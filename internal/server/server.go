// Package server is the local HTTP API plus the embedded web UI.
package server

import (
	"context"
	"embed"
	"encoding/json"
	"io/fs"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"dsa/internal/ai"
	"dsa/internal/config"
	"dsa/internal/problems"
	"dsa/internal/runner"
	"dsa/internal/store"
	"dsa/internal/testgen"
)

//go:embed web
var webFS embed.FS

const randomCases = 50

type Server struct {
	home  string
	lib   *problems.Library
	store *store.Store

	mu    sync.Mutex
	cfg   config.Config
	hosts map[string]bool
}

func New(home string, lib *problems.Library) (*Server, error) {
	st, err := store.New(home)
	if err != nil {
		return nil, err
	}
	cfg, err := config.Load(home)
	if err != nil {
		cfg = config.Default() // a damaged config file should not stop practice
	}
	return &Server{home: home, lib: lib, store: st, cfg: cfg, hosts: map[string]bool{}}, nil
}

// SetAddr records the address we listen on. Requests for any other Host are
// refused, which blocks DNS-rebinding attacks from web pages.
func (s *Server) SetAddr(addr string) {
	_, port, _ := net.SplitHostPort(addr)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hosts = map[string]bool{
		"127.0.0.1:" + port: true,
		"localhost:" + port: true,
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/problems", s.listProblems)
	mux.HandleFunc("GET /api/problems/{id}", s.getProblem)
	mux.HandleFunc("PUT /api/solution", s.saveSolution)
	mux.HandleFunc("POST /api/run", s.run)
	mux.HandleFunc("GET /api/doctor", s.doctor)
	mux.HandleFunc("GET /api/config", s.getConfig)
	mux.HandleFunc("PUT /api/config", s.putConfig)
	mux.HandleFunc("POST /api/ai", s.askAI)

	sub, err := fs.Sub(webFS, "web")
	if err != nil {
		panic(err)
	}
	mux.Handle("/", http.FileServer(http.FS(sub)))
	return s.guard(mux)
}

// guard protects a server that runs code on this computer:
//   - the Host header must be our own loopback address
//   - anything that changes state must carry X-DSA, which a web page on another
//     site cannot add without our permission
func (s *Server) guard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		okHost := s.hosts[r.Host]
		s.mu.Unlock()
		if !okHost {
			http.Error(w, "forbidden host", http.StatusForbidden)
			return
		}
		if r.Method != http.MethodGet && r.Method != http.MethodHead && r.Header.Get("X-DSA") != "1" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/api/") {
			w.Header().Set("Cache-Control", "no-store")
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		writeErr(w, http.StatusBadRequest, "bad request: "+err.Error())
		return false
	}
	return true
}

func validLang(l string) bool { return l == "java" || l == "cpp" }

// ---- problems ----

type listItem struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Difficulty string   `json:"difficulty"`
	Sheet      string   `json:"sheet"`
	Group      string   `json:"group"`
	Order      int      `json:"order"`
	Tags       []string `json:"tags"`
	Status     string   `json:"status"` // "", tried, solved
}

func (s *Server) listProblems(w http.ResponseWriter, _ *http.Request) {
	progress := s.store.Progress()
	items := []listItem{}
	for _, p := range s.lib.List() {
		items = append(items, listItem{
			ID: p.ID, Title: p.Title, Difficulty: p.Difficulty, Sheet: p.Sheet,
			Group: p.Group, Order: p.Order, Tags: p.Tags, Status: progress[p.ID],
		})
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) getProblem(w http.ResponseWriter, r *http.Request) {
	p, ok := s.lib.Get(r.PathValue("id"))
	if !ok {
		writeErr(w, http.StatusNotFound, "no such problem")
		return
	}
	lang := r.URL.Query().Get("lang")
	if !validLang(lang) {
		lang = "java"
	}
	code, saved, err := s.store.Load(p.ID, lang)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !saved {
		code, err = p.File("starter." + lang)
		if err != nil {
			writeErr(w, http.StatusInternalServerError, "this problem has no "+lang+" starter")
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"id": p.ID, "title": p.Title, "difficulty": p.Difficulty, "sheet": p.Sheet,
		"group": p.Group, "tags": p.Tags, "buildsOn": p.BuildsOn, "leadsTo": p.LeadsTo,
		"statement": p.Statement, "examples": p.Examples,
		"lang": lang, "code": code, "saved": saved,
	})
}

func (s *Server) saveSolution(w http.ResponseWriter, r *http.Request) {
	var req struct{ Problem, Lang, Code string }
	if !decode(w, r, &req) {
		return
	}
	if _, ok := s.lib.Get(req.Problem); !ok || !validLang(req.Lang) {
		writeErr(w, http.StatusBadRequest, "unknown problem or language")
		return
	}
	if err := s.store.Save(req.Problem, req.Lang, req.Code); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"saved": true})
}

// ---- running code ----

