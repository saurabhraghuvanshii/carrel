package testgen

import (
	"math"
	"math/rand"
	"slices"
	"strconv"
	"strings"
)

// Generators for the "Dynamic programming" group.

func init() {
	register("stair-ways", genStairWays)
	register("friend-pairs", genFriendPairs)
	register("house-skipper", genHouseSkipper)
	register("fewest-coins", genFewestCoins)
	register("coin-ways", genCoinWays)
	register("grid-path-count", genGridPathCount)
	register("count-search-trees", genCountSearchTrees)
	register("rising-run", genRisingRun)
	register("split-equal", genSplitEqual)
	register("knapsack", genKnapsack)
	register("common-subsequence", genCommonSubsequence)
	register("word-split", genWordSplit)
	register("edit-distance", genEditDistance)
}

const mod = 1_000_000_007

// stair-ways: n stairs, climbing 1 or 2 at a time; the answer is the number
// of ways to reach the top.
func genStairWays(rng *rand.Rand, size int) Case {
	n := 1 + rng.Intn(min(45, max(2, size)))
	a, b := 1, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return Case{Input: strconv.Itoa(n), Expected: strconv.Itoa(b)}
}

// friend-pairs: n people each stay alone or pair with one other; the answer
// is the number of ways, modulo 10^9 + 7.
func genFriendPairs(rng *rand.Rand, size int) Case {
	n := 1 + rng.Intn(min(10_000, max(2, size*5)))
	a, b := 1, 1 // ways for 0 and 1 people
	for i := 2; i <= n; i++ {
		a, b = b, (b+(i-1)*a)%mod
	}
	return Case{Input: strconv.Itoa(n), Expected: strconv.Itoa(b)}
}

// house-skipper: values in a row; the answer is the largest sum with no two
// neighbours taken.
func genHouseSkipper(rng *rand.Rand, size int) Case {
	nums := randomInts(rng, max(1, size), 0, []int{9, 1000}[rng.Intn(2)])
	take, skip := 0, 0
	for _, v := range nums {
		take, skip = skip+v, max(take, skip)
	}
	return Case{Input: arrayInput(nums), Expected: strconv.Itoa(max(take, skip))}
}

// randomCoins makes different coin values. A quarter of the time every coin is even,
// so odd amounts cannot be paid.
func randomCoins(rng *rand.Rand, size int) []int {
	n := 1 + rng.Intn(min(12, max(1, size)))
	top := []int{10, 30, 200}[rng.Intn(3)]
	coins := distinct(rng, min(n, top), 1, top)
	if rng.Intn(4) == 0 {
		for i := range coins {
			coins[i] *= 2
		}
	}
	return coins
}

// fewest-coins: different coin values (each usable any number of times) and
// an amount; the answer is the fewest coins that pay it exactly, or -1.
func genFewestCoins(rng *rand.Rand, size int) Case {
	coins := randomCoins(rng, size)
	amount := rng.Intn(min(10_000, 1+size*20))
	best := make([]int, amount+1)
	for v := 1; v <= amount; v++ {
		best[v] = -1
		for _, c := range coins {
			if c <= v && best[v-c] >= 0 && (best[v] < 0 || best[v-c]+1 < best[v]) {
				best[v] = best[v-c] + 1
			}
		}
	}
	return Case{Input: arrayInput(coins) + "\n" + strconv.Itoa(amount), Expected: strconv.Itoa(best[amount])}
}

// coinWays counts the ways to pay each amount up to top, ignoring order.
// Counts above math.MaxInt32 are capped there.
func coinWays(coins []int, top int) []int {
	ways := make([]int, top+1)
	ways[0] = 1
	for _, c := range coins {
		for v := c; v <= top; v++ {
			ways[v] = min(math.MaxInt32, ways[v]+ways[v-c])
		}
	}
	return ways
}

// coin-ways: different coin values and an amount; the answer is the number
// of ways to pay it, ignoring order. Amounts are chosen so it fits in an int.
func genCoinWays(rng *rand.Rand, size int) Case {
	coins := randomCoins(rng, size)
	top := min(5000, 1+size*10)
	ways := coinWays(coins, top)
	amount := rng.Intn(top + 1)
	for ways[amount] >= math.MaxInt32 {
		amount /= 2
	}
	return Case{Input: arrayInput(coins) + "\n" + strconv.Itoa(amount), Expected: strconv.Itoa(ways[amount])}
}

