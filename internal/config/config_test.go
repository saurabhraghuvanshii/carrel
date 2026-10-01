package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHomeMovesTheOldFolderOnce(t *testing.T) {
	user := t.TempDir()
	t.Setenv("HOME", user)
	t.Setenv("CARREL_HOME", "")
	old := filepath.Join(user, ".dsa", "solutions")
	if err := os.MkdirAll(old, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(old, "two-sum.java"), []byte("mine"), 0o600); err != nil {
		t.Fatal(err)
	}

	home, err := Home()
	if err != nil || home != filepath.Join(user, ".carrel") {
		t.Fatalf("Home = %q, %v", home, err)
	}
	if b, err := os.ReadFile(filepath.Join(home, "solutions", "two-sum.java")); err != nil || string(b) != "mine" {
		t.Fatalf("solution not carried over: %q %v", b, err)
	}
	if _, err := os.Stat(filepath.Join(user, ".dsa")); !os.IsNotExist(err) {
		t.Fatal("the old folder should be gone after the move")
	}

	// A new ~/.dsa appearing later is left alone.
	if err := os.MkdirAll(filepath.Join(user, ".dsa"), 0o700); err != nil {
		t.Fatal(err)
	}
	if _, err := Home(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(user, ".dsa")); err != nil {
		t.Fatal("Home moved ~/.dsa although ~/.carrel already existed")
	}
}

func TestCarrelHomeWins(t *testing.T) {
	user := t.TempDir()
	t.Setenv("HOME", user)
	custom := filepath.Join(t.TempDir(), "data")
	t.Setenv("CARREL_HOME", custom)
	if err := os.MkdirAll(filepath.Join(user, ".dsa"), 0o700); err != nil {
		t.Fatal(err)
	}
	home, err := Home()
	if err != nil || home != custom {
		t.Fatalf("Home = %q, %v", home, err)
	}
	if _, err := os.Stat(filepath.Join(user, ".dsa")); err != nil {
		t.Fatal("~/.dsa was moved although CARREL_HOME is set")
	}
}

func TestSaveIsOwnerOnly(t *testing.T) {
	dir := t.TempDir()
	c := Default()
	c.AI.APIKey = "secret"
	if err := Save(dir, c); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(filepath.Join(dir, "config.json"))
	if err != nil || info.Mode().Perm() != 0o600 {
		t.Fatalf("config.json mode %v, %v", info.Mode().Perm(), err)
	}
	got, err := Load(dir)
	if err != nil || got.AI.APIKey != "secret" {
		t.Fatalf("Load = %+v, %v", got, err)
	}
}
