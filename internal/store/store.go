// Package store keeps the user's solutions and progress as plain files.
package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

var idPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,63}$`)

var exts = map[string]string{"java": ".java", "cpp": ".cpp"}

type Store struct {
	dir string
	mu  sync.Mutex
	now func() time.Time
}

func New(home string) (*Store, error) {
	s := &Store{dir: home, now: time.Now}
	return s, os.MkdirAll(filepath.Join(home, "solutions"), 0o700)
}

func (s *Store) SolutionsDir() string { return filepath.Join(s.dir, "solutions") }

func (s *Store) path(id, lang string) (string, error) {
	ext, ok := exts[lang]
	if !ok {
		return "", fmt.Errorf("unsupported language %q", lang)
	}
	if !idPattern.MatchString(id) {
		return "", fmt.Errorf("invalid problem id %q", id)
	}
	return filepath.Join(s.SolutionsDir(), id+ext), nil
}

// Load returns the saved code, or ok=false when nothing was saved yet.
func (s *Store) Load(id, lang string) (code string, ok bool, err error) {
	p, err := s.path(id, lang)
	if err != nil {
		return "", false, err
	}
	b, err := os.ReadFile(p)
	if errors.Is(err, os.ErrNotExist) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return string(b), true, nil
}

func (s *Store) Save(id, lang, code string) error {
	p, err := s.path(id, lang)
	if err != nil {
		return err
	}
	tmp := p + ".tmp"
	if err := os.WriteFile(tmp, []byte(code), 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, p)
}

// Entry is one problem's progress. A solved problem also carries its review
// schedule: due 3 days after solving, then 10, then 30, then never again.
type Entry struct {
	Status     string     `json:"status"` // tried | solved
	UpdatedAt  time.Time  `json:"updatedAt"`
	ReviewStep int        `json:"reviewStep,omitempty"`
	ReviewDue  *time.Time `json:"reviewDue,omitempty"`
}

var reviewGaps = []time.Duration{3 * 24 * time.Hour, 10 * 24 * time.Hour, 30 * 24 * time.Hour}

// DueForReview reports whether a solved problem should be practised again.
func (e Entry) DueForReview(now time.Time) bool {
	return e.Status == "solved" && e.ReviewDue != nil && !now.Before(*e.ReviewDue)
}

func (s *Store) progressPath() string { return filepath.Join(s.dir, "progress.json") }

// Progress maps problem id to its entry.
func (s *Store) Progress() map[string]Entry {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.readProgress()
}

// readProgress also reads the first format, {id: "tried"}, and dates those
// entries with the file's modification time. Mark writes the new format.
func (s *Store) readProgress() map[string]Entry {
	b, err := os.ReadFile(s.progressPath())
	if err != nil {
		return map[string]Entry{}
	}
	var fileTime time.Time
	if info, err := os.Stat(s.progressPath()); err == nil {
		fileTime = info.ModTime().UTC().Truncate(time.Second)
	}
	m, _ := parseProgress(b, fileTime)
	return m
}

// parseProgress reads either format. Entries in the old format get oldTime.
func parseProgress(b []byte, oldTime time.Time) (map[string]Entry, error) {
	m := map[string]Entry{}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return m, err
	}
	for id, v := range raw {
		var old string
		if json.Unmarshal(v, &old) == nil {
			m[id] = Entry{Status: old, UpdatedAt: oldTime}
			continue
		}
		var e Entry
		if json.Unmarshal(v, &e) == nil {
			m[id] = e
		}
	}
	return m, nil
}

// writeProgress replaces progress.json. The caller holds s.mu.
func (s *Store) writeProgress(m map[string]Entry) error {
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.progressPath() + ".tmp"
	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, s.progressPath())
}

// Mark records progress. A solved problem is never downgraded to tried, but
// its time is updated. Solving schedules the next review.
func (s *Store) Mark(id, status string) error {
	if !idPattern.MatchString(id) {
		return fmt.Errorf("invalid problem id %q", id)
	}
	if status != "tried" && status != "solved" {
		return fmt.Errorf("invalid status %q", status)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	m := s.readProgress()
	e := m[id]
	now := s.now().UTC()
	switch {
	case status == "solved" && e.Status != "solved":
		e.ReviewStep = 0
		e.ReviewDue = after(now, 0)
	case status == "solved" && e.DueForReview(now):
		e.ReviewStep++
		e.ReviewDue = after(now, e.ReviewStep)
	}
	if e.Status != "solved" {
		e.Status = status
	}
	e.UpdatedAt = now
	m[id] = e
	return s.writeProgress(m)
}

func after(now time.Time, step int) *time.Time {
	if step >= len(reviewGaps) {
		return nil
	}
	t := now.Add(reviewGaps[step])
	return &t
}
