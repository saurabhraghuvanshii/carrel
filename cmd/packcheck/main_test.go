package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const (
	realPacks = "../../internal/problems/packs"
	realRefs  = "../../tools/refs"
)

func needCompilers(t *testing.T) {
	t.Helper()
	for _, tool := range []string{"javac", "java", "g++"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skipf("%s is not installed", tool)
		}
	}
}

// copyPack copies one built-in pack into a new packs folder and returns that folder.
func copyPack(t *testing.T, id string) string {
	t.Helper()
	dir := t.TempDir()
	src, dst := filepath.Join(realPacks, id), filepath.Join(dir, id)
	entries, err := os.ReadDir(src)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(dst, 0o700); err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		b, err := os.ReadFile(filepath.Join(src, e.Name()))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dst, e.Name()), b, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func check(args ...string) (int, string) {
	var out strings.Builder
	code := run(append([]string{"-random", "20"}, args...), &out)
	return code, out.String()
}

func TestShippedPacksPass(t *testing.T) {
	needCompilers(t)
	code, out := check("-packs", realPacks, "-refs", realRefs)
	if code != 0 {
		t.Fatalf("exit %d\n%s", code, out)
	}
}

func TestWrongExpectedValueFailsAndNamesTheCase(t *testing.T) {
	needCompilers(t)
	dir := copyPack(t, "pair-with-target-sum")
	path := filepath.Join(dir, "pair-with-target-sum", "tests.json")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	broken := strings.Replace(string(b), `"input": "4\n7 0 2 0\n0", "expected": "1 3"`, `"input": "4\n7 0 2 0\n0", "expected": "1 2"`, 1)
	if broken == string(b) {
		t.Fatal("test data changed; update this test")
	}
	if err := os.WriteFile(path, []byte(broken), 0o600); err != nil {
		t.Fatal(err)
	}
	code, out := check("-packs", dir, "-refs", realRefs)
	if code == 0 {
		t.Fatalf("a wrong expected value passed\n%s", out)
	}
	if !strings.Contains(out, "first failing case: Edge 2: zeros") {
		t.Fatalf("the failing case was not named\n%s", out)
	}
}

func TestMissingReferenceFails(t *testing.T) {
	needCompilers(t)
	dir := copyPack(t, "pair-with-target-sum")
	code, out := check("-packs", dir, "-refs", t.TempDir())
	if code == 0 || !strings.Contains(out, "no reference at") {
		t.Fatalf("exit %d\n%s", code, out)
	}
}

func TestUnknownIDFails(t *testing.T) {
	needCompilers(t)
	code, out := check("-packs", realPacks, "-refs", realRefs, "no-such-problem")
	if code == 0 || !strings.Contains(out, `no pack named "no-such-problem"`) {
		t.Fatalf("exit %d\n%s", code, out)
	}
}
