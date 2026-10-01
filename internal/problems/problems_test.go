package problems

import "testing"

// Every shipped problem must have both languages. Validate checks the rest.
func TestBuiltinPacksAreComplete(t *testing.T) {
	lib, err := Load(Builtin())
	if err != nil {
		t.Fatal(err)
	}
	if len(lib.List()) < 2 {
		t.Fatalf("expected at least two problems, got %d", len(lib.List()))
	}
	for _, p := range lib.List() {
		for _, f := range []string{"starter.java", "starter.cpp", "driver.java", "driver.cpp"} {
			if _, err := p.File(f); err != nil {
				t.Errorf("%s: missing %s", p.ID, f)
			}
		}
		if len(p.Examples) == 0 {
			t.Errorf("%s: no examples", p.ID)
		}
	}
}
