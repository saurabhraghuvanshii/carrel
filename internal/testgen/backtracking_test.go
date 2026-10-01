package testgen

import (
	"slices"
	"strconv"
	"strings"
	"testing"
)

func TestNoAdjacentOnes(t *testing.T) {
	for i, c := range cases(t, "no-adjacent-ones", 40) {
		n, _ := strconv.Atoi(c.Input)
		var want []string
		for m := 0; m < 1<<n; m++ {
			s := strconv.FormatInt(int64(m), 2)
			s = strings.Repeat("0", n-len(s)) + s
			if !strings.Contains(s, "11") {
				want = append(want, s)
			}
		}
		if stringsOutput(want) != c.Expected {
			t.Fatalf("case %d (n=%d) differs", i, n)
		}
	}
}

func TestAllSubsetsAndOrderings(t *testing.T) {
	for i, c := range cases(t, "all-subsets", 40) {
		var want []string
		var walk func(k int, cur string)
		walk = func(k int, cur string) {
			if k == len(c.Input) {
				want = append(want, cur)
				return
			}
			walk(k+1, cur)
			walk(k+1, cur+c.Input[k:k+1])
		}
		walk(0, "")
		if stringsOutput(want) != c.Expected {
			t.Fatalf("subsets case %d differs", i)
		}
	}
	for i, c := range cases(t, "all-orderings", 40) {
		got := strings.Split(strings.Trim(c.Expected, "[]"), ", ")
		fact := 1
		for k := 2; k <= len(c.Input); k++ {
			fact *= k
		}
		seen := map[string]bool{}
		for _, q := range got {
			w := strings.Trim(q, `"`)
			if sortedLetters(w) != sortedLetters(c.Input) || seen[w] {
				t.Fatalf("orderings case %d: bad or repeated %q", i, w)
			}
			seen[w] = true
		}
		if len(got) != fact {
			t.Fatalf("orderings case %d: %d orderings, want %d", i, len(got), fact)
		}
	}
}

func TestComboSum(t *testing.T) {
	for i, c := range cases(t, "combo-sum", 50) {
		cands, rest := parseArray(t, c.Input)
		target, _ := strconv.Atoi(rest[0])
		// Count combinations as in coin change: ways[v] over candidates.
		ways := make([]int, target+1)
		ways[0] = 1
		for _, x := range cands {
			for v := x; v <= target; v++ {
				ways[v] += ways[v-x]
			}
		}
		res := combos(cands, target)
		if len(res) != ways[target] || gridOutput(res) != c.Expected {
			t.Fatalf("case %d: %d combinations, coin change says %d", i, len(res), ways[target])
		}
		for _, r := range res {
			sum := 0
			for _, v := range r {
				sum += v
				if !slices.Contains(cands, v) {
					t.Fatalf("case %d: %d is not a candidate", i, v)
				}
			}
			if sum != target || !slices.IsSorted(r) {
				t.Fatalf("case %d: bad combination %v", i, r)
			}
		}
	}
}

func TestPhoneLetters(t *testing.T) {
	for i, c := range cases(t, "phone-letters", 40) {
		want := 1
		for j := 0; j < len(c.Input); j++ {
			want *= len(keypad[c.Input[j]])
		}
		got := strings.Split(strings.Trim(c.Expected, "[]"), ", ")
		if len(got) != want {
			t.Fatalf("case %d: %d words, want %d", i, len(got), want)
		}
		for _, q := range got {
			w := strings.Trim(q, `"`)
			for j := 0; j < len(w); j++ {
				if !strings.Contains(keypad[c.Input[j]], w[j:j+1]) {
					t.Fatalf("case %d: %q does not match %s", i, w, c.Input)
				}
			}
		}
	}
}

func TestBracketStrings(t *testing.T) {
	for i, c := range cases(t, "bracket-strings", 40) {
		n, _ := strconv.Atoi(c.Input)
		if n > 7 {
			continue
		}
		var want []string
		for m := 0; m < 1<<(2*n); m++ {
			b := make([]byte, 2*n)
			for j := range b {
				b[j] = "()"[m>>j&1]
			}
			if balanced(string(b)) {
				want = append(want, string(b))
			}
		}
		if stringsOutput(want) != c.Expected {
			t.Fatalf("case %d (n=%d) differs", i, n)
		}
	}
}

func TestPalindromeSplits(t *testing.T) {
	for i, c := range cases(t, "palindrome-splits", 50) {
		w := c.Input
		var want [][]string
		for mask := 0; mask < 1<<(len(w)-1); mask++ {
			var pieces []string
			start := 0
			for j := 1; j <= len(w); j++ {
				if j == len(w) || mask>>(j-1)&1 == 1 {
					pieces = append(pieces, w[start:j])
					start = j
				}
			}
			ok := true
			for _, p := range pieces {
				rev := []byte(p)
				slices.Reverse(rev)
				ok = ok && string(rev) == p
			}
			if ok {
				want = append(want, pieces)
			}
		}
		if splitsOutput(want) != c.Expected {
			t.Fatalf("case %d (%s) differs", i, w)
		}
	}
}

func TestGridWord(t *testing.T) {
	seen := map[string]bool{}
	for i, c := range cases(t, "grid-word", 80) {
		lines := strings.Split(c.Input, "\n")
		dims := strings.Fields(lines[0])
		r, _ := strconv.Atoi(dims[0])
		rows, word := lines[1:r+1], lines[r+1]
		// Extend partial paths one letter at a time.
		var paths [][][2]int
		for y := range rows {
			for x := range rows[y] {
				if rows[y][x] == word[0] {
					paths = append(paths, [][2]int{{y, x}})
				}
			}
		}
		for k := 1; k < len(word) && len(paths) > 0; k++ {
			var next [][][2]int
			for _, p := range paths {
				last := p[len(p)-1]
				for _, d := range [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}} {
					q := [2]int{last[0] + d[0], last[1] + d[1]}
					if q[0] >= 0 && q[0] < r && q[1] >= 0 && q[1] < len(rows[0]) && rows[q[0]][q[1]] == word[k] && !slices.Contains(p, q) {
						next = append(next, append(slices.Clone(p), q))
					}
				}
			}
			paths = next
		}
		if want := boolText(len(paths) > 0); want != c.Expected {
			t.Fatalf("case %d: want %s, generator says %s", i, want, c.Expected)
		}
		seen[c.Expected] = true
	}
	if !seen["true"] || !seen["false"] {
		t.Fatal("both answers must appear")
	}
}

func TestQueensHoles(t *testing.T) {
	for i, c := range cases(t, "queens-holes", 40) {
		lines := strings.Split(c.Input, "\n")
		n, _ := strconv.Atoi(lines[0])
		if n > 8 {
			continue
		}
		rows := lines[1:]
		cols := make([]int, n)
		for j := range cols {
			cols[j] = j
		}
		// Try every ordering of columns, one queen per row.
		count := 0
		var perm func(k int)
		perm = func(k int) {
			if k == n {
				for r := 0; r < n; r++ {
					if rows[r][cols[r]] == '#' {
						return
					}
					for q := r + 1; q < n; q++ {
						if q-r == cols[q]-cols[r] || q-r == cols[r]-cols[q] {
							return
						}
					}
				}
				count++
				return
			}
			for j := k; j < n; j++ {
				cols[k], cols[j] = cols[j], cols[k]
				perm(k + 1)
				cols[k], cols[j] = cols[j], cols[k]
			}
		}
		perm(0)
		if strconv.Itoa(count) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, count, c.Expected)
		}
	}
}
