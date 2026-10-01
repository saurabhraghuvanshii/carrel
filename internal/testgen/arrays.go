package testgen

import (
	"math/rand"
	"slices"
	"strconv"
	"strings"
)

// Generators for the "Arrays and hashing" group. Arrays use the shared format:
// "n" then the n values on one line. Lists in answers print as "[a, b]".

func init() {
	register("spot-duplicate", genSpotDuplicate)
	register("same-letters", genSameLetters)
	register("group-letters", genGroupLetters)
	register("most-frequent", genMostFrequent)
	register("product-others", genProductOthers)
	register("longest-run", genLongestRun)
	register("subarray-sum", genSubarraySum)
	register("rotate-grid", genRotateGrid)
	register("missing-positive", genMissingPositive)
}

func arrayInput(nums []int) string {
	return strconv.Itoa(len(nums)) + "\n" + joinInts(nums, " ")
}

func listOutput(nums []int) string { return "[" + joinInts(nums, ", ") + "]" }

func boolText(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// distinct returns n different values from [lo, hi].
func distinct(rng *rand.Rand, n, lo, hi int) []int {
	seen := map[int]bool{}
	out := make([]int, 0, n)
	for len(out) < n {
		v := lo + rng.Intn(hi-lo+1)
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}

// spot-duplicate: does any value appear twice? Half the cases plant one.
func genSpotDuplicate(rng *rand.Rand, size int) Case {
	n := max(size, 1)
	if rng.Intn(10) == 0 {
		n = 1
	}
	nums := distinct(rng, n, -1_000_000_000, 1_000_000_000)
	if n > 1 && rng.Intn(2) == 0 {
		i, j := rng.Intn(n), rng.Intn(n-1)
		if j >= i {
			j++
		}
		nums[j] = nums[i]
	}
	return Case{Input: arrayInput(nums), Expected: boolText(hasDuplicate(nums))}
}

func hasDuplicate(nums []int) bool {
	seen := map[int]bool{}
	for _, v := range nums {
		if seen[v] {
			return true
		}
		seen[v] = true
	}
	return false
}

func randomWord(rng *rand.Rand, n int, letters string) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}

func shuffled(rng *rand.Rand, s string) string {
	b := []byte(s)
	rng.Shuffle(len(b), func(i, j int) { b[i], b[j] = b[j], b[i] })
	return string(b)
}

func sortedLetters(s string) string {
	b := []byte(s)
	slices.Sort(b)
	return string(b)
}

// same-letters: two lowercase words on two lines. Half are rearrangements;
// the rest differ by one letter or by length.
func genSameLetters(rng *rand.Rand, size int) Case {
	n := max(1, min(size, 50_000))
	letters := "abcdefghijklmnopqrstuvwxyz"[:3+rng.Intn(24)]
	a := randomWord(rng, n, letters)
	var b string
	switch rng.Intn(4) {
	case 0, 1:
		b = shuffled(rng, a)
	case 2:
		c := []byte(shuffled(rng, a))
		c[rng.Intn(n)] = letters[rng.Intn(len(letters))]
		b = string(c)
	default:
		b = shuffled(rng, a+randomWord(rng, 1, letters))
	}
	return Case{Input: a + "\n" + b, Expected: boolText(sortedLetters(a) == sortedLetters(b))}
}

// group-letters: n short words; words that are rearrangements of each other
// form a group. The answer lists each group's words sorted, and the groups
// sorted, so it is unique.
func genGroupLetters(rng *rand.Rand, size int) Case {
	n := max(1, min(size, 5000))
	keys := max(1, n/(2+rng.Intn(3)))
	bases := make([]string, keys)
	for i := range bases {
		bases[i] = randomWord(rng, 1+rng.Intn(6), "abcdef")
	}
	words := make([]string, n)
	for i := range words {
		words[i] = shuffled(rng, bases[rng.Intn(keys)])
	}
	return Case{Input: strconv.Itoa(n) + "\n" + strings.Join(words, " "), Expected: groupsOutput(groupLetters(words))}
}

func groupLetters(words []string) [][]string {
	byKey := map[string][]string{}
	for _, w := range words {
		k := sortedLetters(w)
		byKey[k] = append(byKey[k], w)
	}
	var groups [][]string
	for _, g := range byKey {
		groups = append(groups, g)
	}
	return groups
}

// groupsOutput sorts the words in each group, then the groups, and prints
// them as [[a, b], [c]].
func groupsOutput(groups [][]string) string {
	for _, g := range groups {
		slices.Sort(g)
	}
	slices.SortFunc(groups, func(a, b []string) int { return slices.Compare(a, b) })
	parts := make([]string, len(groups))
	for i, g := range groups {
		parts[i] = "[" + strings.Join(g, ", ") + "]"
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

// most-frequent: "n", the values, then k. The k most frequent values are
// unique by construction; the answer is them in ascending order.
func genMostFrequent(rng *rand.Rand, size int) Case {
	target := max(1, size)
	d := max(1, min(target, 1+rng.Intn(max(1, target/2+1))))
	values := distinct(rng, d, -10_000, 10_000)
	counts := make([]int, d)
	total := 0
	for i := range counts {
		counts[i] = 1 + rng.Intn(max(1, 2*target/d))
		total += counts[i]
	}
	order := make([]int, d)
	for i := range order {
		order[i] = i
	}
	slices.SortFunc(order, func(a, b int) int { return counts[b] - counts[a] })
	k := 1 + rng.Intn(d)
	if k < d && counts[order[k-1]] == counts[order[k]] {
		// Lift the top k above the rest so the answer is unique.
		cut := counts[order[k]]
		for _, i := range order[:k] {
			if counts[i] <= cut {
				total += cut + 1 - counts[i]
				counts[i] = cut + 1
			}
		}
	}
	var nums []int
	for i, v := range values {
		for j := 0; j < counts[i]; j++ {
			nums = append(nums, v)
		}
	}
	rng.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })
	return Case{Input: arrayInput(nums) + "\n" + strconv.Itoa(k), Expected: listOutput(mostFrequent(nums, k))}
}

func mostFrequent(nums []int, k int) []int {
	count := map[int]int{}
	for _, v := range nums {
		count[v]++
	}
	var vals []int
	for v := range count {
		vals = append(vals, v)
	}
	slices.SortFunc(vals, func(a, b int) int {
		if count[a] != count[b] {
			return count[b] - count[a]
		}
		return a - b
	})
	top := slices.Clone(vals[:k])
	slices.Sort(top)
	return top
}

// product-others: every product of all but one value fits in 32 bits. Most
// values are 1 or -1, a few are larger, and sometimes there are zeros.
func genProductOthers(rng *rand.Rand, size int) Case {
	n := max(2, size)
	nums := make([]int, n)
	product := 1
	for i := range nums {
		v := 1
		if rng.Intn(2) == 0 {
			v = -1
		}
		if rng.Intn(4) == 0 {
			big := 2 + rng.Intn(29)
			if product*big <= 1_000_000_000 {
				product *= big
				v *= big
			}
		}
		nums[i] = v
	}
	for z := rng.Intn(4) - 1; z > 0; z-- {
		nums[rng.Intn(n)] = 0
	}
	return Case{Input: arrayInput(nums), Expected: listOutput(productOthers(nums))}
}

func productOthers(nums []int) []int {
	out := make([]int, len(nums))
	left := 1
	for i, v := range nums {
		out[i] = left
		left *= v
	}
	right := 1
	for i := len(nums) - 1; i >= 0; i-- {
		out[i] *= right
		right *= nums[i]
	}
	return out
}

// longest-run: the length of the longest set of values that are consecutive
// whole numbers. A few runs are planted among random noise.
func genLongestRun(rng *rand.Rand, size int) Case {
	n := size
	if rng.Intn(12) == 0 {
		n = 0
	}
	var nums []int
	for len(nums) < n {
		if rng.Intn(3) == 0 {
			start := rng.Intn(2_000_000_000-20_000) - 1_000_000_000
			for i := 0; i < 1+rng.Intn(max(1, n/3)) && len(nums) < n; i++ {
				nums = append(nums, start+i)
			}
		} else {
			nums = append(nums, rng.Intn(2_000_000_001)-1_000_000_000)
		}
	}
	for i := 0; i < len(nums)/10; i++ {
		nums[rng.Intn(len(nums))] = nums[rng.Intn(len(nums))] // repeats
	}
	rng.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })
	return Case{Input: arrayInput(nums), Expected: strconv.Itoa(longestRun(nums))}
}

