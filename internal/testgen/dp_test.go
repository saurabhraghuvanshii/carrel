package testgen

import (
	"math/big"
	"math/bits"
	"strconv"
	"strings"
	"testing"
)

func TestStairWays(t *testing.T) {
	for i, c := range cases(t, "stair-ways", 40) {
		n, _ := strconv.Atoi(c.Input)
		if n > 22 {
			continue
		}
		// Each mask marks which of the n-1 gaps a step lands on; a climb is
		// valid when no two neighbouring gaps are skipped.
		count := 0
		for mask := 0; mask < 1<<(n-1); mask++ {
			full := mask | 1<<(n-1)
			ok, last := true, 0
			for p := 1; p <= n; p++ {
				if full>>(p-1)&1 == 1 {
					ok = ok && p-last <= 2
					last = p
				}
			}
			if ok {
				count++
			}
		}
		if strconv.Itoa(count) != c.Expected {
			t.Fatalf("case %d (n=%d): want %d, generator says %s", i, n, count, c.Expected)
		}
	}
}

func TestFriendPairs(t *testing.T) {
	for i, c := range cases(t, "friend-pairs", 40) {
		n, _ := strconv.Atoi(c.Input)
		if n > 300 {
			continue
		}
		// Sum over k pairs: n! / (k! (n-2k)! 2^k).
		total := new(big.Int)
		for k := 0; 2*k <= n; k++ {
			term := new(big.Int).MulRange(1, int64(n))
			term.Div(term, new(big.Int).MulRange(1, int64(k)))
			term.Div(term, new(big.Int).MulRange(1, int64(n-2*k)))
			term.Rsh(term, uint(k))
			total.Add(total, term)
		}
		if want := total.Mod(total, big.NewInt(mod)).String(); want != c.Expected {
			t.Fatalf("case %d (n=%d): want %s, generator says %s", i, n, want, c.Expected)
		}
	}
}

func TestHouseSkipper(t *testing.T) {
	for i, c := range cases(t, "house-skipper", 60) {
		nums, _ := parseArray(t, c.Input)
		if len(nums) > 16 {
			continue
		}
		best := 0
		for mask := 0; mask < 1<<len(nums); mask++ {
			if mask&(mask>>1) != 0 {
				continue
			}
			sum := 0
			for j, v := range nums {
				sum += v * (mask >> j & 1)
			}
			best = max(best, sum)
		}
		if strconv.Itoa(best) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, best, c.Expected)
		}
	}
}

func TestFewestCoins(t *testing.T) {
	seen := map[bool]bool{}
	for i, c := range cases(t, "fewest-coins", 60) {
		coins, rest := parseArray(t, c.Input)
		amount, _ := strconv.Atoi(rest[0])
		// Breadth-first search over amounts paid so far.
		dist := map[int]int{0: 0}
		queue := []int{0}
		for len(queue) > 0 {
			v := queue[0]
			queue = queue[1:]
			for _, x := range coins {
				if _, ok := dist[v+x]; !ok && v+x <= amount {
					dist[v+x] = dist[v] + 1
					queue = append(queue, v+x)
				}
			}
		}
		want, ok := dist[amount]
		if !ok {
			want = -1
		}
		if strconv.Itoa(want) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, want, c.Expected)
		}
		seen[ok] = true
	}
	if !seen[true] || !seen[false] {
		t.Fatal("both payable and unpayable amounts must appear")
	}
}

func TestCoinWays(t *testing.T) {
	for i, c := range cases(t, "coin-ways", 60) {
		coins, rest := parseArray(t, c.Input)
		amount, _ := strconv.Atoi(rest[0])
		// Count by the first coin used, then the ways with the coins after it.
		memo := map[[2]int]int{}
		var count func(k, left int) int
		count = func(k, left int) int {
			if left == 0 {
				return 1
			}
			if k == len(coins) {
				return 0
			}
			if v, ok := memo[[2]int{k, left}]; ok {
				return v
			}
			v := count(k+1, left)
			if coins[k] <= left {
				v += count(k, left-coins[k])
			}
			memo[[2]int{k, left}] = v
			return v
		}
		if want := count(0, amount); strconv.Itoa(want) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, want, c.Expected)
		}
	}
}

func TestGridPathCount(t *testing.T) {
	for i, c := range cases(t, "grid-path-count", 60) {
		lines := strings.Split(c.Input, "\n")
		g := make([][]string, len(lines)-1)
		for r := range g {
			g[r] = strings.Fields(lines[r+1])
		}
		if len(g)+len(g[0]) > 18 {
			continue
		}
		var walk func(r, col int) int
		walk = func(r, col int) int {
			if r >= len(g) || col >= len(g[0]) || g[r][col] == "1" {
				return 0
			}
			if r == len(g)-1 && col == len(g[0])-1 {
				return 1
			}
			return walk(r+1, col) + walk(r, col+1)
		}
		if want := walk(0, 0); strconv.Itoa(want) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, want, c.Expected)
		}
	}
}

func TestCountSearchTrees(t *testing.T) {
	for i, c := range cases(t, "count-search-trees", 40) {
		n, _ := strconv.Atoi(c.Input)
		// Catalan number: C(2n, n) / (n + 1).
		want := new(big.Int).Binomial(int64(2*n), int64(n))
		want.Div(want, big.NewInt(int64(n+1)))
		if want.String() != c.Expected {
			t.Fatalf("case %d (n=%d): want %s, generator says %s", i, n, want, c.Expected)
		}
	}
}

