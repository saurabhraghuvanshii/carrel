package testgen

import (
	"slices"
	"strconv"
	"strings"
	"testing"
)

func parsePairs(t *testing.T, input string) [][2]int {
	t.Helper()
	lines := strings.Split(input, "\n")
	n, _ := strconv.Atoi(lines[0])
	ps := make([][2]int, n)
	for i := range ps {
		f := strings.Fields(lines[i+1])
		ps[i][0], _ = strconv.Atoi(f[0])
		ps[i][1], _ = strconv.Atoi(f[1])
	}
	return ps
}

// permutations calls visit with every ordering of 0 to n-1.
func permutations(n int, visit func(p []int)) {
	p := make([]int, n)
	for i := range p {
		p[i] = i
	}
	var walk func(k int)
	walk = func(k int) {
		if k == n {
			visit(p)
			return
		}
		for i := k; i < n; i++ {
			p[k], p[i] = p[i], p[k]
			walk(k + 1)
			p[k], p[i] = p[i], p[k]
		}
	}
	walk(0)
}

func TestCookieAssignment(t *testing.T) {
	for i, c := range cases(t, "cookie-assignment", 60) {
		greed, rest := parseArray(t, c.Input)
		cookies, _ := parseArray(t, strings.Join(rest, "\n"))
		if len(greed)*len(cookies) > 1_000_000 {
			continue
		}
		// Maximum matching by augmenting paths.
		owner := make([]int, len(cookies))
		for j := range owner {
			owner[j] = -1
		}
		var try func(child int, seen []bool) bool
		try = func(child int, seen []bool) bool {
			for j, size := range cookies {
				if size >= greed[child] && !seen[j] {
					seen[j] = true
					if owner[j] < 0 || try(owner[j], seen) {
						owner[j] = child
						return true
					}
				}
			}
			return false
		}
		fed := 0
		for child := range greed {
			if try(child, make([]bool, len(cookies))) {
				fed++
			}
		}
		if strconv.Itoa(fed) != c.Expected {
			t.Fatalf("case %d: matching finds %d, generator says %s", i, fed, c.Expected)
		}
	}
}

func TestPairGaps(t *testing.T) {
	for i, c := range cases(t, "pair-gaps", 60) {
		a, rest := parseArray(t, c.Input)
		b, _ := parseArray(t, strings.Join(rest, "\n"))
		if len(a) > 7 {
			continue
		}
		best := -1
		permutations(len(a), func(p []int) {
			sum := 0
			for k := range a {
				sum += max(a[k]-b[p[k]], b[p[k]]-a[k])
			}
			if best < 0 || sum < best {
				best = sum
			}
		})
		if strconv.Itoa(best) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, best, c.Expected)
		}
	}
}

func TestPairChain(t *testing.T) {
	for i, c := range cases(t, "pair-chain", 60) {
		ps := parsePairs(t, c.Input)
		if len(ps) > 2000 {
			continue
		}
		// Longest chain ending at each pair, pairs taken by start.
		s := slices.Clone(ps)
		slices.SortFunc(s, func(x, y [2]int) int { return x[0] - y[0] })
		best := make([]int, len(s))
		top := 0
		for k := range s {
			best[k] = 1
			for j := 0; j < k; j++ {
				if s[j][1] < s[k][0] {
					best[k] = max(best[k], best[j]+1)
				}
			}
			top = max(top, best[k])
		}
		if strconv.Itoa(top) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, top, c.Expected)
		}
	}
}

func TestJobDeadlines(t *testing.T) {
	for i, c := range cases(t, "job-deadlines", 60) {
		jobs := parsePairs(t, c.Input)
		if len(jobs) > 14 {
			continue
		}
		// A set of jobs fits when, by deadline, the k-th has deadline >= k.
		best := 0
		for mask := 0; mask < 1<<len(jobs); mask++ {
			var picked [][2]int
			profit := 0
			for j, job := range jobs {
				if mask>>j&1 == 1 {
					picked = append(picked, job)
					profit += job[1]
				}
			}
			slices.SortFunc(picked, func(x, y [2]int) int { return x[0] - y[0] })
			ok := true
			for k, job := range picked {
				ok = ok && job[0] >= k+1
			}
			if ok {
				best = max(best, profit)
			}
		}
		if strconv.Itoa(best) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, best, c.Expected)
		}
	}
}