func longestRun(nums []int) int {
	s := slices.Clone(nums)
	slices.Sort(s)
	s = slices.Compact(s)
	best, cur := 0, 0
	for i, v := range s {
		if i > 0 && v == s[i-1]+1 {
			cur++
		} else {
			cur = 1
		}
		best = max(best, cur)
	}
	return best
}

// subarray-sum: "n", the values, then the target. Small values and many
// zeros make plenty of matching stretches.
func genSubarraySum(rng *rand.Rand, size int) Case {
	n := max(1, size)
	span := []int{3, 10, 1000}[rng.Intn(3)]
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rng.Intn(2*span+1) - span
	}
	i := rng.Intn(n)
	j := i + rng.Intn(n-i)
	target := 0
	for _, v := range nums[i : j+1] {
		target += v
	}
	if rng.Intn(5) == 0 {
		target = rng.Intn(2001) - 1000
	}
	return Case{Input: arrayInput(nums) + "\n" + strconv.Itoa(target), Expected: strconv.Itoa(countSubarrays(nums, target))}
}

func countSubarrays(nums []int, target int) int {
	seen := map[int]int{0: 1}
	sum, count := 0, 0
	for _, v := range nums {
		sum += v
		count += seen[sum-target]
		seen[sum]++
	}
	return count
}

