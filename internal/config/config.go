// Package config reads and writes the user's settings in ~/.carrel/config.json.
// The AI key lives here, on this computer only, never inside a project folder.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type AI struct {
	Provider    string `json:"provider"` // anthropic | openai | ollama
	Model       string `json:"model"`
	APIKey      string `json:"apiKey"`
	ExplainOnly bool   `json:"explainOnly"`
}

type Config struct {
	Theme  string `json:"theme"`  // paper | ink | system
	Accent string `json:"accent"` // brick | forest | amber | slate | plum
	Lang   string `json:"lang"`   // java | cpp
	AI     AI     `json:"ai"`

	RemindReviews bool `json:"remindReviews"` // show solved problems again after 3, 10 and 30 days
}

func Default() Config {
	return Config{
		Theme:  "paper",
		Accent: "brick",
		Lang:   "java",
		AI:     AI{Provider: "anthropic", ExplainOnly: true},
	}
}

// Home returns the data folder: $CARREL_HOME if set, otherwise ~/.carrel.
func Home() (string, error) {
	if h := os.Getenv("CARREL_HOME"); h != "" {
		return h, os.MkdirAll(h, 0o700)
	}
	u, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	h := filepath.Join(u, ".carrel")
	if err := moveOldHome(filepath.Join(u, ".dsa"), h); err != nil {
		return "", err
	}
	return h, os.MkdirAll(h, 0o700)
}

// moveOldHome carries data over from the project's earlier name. If ~/.carrel
// does not exist yet and ~/.dsa does, the old folder is renamed once.
func moveOldHome(old, home string) error {
	if _, err := os.Stat(home); !errors.Is(err, os.ErrNotExist) {
		return nil
	}
	info, err := os.Stat(old)
	if err != nil || !info.IsDir() {
		return nil
	}
	if err := os.Rename(old, home); err != nil {
		return fmt.Errorf("cannot move your data from %s to %s: %w", old, home, err)
	}
	log.Printf("moved your solutions and settings from %s to %s", old, home)
	return nil
}

func Load(home string) (Config, error) {
	c := Default()
	b, err := os.ReadFile(filepath.Join(home, "config.json"))
	if errors.Is(err, os.ErrNotExist) {
		return c, nil
	}
	if err != nil {
		return c, err
	}
	if err := json.Unmarshal(b, &c); err != nil {
		return Default(), err
	}
	return c, nil
}

// Save writes the file with owner-only permissions because it can hold an API key.
func Save(home string, c Config) error {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	tmp := filepath.Join(home, "config.json.tmp")
	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, filepath.Join(home, "config.json"))
}
