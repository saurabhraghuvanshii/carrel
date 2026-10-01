package testgen

import (
	"strconv"
	"strings"
	"testing"
)

func TestBrackets(t *testing.T) {
	seen := map[string]bool{}
	for i, c := range cases(t, "brackets", 80) {
		s := c.Input
		// Remove matched pairs until nothing changes.
		for {
			r := strings.NewReplacer("()", "", "[]", "", "{}", "").Replace(s)
			if r == s {
				break
			}
			s = r
		}
		if want := boolText(s == ""); want != c.Expected || c.Input == "" {
			t.Fatalf("case %d: want %s, generator says %s for %q", i, want, c.Expected, c.Input)
		}
		seen[c.Expected] = true
	}
	if !seen["true"] || !seen["false"] {
		t.Fatal("both answers must appear")
	}
}

func TestAdjacentTwins(t *testing.T) {
	emptied := 0
	for i, c := range cases(t, "adjacent-twins", 60) {
		s := c.Input
		for changed := true; changed; {
			changed = false
			for j := 0; j+1 < len(s); j++ {
				if s[j] == s[j+1] {
					s, changed = s[:j]+s[j+2:], true
					break
				}
			}
		}
		if `"`+s+`"` != c.Expected {
			t.Fatalf("case %d: want %q, generator says %s", i, s, c.Expected)
		}
		if s == "" {
			emptied++
		}
	}
	if emptied == 0 {
		t.Fatal("no case cancels completely")
	}
}

func TestMinStack(t *testing.T) {
	for i, c := range cases(t, "min-stack", 40) {
		lines := strings.Split(c.Input, "\n")
		var stack []int
		var out []string
		for _, line := range lines[1:] {
			f := strings.Fields(line)
			switch f[0] {
			case "new":
				stack = nil
				out = append(out, "null")
			case "push":
				v, _ := strconv.Atoi(f[1])
				stack = append(stack, v)
				out = append(out, "null")
			case "pop":
				stack = stack[:len(stack)-1]
				out = append(out, "null")
			case "top":
				out = append(out, strconv.Itoa(stack[len(stack)-1]))
			case "min":
				m := stack[0]
				for _, v := range stack {
					m = min(m, v)
				}
				out = append(out, strconv.Itoa(m))
			}
		}
		if strings.Join(out, " ") != c.Expected {
			t.Fatalf("case %d: scanning for the minimum gives a different answer", i)
		}
	}
}

func TestPostfix(t *testing.T) {
	for i, c := range cases(t, "postfix", 60) {
		lines := strings.Split(c.Input, "\n")
		tokens := strings.Fields(lines[1])
		var stack []int
		for _, tok := range tokens {
			if len(tok) == 1 && strings.Contains("+-*/", tok) {
				b, a := stack[len(stack)-1], stack[len(stack)-2]
				stack = stack[:len(stack)-2]
				var v int
				switch tok {
				case "+":
					v = a + b
				case "-":
					v = a - b
				case "*":
					v = a * b
				default:
					if b == 0 {
						t.Fatalf("case %d: division by zero", i)
					}
					v = a / b
				}
				if v > exprLimit || v < -exprLimit {
					t.Fatalf("case %d: step value %d too large", i, v)
				}
				stack = append(stack, v)
			} else {
				v, err := strconv.Atoi(tok)
				if err != nil {
					t.Fatalf("case %d: bad token %q", i, tok)
				}
				stack = append(stack, v)
			}
		}
		if len(stack) != 1 || strconv.Itoa(stack[0]) != c.Expected {
			t.Fatalf("case %d: stack evaluation differs", i)
		}
	}
}

// pointlessByDepth checks each bracket pair: it is pointless if nothing sits
// directly inside it, outside any inner bracket pair.
func pointlessByDepth(expr string) bool {
	for i := 0; i < len(expr); i++ {
		if expr[i] != '(' {
			continue
		}
		depth, direct := 0, false
		for j := i + 1; j < len(expr); j++ {
			if expr[j] == '(' {
				depth++
			} else if expr[j] == ')' {
				if depth == 0 {
					break
				}
				depth--
			} else if depth == 0 {
				direct = true
			}
		}
		if !direct {
			return true
		}
	}
	return false
}