// rotate-grid: an n by n grid in the grid format; the answer is the grid
// turned a quarter turn clockwise, printed as [[row], [row]].
func genRotateGrid(rng *rand.Rand, size int) Case {
	n := max(1, size)
	if n > 150 {
		n = 100 + n%101
	}
	grid := make([][]int, n)
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(n) + " " + strconv.Itoa(n))
	for r := range grid {
		grid[r] = make([]int, n)
		for c := range grid[r] {
			grid[r][c] = rng.Intn(2001) - 1000
		}
		sb.WriteString("\n" + joinInts(grid[r], " "))
	}
	return Case{Input: sb.String(), Expected: gridOutput(rotateGrid(grid))}
}

func rotateGrid(grid [][]int) [][]int {
	n := len(grid)
	out := make([][]int, n)
	for r := range out {
		out[r] = make([]int, n)
		for c := range out[r] {
			out[r][c] = grid[n-1-c][r]
		}
	}
	return out
}

func gridOutput(grid [][]int) string {
	rows := make([]string, len(grid))
	for i, row := range grid {
		rows[i] = listOutput(row)
	}
	return "[" + strings.Join(rows, ", ") + "]"
}

// missing-positive: most of 1..m present, often with one gap, mixed with
// zeros, negatives, large values and repeats.
func genMissingPositive(rng *rand.Rand, size int) Case {
	n := max(1, size)
	m := rng.Intn(n + 1)
	var nums []int
	gap := 1 + rng.Intn(m+1)
	for v := 1; v <= m; v++ {
		if v != gap || rng.Intn(4) == 0 {
			nums = append(nums, v)
		}
	}
	for len(nums) < n {
		switch rng.Intn(4) {
		case 0:
			nums = append(nums, -rng.Intn(1_000_000_001))
		case 1:
			nums = append(nums, 1+rng.Intn(1_000_000_000))
		default:
			if len(nums) > 0 {
				nums = append(nums, nums[rng.Intn(len(nums))])
			} else {
				nums = append(nums, 0)
			}
		}
	}
	nums = nums[:n]
	rng.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })
	return Case{Input: arrayInput(nums), Expected: strconv.Itoa(smallestMissing(nums))}
}

func smallestMissing(nums []int) int {
	seen := map[int]bool{}
	for _, v := range nums {
		seen[v] = true
	}
	for v := 1; ; v++ {
		if !seen[v] {
			return v
		}
	}
}
