package testgen

import (
	"math"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func unquote(t *testing.T, s string) string {
	t.Helper()
	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		t.Fatalf("not in quotes: %q", s)
	}
	return s[1 : len(s)-1]
}

func TestSquashRepeats(t *testing.T) {
	for i, c := range cases(t, "squash-repeats", 60) {
		// Remove one letter of any doubled pair until none is left.
		w := c.Input
		for changed := true; changed; {
			changed = false
			for ch := 'a'; ch <= 'z'; ch++ {
				d := string(ch) + string(ch)
				if strings.Contains(w, d) {
					w = strings.ReplaceAll(w, d, string(ch))
					changed = true
				}
			}
		}
		if `"`+w+`"` != c.Expected {
			t.Fatalf("case %d: want %q, generator says %s", i, w, c.Expected)
		}
	}
}

func TestRunLength(t *testing.T) {
	groups := regexp.MustCompile(`([a-z])([0-9]*)`)
	for i, c := range cases(t, "run-length", 60) {
		code := unquote(t, c.Expected)
		var decoded strings.Builder
		prev := byte(0)
		for _, m := range groups.FindAllStringSubmatch(code, -1) {
			n := 1
			if m[2] != "" {
				n, _ = strconv.Atoi(m[2])
				if n < 2 {
					t.Fatalf("case %d: count %d in %s", i, n, code)
				}
			}
			if m[1][0] == prev {
				t.Fatalf("case %d: two groups of %s side by side", i, m[1])
			}
			prev = m[1][0]
			decoded.WriteString(strings.Repeat(m[1], n))
		}
		if decoded.String() != c.Input {
			t.Fatalf("case %d: %s does not decode to the input", i, code)
		}
	}
}

func TestWalkDistance(t *testing.T) {
	wide := false
	for i, c := range cases(t, "walk-distance", 80) {
		x := strings.Count(c.Input, "E") - strings.Count(c.Input, "W")
		y := strings.Count(c.Input, "N") - strings.Count(c.Input, "S")
		want := x*x + y*y
		if strconv.Itoa(want) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, want, c.Expected)
		}
		wide = wide || want > math.MaxInt32
	}
	if !wide {
		t.Fatal("no case passes 32 bits")
	}
}

func TestCommonStart(t *testing.T) {
	for i, c := range cases(t, "common-start", 60) {
		words := strings.Fields(strings.Split(c.Input, "\n")[1])
		got := unquote(t, c.Expected)
		for _, w := range words {
			if !strings.HasPrefix(w, got) {
				t.Fatalf("case %d: %q does not start with %q", i, w, got)
			}
		}
		if len(got) < len(words[0]) {
			longer := words[0][:len(got)+1]
			all := true
			for _, w := range words {
				all = all && strings.HasPrefix(w, longer)
			}
			if all {
				t.Fatalf("case %d: %q is shared too", i, longer)
			}
		}
	}
}

func TestReverseWords(t *testing.T) {
	for i, c := range cases(t, "reverse-words", 60) {
		text := unquote(t, c.Input)
		// Read words from the right end.
		var out []string
		end := len(text)
		for k := len(text) - 1; k >= -1; k-- {
			if k == -1 || text[k] == ' ' {
				if k+1 < end {
					out = append(out, text[k+1:end])
				}
				end = k
			}
		}
		if want := `"` + strings.Join(out, " ") + `"`; want != c.Expected {
			t.Fatalf("case %d: want %s, generator says %s", i, want, c.Expected)
		}
	}
}

func TestFirstUnique(t *testing.T) {
	for i, c := range cases(t, "first-unique", 60) {
		w := c.Input
		if len(w) > 400 {
			continue
		}
		out := make([]byte, len(w))
		for p := range w {
			out[p] = '#'
			for k := 0; k <= p; k++ {
				if strings.Count(w[:p+1], w[k:k+1]) == 1 {
					out[p] = w[k]
					break
				}
			}
		}
		if `"`+string(out)+`"` != c.Expected {
			t.Fatalf("case %d differs", i)
		}
	}
}

func TestMirrorStretch(t *testing.T) {
	for i, c := range cases(t, "mirror-stretch", 60) {
		w := c.Input
		if len(w) > 80 {
			continue
		}
		best := ""
		for size := len(w); size > 0 && best == ""; size-- {
			for from := 0; from+size <= len(w) && best == ""; from++ {
				if isPal(w[from : from+size]) {
					best = w[from : from+size]
				}
			}
		}
		if `"`+best+`"` != c.Expected {
			t.Fatalf("case %d: want %q, generator says %s", i, best, c.Expected)
		}
	}
}

func TestTextToNumber(t *testing.T) {
	lead := regexp.MustCompile(`^ *([+-]?)([0-9]*)`)
	seen := map[string]bool{}
	for i, c := range cases(t, "text-to-number", 120) {
		m := lead.FindStringSubmatch(unquote(t, c.Input))
		v := new(big.Int)
		if m[2] != "" {
			v.SetString(m[2], 10)
		}
		if m[1] == "-" {
			v.Neg(v)
		}
		if v.Cmp(big.NewInt(math.MaxInt32)) > 0 {
			v.SetInt64(math.MaxInt32)
		}
		if v.Cmp(big.NewInt(math.MinInt32)) < 0 {
			v.SetInt64(math.MinInt32)
		}
		if v.String() != c.Expected {
			t.Fatalf("case %d (%s): want %s, generator says %s", i, c.Input, v, c.Expected)
		}
		seen[c.Expected] = true
	}
	for _, want := range []string{"0", "2147483647", "-2147483648"} {
		if !seen[want] {
			t.Fatalf("no case gives %s", want)
		}
	}
}
