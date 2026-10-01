package testgen

import (
	"math/rand"
	"strconv"
)

// Generators for the "Sliding window" group.

func init() {
	register("window-sum", genWindowSum)
	register("unique-stretch", genUniqueStretch)
	register("two-kinds", genTwoKinds)
	register("k-changes", genKChanges)
	register("sum-at-least", genSumAtLeast)
	register("contains-rearrangement", genContainsRearrangement)
	register("covering-window", genCoveringWindow)
	register("window-max", genWindowMax)
}

func randomInts(rng *rand.Rand, n, lo, hi int) []int {
	nums := make([]int, n)
	for i := range nums {
		nums[i] = lo + rng.Intn(hi-lo+1)
	}
	return nums
}

// window-sum: the largest sum of k neighbouring values.
func genWindowSum(rng *rand.Rand, size int) Case {
	n := max(1, size)
	nums := randomInts(rng, n, -10_000, 10_000)
	k := 1 + rng.Intn(n)
	return Case{Input: arrayInput(nums) + "\n" + strconv.Itoa(k), Expected: strconv.Itoa(windowSum(nums, k))}
}

func windowSum(nums []int, k int) int {
	sum := 0
	for _, v := range nums[:k] {
		sum += v
	}
	best := sum
	for i := k; i < len(nums); i++ {
		sum += nums[i] - nums[i-k]
		best = max(best, sum)
	}
	return best
}

const wordChars = "abcdefghijklmnopqrstuvwxyz0123456789"

// unique-stretch: a word of lowercase letters and digits; the answer is the
// longest stretch with no repeated character.
func genUniqueStretch(rng *rand.Rand, size int) Case {
	text := randomWord(rng, max(1, size), wordChars[:1+rng.Intn(len(wordChars))])
	return Case{Input: text, Expected: strconv.Itoa(uniqueStretch(text))}
}

func uniqueStretch(text string) int {
	last := map[byte]int{}
	best, start := 0, 0
	for i := 0; i < len(text); i++ {
		if j, ok := last[text[i]]; ok && j >= start {
			start = j + 1
		}
		last[text[i]] = i
		best = max(best, i-start+1)
	}
	return best
}

// two-kinds: kinds of plants in a row, often in runs; the answer is the
// longest stretch with at most two kinds.
func genTwoKinds(rng *rand.Rand, size int) Case {
	n := max(1, size)
	kinds := distinct(rng, 1+rng.Intn(6), 0, 10_000)
	nums := make([]int, 0, n)
	for len(nums) < n {
		k := kinds[rng.Intn(len(kinds))]
		for r := 1 + rng.Intn(4); r > 0 && len(nums) < n; r-- {
			nums = append(nums, k)
		}
	}
	return Case{Input: arrayInput(nums), Expected: strconv.Itoa(twoKinds(nums))}
}

func twoKinds(nums []int) int {
	count := map[int]int{}
	best, start := 0, 0
	for i, v := range nums {
		count[v]++
		for len(count) > 2 {
			count[nums[start]]--
			if count[nums[start]] == 0 {
				delete(count, nums[start])
			}
			start++
		}
		best = max(best, i-start+1)
	}
	return best
}

// k-changes: an uppercase word and k; the answer is the longest run of one
// letter after changing at most k letters.
func genKChanges(rng *rand.Rand, size int) Case {
	n := max(1, size)
	text := randomWord(rng, n, "ABCDEFGHIJKLMNOPQRSTUVWXYZ"[:1+rng.Intn(min(26, 1+rng.Intn(8)))])
	k := rng.Intn(n + 1)
	if rng.Intn(2) == 0 {
		k = rng.Intn(min(n, 5) + 1)
	}
	return Case{Input: text + "\n" + strconv.Itoa(k), Expected: strconv.Itoa(kChanges(text, k))}
}

func kChanges(text string, k int) int {
	var count [26]int
	best, most, start := 0, 0, 0
	for i := 0; i < len(text); i++ {
		count[text[i]-'A']++
		most = max(most, count[text[i]-'A'])
		for i-start+1-most > k {
			count[text[start]-'A']--
			start++
		}
		best = max(best, i-start+1)
	}
	return best
}

// sum-at-least: positive values and a target; the answer is the shortest
// stretch whose sum reaches the target, or 0.
func genSumAtLeast(rng *rand.Rand, size int) Case {
	n := max(1, size)
	top := []int{10, 10_000}[rng.Intn(2)]
	nums := randomInts(rng, n, 1, top)
	total := 0
	for _, v := range nums {
		total += v
	}
	target := 1 + rng.Intn(total+total/5+1)
	return Case{Input: arrayInput(nums) + "\n" + strconv.Itoa(target), Expected: strconv.Itoa(sumAtLeast(nums, target))}
}

