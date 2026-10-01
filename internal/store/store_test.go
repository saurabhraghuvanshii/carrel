package store

import "testing"

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
	_ = s.Mark("two-sum", "solved")
	_ = s.Mark("two-sum", "tried") // must not downgrade
	if s.Progress()["two-sum"] != "solved" {
		t.Error("a solved problem was downgraded")
	}
}
