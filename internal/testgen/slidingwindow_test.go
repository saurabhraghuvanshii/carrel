package testgen

import (
	"slices"
	"strconv"
	"strings"
	"testing"
)

func TestWindowSum(t *testing.T) {
	for i, c := range cases(t, "window-sum", 50) {
		nums, rest := parseArray(t, c.Input)
		k, _ := strconv.Atoi(rest[0])
		if len(nums) > 1500 {
			continue
		}
		best := 0
		for s := 0; s+k <= len(nums); s++ {
			sum := 0
			for _, v := range nums[s : s+k] {
				sum += v
			}
			if s == 0 || sum > best {
				best = sum
			}
		}
		if strconv.Itoa(best) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, best, c.Expected)
		}
	}
}

// allDifferent reports whether no byte repeats in s.
func allDifferent(s string) bool {
	seen := map[byte]bool{}
	for i := 0; i < len(s); i++ {
		if seen[s[i]] {
			return false
		}
		seen[s[i]] = true
	}
	return true
}

func TestUniqueStretch(t *testing.T) {
	for i, c := range cases(t, "unique-stretch", 50) {
		text := c.Input
		if len(text) > 800 {
			continue
		}
		best := 0
		for a := 0; a < len(text); a++ {
			for b := a + 1; b <= len(text) && b-a <= 36; b++ {
				if allDifferent(text[a:b]) {
					best = max(best, b-a)
				}
			}
		}
		if strconv.Itoa(best) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, best, c.Expected)
		}
	}
}

func TestTwoKinds(t *testing.T) {
	for i, c := range cases(t, "two-kinds", 50) {
		nums, _ := parseArray(t, c.Input)
		if len(nums) > 800 {
			continue
		}
		best := 0
		for a := range nums {
			kinds := map[int]bool{}
			for b := a; b < len(nums); b++ {
				kinds[nums[b]] = true
				if len(kinds) > 2 {
					break
				}
				best = max(best, b-a+1)
			}
		}
		if strconv.Itoa(best) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, best, c.Expected)
		}
	}
}

func TestKChanges(t *testing.T) {
	for i, c := range cases(t, "k-changes", 50) {
		lines := strings.Split(c.Input, "\n")
		text := lines[0]
		k, _ := strconv.Atoi(lines[1])
		if len(text) > 400 {
			continue
		}
		best := 0
		for a := 0; a < len(text); a++ {
			for b := a + 1; b <= len(text); b++ {
				for letter := byte('A'); letter <= 'Z'; letter++ {
					changes := 0
					for j := a; j < b; j++ {
						if text[j] != letter {
							changes++
						}
					}
					if changes <= k {
						best = max(best, b-a)
					}
				}
			}
		}
		if strconv.Itoa(best) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, best, c.Expected)
		}
	}
}

func TestSumAtLeast(t *testing.T) {
	zero := 0
	for i, c := range cases(t, "sum-at-least", 60) {
		nums, rest := parseArray(t, c.Input)
		target, _ := strconv.Atoi(rest[0])
		if len(nums) > 1500 {
			continue
		}
		best := 0
		for a := range nums {
			sum := 0
			for b := a; b < len(nums); b++ {
				sum += nums[b]
				if sum >= target {
					if best == 0 || b-a+1 < best {
						best = b - a + 1
					}
					break
				}
			}
		}
		if strconv.Itoa(best) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, best, c.Expected)
		}
		if best == 0 {
			zero++
		}
	}
	if zero == 0 {
		t.Fatal("no case where the target cannot be reached")
	}
}

func TestContainsRearrangement(t *testing.T) {
	seen := map[string]bool{}
	for i, c := range cases(t, "contains-rearrangement", 60) {
		lines := strings.Split(c.Input, "\n")
		pattern, text := lines[0], lines[1]
		want := false
		key := sortedLetters(pattern)
		for a := 0; a+len(pattern) <= len(text) && !want; a++ {
			want = sortedLetters(text[a:a+len(pattern)]) == key
		}
		if boolText(want) != c.Expected {
			t.Fatalf("case %d: want %v, generator says %s", i, want, c.Expected)
		}
		seen[c.Expected] = true
	}
	if !seen["true"] || !seen["false"] {
		t.Fatal("both answers must appear")
	}
}

func covers(window, pattern string) bool {
	count := map[byte]int{}
	for i := 0; i < len(window); i++ {
		count[window[i]]++
	}
	for i := 0; i < len(pattern); i++ {
		count[pattern[i]]--
		if count[pattern[i]] < 0 {
			return false
		}
	}
	return true
}

func TestCoveringWindow(t *testing.T) {
	empty := 0
	for i, c := range cases(t, "covering-window", 60) {
		lines := strings.Split(c.Input, "\n")
		text, pattern := lines[0], lines[1]
		if len(text) > 400 {
			continue
		}
		want := ""
		found := false
		for length := 1; length <= len(text) && !found; length++ {
			for a := 0; a+length <= len(text); a++ {
				if covers(text[a:a+length], pattern) {
					want, found = text[a:a+length], true
					break
				}
			}
		}
		if `"`+want+`"` != c.Expected {
			t.Fatalf("case %d: want %q, generator says %s", i, want, c.Expected)
		}
		if !found {
			empty++
		}
	}
	if empty == 0 {
		t.Fatal("no case without a window")
	}
}

func TestWindowMax(t *testing.T) {
	for i, c := range cases(t, "window-max", 50) {
		nums, rest := parseArray(t, c.Input)
		k, _ := strconv.Atoi(rest[0])
		if len(nums) > 1500 {
			continue
		}
		var want []int
		for a := 0; a+k <= len(nums); a++ {
			want = append(want, slices.Max(nums[a:a+k]))
		}
		if listOutput(want) != c.Expected {
			t.Fatalf("case %d: window maximums differ", i)
		}
	}
}
