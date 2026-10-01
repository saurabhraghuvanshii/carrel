package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

const testKey = "sk-test-key-0123456789"

// fake is a provider that records each request and answers with reply(n),
// where n counts calls from 1.
type fake struct {
	mu       sync.Mutex
	headers  []http.Header
	bodies   []map[string]any
	status   int
	reply    func(n int) string
	delay    time.Duration
	provider string
}

func newFake(t *testing.T, provider string, reply func(n int) string) *fake {
	t.Helper()
	f := &fake{status: 200, reply: reply, provider: provider}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		f.mu.Lock()
		f.headers = append(f.headers, r.Header.Clone())
		f.bodies = append(f.bodies, body)
		n := len(f.bodies)
		f.mu.Unlock()
		select {
		case <-time.After(f.delay):
		case <-r.Context().Done():
			return
		}
		w.WriteHeader(f.status)
		_, _ = w.Write([]byte(f.reply(n)))
	}))
	t.Cleanup(srv.Close)
	urls := map[string]*string{"anthropic": &anthropicURL, "openai": &openaiURL, "ollama": &ollamaURL}
	old := *urls[provider]
	*urls[provider] = srv.URL
	t.Cleanup(func() { *urls[provider] = old })
	return f
}

// wrap puts text in each provider's success shape.
func wrap(provider, text string) string {
	t, _ := json.Marshal(text)
	switch provider {
	case "anthropic":
		return `{"content":[{"type":"text","text":` + string(t) + `}]}`
	case "openai":
		return `{"output":[{"type":"reasoning","summary":[]},{"type":"message","content":[{"type":"output_text","text":` + string(t) + `}]}]}`
	}
	return `{"message":{"role":"assistant","content":` + string(t) + `}}`
}

func request(provider string) Request {
	return Request{Provider: provider, Model: "test-model", APIKey: testKey, ExplainOnly: true,
		Statement: "Statement X", Language: "java", Code: "class Solution {}", Question: "Why does it fail?"}
}

var providers = []string{"anthropic", "openai", "ollama"}

func TestRequestShapesAndGoodReplies(t *testing.T) {
	for _, p := range providers {
		f := newFake(t, p, func(int) string { return wrap(p, "Try a hash map.") })
		got, err := Ask(context.Background(), request(p))
		if err != nil || got != "Try a hash map." {
			t.Fatalf("%s: %q, %v", p, got, err)
		}
		h, b := f.headers[0], f.bodies[0]
		if b["model"] != "test-model" || h.Get("Content-Type") != "application/json" {
			t.Fatalf("%s: model or content type wrong: %v %v", p, b["model"], h.Get("Content-Type"))
		}
		var system, firstUser string
		switch p {
		case "anthropic":
			if h.Get("Authorization") != "Bearer "+testKey || h.Get("anthropic-version") != "2023-06-01" || b["max_tokens"] != float64(1024) {
				t.Fatalf("anthropic headers or max_tokens wrong: %v %v", h, b["max_tokens"])
			}
			system = b["system"].(string)
			firstUser = b["messages"].([]any)[0].(map[string]any)["content"].(string)
		case "openai":
			if h.Get("Authorization") != "Bearer "+testKey {
				t.Fatal("openai: missing bearer key")
			}
			system = b["instructions"].(string)
			firstUser = b["input"].([]any)[0].(map[string]any)["content"].(string)
		case "ollama":
			if h.Get("Authorization") != "" || b["stream"] != false {
				t.Fatal("ollama: should send no key and stream false")
			}
			msgs := b["messages"].([]any)
			system = msgs[0].(map[string]any)["content"].(string)
			firstUser = msgs[1].(map[string]any)["content"].(string)
		}
		if !strings.Contains(system, "Never write the full solution") {
			t.Fatalf("%s: explain-only rule missing from the system prompt", p)
		}
		for _, want := range []string{"Statement X", "class Solution {}", "Why does it fail?"} {
			if !strings.Contains(firstUser, want) {
				t.Fatalf("%s: user message lacks %q", p, want)
			}
		}
	}
}

func TestErrorsNeverContainTheKey(t *testing.T) {
	for _, p := range providers {
		f := newFake(t, p, func(int) string {
			return `{"type":"error","error":{"type":"authentication_error","message":"invalid key ` + testKey + `"}}`
		})
		f.status = 401
		_, err := Ask(context.Background(), request(p))
		if err == nil || !strings.Contains(err.Error(), "401") || !strings.Contains(err.Error(), "invalid key") {
			t.Fatalf("%s: %v", p, err)
		}
		if strings.Contains(err.Error(), testKey) {
			t.Fatalf("%s: the error leaked the key: %v", p, err)
		}
	}
}

func TestMalformedAndEmptyReplies(t *testing.T) {
	for _, p := range providers {
		newFake(t, p, func(int) string { return "<html>oops" })
		if _, err := Ask(context.Background(), request(p)); err == nil || !strings.Contains(err.Error(), "could not be read") {
			t.Fatalf("%s malformed: %v", p, err)
		}
		newFake(t, p, func(int) string { return "{}" })
		if _, err := Ask(context.Background(), request(p)); err == nil || !strings.Contains(err.Error(), "no answer") {
			t.Fatalf("%s empty: %v", p, err)
		}
	}
}

