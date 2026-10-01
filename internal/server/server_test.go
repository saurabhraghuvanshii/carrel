package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"carrel/internal/problems"
)

func builtinServer(t *testing.T) *Server {
	t.Helper()
	lib, err := problems.Load(problems.Builtin())
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

func newTestServer(t *testing.T) http.Handler { return builtinServer(t).Handler() }

func do(h http.Handler, method, url, body string, header map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, url, strings.NewReader(body))
	for k, v := range header {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestGuardBlocksForeignHostAndMissingHeader(t *testing.T) {
	h := newTestServer(t)

	if rec := do(h, "GET", "http://127.0.0.1:7777/api/problems", "", nil); rec.Code != 200 {
		t.Fatalf("own host should work, got %d", rec.Code)
	}
	if rec := do(h, "GET", "http://evil.example:7777/api/problems", "", nil); rec.Code != 403 {
		t.Fatalf("foreign host should be refused, got %d", rec.Code)
	}
	body := `{"problem":"pair-with-target-sum","lang":"java","code":"x"}`
	if rec := do(h, "PUT", "http://127.0.0.1:7777/api/solution", body, nil); rec.Code != 403 {
		t.Fatalf("a write without X-Carrel should be refused, got %d", rec.Code)
	}
	if rec := do(h, "PUT", "http://127.0.0.1:7777/api/solution", body, map[string]string{"X-Carrel": "1"}); rec.Code != 200 {
		t.Fatalf("a write with X-Carrel should work, got %d", rec.Code)
	}
}

func TestConfigNeverReturnsTheKey(t *testing.T) {
	h := newTestServer(t)
	put := `{"theme":"ink","accent":"forest","lang":"cpp","ai":{"provider":"anthropic","model":"m","apiKey":"secret-key","explainOnly":true}}`
	rec := do(h, "PUT", "http://127.0.0.1:7777/api/config", put, map[string]string{"X-Carrel": "1"})
	if rec.Code != 200 {
		t.Fatalf("put config: %d %s", rec.Code, rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), "secret-key") {
		t.Fatal("the response leaked the API key")
	}
	get := do(h, "GET", "http://127.0.0.1:7777/api/config", "", nil)
	if strings.Contains(get.Body.String(), "secret-key") || !strings.Contains(get.Body.String(), `"hasKey":true`) {
		t.Fatalf("unexpected config view: %s", get.Body.String())
	}
}