// grid-path-count: a grid of 0 (open) and 1 (wall); the answer is the number
// of paths from the top-left to the bottom-right moving only right or down.
func genGridPathCount(rng *rand.Rand, size int) Case {
	r, c := 1+rng.Intn(min(16, max(1, size))), 1+rng.Intn(min(16, max(1, size)))
	walls := []float64{0, 0.1, 0.25}[rng.Intn(3)]
	g := make([][]int, r)
	for i := range g {
		g[i] = make([]int, c)
		for j := range g[i] {
			if rng.Float64() < walls {
				g[i][j] = 1
			}
		}
	}
	if rng.Intn(4) != 0 { // keep the ends open most of the time
		g[0][0], g[r-1][c-1] = 0, 0
	}
	ways := make([][]int, r)
	for i := range ways {
		ways[i] = make([]int, c)
		for j := range ways[i] {
			switch {
			case g[i][j] == 1:
			case i == 0 && j == 0:
				ways[i][j] = 1
			default:
				if i > 0 {
					ways[i][j] += ways[i-1][j]
				}
				if j > 0 {
					ways[i][j] += ways[i][j-1]
				}
			}
		}
	}
	return Case{Input: gridInput(g), Expected: strconv.Itoa(ways[r-1][c-1])}
}

func gridInput(g [][]int) string {
	rows := make([]string, len(g))
	for i, row := range g {
		rows[i] = joinInts(row, " ")
	}
	return strconv.Itoa(len(g)) + " " + strconv.Itoa(len(g[0])) + "\n" + strings.Join(rows, "\n")
}

// count-search-trees: n; the answer is the number of differently shaped
// search trees holding the keys 1 to n.
func genCountSearchTrees(rng *rand.Rand, size int) Case {
	n := 1 + rng.Intn(min(19, max(2, size)))
	trees := make([]int, n+1)
	trees[0] = 1
	for k := 1; k <= n; k++ {
		for root := 1; root <= k; root++ {
			trees[k] += trees[root-1] * trees[k-root]
		}
	}
	return Case{Input: strconv.Itoa(n), Expected: strconv.Itoa(trees[n])}
}

// rising-run: the answer is the length of the longest strictly increasing
// subsequence. Narrow value ranges give many repeats.
func genRisingRun(rng *rand.Rand, size int) Case {
	n := max(1, min(2500, size))
	top := []int{3, 50, 10_000}[rng.Intn(3)]
	nums := randomInts(rng, n, -top, top)
	if rng.Intn(4) == 0 { // a rising trend with noise
		for i := range nums {
			nums[i] = i + rng.Intn(top) - top/2
		}
	}
	var tails []int
	for _, v := range nums {
		i, _ := slices.BinarySearch(tails, v)
		if i == len(tails) {
			tails = append(tails, v)
		} else {
			tails[i] = v
		}
	}
	return Case{Input: arrayInput(nums), Expected: strconv.Itoa(len(tails))}
}

// split-equal: positive values; the answer is whether they split into two
// groups with the same sum. A third of the cases are built from two equal
// halves, and a third have an even total that may still not split.
func genSplitEqual(rng *rand.Rand, size int) Case {
	n := max(1, min(200, size))
	top := []int{5, 100}[rng.Intn(2)]
	nums := randomInts(rng, n, 1, top)
	switch mode := rng.Intn(3); {
	case mode == 0 && n > 1:
		k := 1 + rng.Intn(n-1)
		a, b := 0, 0
		for _, v := range nums[:k] {
			a += v
		}
		for _, v := range nums[k:] {
			b += v
		}
		// Grow the smaller side with values up to top until the sums match.
		for a != b {
			d := min(top, max(a, b)-min(a, b))
			nums = append(nums, d)
			if a < b {
				a += d
			} else {
				b += d
			}
		}
		rng.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })
		nums = nums[:min(len(nums), 200)]
	case mode == 1:
		total := 0
		for _, v := range nums {
			total += v
		}
		if total%2 == 1 && nums[0] > 1 {
			nums[0]--
		} else if total%2 == 1 {
			nums[0]++
		}
	}
	return Case{Input: arrayInput(nums), Expected: boolText(canSplit(nums))}
}

func canSplit(nums []int) bool {
	total := 0
	for _, v := range nums {
		total += v
	}
	if total%2 == 1 {
		return false
	}
	reach := make([]bool, total/2+1)
	reach[0] = true
	for _, v := range nums {
		for s := total / 2; s >= v; s-- {
			reach[s] = reach[s] || reach[s-v]
		}
	}
	return reach[total/2]
}

// knapsack: weights, values and a capacity; each item is taken once or not
// at all. The answer is the largest total value that fits.
func genKnapsack(rng *rand.Rand, size int) Case {
	n := max(1, min(100, size))
	wTop := []int{10, 1000}[rng.Intn(2)]
	weights := randomInts(rng, n, 1, wTop)
	values := randomInts(rng, n, 1, 1000)
	total := 0
	for _, w := range weights {
		total += w
	}
	capacity := rng.Intn(min(10_000, total+5) + 1)
	best := make([]int, capacity+1)
	for i, w := range weights {
		for c := capacity; c >= w; c-- {
			best[c] = max(best[c], best[c-w]+values[i])
		}
	}
	return Case{
		Input:    arrayInput(weights) + "\n" + arrayInput(values) + "\n" + strconv.Itoa(capacity),
		Expected: strconv.Itoa(best[capacity]),
	}
}

