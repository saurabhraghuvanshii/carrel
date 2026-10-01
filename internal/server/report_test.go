package server

import (
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

type gridReport struct {
	Status  string
	Passed  int
	Total   int
	Seed    int64
	Results []struct {
		Kind, Label, Input, Expected, Got string
		Passed                            *bool
	}
}

func submit(t *testing.T, h http.Handler, problem, code string) (gridReport, int) {
	t.Helper()
	body, _ := json.Marshal(map[string]string{"problem": problem, "lang": "java", "code": code, "mode": "submit"})
	rec := do(h, "POST", "http://127.0.0.1:7777/api/run", string(body), map[string]string{"X-Carrel": "1"})
	if rec.Code != http.StatusOK {
		t.Fatalf("run: %d %s", rec.Code, rec.Body.String())
	}
	var rep gridReport
	if err := json.Unmarshal(rec.Body.Bytes(), &rep); err != nil {
		t.Fatal(err)
	}
	return rep, rec.Body.Len()
}

// The test grid draws one square per case and opens any of them, so every
// result needs kind, label, passed, input, expected and got.
func TestSubmitReportHasWhatTheGridNeeds(t *testing.T) {
	for _, tool := range []string{"javac", "java"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skipf("%s is not installed", tool)
		}
	}
	ref, err := os.ReadFile(filepath.Join("..", "..", "tools", "refs", "pair-with-target-sum", "Solution.java"))
	if err != nil {
		t.Fatal(err)
	}
	h := newTestServer(t)

	rep, size := submit(t, h, "pair-with-target-sum", string(ref))
	if rep.Status != "ok" || rep.Passed != rep.Total || rep.Total != len(rep.Results) || rep.Total <= randomCases {
		t.Fatalf("status %s, passed %d of %d, %d results", rep.Status, rep.Passed, rep.Total, len(rep.Results))
	}
	if rep.Seed < 1 || rep.Seed > seedRange {
		t.Fatalf("seed %d is outside 1..%d", rep.Seed, seedRange)
	}
	if size > 1<<20 {
		t.Fatalf("the report is %d bytes, over 1 MB", size)
	}
	random := false
	for i, r := range rep.Results {
		if r.Kind != "example" && r.Kind != "edge" && r.Kind != "random" {
			t.Fatalf("case %d: kind %q", i+1, r.Kind)
		}
		if random && r.Kind != "random" {
			t.Fatalf("case %d: a fixed case comes after a random one", i+1)
		}
		random = r.Kind == "random"
		if r.Label == "" || r.Passed == nil || r.Input == "" || r.Expected == "" || r.Got == "" {
			t.Fatalf("case %d is missing a field: %+v", i+1, r)
		}
	}
	if !random {
		t.Fatal("the report has no random cases")
	}
}

func TestFocusLayoutIsSaved(t *testing.T) {
	h := newTestServer(t)
	put := `{"theme":"paper","accent":"brick","lang":"java","focusLayout":true,"ai":{"provider":"anthropic"}}`
	rec := do(h, "PUT", "http://127.0.0.1:7777/api/config", put, map[string]string{"X-Carrel": "1"})
	if rec.Code != http.StatusOK {
		t.Fatalf("put config: %d %s", rec.Code, rec.Body.String())
	}
	var got struct{ FocusLayout bool }
	if code := getJSON(t, h, "/api/config", &got); code != http.StatusOK || !got.FocusLayout {
		t.Fatalf("focusLayout not saved: %d %+v", code, got)
	}
}