func sumAtLeast(nums []int, target int) int {
	best, sum, start := 0, 0, 0
	for i, v := range nums {
		sum += v
		for sum >= target {
			if best == 0 || i-start+1 < best {
				best = i - start + 1
			}
			sum -= nums[start]
			start++
		}
	}
	return best
}

// contains-rearrangement: a pattern and a text on two lines. Half the cases
// hide a rearranged copy of the pattern in the text.
func genContainsRearrangement(rng *rand.Rand, size int) Case {
	n := max(1, size)
	letters := "abcdefghijklmnopqrstuvwxyz"[:2+rng.Intn(10)]
	text := []byte(randomWord(rng, n, letters))
	m := 1 + rng.Intn(min(n+2, 30))
	pattern := randomWord(rng, m, letters)
	if m <= n && rng.Intn(2) == 0 {
		at := rng.Intn(n - m + 1)
		copy(text[at:], shuffled(rng, pattern))
	}
	return Case{Input: pattern + "\n" + string(text), Expected: boolText(containsRearrangement(pattern, string(text)))}
}

func containsRearrangement(pattern, text string) bool {
	if len(pattern) > len(text) {
		return false
	}
	var need, have [26]int
	for i := 0; i < len(pattern); i++ {
		need[pattern[i]-'a']++
		have[text[i]-'a']++
	}
	for i := len(pattern); ; i++ {
		if need == have {
			return true
		}
		if i == len(text) {
			return false
		}
		have[text[i]-'a']++
		have[text[i-len(pattern)]-'a']--
	}
}

// covering-window: a text and a pattern of letters (case matters). The answer
// is the shortest stretch of text holding every pattern letter as often as
// the pattern does, the first one if several tie, in quotes; "" if none.
func genCoveringWindow(rng *rand.Rand, size int) Case {
	n := max(1, size)
	letters := "abcdeABCDExyzXYZ"[:2+rng.Intn(15)]
	text := randomWord(rng, n, letters)
	m := 1 + rng.Intn(min(n, 8))
	var pattern string
	if rng.Intn(4) == 0 {
		pattern = randomWord(rng, m, letters)
	} else {
		at := rng.Intn(n - m + 1)
		pattern = shuffled(rng, text[at:at+m])
		if rng.Intn(3) == 0 {
			pattern = randomWord(rng, 1, text) + pattern[1:]
		}
	}
	return Case{Input: text + "\n" + pattern, Expected: `"` + coveringWindow(text, pattern) + `"`}
}

func coveringWindow(text, pattern string) string {
	need := map[byte]int{}
	for i := 0; i < len(pattern); i++ {
		need[pattern[i]]++
	}
	missing := len(pattern)
	bestStart, bestLen := 0, -1
	start := 0
	for i := 0; i < len(text); i++ {
		if need[text[i]] > 0 {
			missing--
		}
		need[text[i]]--
		for missing == 0 {
			if bestLen < 0 || i-start+1 < bestLen {
				bestStart, bestLen = start, i-start+1
			}
			need[text[start]]++
			if need[text[start]] > 0 {
				missing++
			}
			start++
		}
	}
	if bestLen < 0 {
		return ""
	}
	return text[bestStart : bestStart+bestLen]
}

// window-max: values that rise and fall in trends; the answer is the largest
// value in each window of k.
func genWindowMax(rng *rand.Rand, size int) Case {
	n := max(1, size)
	nums := make([]int, n)
	v := rng.Intn(20_001) - 10_000
	for i := range nums {
		v = max(-10_000, min(10_000, v+rng.Intn(2001)-1000))
		nums[i] = v
	}
	k := 1 + rng.Intn(n)
	if rng.Intn(2) == 0 {
		k = 1 + rng.Intn(min(n, 10))
	}
	return Case{Input: arrayInput(nums) + "\n" + strconv.Itoa(k), Expected: listOutput(windowMax(nums, k))}
}

func windowMax(nums []int, k int) []int {
	var queue, out []int // queue holds positions with falling values
	for i, v := range nums {
		for len(queue) > 0 && nums[queue[len(queue)-1]] <= v {
			queue = queue[:len(queue)-1]
		}
		queue = append(queue, i)
		if queue[0] <= i-k {
			queue = queue[1:]
		}
		if i >= k-1 {
			out = append(out, nums[queue[0]])
		}
	}
	return out
}
