package store

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type zipEntry struct {
	name string
	data string
	mode os.FileMode
}

func makeZip(t *testing.T, entries ...zipEntry) *bytes.Reader {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for _, e := range entries {
		h := &zip.FileHeader{Name: e.name, Method: zip.Deflate}
		if e.mode != 0 {
			h.SetMode(e.mode)
		}
		w, err := zw.CreateHeader(h)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write([]byte(e.data)); err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return bytes.NewReader(buf.Bytes())
}

func importZip(t *testing.T, s *Store, r *bytes.Reader, overwrite bool) ImportReport {
	t.Helper()
	rep, err := s.ImportZip(r, r.Size(), overwrite)
	if err != nil {
		t.Fatal(err)
	}
	return rep
}

func TestExportImportRoundTrip(t *testing.T) {
	src, _ := newClocked(t)
	_ = src.Save("two-sum", "java", "class Solution {}")
	_ = src.Save("two-sum", "cpp", "int x;")
	_ = src.Save("trees", "java", "class Solution { /* trees */ }")
	_ = src.Mark("two-sum", "solved")
	_ = src.Mark("trees", "tried")

	var buf bytes.Buffer
	if err := src.ExportZip(&buf); err != nil {
		t.Fatal(err)
	}
	zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, f := range zr.File {
		names = append(names, f.Name)
	}
	if got := strings.Join(names, ","); got != "README.txt,solutions/trees.java,solutions/two-sum.cpp,solutions/two-sum.java,progress.json" {
		t.Fatalf("export holds %s", got)
	}

	dst, _ := newClocked(t)
	rep := importZip(t, dst, bytes.NewReader(buf.Bytes()), false)
	if len(rep.Imported) != 3 || len(rep.Skipped) != 0 || len(rep.Rejected) != 0 || rep.ProgressMerged != 2 {
		t.Fatalf("report: %+v", rep)
	}
	code, ok, _ := dst.Load("trees", "java")
	if !ok || code != "class Solution { /* trees */ }" {
		t.Fatalf("trees.java = %q, %v", code, ok)
	}
	if p := dst.Progress(); p["two-sum"].Status != "solved" || p["trees"].Status != "tried" {
		t.Fatalf("progress: %+v", p)
	}
}

func TestImportSkipsExistingUnlessOverwrite(t *testing.T) {
	s, _ := newClocked(t)
	_ = s.Save("two-sum", "java", "mine")
	zipped := func() *bytes.Reader {
		return makeZip(t, zipEntry{name: "solutions/two-sum.java", data: "theirs"}, zipEntry{name: "solutions/new-one.cpp", data: "new"})
	}

	rep := importZip(t, s, zipped(), false)
	if strings.Join(rep.Skipped, ",") != "solutions/two-sum.java" || strings.Join(rep.Imported, ",") != "solutions/new-one.cpp" {
		t.Fatalf("report: %+v", rep)
	}
	if code, _, _ := s.Load("two-sum", "java"); code != "mine" {
		t.Fatal("an existing file was replaced without overwrite")
	}

	rep = importZip(t, s, zipped(), true)
	if len(rep.Imported) != 2 || len(rep.Skipped) != 0 {
		t.Fatalf("report with overwrite: %+v", rep)
	}
	if code, _, _ := s.Load("two-sum", "java"); code != "theirs" {
		t.Fatal("overwrite did not replace the file")
	}
}

func TestImportRejectsUnsafeEntries(t *testing.T) {
	parent := t.TempDir()
	home := filepath.Join(parent, "home")
	s, err := New(home)
	if err != nil {
		t.Fatal(err)
	}
	r := makeZip(t,
		zipEntry{name: "../evil.java", data: "x"},
		zipEntry{name: "solutions/../../evil.java", data: "x"},
		zipEntry{name: "/abs.java", data: "x"},
		zipEntry{name: `solutions\back.java`, data: "x"},
		zipEntry{name: "solutions/UPPER.java", data: "x"},
		zipEntry{name: "solutions/notes.txt", data: "x"},
		zipEntry{name: "solutions/sub/deep.java", data: "x"},
		zipEntry{name: "solutions/link.java", data: "/etc/passwd", mode: os.ModeSymlink | 0o777},
		zipEntry{name: "solutions/huge.java", data: strings.Repeat("a", maxImportFile+1)},
		zipEntry{name: "config.json", data: `{"ai":{"apiKey":"x"}}`},
		zipEntry{name: "solutions/fine.java", data: "ok"},
	)
	rep := importZip(t, s, r, true)
	if len(rep.Rejected) != 10 || strings.Join(rep.Imported, ",") != "solutions/fine.java" {
		t.Fatalf("report: %+v", rep)
	}
	for _, rj := range rep.Rejected {
		if rj.Reason == "" {
			t.Errorf("%s rejected without a reason", rj.Name)
		}
	}

	var written []string
	_ = filepath.WalkDir(parent, func(path string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			rel, _ := filepath.Rel(parent, path)
			written = append(written, rel)
		}
		return nil
	})
	if got := strings.Join(written, ","); got != filepath.Join("home", "solutions", "fine.java") {
		t.Fatalf("files on disk after import: %s", got)
	}
}

