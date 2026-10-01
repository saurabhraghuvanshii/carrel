package server

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func doBytes(h http.Handler, method, url string, body []byte, header map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, url, bytes.NewReader(body))
	for k, v := range header {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestExportReturnsAZip(t *testing.T) {
	h := newTestServer(t)
	save := `{"problem":"pair-with-target-sum","lang":"java","code":"class Solution {}"}`
	if rec := do(h, "PUT", "http://127.0.0.1:7777/api/solution", save, map[string]string{"X-DSA": "1"}); rec.Code != 200 {
		t.Fatalf("save: %d", rec.Code)
	}

	rec := do(h, "GET", "http://127.0.0.1:7777/api/export", "", nil)
	if rec.Code != 200 || rec.Header().Get("Content-Type") != "application/zip" {
		t.Fatalf("export: %d %s", rec.Code, rec.Header().Get("Content-Type"))
	}
	if cd := rec.Header().Get("Content-Disposition"); !regexp.MustCompile(`^attachment; filename="dsa-solutions-\d{4}-\d{2}-\d{2}\.zip"$`).MatchString(cd) {
		t.Fatalf("Content-Disposition = %q", cd)
	}
	zr, err := zip.NewReader(bytes.NewReader(rec.Body.Bytes()), int64(rec.Body.Len()))
	if err != nil {
		t.Fatalf("not a valid zip: %v", err)
	}
	found := false
	for _, f := range zr.File {
		found = found || f.Name == "solutions/pair-with-target-sum.java"
	}
	if !found {
		t.Fatal("the saved solution is not in the export")
	}
}

func TestImport(t *testing.T) {
	h := newTestServer(t)
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for name, data := range map[string]string{
		"solutions/pair-with-target-sum.java": "class Solution {}",
		"solutions/reverse-list.cpp":          "// mine",
		"../evil.java":                        "x",
	} {
		w, _ := zw.Create(name)
		_, _ = w.Write([]byte(data))
	}
	_ = zw.Close()
	zipBody := buf.Bytes()
	url := "http://127.0.0.1:7777/api/import"
	zipType := map[string]string{"Content-Type": "application/zip"}

	if rec := doBytes(h, "POST", url, zipBody, zipType); rec.Code != http.StatusForbidden {
		t.Fatalf("import without X-DSA: %d, want 403", rec.Code)
	}
	if rec := doBytes(h, "POST", url, zipBody, map[string]string{"X-DSA": "1", "Content-Type": "application/json"}); rec.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("import as JSON: %d, want 415", rec.Code)
	}

	withHeader := map[string]string{"X-DSA": "1", "Content-Type": "application/zip"}
	var rep struct {
		Imported, Skipped []string
		Rejected          []struct{ Name, Reason string }
	}
	rec := doBytes(h, "POST", url, zipBody, withHeader)
	if rec.Code != 200 {
		t.Fatalf("import: %d %s", rec.Code, rec.Body.String())
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &rep)
	if len(rep.Imported) != 2 || len(rep.Skipped) != 0 || len(rep.Rejected) != 1 {
		t.Fatalf("first import: %s", rec.Body.String())
	}

	rec = doBytes(h, "POST", url, zipBody, withHeader)
	_ = json.Unmarshal(rec.Body.Bytes(), &rep)
	if len(rep.Imported) != 0 || len(rep.Skipped) != 2 {
		t.Fatalf("second import should skip both: %s", rec.Body.String())
	}

	rec = doBytes(h, "POST", url+"?overwrite=true", zipBody, withHeader)
	_ = json.Unmarshal(rec.Body.Bytes(), &rep)
	if len(rep.Imported) != 2 || len(rep.Skipped) != 0 {
		t.Fatalf("overwrite import: %s", rec.Body.String())
	}

	if rec := doBytes(h, "POST", url, []byte("not a zip"), withHeader); rec.Code != http.StatusBadRequest {
		t.Fatalf("non-zip body: %d, want 400", rec.Code)
	}
}
