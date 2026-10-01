package testgen

import (
	"slices"
	"strconv"
	"strings"
	"testing"
)

// parseArray reads "n" and a line of n values, returning the values and the
// lines after them.
func parseArray(t *testing.T, input string) ([]int, []string) {
	t.Helper()
	lines := strings.Split(input, "\n")
	n, _ := strconv.Atoi(lines[0])
	fields := strings.Fields(lines[1])
	if len(fields) != n {
		t.Fatalf("array says %d values, has %d", n, len(fields))
	}
	nums := make([]int, n)
	for i, f := range fields {
		nums[i], _ = strconv.Atoi(f)
	}
	return nums, lines[2:]
}

func cases(t *testing.T, name string, n int) []Case {
	t.Helper()
	cs, err := Generate(name, 21, n)
	if err != nil {
		t.Fatal(err)
	}
	return cs
}

func TestSpotDuplicate(t *testing.T) {
	seen := map[string]bool{}
	for i, c := range cases(t, "spot-duplicate", 60) {
		nums, _ := parseArray(t, c.Input)
		want := false
		for p := 0; p < len(nums) && len(nums) <= 2000; p++ {
			for q := p + 1; q < len(nums); q++ {
				want = want || nums[p] == nums[q]
			}
		}
		if len(nums) > 2000 {
			s := slices.Clone(nums)
			slices.Sort(s)
			want = len(slices.Compact(s)) != len(nums)
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

func TestSameLetters(t *testing.T) {
	seen := map[string]bool{}
	for i, c := range cases(t, "same-letters", 60) {
		w := strings.Split(c.Input, "\n")
		var count [26]int
		for _, ch := range w[0] {
			count[ch-'a']++
		}
		for _, ch := range w[1] {
			count[ch-'a']--
		}
		want := count == [26]int{}
		if boolText(want) != c.Expected || w[0] == "" || w[1] == "" {
			t.Fatalf("case %d: want %v, generator says %s", i, want, c.Expected)
		}
		seen[c.Expected] = true
	}
	if !seen["true"] || !seen["false"] {
		t.Fatal("both answers must appear")
	}
}

func TestGroupLetters(t *testing.T) {
	for i, c := range cases(t, "group-letters", 40) {
		lines := strings.Split(c.Input, "\n")
		words := strings.Fields(lines[1])
		byCount := map[[26]int][]string{}
		for _, w := range words {
			var k [26]int
			for _, ch := range w {
				k[ch-'a']++
			}
			byCount[k] = append(byCount[k], w)
		}
		var groups [][]string
		for _, g := range byCount {
			groups = append(groups, g)
		}
		if want := groupsOutput(groups); want != c.Expected {
			t.Fatalf("case %d: want %s, generator says %s", i, want, c.Expected)
		}
	}
}

func TestMostFrequentIsUnique(t *testing.T) {
	for i, c := range cases(t, "most-frequent", 60) {
		nums, rest := parseArray(t, c.Input)
		k, _ := strconv.Atoi(rest[0])
		count := map[int]int{}
		for _, v := range nums {
			count[v]++
		}
		var freqs []int
		for _, f := range count {
			freqs = append(freqs, f)
		}
		slices.Sort(freqs)
		slices.Reverse(freqs)
		if k < 1 || k > len(freqs) || (k < len(freqs) && freqs[k-1] == freqs[k]) {
			t.Fatalf("case %d: k=%d does not give a unique answer %v", i, k, freqs)
		}
		var top []int
		for v, f := range count {
			if f >= freqs[k-1] {
				top = append(top, v)
			}
		}
		slices.Sort(top)
		if want := listOutput(top); want != c.Expected {
			t.Fatalf("case %d: want %s, generator says %s", i, want, c.Expected)
		}
	}
}

func TestProductOthers(t *testing.T) {
	for i, c := range cases(t, "product-others", 50) {
		nums, _ := parseArray(t, c.Input)
		if len(nums) > 400 {
			continue
		}
		out := make([]int, len(nums))
		for p := range nums {
			out[p] = 1
			for q, v := range nums {
				if q != p {
					out[p] *= v
				}
			}
			if out[p] > 1<<31-1 || out[p] < -1<<31 {
				t.Fatalf("case %d: product %d does not fit in 32 bits", i, out[p])
			}
		}
		if want := listOutput(out); want != c.Expected {
			t.Fatalf("case %d: brute force differs", i)
		}
	}
}

func TestLongestRun(t *testing.T) {
	for i, c := range cases(t, "longest-run", 60) {
		nums, _ := parseArray(t, c.Input)
		set := map[int]bool{}
		for _, v := range nums {
			set[v] = true
		}
		best := 0
		for v := range set {
			if set[v-1] {
				continue
			}
			l := 1
			for set[v+l] {
				l++
			}
			best = max(best, l)
		}
		if strconv.Itoa(best) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, best, c.Expected)
		}
	}
}

func TestSubarraySum(t *testing.T) {
	for i, c := range cases(t, "subarray-sum", 50) {
		nums, rest := parseArray(t, c.Input)
		if len(nums) > 1500 {
			continue
		}
		target, _ := strconv.Atoi(rest[0])
		count := 0
		for p := range nums {
			sum := 0
			for q := p; q < len(nums); q++ {
				sum += nums[q]
				if sum == target {
					count++
				}
			}
		}
		if strconv.Itoa(count) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, count, c.Expected)
		}
	}
}

func TestRotateGrid(t *testing.T) {
	for i, c := range cases(t, "rotate-grid", 50) {
		lines := strings.Split(c.Input, "\n")
		dims := strings.Fields(lines[0])
		n, _ := strconv.Atoi(dims[0])
		if dims[0] != dims[1] || len(lines) != n+1 || n > 200 {
			t.Fatalf("case %d: bad grid shape", i)
		}
		grid := make([][]int, n)
		for r := range grid {
			for _, f := range strings.Fields(lines[r+1]) {
				v, _ := strconv.Atoi(f)
				grid[r] = append(grid[r], v)
			}
		}
		// Transpose, then reverse each row.
		for r := 0; r < n; r++ {
			for c := r + 1; c < n; c++ {
				grid[r][c], grid[c][r] = grid[c][r], grid[r][c]
			}
		}
		for _, row := range grid {
			slices.Reverse(row)
		}
		if want := gridOutput(grid); want != c.Expected {
			t.Fatalf("case %d: rotation differs", i)
		}
	}
}

func TestMissingPositive(t *testing.T) {
	for i, c := range cases(t, "missing-positive", 60) {
		nums, _ := parseArray(t, c.Input)
		s := slices.Clone(nums)
		slices.Sort(s)
		want := 1
		for _, v := range s {
			if v == want {
				want++
			}
		}
		if strconv.Itoa(want) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, want, c.Expected)
		}
	}
}
