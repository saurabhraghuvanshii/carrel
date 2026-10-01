package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRejectsUnsafeIDs(t *testing.T) {
	s, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"../evil", "a/b", "", "UPPER", "a b"} {
		if err := s.Save(id, "java", "x"); err == nil {
			t.Errorf("Save accepted id %q", id)
		}
	}
	if err := s.Save("pair-with-target-sum", "rust", "x"); err == nil {
		t.Error("Save accepted an unsupported language")
	}
}

func TestSaveLoadAndProgress(t *testing.T) {
	s, _ := New(t.TempDir())
	if err := s.Save("two-sum", "cpp", "code"); err != nil {
		t.Fatal(err)
	}
	got, ok, err := s.Load("two-sum", "cpp")
	if err != nil || !ok || got != "code" {
		t.Fatalf("Load = %q, %v, %v", got, ok, err)
	}
}

type clock struct{ t time.Time }

func (c *clock) now() time.Time      { return c.t }
func (c *clock) add(d time.Duration) { c.t = c.t.Add(d) }

func newClocked(t *testing.T) (*Store, *clock) {
	t.Helper()
	s, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	c := &clock{time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC)}
	s.now = c.now
	return s, c
}

const day = 24 * time.Hour

func TestMarkNeverDowngradesAndUpdatesTime(t *testing.T) {
	s, c := newClocked(t)
	if err := s.Mark("two-sum", "tried"); err != nil {
		t.Fatal(err)
	}
	if e := s.Progress()["two-sum"]; e.Status != "tried" || !e.UpdatedAt.Equal(c.t) {
		t.Fatalf("after tried: %+v", e)
	}
	c.add(time.Hour)
	_ = s.Mark("two-sum", "solved")
	c.add(time.Hour)
	_ = s.Mark("two-sum", "tried")
	e := s.Progress()["two-sum"]
	if e.Status != "solved" {
		t.Error("a solved problem was downgraded")
	}
	if !e.UpdatedAt.Equal(c.t) {
		t.Errorf("updatedAt = %v, want %v", e.UpdatedAt, c.t)
	}
	if err := s.Mark("two-sum", "done"); err == nil {
		t.Error("Mark accepted an unknown status")
	}
}

func TestOldProgressFormatIsReadAndConverted(t *testing.T) {
	s, _ := newClocked(t)
	old := `{"two-sum": "solved", "trees": "tried"}`
	if err := os.WriteFile(filepath.Join(s.dir, "progress.json"), []byte(old), 0o600); err != nil {
		t.Fatal(err)
	}
	p := s.Progress()
	if p["two-sum"].Status != "solved" || p["trees"].Status != "tried" || p["trees"].UpdatedAt.IsZero() {
		t.Fatalf("old format not read: %+v", p)
	}
	if err := s.Mark("graphs", "tried"); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(filepath.Join(s.dir, "progress.json"))
	var raw map[string]Entry
	if err := json.Unmarshal(b, &raw); err != nil {
		t.Fatalf("file was not converted: %v\n%s", err, b)
	}
	if raw["two-sum"].Status != "solved" || raw["trees"].Status != "tried" || raw["graphs"].Status != "tried" {
		t.Fatalf("entries lost in conversion: %s", b)
	}
}

func TestReviewSchedule(t *testing.T) {
	s, c := newClocked(t)
	due := func() bool { return s.Progress()["two-sum"].DueForReview(c.t) }

	_ = s.Mark("two-sum", "solved")
	for i, gap := range []time.Duration{3 * day, 10 * day, 30 * day} {
		c.add(gap - time.Hour)
		if due() {
			t.Fatalf("review %d is due too early", i+1)
		}
		_ = s.Mark("two-sum", "solved") // solving early does not move the schedule
		c.add(time.Hour)
		if !due() {
			t.Fatalf("review %d is not due after %v", i+1, gap)
		}
		_ = s.Mark("two-sum", "solved")
		if due() {
			t.Fatalf("review %d is still due after solving again", i+1)
		}
	}
	c.add(365 * day)
	if due() {
		t.Fatal("no review should follow the 30-day one")
	}
}
