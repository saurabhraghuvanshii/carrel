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
)

var idPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,63}$`)

var exts = map[string]string{"java": ".java", "cpp": ".cpp"}

type Store struct {
	dir string
	mu  sync.Mutex
}

func New(home string) (*Store, error) {
	s := &Store{dir: home}
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

// Progress maps problem id to "tried" or "solved".
func (s *Store) Progress() map[string]string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.readProgress()
}

func (s *Store) readProgress() map[string]string {
	m := map[string]string{}
	b, err := os.ReadFile(filepath.Join(s.dir, "progress.json"))
	if err == nil {
		_ = json.Unmarshal(b, &m)
	}
	return m
}

// Mark records progress. A solved problem is never downgraded to tried.
func (s *Store) Mark(id, status string) error {
	if !idPattern.MatchString(id) {
		return fmt.Errorf("invalid problem id %q", id)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	m := s.readProgress()
	if m[id] == "solved" {
		return nil
	}
	m[id] = status
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	tmp := filepath.Join(s.dir, "progress.json.tmp")
	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, filepath.Join(s.dir, "progress.json"))
}
