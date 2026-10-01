// Package testgen makes fresh random test cases from a seed, so nothing large
// has to be preloaded and a different set can run every time.
//
// Each generator plants a known answer, or computes one with a reference
// solution, and returns the input text plus the expected output line.
package testgen

import (
	"fmt"
	"math/rand"
)

type Case struct {
	Input    string
	Expected string
}

// Generator builds one case of roughly the given size.
type Generator func(rng *rand.Rand, size int) Case

var registry = map[string]Generator{}

func register(name string, g Generator) { registry[name] = g }

func Has(name string) bool { _, ok := registry[name]; return ok }

// Generate returns n cases for the named generator. The same seed always gives
// the same cases, which makes a failing run easy to reproduce.
func Generate(name string, seed int64, n int) ([]Case, error) {
	g, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("no generator named %q", name)
	}
	rng := rand.New(rand.NewSource(seed))
	out := make([]Case, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, g(rng, sizeFor(i, n)))
	}
	return out, nil
}

// sizeFor goes from many tiny cases, through medium ones, to a few large ones.
func sizeFor(i, n int) int {
	switch {
	case i < n/2:
		return 2 + i%10
	case i < n*9/10:
		return 20 + (i*37)%300
	default:
		return 2000 + (i*911)%8000
	}
}