func (s *Server) run(w http.ResponseWriter, r *http.Request) {
	var req struct{ Problem, Lang, Code, Mode string }
	if !decode(w, r, &req) {
		return
	}
	p, ok := s.lib.Get(req.Problem)
	if !ok || !validLang(req.Lang) || (req.Mode != "examples" && req.Mode != "submit") {
		writeErr(w, http.StatusBadRequest, "unknown problem, language or mode")
		return
	}
	driver, err := p.File("driver." + req.Lang)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "this problem has no "+req.Lang+" driver")
		return
	}
	_ = s.store.Save(p.ID, req.Lang, req.Code) // running always saves

	var cases []runner.Case
	for _, c := range p.Examples {
		cases = append(cases, runner.Case{Kind: "example", Label: c.Label, Input: c.Input, Expected: c.Expected})
	}
	var seed int64
	if req.Mode == "submit" {
		for _, c := range p.Edge {
			cases = append(cases, runner.Case{Kind: "edge", Label: c.Label, Input: c.Input, Expected: c.Expected})
		}
		if testgen.Has(p.Generator) {
			seed = time.Now().UnixNano()
			gen, err := testgen.Generate(p.Generator, seed, randomCases)
			if err != nil {
				writeErr(w, http.StatusInternalServerError, err.Error())
				return
			}
			for i, c := range gen {
				cases = append(cases, runner.Case{Kind: "random", Label: "Random " + strconv.Itoa(i+1), Input: c.Input, Expected: c.Expected})
			}
		}
	}

	rep := runner.Run(r.Context(), runner.Request{
		Lang: req.Lang, Code: req.Code, Driver: driver, Cases: cases,
		CompileTimeout: 30 * time.Second, RunTimeout: 10 * time.Second,
	})

	if rep.Status == "ok" || rep.Status == "runtime_error" || rep.Status == "timeout" {
		if req.Mode == "submit" && rep.Status == "ok" && rep.Passed == rep.Total {
			_ = s.store.Mark(p.ID, "solved")
		} else {
			_ = s.store.Mark(p.ID, "tried")
		}
	}
	writeJSON(w, http.StatusOK, struct {
		runner.Report
		Seed int64 `json:"seed,omitempty"`
	}{rep, seed})
}

func (s *Server) doctor(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"tools":        runner.Doctor(),
		"solutionsDir": s.store.SolutionsDir(),
	})
}

// ---- settings ----

type aiView struct {
	Provider    string `json:"provider"`
	Model       string `json:"model"`
	HasKey      bool   `json:"hasKey"`
	ExplainOnly bool   `json:"explainOnly"`
}

type configView struct {
	Theme  string `json:"theme"`
	Accent string `json:"accent"`
	Lang   string `json:"lang"`
	AI     aiView `json:"ai"`
}

func view(c config.Config) configView {
	return configView{
		Theme: c.Theme, Accent: c.Accent, Lang: c.Lang,
		AI: aiView{Provider: c.AI.Provider, Model: c.AI.Model, HasKey: c.AI.APIKey != "", ExplainOnly: c.AI.ExplainOnly},
	}
}

// getConfig never sends the API key back to the browser.
func (s *Server) getConfig(w http.ResponseWriter, _ *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	writeJSON(w, http.StatusOK, view(s.cfg))
}

func oneOf(v string, allowed ...string) bool {
	for _, a := range allowed {
		if v == a {
			return true
		}
	}
	return false
}

func (s *Server) putConfig(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Theme  string `json:"theme"`
		Accent string `json:"accent"`
		Lang   string `json:"lang"`
		AI     struct {
			Provider    string `json:"provider"`
			Model       string `json:"model"`
			APIKey      string `json:"apiKey"` // empty means keep the saved key
			ClearKey    bool   `json:"clearKey"`
			ExplainOnly bool   `json:"explainOnly"`
		} `json:"ai"`
	}
	if !decode(w, r, &req) {
		return
	}
	if !oneOf(req.Theme, "paper", "ink", "system") ||
		!oneOf(req.Accent, "brick", "forest", "amber", "slate", "plum") ||
		!validLang(req.Lang) ||
		!oneOf(req.AI.Provider, "anthropic", "openai", "ollama") {
		writeErr(w, http.StatusBadRequest, "invalid setting")
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cfg.Theme, s.cfg.Accent, s.cfg.Lang = req.Theme, req.Accent, req.Lang
	s.cfg.AI.Provider = req.AI.Provider
	s.cfg.AI.Model = strings.TrimSpace(req.AI.Model)
	s.cfg.AI.ExplainOnly = req.AI.ExplainOnly
	if req.AI.ClearKey {
		s.cfg.AI.APIKey = ""
	} else if k := strings.TrimSpace(req.AI.APIKey); k != "" {
		s.cfg.AI.APIKey = k
	}
	if err := config.Save(s.home, s.cfg); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, view(s.cfg))
}

// ---- AI help ----

func (s *Server) askAI(w http.ResponseWriter, r *http.Request) {
	var req struct{ Problem, Lang, Code, Question, LastReport string }
	if !decode(w, r, &req) {
		return
	}
	p, ok := s.lib.Get(req.Problem)
	if !ok {
		writeErr(w, http.StatusBadRequest, "unknown problem")
		return
	}
	s.mu.Lock()
	c := s.cfg.AI
	s.mu.Unlock()
	if c.APIKey == "" && c.Provider != "ollama" {
		writeErr(w, http.StatusBadRequest, "Add your API key in Settings to use AI help.")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 90*time.Second)
	defer cancel()
	text, err := ai.Ask(ctx, ai.Request{
		Provider: c.Provider, Model: c.Model, APIKey: c.APIKey, ExplainOnly: c.ExplainOnly,
		Statement: p.Statement, Language: req.Lang, Code: req.Code,
		LastReport: req.LastReport, Question: req.Question,
	})
	if err != nil {
		writeErr(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"answer": text})
}