func TestRisingRun(t *testing.T) {
	for i, c := range cases(t, "rising-run", 60) {
		nums, _ := parseArray(t, c.Input)
		if len(nums) > 15 {
			continue
		}
		best := 0
		for mask := 1; mask < 1<<len(nums); mask++ {
			ok, last := true, -1<<62
			for j, v := range nums {
				if mask>>j&1 == 1 {
					ok = ok && v > last
					last = v
				}
			}
			if ok {
				best = max(best, bits.OnesCount(uint(mask)))
			}
		}
		if strconv.Itoa(best) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, best, c.Expected)
		}
	}
}

func TestSplitEqual(t *testing.T) {
	seen := map[string]bool{}
	for i, c := range cases(t, "split-equal", 80) {
		nums, _ := parseArray(t, c.Input)
		seen[c.Expected] = true
		if len(nums) > 16 {
			continue
		}
		total, ok := 0, false
		for _, v := range nums {
			total += v
		}
		for mask := 0; mask < 1<<len(nums); mask++ {
			sum := 0
			for j, v := range nums {
				sum += v * (mask >> j & 1)
			}
			ok = ok || 2*sum == total
		}
		if boolText(ok) != c.Expected {
			t.Fatalf("case %d: want %v, generator says %s", i, ok, c.Expected)
		}
	}
	if !seen["true"] || !seen["false"] {
		t.Fatal("both answers must appear")
	}
}

func TestKnapsack(t *testing.T) {
	for i, c := range cases(t, "knapsack", 60) {
		weights, rest := parseArray(t, c.Input)
		values, rest := parseArray(t, strings.Join(rest, "\n"))
		capacity, _ := strconv.Atoi(rest[0])
		if len(weights) > 15 {
			continue
		}
		best := 0
		for mask := 0; mask < 1<<len(weights); mask++ {
			w, v := 0, 0
			for j := range weights {
				w += weights[j] * (mask >> j & 1)
				v += values[j] * (mask >> j & 1)
			}
			if w <= capacity {
				best = max(best, v)
			}
		}
		if strconv.Itoa(best) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, best, c.Expected)
		}
	}
}

func isSubsequence(s, of string) bool {
	k := 0
	for i := 0; i < len(of) && k < len(s); i++ {
		if of[i] == s[k] {
			k++
		}
	}
	return k == len(s)
}

func TestCommonSubsequence(t *testing.T) {
	for i, c := range cases(t, "common-subsequence", 60) {
		w := strings.Split(c.Input, "\n")
		a, b := w[0], w[1]
		if len(a) > 14 {
			continue
		}
		best := 0
		for mask := 0; mask < 1<<len(a); mask++ {
			var s []byte
			for j := range a {
				if mask>>j&1 == 1 {
					s = append(s, a[j])
				}
			}
			if len(s) > best && isSubsequence(string(s), b) {
				best = len(s)
			}
		}
		if strconv.Itoa(best) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, best, c.Expected)
		}
	}
}

func TestWordSplit(t *testing.T) {
	seen := map[string]bool{}
	for i, c := range cases(t, "word-split", 80) {
		lines := strings.Split(c.Input, "\n")
		text, words := lines[0], strings.Fields(lines[2])
		if n, _ := strconv.Atoi(lines[1]); n != len(words) {
			t.Fatalf("case %d: says %d words, has %d", i, n, len(words))
		}
		// Recursion from the end, remembering texts already tried.
		memo := map[int]bool{}
		var fits func(end int) bool
		fits = func(end int) bool {
			if end == 0 {
				return true
			}
			if v, ok := memo[end]; ok {
				return v
			}
			found := false
			for _, w := range words {
				found = found || strings.HasSuffix(text[:end], w) && fits(end-len(w))
			}
			memo[end] = found
			return found
		}
		if boolText(fits(len(text))) != c.Expected {
			t.Fatalf("case %d: want %v, generator says %s", i, fits(len(text)), c.Expected)
		}
		seen[c.Expected] = true
	}
	if !seen["true"] || !seen["false"] {
		t.Fatal("both answers must appear")
	}
}

func TestEditDistance(t *testing.T) {
	for i, c := range cases(t, "edit-distance", 60) {
		w := strings.Split(c.Input, "\n")
		a, b := w[0], w[1]
		if len(a)*len(b) > 40_000 {
			continue
		}
		// Recursion from the front over both suffixes.
		memo := map[[2]int]int{}
		var dist func(i, j int) int
		dist = func(i, j int) int {
			if i == len(a) || j == len(b) {
				return len(a) - i + len(b) - j
			}
			if v, ok := memo[[2]int{i, j}]; ok {
				return v
			}
			v := 1 + min(dist(i+1, j), dist(i, j+1), dist(i+1, j+1))
			if a[i] == b[j] {
				v = dist(i+1, j+1)
			}
			memo[[2]int{i, j}] = v
			return v
		}
		if want := dist(0, 0); strconv.Itoa(want) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, want, c.Expected)
		}
	}
}