func TestImportLimits(t *testing.T) {
	s, _ := newClocked(t)
	var many []zipEntry
	for i := 0; i <= maxImportEntries; i++ {
		many = append(many, zipEntry{name: "solutions/a.java", data: "x"})
	}
	r := makeZip(t, many...)
	if _, err := s.ImportZip(r, r.Size(), false); err == nil || !strings.Contains(err.Error(), "entries") {
		t.Fatalf("too many entries: %v", err)
	}

	var big []zipEntry
	for i := 0; i < 21; i++ {
		big = append(big, zipEntry{name: "solutions/p" + string(rune('a'+i)) + ".java", data: strings.Repeat("b", maxImportFile)})
	}
	r = makeZip(t, big...)
	if _, err := s.ImportZip(r, r.Size(), false); err == nil || !strings.Contains(err.Error(), "MB") {
		t.Fatalf("too large in total: %v", err)
	}
	if entries, _ := os.ReadDir(s.SolutionsDir()); len(entries) != 0 {
		t.Fatalf("a rejected zip wrote %d files", len(entries))
	}

	notZip := bytes.NewReader([]byte("hello"))
	if _, err := s.ImportZip(notZip, notZip.Size(), false); err == nil {
		t.Fatal("a non-zip was accepted")
	}
}

func TestImportMergesProgress(t *testing.T) {
	s, c := newClocked(t)
	_ = s.Mark("kept-solved", "solved")
	_ = s.Mark("older-tried", "tried")
	later := c.t.Add(time.Hour).Format(time.RFC3339)
	earlier := c.t.Add(-time.Hour).Format(time.RFC3339)
	progress := `{
		"kept-solved": {"status": "tried", "updatedAt": "` + later + `"},
		"older-tried": {"status": "tried", "updatedAt": "` + later + `"},
		"new-one": {"status": "solved", "updatedAt": "` + earlier + `"},
		"old-format": "tried",
		"BAD ID": {"status": "solved"},
		"bad-status": {"status": "great"}
	}`
	rep := importZip(t, s, makeZip(t, zipEntry{name: "progress.json", data: progress}), false)
	if rep.ProgressMerged != 3 {
		t.Fatalf("merged %d entries, want 3 (older-tried, new-one, old-format)", rep.ProgressMerged)
	}
	p := s.Progress()
	if p["kept-solved"].Status != "solved" {
		t.Error("import downgraded a solved problem")
	}
	if !p["older-tried"].UpdatedAt.After(c.t) || p["new-one"].Status != "solved" || p["old-format"].Status != "tried" {
		t.Errorf("progress after merge: %+v", p)
	}
	if _, ok := p["BAD ID"]; ok {
		t.Error("an invalid id was imported")
	}
	if _, ok := p["bad-status"]; ok {
		t.Error("an invalid status was imported")
	}

	rep = importZip(t, s, makeZip(t, zipEntry{name: "progress.json", data: "not json"}), false)
	if len(rep.Rejected) != 1 {
		t.Fatalf("bad progress.json: %+v", rep)
	}
}
