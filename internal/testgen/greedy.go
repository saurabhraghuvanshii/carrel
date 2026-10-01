package testgen

import (
	"math/rand"
	"slices"
	"strconv"
	"strings"
)

// Generators for the "Greedy" group.

func init() {
	register("cookie-assignment", genCookieAssignment)
	register("pair-gaps", genPairGaps)
	register("pair-chain", genPairChain)
	register("job-deadlines", genJobDeadlines)
	register("can-reach-end", genCanReachEnd)
	register("fuel-circuit", genFuelCircuit)
	register("label-partitions", genLabelPartitions)
	register("chocolate-cuts", genChocolateCuts)
	register("candy-handout", genCandyHandout)
}

func sortedCopy(nums []int) []int {
	s := slices.Clone(nums)
	slices.Sort(s)
	return s
}

// cookie-assignment: each child needs a cookie at least their greed size;
// the answer is the most children who can get one cookie each.
func genCookieAssignment(rng *rand.Rand, size int) Case {
	top := []int{5, 1000, 1_000_000_000}[rng.Intn(3)]
	greed := randomInts(rng, 1+rng.Intn(max(1, size)), 1, top)
	cookies := randomInts(rng, 1+rng.Intn(max(1, size)), 1, top)
	g, c := sortedCopy((greed)), sortedCopy((cookies))
	fed := 0
	for _, s := range c {
		if fed < len(g) && s >= g[fed] {
			fed++
		}
	}
	return Case{Input: arrayInput(greed) + "\n" + arrayInput(cookies), Expected: strconv.Itoa(fed)}
}

// pair-gaps: two lists of the same length; pair every value of the first
// with one of the second. The answer is the smallest sum of the gaps.
func genPairGaps(rng *rand.Rand, size int) Case {
	n := max(1, size)
	top := []int{10, 1_000_000}[rng.Intn(2)]
	a, b := randomInts(rng, n, 0, top), randomInts(rng, n, 0, top)
	sa, sb := sortedCopy((a)), sortedCopy((b))
	total := 0
	for i := range sa {
		total += max(sa[i]-sb[i], sb[i]-sa[i])
	}
	return Case{Input: arrayInput(a) + "\n" + arrayInput(b), Expected: strconv.Itoa(total)}
}

// pairsInput prints "n" then n lines "a b".
func pairsInput(ps [][2]int) string {
	var b strings.Builder
	b.WriteString(strconv.Itoa(len(ps)))
	for _, p := range ps {
		b.WriteString("\n" + strconv.Itoa(p[0]) + " " + strconv.Itoa(p[1]))
	}
	return b.String()
}

// pair-chain: pairs (a, b) with a < b; (c, d) may follow (a, b) when c > b.
// The answer is the longest chain, using pairs in any order.
func genPairChain(rng *rand.Rand, size int) Case {
	n := max(1, size)
	room := max(10, n*[]int{1, 5, 20}[rng.Intn(3)])
	ps := make([][2]int, n)
	for i := range ps {
		a := rng.Intn(room) - room/2
		ps[i] = [2]int{a, a + 1 + rng.Intn(max(1, room/max(1, n)*[]int{1, 4}[rng.Intn(2)]))}
	}
	if rng.Intn(3) == 0 { // pairs that touch end to start do not chain
		for i := 1; i < n; i++ {
			if rng.Intn(2) == 0 {
				p := ps[rng.Intn(i)]
				ps[i] = [2]int{p[1], p[1] + 1 + rng.Intn(5)}
			}
		}
	}
	s := slices.Clone(ps)
	slices.SortFunc(s, func(x, y [2]int) int { return x[1] - y[1] })
	chain, end := 0, -1<<62
	for _, p := range s {
		if p[0] > end {
			chain++
			end = p[1]
		}
	}
	return Case{Input: pairsInput(ps), Expected: strconv.Itoa(chain)}
}

// job-deadlines: jobs that each take one unit of time, with a deadline (the
// job must finish by then) and a profit. The answer is the largest profit.
func genJobDeadlines(rng *rand.Rand, size int) Case {
	n := max(1, size)
	dTop := max(1, n/[]int{1, 3, 10}[rng.Intn(3)])
	jobs := make([][2]int, n)
	for i := range jobs {
		jobs[i] = [2]int{1 + rng.Intn(dTop), 1 + rng.Intn(1000)}
	}
	// Richest first, each into the latest free slot by its deadline. free[t]
	// points at the latest free slot at or before t (0 means none).
	s := slices.Clone(jobs)
	slices.SortFunc(s, func(x, y [2]int) int { return y[1] - x[1] })
	free := make([]int, dTop+1)
	for t := range free {
		free[t] = t
	}
	var find func(t int) int
	find = func(t int) int {
		for free[t] != t {
			free[t] = free[free[t]]
			t = free[t]
		}
		return t
	}
	total := 0
	for _, j := range s {
		if t := find(j[0]); t > 0 {
			total += j[1]
			free[t] = t - 1
		}
	}
	return Case{Input: pairsInput(jobs), Expected: strconv.Itoa(total)}
}

