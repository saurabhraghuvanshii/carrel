package server

import (
	"encoding/json"
	"net/http"
	"testing"
	"testing/fstest"
	"time"

	"carrel/internal/problems"
)

// pathLib: patterns a -> c, b on its own; real r1 builds on c.
func pathLib(t *testing.T) *Server {
	t.Helper()
	fsys := fstest.MapFS{}
	add := func(id, sheet, difficulty string, order int, buildsOn, leadsTo []string) {
		m := problems.Meta{ID: id, Title: "Title " + id, Difficulty: difficulty, Sheet: sheet, Group: "G",
			Order: order, BuildsOn: buildsOn, LeadsTo: leadsTo, Generator: "pair-sum"}
		b, _ := json.Marshal(m)
		fsys[id+"/meta.json"] = &fstest.MapFile{Data: b}
		fsys[id+"/statement.md"] = &fstest.MapFile{Data: []byte("x")}
		fsys[id+"/tests.json"] = &fstest.MapFile{Data: []byte(`{"examples":[],"edge":[]}`)}
	}
	add("a", "patterns", "easy", 101, nil, []string{"c"})
	add("b", "patterns", "easy", 102, nil, nil)
	add("c", "patterns", "medium", 103, []string{"a"}, nil)
	add("r1", "real", "hard", 101, []string{"c"}, nil)
	lib, err := problems.Load(fsys)
	if err != nil {
		t.Fatal(err)
	}
	s, err := New(t.TempDir(), lib)
	if err != nil {
		t.Fatal(err)
	}
	s.SetAddr("127.0.0.1:7777")
	return s
}

func getJSON(t *testing.T, h http.Handler, url string, v any) int {
	t.Helper()
	rec := do(h, "GET", "http://127.0.0.1:7777"+url, "", nil)
	if rec.Code == http.StatusOK && v != nil {
		if err := json.Unmarshal(rec.Body.Bytes(), v); err != nil {
			t.Fatalf("%s: %v", url, err)
		}
	}
	return rec.Code
}

func wantNext(t *testing.T, h http.Handler, url, id, reason string) {
	t.Helper()
	var got struct{ ID, Reason string }
	if code := getJSON(t, h, url, &got); code != http.StatusOK {
		t.Fatalf("%s: status %d", url, code)
	}
	if got.ID != id || got.Reason != reason {
		t.Fatalf("%s: got %s (%s), want %s (%s)", url, got.ID, got.Reason, id, reason)
	}
}

func TestNextSuggestsInLearningOrder(t *testing.T) {
	s := pathLib(t)
	h := s.Handler()

	wantNext(t, h, "/api/next?sheet=patterns", "a", "ready")
	wantNext(t, h, "/api/next?sheet=real", "r1", "any") // c is not solved, so r1 is locked

	_ = s.store.Mark("b", "tried")
	wantNext(t, h, "/api/next?sheet=patterns", "b", "continue")

	_ = s.store.Mark("a", "solved")
	wantNext(t, h, "/api/next?sheet=patterns&after=a", "c", "leads")
	wantNext(t, h, "/api/next?sheet=patterns", "b", "continue")

	_ = s.store.Mark("c", "tried")
	wantNext(t, h, "/api/next?sheet=patterns", "c", "continue") // the most recent tried wins

	_ = s.store.Mark("b", "solved")
	_ = s.store.Mark("c", "solved")
	wantNext(t, h, "/api/next?sheet=real", "r1", "ready")
	if code := getJSON(t, h, "/api/next?sheet=patterns", nil); code != http.StatusNotFound {
		t.Fatalf("all solved: status %d, want 404", code)
	}
	if code := getJSON(t, h, "/api/next?sheet=other", nil); code != http.StatusBadRequest {
		t.Fatalf("bad sheet: status %d, want 400", code)
	}
}

func TestProgressCounts(t *testing.T) {
	s := pathLib(t)
	_ = s.store.Mark("a", "solved")
	_ = s.store.Mark("c", "tried")
	var got map[string]struct {
		Solved, Tried, Total int
		ByDifficulty         map[string]struct{ Solved, Tried, Total int }
	}
	if code := getJSON(t, s.Handler(), "/api/progress", &got); code != http.StatusOK {
		t.Fatalf("status %d", code)
	}
	p, r := got["patterns"], got["real"]
	if p.Solved != 1 || p.Tried != 1 || p.Total != 3 || r.Total != 1 || r.Solved != 0 {
		t.Fatalf("sheet counts: %+v %+v", p, r)
	}
	if e := p.ByDifficulty["easy"]; e.Solved != 1 || e.Total != 2 {
		t.Fatalf("easy: %+v", e)
	}
	if m := p.ByDifficulty["medium"]; m.Tried != 1 || m.Total != 1 {
		t.Fatalf("medium: %+v", m)
	}
	if hd := r.ByDifficulty["hard"]; hd.Total != 1 {
		t.Fatalf("hard: %+v", hd)
	}
}

type item struct {
	ID     string
	Status string
	Ready  bool
	Needs  []string
	Review bool
}

func listByID(t *testing.T, h http.Handler) map[string]item {
	t.Helper()
	var items []item
	getJSON(t, h, "/api/problems", &items)
	m := map[string]item{}
	for _, it := range items {
		m[it.ID] = it
	}
	return m
}

func TestListShowsReadyNeedsAndReview(t *testing.T) {
	s := pathLib(t)
	h := s.Handler()

	l := listByID(t, h)
	if c := l["c"]; c.Ready || len(c.Needs) != 1 || c.Needs[0] != "a" {
		t.Fatalf("c before a is solved: %+v", c)
	}
	if a := l["a"]; a.Ready || a.Needs != nil {
		t.Fatalf("a has no prerequisites, so it is neither ready nor locked: %+v", a)
	}

	_ = s.store.Mark("a", "solved")
	if c := listByID(t, h)["c"]; !c.Ready || c.Needs != nil {
		t.Fatalf("c after a is solved: %+v", c)
	}

	s.now = func() time.Time { return time.Now().Add(4 * 24 * time.Hour) }
	if listByID(t, h)["a"].Review {
		t.Fatal("review shown while reminders are off")
	}
	put := `{"theme":"paper","accent":"brick","lang":"java","remindReviews":true,"ai":{"provider":"anthropic"}}`
	if rec := do(h, "PUT", "http://127.0.0.1:7777/api/config", put, map[string]string{"X-Carrel": "1"}); rec.Code != 200 {
		t.Fatalf("put config: %d %s", rec.Code, rec.Body.String())
	}
	if !listByID(t, h)["a"].Review {
		t.Fatal("a is due after 3 days but not marked for review")
	}
	s.now = time.Now
	if listByID(t, h)["a"].Review {
		t.Fatal("a is marked for review before it is due")
	}
}
