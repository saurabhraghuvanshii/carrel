// Package problems loads problem packs: one folder per problem.
//
//	packs/<id>/meta.json       title, difficulty, sheet, order, tags, generator
//	packs/<id>/statement.md    the written statement
//	packs/<id>/tests.json      visible examples and fixed edge cases
//	packs/<id>/starter.java    starter code, one file per language
//	packs/<id>/driver.java     reads test cases from stdin, calls the user's code
package problems

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"regexp"
	"sort"
)

//go:embed all:packs
var builtin embed.FS

// Builtin returns the problem packs that ship inside the binary.
func Builtin() fs.FS {
	sub, err := fs.Sub(builtin, "packs")
	if err != nil {
		panic(err)
	}
	return sub
}

var idPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,63}$`)

type Meta struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Difficulty string   `json:"difficulty"` // easy | medium | hard
	Sheet      string   `json:"sheet"`      // patterns | real
	Group      string   `json:"group"`
	Order      int      `json:"order"`
	Tags       []string `json:"tags"`
	BuildsOn   []string `json:"buildsOn"`
	LeadsTo    []string `json:"leadsTo"`
	Generator  string   `json:"generator"` // name registered in package testgen
}

// Case is one test. Input is what the driver reads, Expected is one output line.
// Display is a friendly form shown to the user for examples.
type Case struct {
	Label    string `json:"label"`
	Input    string `json:"input"`
	Expected string `json:"expected"`
	Display  string `json:"display,omitempty"`
}

type Problem struct {
	Meta
	Statement string
	Examples  []Case
	Edge      []Case
	fsys      fs.FS
}

type Library struct {
	byID map[string]*Problem
	list []*Problem
}

// Load reads every pack from the given file systems. Later ones override earlier
// ones, so a user's own pack can replace a built-in problem.
func Load(fsyss ...fs.FS) (*Library, error) {
	lib := &Library{byID: map[string]*Problem{}}
	for _, fsys := range fsyss {
		entries, err := fs.ReadDir(fsys, ".")
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, err
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			p, err := loadOne(fsys, e.Name())
			if err != nil {
				return nil, fmt.Errorf("pack %q: %w", e.Name(), err)
			}
			lib.byID[p.ID] = p
		}
	}
	for _, p := range lib.byID {
		lib.list = append(lib.list, p)
	}
	sort.Slice(lib.list, func(i, j int) bool {
		a, b := lib.list[i], lib.list[j]
		if a.Sheet != b.Sheet {
			return a.Sheet > b.Sheet // "patterns" before "real"
		}
		if a.Order != b.Order {
			return a.Order < b.Order
		}
		return a.ID < b.ID
	})
	return lib, nil
}

func loadOne(fsys fs.FS, dir string) (*Problem, error) {
	sub, err := fs.Sub(fsys, dir)
	if err != nil {
		return nil, err
	}
	var p Problem
	p.fsys = sub

	b, err := fs.ReadFile(sub, "meta.json")
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &p.Meta); err != nil {
		return nil, fmt.Errorf("meta.json: %w", err)
	}
	if !idPattern.MatchString(p.ID) || p.ID != dir {
		return nil, fmt.Errorf("meta.json id %q must match the folder name", p.ID)
	}
	switch p.Difficulty {
	case "easy", "medium", "hard":
	default:
		return nil, fmt.Errorf("difficulty must be easy, medium or hard")
	}

	st, err := fs.ReadFile(sub, "statement.md")
	if err != nil {
		return nil, err
	}
	p.Statement = string(st)

	tb, err := fs.ReadFile(sub, "tests.json")
	if err != nil {
		return nil, err
	}
	var tests struct {
		Examples []Case `json:"examples"`
		Edge     []Case `json:"edge"`
	}
	if err := json.Unmarshal(tb, &tests); err != nil {
		return nil, fmt.Errorf("tests.json: %w", err)
	}
	p.Examples, p.Edge = tests.Examples, tests.Edge
	return &p, nil
}

func (l *Library) Get(id string) (*Problem, bool) {
	p, ok := l.byID[id]
	return p, ok
}

func (l *Library) List() []*Problem { return l.list }

// File reads a file from the problem's pack, such as "driver.java".
func (p *Problem) File(name string) (string, error) {
	b, err := fs.ReadFile(p.fsys, name)
	return string(b), err
}