func TestTimeout(t *testing.T) {
	for _, p := range providers {
		f := newFake(t, p, func(int) string { return wrap(p, "late") })
		f.delay = 2 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		start := time.Now()
		_, err := Ask(ctx, request(p))
		cancel()
		if err == nil || !strings.Contains(err.Error(), "did not answer in time") || time.Since(start) > time.Second {
			t.Fatalf("%s: %v after %v", p, err, time.Since(start))
		}
	}
}

func longSolution() string {
	var b strings.Builder
	b.WriteString("Here is the whole thing:\n```java\n")
	for i := 0; i < 40; i++ {
		fmt.Fprintf(&b, "int line%d = %d;\n", i, i)
	}
	b.WriteString("```\nThat should pass.")
	return b.String()
}

func TestExplainOnlyIsEnforced(t *testing.T) {
	for _, p := range providers {
		f := newFake(t, p, func(int) string { return wrap(p, longSolution()) })
		got, err := Ask(context.Background(), request(p))
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(got, "line20") || !strings.Contains(got, codeRemoved) || !strings.Contains(got, "That should pass.") {
			t.Fatalf("%s: long code reached the learner:\n%s", p, got)
		}
		if len(f.bodies) != 2 {
			t.Fatalf("%s: want one retry, got %d calls", p, len(f.bodies))
		}
		if retry, _ := json.Marshal(f.bodies[1]); !strings.Contains(string(retry), "too much code") {
			t.Fatalf("%s: the retry did not ask for words instead of code", p)
		}

		f = newFake(t, p, func(int) string { return wrap(p, longSolution()) })
		r := request(p)
		r.ExplainOnly = false
		got, _ = Ask(context.Background(), r)
		if got != longSolution() || len(f.bodies) != 1 {
			t.Fatalf("%s: with Explain only off the full reply should pass through once", p)
		}
		if sys, _ := json.Marshal(f.bodies[0]); strings.Contains(string(sys), "Never write the full solution") {
			t.Fatalf("%s: explain-only rule sent while it is off", p)
		}
	}
}

func TestExplainOnlyKeepsAGoodRetry(t *testing.T) {
	newFake(t, "anthropic", func(n int) string {
		if n == 1 {
			return wrap("anthropic", longSolution())
		}
		return wrap("anthropic", "Think about what you have already seen.")
	})
	got, _ := Ask(context.Background(), request("anthropic"))
	if got != "Think about what you have already seen." {
		t.Fatalf("got %q", got)
	}
}

func TestTrimLongCode(t *testing.T) {
	block := func(fence string, lines int, closed bool) string {
		s := fence + "\n" + strings.Repeat("x\n", lines)
		if closed {
			s += fence
		}
		return "before\n" + s
	}
	cases := []struct {
		text    string
		removed bool
	}{
		{block("```", 12, true), false},
		{block("```", 13, true), true},
		{block("~~~", 13, true), true},
		{block("```", 13, false), true},
		{block("```", 3, false), false},
		{"no code at all", false},
	}
	for i, c := range cases {
		got, found := trimLongCode(c.text)
		if found != c.removed || strings.Contains(got, codeRemoved) != c.removed {
			t.Fatalf("case %d: found=%v\n%s", i, found, got)
		}
	}
}

func TestWhatIsSentIsCut(t *testing.T) {
	f := newFake(t, "anthropic", func(int) string { return wrap("anthropic", "ok") })
	r := request("anthropic")
	r.Code = strings.Repeat("c", 30000)
	r.Question = strings.Repeat("q", 5000)
	r.LastReport = strings.Repeat("r", 5000)
	if _, err := Ask(context.Background(), r); err != nil {
		t.Fatal(err)
	}
	sent := f.bodies[0]["messages"].([]any)[0].(map[string]any)["content"].(string)
	for letter, max := range map[string]int{"c": maxCode, "q": maxQuestion, "r": maxReport} {
		if n := strings.Count(sent, letter); n > max+20 {
			t.Fatalf("%q sent %d times, limit %d", letter, n, max)
		}
	}
}

func TestPingAndMissingModel(t *testing.T) {
	f := newFake(t, "openai", func(int) string { return `{"output":[{"type":"reasoning"}]}` })
	if _, err := Ask(context.Background(), Request{Provider: "openai", Model: "m", APIKey: testKey, Ping: true}); err != nil {
		t.Fatalf("a ping that reached the model should succeed: %v", err)
	}
	if b, _ := json.Marshal(f.bodies[0]); strings.Contains(string(b), "Statement") {
		t.Fatal("a ping should not send the problem")
	}
	if _, err := Ask(context.Background(), Request{Provider: "openai", APIKey: testKey}); err == nil || !strings.Contains(err.Error(), "model") {
		t.Fatalf("missing model: %v", err)
	}
}
