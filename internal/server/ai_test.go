package server

import (
	"bytes"
	"context"
	"errors"
	"log"
	"os"
	"strings"
	"testing"

	"dsa/internal/ai"
)

const secretKey = "sk-very-secret-0123456789"

func TestAIErrorsAndLogsNeverContainTheKey(t *testing.T) {
	s := builtinServer(t)
	var asked []ai.Request
	s.ask = func(_ context.Context, r ai.Request) (string, error) {
		asked = append(asked, r)
		return "", errors.New("the provider answered 401: invalid key " + r.APIKey)
	}
	var logs bytes.Buffer
	log.SetOutput(&logs)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })
	h := s.Handler()

	put := `{"theme":"paper","accent":"brick","lang":"java","ai":{"provider":"anthropic","model":"m","apiKey":"` + secretKey + `","explainOnly":true}}`
	if rec := do(h, "PUT", "http://127.0.0.1:7777/api/config", put, map[string]string{"X-DSA": "1"}); rec.Code != 200 {
		t.Fatalf("put config: %d", rec.Code)
	}
	ask := `{"problem":"pair-with-target-sum","lang":"java","code":"x","question":"why?"}`
	for _, url := range []string{"/api/ai", "/api/ai/test"} {
		rec := do(h, "POST", "http://127.0.0.1:7777"+url, ask, map[string]string{"X-DSA": "1"})
		if rec.Code != 502 || !strings.Contains(rec.Body.String(), "invalid key [your key]") {
			t.Fatalf("%s: %d %s", url, rec.Code, rec.Body.String())
		}
		if strings.Contains(rec.Body.String(), secretKey) {
			t.Fatalf("%s: the response leaked the key", url)
		}
	}
	if strings.Contains(logs.String(), secretKey) {
		t.Fatal("the log contains the key")
	}
	if len(asked) != 2 || !asked[0].ExplainOnly || asked[0].Statement == "" || !asked[1].Ping {
		t.Fatalf("requests passed to ai: %+v", asked)
	}
}

func TestAIWithoutKeyAsksForOne(t *testing.T) {
	s := builtinServer(t)
	s.ask = func(context.Context, ai.Request) (string, error) {
		t.Fatal("should not call the provider")
		return "", nil
	}
	rec := do(s.Handler(), "POST", "http://127.0.0.1:7777/api/ai/test", "{}", map[string]string{"X-DSA": "1"})
	if rec.Code != 400 || !strings.Contains(rec.Body.String(), "API key") {
		t.Fatalf("%d %s", rec.Code, rec.Body.String())
	}
}
