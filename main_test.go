package main

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestVersionIsStampedAtBuild(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "carrel")
	build := exec.Command("go", "build", "-ldflags", "-X main.version=v9.8.7", "-o", bin, ".")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}
	out, err := exec.Command(bin, "version").Output()
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(string(out)); got != "carrel v9.8.7" {
		t.Fatalf("version printed %q", got)
	}
}