func TestCanReachEnd(t *testing.T) {
	seen := map[string]bool{}
	for i, c := range cases(t, "can-reach-end", 80) {
		nums, _ := parseArray(t, c.Input)
		// Breadth-first search over positions.
		visited := make([]bool, len(nums))
		visited[0] = true
		queue := []int{0}
		for len(queue) > 0 {
			p := queue[0]
			queue = queue[1:]
			for step := 1; step <= nums[p] && p+step < len(nums); step++ {
				if !visited[p+step] {
					visited[p+step] = true
					queue = append(queue, p+step)
				}
			}
		}
		if boolText(visited[len(nums)-1]) != c.Expected {
			t.Fatalf("case %d: want %v, generator says %s", i, visited[len(nums)-1], c.Expected)
		}
		seen[c.Expected] = true
	}
	if !seen["true"] || !seen["false"] {
		t.Fatal("both answers must appear")
	}
}

func TestFuelCircuit(t *testing.T) {
	found := 0
	for i, c := range cases(t, "fuel-circuit", 60) {
		gas, rest := parseArray(t, c.Input)
		cost, _ := parseArray(t, strings.Join(rest, "\n"))
		n := len(gas)
		if n > 1000 {
			continue
		}
		want := -1
		for s := 0; s < n && want < 0; s++ {
			tank, ok := 0, true
			for k := 0; k < n && ok; k++ {
				tank += gas[(s+k)%n] - cost[(s+k)%n]
				ok = tank >= 0
			}
			if ok {
				want = s
			}
		}
		if strconv.Itoa(want) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, want, c.Expected)
		}
		if want > 0 {
			found++
		}
	}
	if found == 0 {
		t.Fatal("no case starts past station 0")
	}
}

func TestLabelPartitions(t *testing.T) {
	for i, c := range cases(t, "label-partitions", 60) {
		w := c.Input
		if len(w) > 14 {
			continue
		}
		// Try every set of cuts; keep the valid one with the most pieces.
		var best []int
		for mask := 0; mask < 1<<(len(w)-1); mask++ {
			var sizes []int
			piece := map[byte]int{}
			start, ok := 0, true
			for j := 1; j <= len(w); j++ {
				if j == len(w) || mask>>(j-1)&1 == 1 {
					for k := start; k < j; k++ {
						if p, seen := piece[w[k]]; seen && p != len(sizes) {
							ok = false
						}
						piece[w[k]] = len(sizes)
					}
					sizes = append(sizes, j-start)
					start = j
				}
			}
			if ok && len(sizes) > len(best) {
				best = sizes
			}
		}
		if listOutput(best) != c.Expected {
			t.Fatalf("case %d (%s): want %v, generator says %s", i, w, best, c.Expected)
		}
	}
}

func TestChocolateCuts(t *testing.T) {
	for i, c := range cases(t, "chocolate-cuts", 60) {
		lines := strings.Split(c.Input, "\n")
		var h, v []int
		for _, f := range strings.Fields(lines[1]) {
			x, _ := strconv.Atoi(f)
			h = append(h, x)
		}
		for _, f := range strings.Fields(lines[2]) {
			x, _ := strconv.Atoi(f)
			v = append(v, x)
		}
		if len(h)+len(v) > 8 {
			continue
		}
		// Try every order of the cuts; a cut crosses one more piece than
		// the cuts already made the other way.
		best := -1
		permutations(len(h)+len(v), func(p []int) {
			total, hDone, vDone := 0, 0, 0
			for _, k := range p {
				if k < len(h) {
					total += h[k] * (vDone + 1)
					hDone++
				} else {
					total += v[k-len(h)] * (hDone + 1)
					vDone++
				}
			}
			if best < 0 || total < best {
				best = total
			}
		})
		if strconv.Itoa(best) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, best, c.Expected)
		}
	}
}

func TestCandyHandout(t *testing.T) {
	for i, c := range cases(t, "candy-handout", 60) {
		r, _ := parseArray(t, c.Input)
		if len(r) > 3000 {
			continue
		}
		// Raise counts until every rule holds.
		sweets := make([]int, len(r))
		for k := range sweets {
			sweets[k] = 1
		}
		for changed := true; changed; {
			changed = false
			for k := range r {
				if k > 0 && r[k] > r[k-1] && sweets[k] <= sweets[k-1] {
					sweets[k], changed = sweets[k-1]+1, true
				}
				if k+1 < len(r) && r[k] > r[k+1] && sweets[k] <= sweets[k+1] {
					sweets[k], changed = sweets[k+1]+1, true
				}
			}
		}
		total := 0
		for _, s := range sweets {
			total += s
		}
		if strconv.Itoa(total) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, total, c.Expected)
		}
	}
}