// common-subsequence: two words; the answer is the length of the longest
// sequence of letters found, in order, in both.
func genCommonSubsequence(rng *rand.Rand, size int) Case {
	letters := "abcdefghijklmnopqrstuvwxyz"[:[]int{2, 4, 26}[rng.Intn(3)]]
	a := randomWord(rng, 1+rng.Intn(min(1000, max(1, size))), letters)
	b := randomWord(rng, 1+rng.Intn(min(1000, max(1, size))), letters)
	return Case{Input: a + "\n" + b, Expected: strconv.Itoa(lcsLength(a, b))}
}

func lcsLength(a, b string) int {
	prev, cur := make([]int, len(b)+1), make([]int, len(b)+1)
	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			if a[i-1] == b[j-1] {
				cur[j] = prev[j-1] + 1
			} else {
				cur[j] = max(prev[j], cur[j-1])
			}
		}
		prev, cur = cur, prev
	}
	return prev[len(b)]
}

// word-split: a text and a list of words; the answer is whether the text cuts
// into pieces that are all on the list (words may repeat). Some cases repeat
// one letter with a stray letter at the end, which defeats plain recursion.
func genWordSplit(rng *rand.Rand, size int) Case {
	var text string
	var words []string
	switch k := rng.Intn(4); {
	case k == 0 && size > 15:
		for i := 1; i <= 1+rng.Intn(10); i++ {
			words = append(words, strings.Repeat("a", i))
		}
		text = strings.Repeat("a", min(299, size)) + []string{"b", ""}[rng.Intn(2)]
	default:
		letters := "abc"[:2+rng.Intn(2)]
		for i := 1 + rng.Intn(min(20, max(1, size))); i > 0; i-- {
			words = append(words, randomWord(rng, 1+rng.Intn(5), letters))
		}
		slices.Sort(words)
		words = slices.Compact(words)
		length := 1 + rng.Intn(min(300, max(1, size*3)))
		if k%2 == 0 { // glue words together, then maybe break one letter
			var b strings.Builder
			for b.Len() < length {
				w := words[rng.Intn(len(words))]
				if b.Len() > 0 && b.Len()+len(w) > 300 {
					break
				}
				b.WriteString(w)
			}
			text = b.String()
			if rng.Intn(2) == 0 {
				p := rng.Intn(len(text))
				text = text[:p] + letters[rng.Intn(len(letters)):][:1] + text[p+1:]
			}
		} else {
			text = randomWord(rng, length, letters)
		}
		rng.Shuffle(len(words), func(i, j int) { words[i], words[j] = words[j], words[i] })
	}
	return Case{
		Input:    text + "\n" + strconv.Itoa(len(words)) + "\n" + strings.Join(words, " "),
		Expected: boolText(canSplitWords(text, words)),
	}
}

func canSplitWords(text string, words []string) bool {
	ok := make([]bool, len(text)+1)
	ok[0] = true
	for i := 0; i < len(text); i++ {
		if !ok[i] {
			continue
		}
		for _, w := range words {
			if strings.HasPrefix(text[i:], w) {
				ok[i+len(w)] = true
			}
		}
	}
	return ok[len(text)]
}

// edit-distance: two words; the answer is the fewest single-letter inserts,
// deletes and replacements that turn the first into the second.
func genEditDistance(rng *rand.Rand, size int) Case {
	letters := "abcdefghijklmnopqrstuvwxyz"[:[]int{2, 4, 26}[rng.Intn(3)]]
	a := randomWord(rng, 1+rng.Intn(min(496, max(1, size))), letters)
	b := randomWord(rng, 1+rng.Intn(min(500, max(1, size))), letters)
	if rng.Intn(2) == 0 { // a few changes to the first word
		w := []byte(a)
		for k := 1 + rng.Intn(4); k > 0; k-- {
			p := rng.Intn(len(w) + 1)
			ch := letters[rng.Intn(len(letters))]
			switch {
			case rng.Intn(3) == 0 || len(w) < 2:
				w = slices.Insert(w, p, ch)
			case rng.Intn(2) == 0:
				w = slices.Delete(w, min(p, len(w)-1), min(p, len(w)-1)+1)
			default:
				w[min(p, len(w)-1)] = ch
			}
		}
		b = string(w)
	}
	return Case{Input: a + "\n" + b, Expected: strconv.Itoa(editDistance(a, b))}
}

func editDistance(a, b string) int {
	prev, cur := make([]int, len(b)+1), make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		cur[0] = i
		for j := 1; j <= len(b); j++ {
			if a[i-1] == b[j-1] {
				cur[j] = prev[j-1]
			} else {
				cur[j] = 1 + min(prev[j-1], prev[j], cur[j-1])
			}
		}
		prev, cur = cur, prev
	}
	return prev[len(b)]
}
