package testgen

import (
	"slices"
	"strconv"
	"strings"
	"testing"
)

func TestKthLargest(t *testing.T) {
	for i, c := range cases(t, "kth-largest", 50) {
		nums, rest := parseArray(t, c.Input)
		k, _ := strconv.Atoi(rest[0])
		s := slices.Clone(nums)
		slices.Sort(s)
		slices.Reverse(s)
		if strconv.Itoa(s[k-1]) != c.Expected {
			t.Fatalf("case %d: k-th largest differs", i)
		}
	}
}

func TestSmashStones(t *testing.T) {
	for i, c := range cases(t, "smash-stones", 50) {
		stones, _ := parseArray(t, c.Input)
		if len(stones) > 1500 {
			continue
		}
		s := slices.Clone(stones)
		for len(s) > 1 {
			slices.Sort(s)
			y, x := s[len(s)-1], s[len(s)-2]
			s = s[:len(s)-2]
			if y != x {
				s = append(s, y-x)
			}
		}
		want := 0
		if len(s) == 1 {
			want = s[0]
		}
		if strconv.Itoa(want) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, want, c.Expected)
		}
	}
}

func TestRopeCost(t *testing.T) {
	for i, c := range cases(t, "rope-cost", 50) {
		ropes, _ := parseArray(t, c.Input)
		if len(ropes) > 1500 {
			continue
		}
		s := slices.Clone(ropes)
		cost := 0
		for len(s) > 1 {
			slices.Sort(s)
			j := s[0] + s[1]
			cost += j
			s = append([]int{j}, s[2:]...)
		}
		if strconv.Itoa(cost) != c.Expected || cost > 1<<31-1 {
			t.Fatalf("case %d: want %d, generator says %s", i, cost, c.Expected)
		}
	}
}

func TestClosestPoints(t *testing.T) {
	for i, c := range cases(t, "closest-points", 60) {
		lines := strings.Split(c.Input, "\n")
		n, _ := strconv.Atoi(lines[0])
		k, _ := strconv.Atoi(lines[n+1])
		type pt struct{ d, i int }
		var pts []pt
		for j := 0; j < n; j++ {
			f := strings.Fields(lines[j+1])
			x, _ := strconv.Atoi(f[0])
			y, _ := strconv.Atoi(f[1])
			pts = append(pts, pt{x*x + y*y, j})
		}
		slices.SortFunc(pts, func(a, b pt) int { return a.d - b.d })
		if k < n && pts[k-1].d == pts[k].d {
			t.Fatalf("case %d: tie at the cut", i)
		}
		var idx []int
		for _, p := range pts[:k] {
			idx = append(idx, p.i)
		}
		slices.Sort(idx)
		if listOutput(idx) != c.Expected {
			t.Fatalf("case %d: closest points differ", i)
		}
	}
}

// simulate runs the most-remaining task that is not cooling down, else idles.
func simulate(tasks string, n int) int {
	var left [26]int
	for i := 0; i < len(tasks); i++ {
		left[tasks[i]-'A']++
	}
	var ready [26]int
	remaining := len(tasks)
	for time := 0; ; time++ {
		if remaining == 0 {
			return time
		}
		best := -1
		for t := 0; t < 26; t++ {
			if left[t] > 0 && ready[t] <= time && (best < 0 || left[t] > left[best]) {
				best = t
			}
		}
		if best >= 0 {
			left[best]--
			remaining--
			ready[best] = time + n + 1
		}
	}
}

func TestCooldown(t *testing.T) {
	for i, c := range cases(t, "cooldown", 60) {
		lines := strings.Split(c.Input, "\n")
		n, _ := strconv.Atoi(lines[1])
		if len(lines[0]) > 3000 {
			continue
		}
		if want := simulate(lines[0], n); strconv.Itoa(want) != c.Expected {
			t.Fatalf("case %d: simulation gives %d, generator says %s", i, want, c.Expected)
		}
	}
}

func TestMergeK(t *testing.T) {
	for i, c := range cases(t, "merge-k", 50) {
		lines := strings.Split(c.Input, "\n")
		k, _ := strconv.Atoi(lines[0])
		rest := lines[1:]
		var all []int
		for j := 0; j < k; j++ {
			list, after := parseArray(t, strings.Join(rest, "\n"))
			if !slices.IsSorted(list) {
				t.Fatalf("case %d: list %d not sorted", i, j)
			}
			all = append(all, list...)
			rest = after
		}
		slices.Sort(all)
		if len(rest) != 0 || listOutput(all) != c.Expected {
			t.Fatalf("case %d: merge differs", i)
		}
	}
}

func TestRunningMedian(t *testing.T) {
	for i, c := range cases(t, "running-median", 40) {
		var vals []int
		var out []string
		for _, line := range strings.Split(c.Input, "\n")[1:] {
			f := strings.Fields(line)
			switch f[0] {
			case "new":
				vals = nil
				out = append(out, "null")
			case "add":
				v, _ := strconv.Atoi(f[1])
				vals = append(vals, v)
				out = append(out, "null")
			default:
				s := slices.Clone(vals)
				slices.Sort(s)
				m := float64(s[len(s)/2])
				if len(s)%2 == 0 {
					m = float64(s[len(s)/2-1]+s[len(s)/2]) / 2
				}
				out = append(out, strconv.FormatFloat(m, 'f', 1, 64))
			}
		}
		if strings.Join(out, " ") != c.Expected {
			t.Fatalf("case %d: running median differs", i)
		}
	}
}