func TestRedundantBrackets(t *testing.T) {
	seen := map[string]bool{}
	for i, c := range cases(t, "redundant-brackets", 80) {
		if want := boolText(pointlessByDepth(c.Input)); want != c.Expected {
			t.Fatalf("case %d: want %s, generator says %s for %q", i, want, c.Expected, c.Input)
		}
		seen[c.Expected] = true
	}
	if !seen["true"] || !seen["false"] {
		t.Fatal("both answers must appear")
	}
}

func TestWarmerDay(t *testing.T) {
	for i, c := range cases(t, "warmer-day", 50) {
		temps, _ := parseArray(t, c.Input)
		if len(temps) > 2000 {
			continue
		}
		want := make([]int, len(temps))
		for a := range temps {
			for b := a + 1; b < len(temps); b++ {
				if temps[b] > temps[a] {
					want[a] = b - a
					break
				}
			}
		}
		if listOutput(want) != c.Expected {
			t.Fatalf("case %d: brute force differs", i)
		}
	}
}

func TestPriceSpan(t *testing.T) {
	for i, c := range cases(t, "price-span", 50) {
		prices, _ := parseArray(t, c.Input)
		if len(prices) > 2000 {
			continue
		}
		want := make([]int, len(prices))
		for a := range prices {
			for b := a; b >= 0 && prices[b] <= prices[a]; b-- {
				want[a]++
			}
		}
		if listOutput(want) != c.Expected {
			t.Fatalf("case %d: brute force differs", i)
		}
	}
}

func TestDecodeText(t *testing.T) {
	for i, c := range cases(t, "decode-text", 60) {
		// Expand the innermost k[...] again and again.
		s := c.Input
		for strings.Contains(s, "[") {
			close := strings.IndexByte(s, ']')
			open := strings.LastIndexByte(s[:close], '[')
			start := open
			for start > 0 && s[start-1] >= '0' && s[start-1] <= '9' {
				start--
			}
			k, _ := strconv.Atoi(s[start:open])
			s = s[:start] + strings.Repeat(s[open+1:close], k) + s[close+1:]
		}
		if s != c.Expected || len(s) > decodedLimit || s == "" {
			t.Fatalf("case %d: expansion differs or is too long (%d)", i, len(s))
		}
	}
}

func TestRobotRail(t *testing.T) {
	for i, c := range cases(t, "robot-rail", 50) {
		robots, _ := parseArray(t, c.Input)
		if len(robots) > 1500 {
			continue
		}
		// Collide the first meeting pair again and again.
		r := robots
		for {
			hit := -1
			for j := 0; j+1 < len(r); j++ {
				if r[j] > 0 && r[j+1] < 0 {
					hit = j
					break
				}
			}
			if hit < 0 {
				break
			}
			a, b := r[hit], -r[hit+1]
			next := append([]int{}, r[:hit]...)
			if a > b {
				next = append(next, r[hit])
			} else if b > a {
				next = append(next, r[hit+1])
			}
			r = append(next, r[hit+2:]...)
		}
		if listOutput(r) != c.Expected {
			t.Fatalf("case %d: want %v, generator says %s", i, r, c.Expected)
		}
	}
}

func TestHistogram(t *testing.T) {
	for i, c := range cases(t, "histogram", 50) {
		h, _ := parseArray(t, c.Input)
		if len(h) > 1500 {
			continue
		}
		best := 0
		for a := range h {
			low := h[a]
			for b := a; b < len(h); b++ {
				low = min(low, h[b])
				best = max(best, low*(b-a+1))
			}
		}
		if strconv.Itoa(best) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, best, c.Expected)
		}
	}
}
