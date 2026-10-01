package testgen

import (
	"slices"
	"strconv"
	"strings"
	"testing"
)

func TestSortedPair(t *testing.T) {
	for i, c := range cases(t, "sorted-pair", 50) {
		nums, rest := parseArray(t, c.Input)
		target, _ := strconv.Atoi(rest[0])
		if !slices.IsSorted(nums) || target > 1<<31-1 || target < -1<<31 {
			t.Fatalf("case %d: not sorted or target too large", i)
		}
		var found []string
		for p := 0; p < len(nums) && len(nums) <= 2000; p++ {
			for q := p + 1; q < len(nums); q++ {
				if nums[p]+nums[q] == target {
					found = append(found, strconv.Itoa(p)+" "+strconv.Itoa(q))
				}
			}
		}
		if len(nums) <= 2000 && (len(found) != 1 || found[0] != c.Expected) {
			t.Fatalf("case %d: brute force finds %v, generator says %s", i, found, c.Expected)
		}
	}
}

func TestPalindrome(t *testing.T) {
	seen := map[string]bool{}
	for i, c := range cases(t, "palindrome", 80) {
		text := c.Input
		if text == "" || strings.Contains(text, "\n") || strings.TrimSpace(text) != text {
			t.Fatalf("case %d: bad line %q", i, text)
		}
		var kept []byte
		for j := 0; j < len(text); j++ {
			ch := text[j]
			if ch >= 'A' && ch <= 'Z' {
				ch += 'a' - 'A'
			}
			if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') {
				kept = append(kept, ch)
			}
		}
		rev := slices.Clone(kept)
		slices.Reverse(rev)
		if want := boolText(string(kept) == string(rev)); want != c.Expected {
			t.Fatalf("case %d: want %s, generator says %s for %q", i, want, c.Expected, text)
		}
		seen[c.Expected] = true
	}
	if !seen["true"] || !seen["false"] {
		t.Fatal("both answers must appear")
	}
}

func TestRemoveRepeats(t *testing.T) {
	for i, c := range cases(t, "remove-repeats", 60) {
		nums, _ := parseArray(t, c.Input)
		var want []int
		seen := map[int]bool{}
		for _, v := range nums {
			if !seen[v] {
				seen[v] = true
				want = append(want, v)
			}
		}
		if !slices.IsSorted(nums) || listOutput(want) != c.Expected {
			t.Fatalf("case %d: want %v, generator says %s", i, want, c.Expected)
		}
	}
}

func TestMergeSorted(t *testing.T) {
	for i, c := range cases(t, "merge-sorted", 60) {
		a, rest := parseArray(t, c.Input)
		b, _ := parseArray(t, strings.Join(rest, "\n"))
		var merged []int
		p, q := 0, 0
		for p < len(a) || q < len(b) {
			if q == len(b) || (p < len(a) && a[p] <= b[q]) {
				merged = append(merged, a[p])
				p++
			} else {
				merged = append(merged, b[q])
				q++
			}
		}
		if !slices.IsSorted(a) || !slices.IsSorted(b) || listOutput(merged) != c.Expected {
			t.Fatalf("case %d: merge differs", i)
		}
	}
}

func TestThreeWay(t *testing.T) {
	for i, c := range cases(t, "three-way", 50) {
		nums, _ := parseArray(t, c.Input)
		var count [3]int
		for _, v := range nums {
			count[v]++
		}
		var want []int
		for v := 0; v < 3; v++ {
			for j := 0; j < count[v]; j++ {
				want = append(want, v)
			}
		}
		if listOutput(want) != c.Expected {
			t.Fatalf("case %d: counting sort differs", i)
		}
	}
}

func TestTripleSum(t *testing.T) {
	nonEmpty := 0
	for i, c := range cases(t, "triple-sum", 60) {
		nums, rest := parseArray(t, c.Input)
		target, _ := strconv.Atoi(rest[0])
		if len(nums) > 500 {
			t.Fatalf("case %d: %d values, limit 500", i, len(nums))
		}
		if len(nums) > 80 {
			continue
		}
		set := map[[3]int]bool{}
		for p := range nums {
			for q := p + 1; q < len(nums); q++ {
				for r := q + 1; r < len(nums); r++ {
					if nums[p]+nums[q]+nums[r] == target {
						tr := []int{nums[p], nums[q], nums[r]}
						slices.Sort(tr)
						set[[3]int(tr)] = true
					}
				}
			}
		}
		var lists [][]int
		for tr := range set {
			lists = append(lists, []int{tr[0], tr[1], tr[2]})
		}
		if want := triplesOutput(lists); want != c.Expected {
			t.Fatalf("case %d: want %s, generator says %s", i, want, c.Expected)
		}
		if len(set) > 0 {
			nonEmpty++
		}
	}
	if nonEmpty == 0 {
		t.Fatal("no case has any triple")
	}
}

func TestContainer(t *testing.T) {
	for i, c := range cases(t, "container", 50) {
		h, _ := parseArray(t, c.Input)
		if len(h) > 1500 {
			continue
		}
		best := 0
		for p := range h {
			for q := p + 1; q < len(h); q++ {
				best = max(best, min(h[p], h[q])*(q-p))
			}
		}
		if strconv.Itoa(best) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, best, c.Expected)
		}
	}
}

func TestRainTrap(t *testing.T) {
	positive := 0
	for i, c := range cases(t, "rain-trap", 50) {
		h, _ := parseArray(t, c.Input)
		if len(h) > 1500 {
			continue
		}
		total := 0
		for p := range h {
			left, right := 0, 0
			for q := 0; q <= p; q++ {
				left = max(left, h[q])
			}
			for q := p; q < len(h); q++ {
				right = max(right, h[q])
			}
			total += min(left, right) - h[p]
		}
		if strconv.Itoa(total) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, total, c.Expected)
		}
		if total > 0 {
			positive++
		}
	}
	if positive == 0 {
		t.Fatal("no case traps any water")
	}
}
