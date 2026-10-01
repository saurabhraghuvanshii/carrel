package store

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	maxImportTotal   = 20 << 20
	maxImportEntries = 2000
	maxImportFile    = 1 << 20
)

const exportReadme = `Solutions exported from Carrel.

solutions/<problem>.java   your Java code for each problem
solutions/<problem>.cpp    your C++ code for each problem
progress.json              which problems you tried or solved

They are plain files. Open them in any editor, or bring them back with
"Import solutions" in Carrel (or: carrel import <this file>).
`

// ExportZip writes every saved solution and the progress file as a zip.
func (s *Store) ExportZip(w io.Writer) error {
	zw := zip.NewWriter(w)
	add := func(name string, data []byte) error {
		f, err := zw.CreateHeader(&zip.FileHeader{Name: name, Method: zip.Deflate, Modified: time.Now()})
		if err != nil {
			return err
		}
		_, err = f.Write(data)
		return err
	}
	if err := add("README.txt", []byte(exportReadme)); err != nil {
		return err
	}

	entries, err := os.ReadDir(s.SolutionsDir())
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	for _, e := range entries {
		if _, _, ok := solutionName(e.Name()); !ok || !e.Type().IsRegular() {
			continue
		}
		b, err := os.ReadFile(filepath.Join(s.SolutionsDir(), e.Name()))
		if err != nil {
			return err
		}
		if err := add("solutions/"+e.Name(), b); err != nil {
			return err
		}
	}

	s.mu.Lock()
	b, err := os.ReadFile(s.progressPath())
	s.mu.Unlock()
	if err == nil {
		if err := add("progress.json", b); err != nil {
			return err
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return zw.Close()
}

// solutionName splits "two-sum.java" into id and language, if both are valid.
func solutionName(name string) (id, lang string, ok bool) {
	for l, ext := range exts {
		if base, found := strings.CutSuffix(name, ext); found && idPattern.MatchString(base) {
			return base, l, true
		}
	}
	return "", "", false
}

type Rejected struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

type ImportReport struct {
	Imported []string   `json:"imported"`
	Skipped  []string   `json:"skipped"`
	Rejected []Rejected `json:"rejected"`
	// ProgressMerged counts progress entries that were added or improved.
	ProgressMerged int `json:"progressMerged"`
}

type pending struct {
	name, id, lang string
	data           []byte
}

// ImportZip reads solutions and progress from a zip made by ExportZip (or by
// hand). Only solutions/<id>.java|cpp and progress.json are accepted; every
// other entry is reported as rejected and never touches the disk. The whole
// zip is read and checked before anything is written.
//
// Existing solutions are skipped unless overwrite is set. Progress is merged:
// solved beats tried, and between equal statuses the newer entry wins.
func (s *Store) ImportZip(r io.ReaderAt, size int64, overwrite bool) (ImportReport, error) {
	rep := ImportReport{Imported: []string{}, Skipped: []string{}, Rejected: []Rejected{}}
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return rep, errors.New("this is not a zip file")
	}
	if len(zr.File) > maxImportEntries {
		return rep, fmt.Errorf("the zip has %d entries; the limit is %d", len(zr.File), maxImportEntries)
	}

	reject := func(name, reason string) { rep.Rejected = append(rep.Rejected, Rejected{name, reason}) }
	var files []pending
	var progress map[string]Entry
	total := 0
	for _, f := range zr.File {
		name := f.Name
		switch {
		case name == "README.txt" || name == "solutions/":
			continue
		case f.Mode()&os.ModeSymlink != 0:
			reject(name, "links are not allowed")
			continue
		case strings.Contains(name, `\`), strings.HasPrefix(name, "/"), strings.Contains(name, ".."):
			reject(name, "unsafe path")
			continue
		case f.Mode().IsDir():
			reject(name, "unexpected folder")
			continue
		}

		var id, lang string
		if name != "progress.json" {
			base, found := strings.CutPrefix(name, "solutions/")
			var ok bool
			if id, lang, ok = solutionName(base); !found || !ok {
				reject(name, "not a solution file (expected solutions/<problem>.java or .cpp)")
				continue
			}
		}

		data, err := readLimited(f, maxImportFile)
		if err != nil {
			reject(name, err.Error())
			continue
		}
		total += len(data)
		if total > maxImportTotal {
			return rep, fmt.Errorf("the zip holds more than %d MB", maxImportTotal>>20)
		}

		if name == "progress.json" {
			if progress, err = parseProgress(data, f.Modified.UTC()); err != nil {
				reject(name, "not a valid progress file")
			}
			continue
		}
		files = append(files, pending{name, id, lang, data})
	}

	for _, p := range files {
		if _, exists, _ := s.Load(p.id, p.lang); exists && !overwrite {
			rep.Skipped = append(rep.Skipped, p.name)
			continue
		}
		if err := s.Save(p.id, p.lang, string(p.data)); err != nil {
			return rep, err
		}
		rep.Imported = append(rep.Imported, p.name)
	}
	if progress != nil {
		n, err := s.mergeProgress(progress)
		if err != nil {
			return rep, err
		}
		rep.ProgressMerged = n
	}
	sort.Strings(rep.Imported)
	sort.Strings(rep.Skipped)
	return rep, nil
}

func readLimited(f *zip.File, limit int64) ([]byte, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, errors.New("cannot read this entry")
	}
	defer rc.Close()
	var buf bytes.Buffer
	n, err := io.Copy(&buf, io.LimitReader(rc, limit+1))
	if err != nil {
		return nil, errors.New("cannot read this entry")
	}
	if n > limit {
		return nil, fmt.Errorf("larger than %d MB", limit>>20)
	}
	return buf.Bytes(), nil
}

func rank(status string) int {
	switch status {
	case "solved":
		return 2
	case "tried":
		return 1
	}
	return 0
}

// mergeProgress adds imported entries that are new or better than ours.
func (s *Store) mergeProgress(in map[string]Entry) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	m := s.readProgress()
	changed := 0
	for id, e := range in {
		if !idPattern.MatchString(id) || rank(e.Status) == 0 {
			continue
		}
		cur, ok := m[id]
		if !ok || rank(e.Status) > rank(cur.Status) ||
			(rank(e.Status) == rank(cur.Status) && e.UpdatedAt.After(cur.UpdatedAt)) {
			m[id] = e
			changed++
		}
	}
	if changed == 0 {
		return 0, nil
	}
	return changed, s.writeProgress(m)
}
