package problems

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"dsa/internal/testgen"
)

// Validate checks the structure of every pack and the links between them.
// It returns all problems found, joined, or nil.
func Validate(lib *Library) error {
	var errs []error
	fail := func(id, format string, args ...any) {
		errs = append(errs, fmt.Errorf("%s: "+format, append([]any{id}, args...)...))
	}

	for _, p := range lib.List() {
		switch p.Difficulty {
		case "easy", "medium", "hard":
		default:
			fail(p.ID, "difficulty must be easy, medium or hard")
		}
		if p.Sheet != "patterns" && p.Sheet != "real" {
			fail(p.ID, "sheet must be patterns or real")
		}
		if strings.TrimSpace(p.Group) == "" {
			fail(p.ID, "group is not set")
		}
		if p.Order <= 0 {
			fail(p.ID, "order must be 1 or more")
		}
		if !testgen.Has(p.Generator) {
			fail(p.ID, "generator %q is not registered in internal/testgen", p.Generator)
		}
		if !slices.Contains(strings.Split(strings.ReplaceAll(p.Statement, "\r", ""), "\n"), "## Constraints") {
			fail(p.ID, "statement has no \"## Constraints\" heading")
		}
		if len(p.Examples) < 2 {
			fail(p.ID, "needs at least 2 examples, has %d", len(p.Examples))
		}
		if len(p.Edge) < 3 {
			fail(p.ID, "needs at least 3 edge cases, has %d", len(p.Edge))
		}
		for _, c := range append(slices.Clone(p.Examples), p.Edge...) {
			if msg := inputProblem(c.Input); msg != "" {
				fail(p.ID, "%q: input %s", c.Label, msg)
			}
		}
		for _, id := range p.BuildsOn {
			if _, ok := lib.Get(id); !ok {
				fail(p.ID, "buildsOn names %q, which does not exist", id)
			}
		}
		for _, id := range p.LeadsTo {
			next, ok := lib.Get(id)
			if !ok {
				fail(p.ID, "leadsTo names %q, which does not exist", id)
			} else if !slices.Contains(next.BuildsOn, p.ID) {
				fail(p.ID, "leads to %q, but %q does not build on it", id, id)
			}
		}
	}
	if cycle := findCycle(lib); cycle != nil {
		errs = append(errs, fmt.Errorf("buildsOn has a cycle: %s", strings.Join(cycle, " -> ")))
	}
	return errors.Join(errs...)
}

// inputProblem rejects trailing spaces and leading or trailing blank lines.
// Blank lines inside an input are allowed, because an empty array is "0"
// followed by an empty line.
func inputProblem(in string) string {
	if strings.TrimSpace(in) == "" {
		return "is empty"
	}
	if strings.HasPrefix(in, "\n") || strings.HasSuffix(in, "\n") {
		return "starts or ends with a blank line"
	}
	for _, line := range strings.Split(in, "\n") {
		if strings.TrimRight(line, " \t\r") != line {
			return "has trailing spaces"
		}
	}
	return ""
}

// findCycle returns the ids along one buildsOn cycle, or nil.
func findCycle(lib *Library) []string {
	const (
		unseen = iota
		active
		finished
	)
	state := map[string]int{}
	var path []string
	var visit func(id string) []string
	visit = func(id string) []string {
		switch state[id] {
		case active:
			i := slices.Index(path, id)
			return append(slices.Clone(path[i:]), id)
		case finished:
			return nil
		}
		state[id] = active
		path = append(path, id)
		if p, ok := lib.Get(id); ok {
			for _, dep := range p.BuildsOn {
				if c := visit(dep); c != nil {
					return c
				}
			}
		}
		path = path[:len(path)-1]
		state[id] = finished
		return nil
	}
	for _, p := range lib.List() {
		if c := visit(p.ID); c != nil {
			return c
		}
	}
	return nil
}
