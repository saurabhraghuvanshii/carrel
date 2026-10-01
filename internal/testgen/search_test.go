package testgen

import (
	"slices"
	"strconv"
	"strings"
	"testing"
)

func TestSuggestionCounts(t *testing.T) {
	for i, c := range cases(t, "suggestion-counts", 60) {
		lines := strings.Split(c.Input, "\n")
		words, prefixes := strings.Fields(lines[1]), strings.Fields(lines[2])
		if len(words)*len(prefixes) > 4_000_000 {
			continue
		}
		counts := make([]int, len(prefixes))
		for k, p := range prefixes {
			for _, w := range words {
				if strings.HasPrefix(w, p) {
					counts[k]++
				}
			}
		}
		if listOutput(counts) != c.Expected {
			t.Fatalf("case %d differs", i)
		}
	}
}

func TestNearMatches(t *testing.T) {
	for i, c := range cases(t, "near-matches", 60) {
		lines := strings.Split(c.Input, "\n")
		words, queries := strings.Fields(lines[1]), strings.Fields(lines[2])
		seen := map[string]bool{}
		for _, w := range words {
			if seen[w] {
				t.Fatalf("case %d: %q is repeated", i, w)
			}
			seen[w] = true
		}
		counts := make([]int, len(queries))
		for k, q := range queries {
			for _, w := range words {
				if len(w) != len(q) {
					continue
				}
				diff := 0
				for j := range w {
					if w[j] != q[j] {
						diff++
					}
				}
				if diff <= 1 {
					counts[k]++
				}
			}
		}
		if listOutput(counts) != c.Expected {
			t.Fatalf("case %d differs", i)
		}
	}
}

func TestTrendingTags(t *testing.T) {
	empty := 0
	for i, c := range cases(t, "trending-tags", 80) {
		lines := strings.Split(c.Input, "\n")
		head := intsOf(lines[0])
		rows := rowsOf(lines, 1, head[0])
		last := rows[len(rows)-1][0]
		growth := func(tag int) int {
			g := 0
			for _, r := range rows {
				if r[1] != tag {
					continue
				}
				if age := last - r[0]; age < head[1] {
					g++
				} else if age < 2*head[1] {
					g--
				}
			}
			return g
		}
		// Pick the best remaining tag k times.
		var got []int
		for len(got) < head[2] {
			best, bestGrowth := 0, 0
			for tag := 1; tag <= 10; tag++ {
				if g := growth(tag); g > bestGrowth && !slices.Contains(got, tag) {
					best, bestGrowth = tag, g
				}
			}
			if bestGrowth == 0 {
				break
			}
			got = append(got, best)
		}
		if listOutput(got) != c.Expected {
			t.Fatalf("case %d: want %v, generator says %s", i, got, c.Expected)
		}
		if len(got) == 0 {
			empty++
		}
	}
	if empty > 30 {
		t.Fatalf("%d of 80 answers are empty", empty)
	}
}

func TestAutocomplete(t *testing.T) {
	for i, c := range cases(t, "autocomplete", 50) {
		lines := strings.Split(c.Input, "\n")
		total := map[string]int{}
		got := []string{"null"}
		for _, l := range lines[2:] {
			f := strings.Fields(l)
			if f[0] == "add" {
				score, _ := strconv.Atoi(f[2])
				total[f[1]] += score
				got = append(got, "null")
				continue
			}
			// Take the best remaining word up to three times.
			var picked []string
			for len(picked) < 3 {
				best := ""
				for w, s := range total {
					if !strings.HasPrefix(w, f[1]) || strings.Contains(" "+strings.Join(picked, " ")+" ", " "+w+" ") {
						continue
					}
					if best == "" || s > total[best] || (s == total[best] && w < best) {
						best = w
					}
				}
				if best == "" {
					break
				}
				picked = append(picked, best)
			}
			got = append(got, quoted(picked))
		}
		if strings.Join(got, " ") != c.Expected {
			t.Fatalf("case %d differs", i)
		}
	}
}

func TestMergeContacts(t *testing.T) {
	merged := false
	for i, c := range cases(t, "merge-contacts", 60) {
		lines := strings.Split(c.Input, "\n")
		n, _ := strconv.Atoi(lines[0])
		if n > 300 {
			continue
		}
		emails := make([][]string, n)
		for k := range emails {
			f := strings.Fields(lines[k+1])
			if cnt, _ := strconv.Atoi(f[0]); cnt != len(f)-1 {
				t.Fatalf("case %d: contact %d says %s emails, has %d", i, k, f[0], len(f)-1)
			}
			emails[k] = f[1:]
		}
		// Spread the lowest number between contacts that share an email
		// until nothing changes.
		owner := make([]int, n)
		for k := range owner {
			owner[k] = k
		}
		for changed := true; changed; {
			changed = false
			for a := 0; a < n; a++ {
				for b := a + 1; b < n; b++ {
					if owner[a] == owner[b] || !shareAny(emails[a], emails[b]) {
						continue
					}
					low := min(owner[a], owner[b])
					owner[a], owner[b], changed = low, low, true
				}
			}
		}
		if listOutput(owner) != c.Expected {
			t.Fatalf("case %d differs", i)
		}
		for k, o := range owner {
			merged = merged || o != k
		}
	}
	if !merged {
		t.Fatal("no contacts were merged")
	}
}

func shareAny(a, b []string) bool {
	for _, x := range a {
		for _, y := range b {
			if x == y {
				return true
			}
		}
	}
	return false
}