// can-reach-end: nums[i] is the longest jump from position i; the answer is
// whether the last position can be reached from the first.
func genCanReachEnd(rng *rand.Rand, size int) Case {
	n := max(1, size)
	nums := randomInts(rng, n, 0, []int{2, 3, 10}[rng.Intn(3)])
	if rng.Intn(2) == 0 { // a run of zeros that a long jump must clear
		at := rng.Intn(n)
		for i := at; i < min(n-1, at+1+rng.Intn(4)); i++ {
			nums[i] = 0
		}
	}
	if rng.Intn(3) == 0 { // plant a route of jumps; zeros stay elsewhere
		for p := 0; p < n-1; {
			step := 1 + rng.Intn(min(10, n-1-p))
			nums[p] = max(nums[p], step)
			p += step
		}
	}
	reach := 0
	for i := 0; i < n && i <= reach; i++ {
		reach = max(reach, i+nums[i])
	}
	return Case{Input: arrayInput(nums), Expected: boolText(reach >= n-1)}
}

// fuel-circuit: stations in a circle; station i gives gas[i] and driving on
// to the next costs cost[i]. The answer is the smallest start index from
// which the tank never goes below zero all the way round, or -1.
func genFuelCircuit(rng *rand.Rand, size int) Case {
	n := max(1, size)
	gas, cost := randomInts(rng, n, 0, 100), randomInts(rng, n, 0, 100)
	if rng.Intn(2) == 0 { // move gas until there is enough in total
		diff := 0
		for i := range gas {
			diff += gas[i] - cost[i]
		}
		for i := 0; diff < 0; i = (i + 1) % n {
			d := min(-diff, cost[i])
			cost[i] -= d
			diff += d
		}
	}
	start, tank, total := 0, 0, 0
	for i := range gas {
		tank += gas[i] - cost[i]
		total += gas[i] - cost[i]
		if tank < 0 {
			start, tank = i+1, 0
		}
	}
	if total < 0 {
		start = -1
	}
	return Case{Input: arrayInput(gas) + "\n" + arrayInput(cost), Expected: strconv.Itoa(start)}
}

// label-partitions: cut a word into as many pieces as possible so that every
// letter appears in one piece only. The answer lists the piece lengths.
func genLabelPartitions(rng *rand.Rand, size int) Case {
	n := max(1, size)
	var word string
	if rng.Intn(2) == 0 {
		word = randomWord(rng, n, "abcdefghijklmnopqrstuvwxyz"[:1+rng.Intn(26)])
	} else { // blocks with their own letters, lightly mixed
		var b strings.Builder
		letters := shuffled(rng, "abcdefghijklmnopqrstuvwxyz")
		for b.Len() < n {
			k := 1 + rng.Intn(3)
			from := rng.Intn(26 - k + 1)
			b.WriteString(randomWord(rng, 1+rng.Intn(6), letters[from:from+k]))
		}
		word = b.String()[:n]
	}
	last := map[byte]int{}
	for i := 0; i < len(word); i++ {
		last[word[i]] = i
	}
	var sizes []int
	start, end := 0, 0
	for i := 0; i < len(word); i++ {
		end = max(end, last[word[i]])
		if i == end {
			sizes = append(sizes, end-start+1)
			start = i + 1
		}
	}
	return Case{Input: word, Expected: listOutput(sizes)}
}

// chocolate-cuts: a bar of rows by cols squares. Cutting along horizontal
// line i costs h[i] for each piece it passes through, and the same for
// vertical lines. The answer is the cheapest total to cut it into squares.
func genChocolateCuts(rng *rand.Rand, size int) Case {
	rows, cols := 2+rng.Intn(min(999, max(1, size))), 2+rng.Intn(min(999, max(1, size)))
	top := []int{3, 10_000}[rng.Intn(2)]
	h, v := randomInts(rng, rows-1, 1, top), randomInts(rng, cols-1, 1, top)
	sh, sv := sortedCopy((h)), sortedCopy((v))
	slices.Reverse(sh)
	slices.Reverse(sv)
	// Dearest line first; each cut crosses one more than the cuts made
	// the other way.
	total, i, j := 0, 0, 0
	for i < len(sh) || j < len(sv) {
		if j == len(sv) || (i < len(sh) && sh[i] >= sv[j]) {
			total += sh[i] * (j + 1)
			i++
		} else {
			total += sv[j] * (i + 1)
			j++
		}
	}
	return Case{
		Input:    strconv.Itoa(rows) + " " + strconv.Itoa(cols) + "\n" + joinInts(h, " ") + "\n" + joinInts(v, " "),
		Expected: strconv.Itoa(total),
	}
}

// candy-handout: children in a row with ratings; each gets at least one
// sweet, and a child rated higher than a neighbour gets more than that
// neighbour. The answer is the fewest sweets in total.
func genCandyHandout(rng *rand.Rand, size int) Case {
	n := max(1, size)
	ratings := randomInts(rng, n, 0, []int{2, 5, 20_000}[rng.Intn(3)])
	if rng.Intn(3) == 0 { // long rising and falling slopes
		v, dir := 10_000, 1
		for i := range ratings {
			if rng.Intn(8) == 0 {
				dir = -dir
			}
			v = min(20_000, max(0, v+dir*rng.Intn(3)))
			ratings[i] = v
		}
	}
	sweets := make([]int, n)
	for i := range sweets {
		sweets[i] = 1
		if i > 0 && ratings[i] > ratings[i-1] {
			sweets[i] = sweets[i-1] + 1
		}
	}
	total := sweets[n-1]
	for i := n - 2; i >= 0; i-- {
		if ratings[i] > ratings[i+1] {
			sweets[i] = max(sweets[i], sweets[i+1]+1)
		}
		total += sweets[i]
	}
	return Case{Input: arrayInput(ratings), Expected: strconv.Itoa(total)}
}
