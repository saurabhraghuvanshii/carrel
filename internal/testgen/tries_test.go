package testgen

import (
	"strconv"
	"strings"
	"testing"
)

func TestPrefixTree(t *testing.T) {
	seen := map[string]bool{}
	for i, c := range cases(t, "prefix-tree", 60) {
		lines := strings.Split(c.Input, "\n")
		words := map[string]bool{}
		var got []string
		for _, l := range lines[1:] {
			f := strings.Fields(l)
			switch f[0] {
			case "insert":
				words[f[1]] = true
			case "search":
				got = append(got, boolText(words[f[1]]))
			default:
				ok := false
				for w := range words {
					ok = ok || strings.HasPrefix(w, f[1])
				}
				got = append(got, boolText(ok))
			}
		}
		if "["+strings.Join(got, ", ")+"]" != c.Expected {
			t.Fatalf("case %d differs", i)
		}
		for _, g := range got {
			seen[g] = true
		}
	}
	if !seen["true"] || !seen["false"] {
		t.Fatal("both answers must appear")
	}
}

func TestUniquePrefix(t *testing.T) {
	for i, c := range cases(t, "unique-prefix", 60) {
		words := strings.Fields(strings.Split(c.Input, "\n")[1])
		out := make([]string, len(words))
		for k, w := range words {
			for size := 1; size <= len(w) && out[k] == ""; size++ {
				shared := false
				for j, o := range words {
					shared = shared || (j != k && strings.HasPrefix(o, w[:size]))
				}
				if !shared {
					out[k] = w[:size]
				}
			}
			if out[k] == "" {
				t.Fatalf("case %d: %q has no unique beginning", i, w)
			}
		}
		if quoted(out) != c.Expected {
			t.Fatalf("case %d differs", i)
		}
	}
}

func TestBuildableWord(t *testing.T) {
	empty := 0
	for i, c := range cases(t, "buildable-word", 80) {
		words := strings.Fields(strings.Split(c.Input, "\n")[1])
		set := map[string]bool{}
		for _, w := range words {
			set[w] = true
		}
		best := ""
		for _, w := range words {
			ok := true
			for p := 1; p < len(w); p++ {
				ok = ok && set[w[:p]]
			}
			if ok && (len(w) > len(best) || (len(w) == len(best) && w < best)) {
				best = w
			}
		}
		if `"`+best+`"` != c.Expected {
			t.Fatalf("case %d: want %q, generator says %s", i, best, c.Expected)
		}
		if best == "" {
			empty++
		}
	}
	if empty > 20 {
		t.Fatalf("%d of 80 answers are empty", empty)
	}
}

func TestDistinctSubstrings(t *testing.T) {
	for i, c := range cases(t, "distinct-substrings", 60) {
		w := c.Input
		if len(w) > 150 {
			continue
		}
		set := map[string]bool{}
		for a := 0; a < len(w); a++ {
			for b := a + 1; b <= len(w); b++ {
				set[w[a:b]] = true
			}
		}
		if strconv.Itoa(len(set)) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, len(set), c.Expected)
		}
	}
}
