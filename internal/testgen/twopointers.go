package testgen

import (
	"math/rand"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

// Generators for the "Two pointers" group.

func init() {
	register("sorted-pair", genSortedPair)
	register("palindrome", genPalindrome)
	register("remove-repeats", genRemoveRepeats)
	register("merge-sorted", genMergeSorted)
	register("three-way", genThreeWay)
	register("triple-sum", genTripleSum)
	register("container", genContainer)
	register("rain-trap", genRainTrap)
}

// sorted-pair: a sorted list and a target; exactly one pair of positions adds
// up to it. Values stay within ±5×10⁸ so the target fits in 32 bits.
func genSortedPair(rng *rand.Rand, size int) Case {
	n := max(2, size)
	span := min(1_000_000_000, max(100, n*n*8))
	for {
		nums := make([]int, n)
		for i := range nums {
			nums[i] = rng.Intn(span+1) - span/2
		}
		slices.Sort(nums)
		i := rng.Intn(n - 1)
		j := i + 1 + rng.Intn(n-1-i)
		target := nums[i] + nums[j]
		if countPairs(nums, target) == 1 {
			return Case{Input: arrayInput(nums) + "\n" + strconv.Itoa(target), Expected: strconv.Itoa(i) + " " + strconv.Itoa(j)}
		}
	}
}

const palindromeNoise = " ,.;:!?'-_()"

// palindrome: one line of text. Letters and digits must read the same both
// ways, ignoring case; everything else is ignored. Half are palindromes.
func genPalindrome(rng *rand.Rand, size int) Case {
	alnum := "abcdefghijklmnopqrstuvwxyz0123456789"[:2+rng.Intn(35)]
	half := randomWord(rng, rng.Intn(max(1, size/2)+1), alnum)
	core := []byte(half)
	if rng.Intn(2) == 0 {
		core = append(core, alnum[rng.Intn(len(alnum))])
	}
	for i := len(half) - 1; i >= 0; i-- {
		core = append(core, half[i])
	}
	if len(core) > 0 && rng.Intn(2) == 0 {
		i := rng.Intn(len(core))
		core[i] = alnum[rng.Intn(len(alnum))]
	}
	var b strings.Builder
	for _, ch := range core {
		for rng.Intn(4) == 0 {
			b.WriteByte(palindromeNoise[rng.Intn(len(palindromeNoise))])
		}
		if rng.Intn(2) == 0 {
			ch = byte(unicode.ToUpper(rune(ch)))
		}
		b.WriteByte(ch)
	}
	text := strings.TrimSpace(b.String())
	if text == "" {
		text = "?!"
	}
	return Case{Input: text, Expected: boolText(isPalindrome(text))}
}

func isPalindrome(text string) bool {
	var keep []rune
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			keep = append(keep, unicode.ToLower(r))
		}
	}
	for i, j := 0, len(keep)-1; i < j; i, j = i+1, j-1 {
		if keep[i] != keep[j] {
			return false
		}
	}
	return true
}

// remove-repeats: a sorted list with many repeats; the answer is the list
// with each value once.
func genRemoveRepeats(rng *rand.Rand, size int) Case {
	n := size
	if rng.Intn(12) == 0 {
		n = 0
	}
	nums := make([]int, n)
	v := rng.Intn(20_001) - 10_000
	for i := range nums {
		if i > 0 && rng.Intn(2) == 0 {
			v += 1 + rng.Intn(5)
		}
		nums[i] = v
	}
	return Case{Input: arrayInput(nums), Expected: listOutput(slices.Compact(slices.Clone(nums)))}
}

func sortedRandom(rng *rand.Rand, n, span int) []int {
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rng.Intn(2*span+1) - span
	}
	slices.Sort(nums)
	return nums
}

// merge-sorted: two sorted lists, either may be empty; the answer is them
// merged in order.
func genMergeSorted(rng *rand.Rand, size int) Case {
	total := max(1, size)
	m := rng.Intn(total + 1)
	span := []int{5, 1000, 1_000_000_000}[rng.Intn(3)]
	a, b := sortedRandom(rng, m, span), sortedRandom(rng, total-m, span)
	merged := append(slices.Clone(a), b...)
	slices.Sort(merged)
	return Case{Input: arrayInput(a) + "\n" + arrayInput(b), Expected: listOutput(merged)}
}

// three-way: values 0, 1 and 2 in mixed proportions; the answer is them sorted.
func genThreeWay(rng *rand.Rand, size int) Case {
	n := max(1, size)
	weights := []int{rng.Intn(5), rng.Intn(5), rng.Intn(5)}
	weights[rng.Intn(3)]++
	sum := weights[0] + weights[1] + weights[2]
	nums := make([]int, n)
	for i := range nums {
		r := rng.Intn(sum)
		switch {
		case r < weights[0]:
			nums[i] = 0
		case r < weights[0]+weights[1]:
			nums[i] = 1
		default:
			nums[i] = 2
		}
	}
	sorted := slices.Clone(nums)
	slices.Sort(sorted)
	return Case{Input: arrayInput(nums), Expected: listOutput(sorted)}
}

// triple-sum: up to 500 values and a target; the answer is every different
// triple of values (from different positions) that adds up to the target,
// each triple ascending and the list sorted.
func genTripleSum(rng *rand.Rand, size int) Case {
	n := max(3, size)
	if n > 500 {
		n = 300 + n%201
	}
	span := []int{10, 1000, 100_000}[rng.Intn(3)]
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rng.Intn(2*span+1) - span
	}
	i, j, k := rng.Intn(n), rng.Intn(n), rng.Intn(n)
	target := nums[i] + nums[j] + nums[k]
	if i == j || j == k || i == k || rng.Intn(5) == 0 {
		target = rng.Intn(2*span+1) - span
	}
	return Case{Input: arrayInput(nums) + "\n" + strconv.Itoa(target), Expected: triplesOutput(tripleSum(nums, target))}
}

func tripleSum(nums []int, target int) [][]int {
	s := slices.Clone(nums)
	slices.Sort(s)
	var out [][]int
	for i := 0; i < len(s)-2; i++ {
		if i > 0 && s[i] == s[i-1] {
			continue
		}
		lo, hi := i+1, len(s)-1
		for lo < hi {
			sum := s[i] + s[lo] + s[hi]
			switch {
			case sum < target:
				lo++
			case sum > target:
				hi--
			default:
				out = append(out, []int{s[i], s[lo], s[hi]})
				for lo < hi && s[lo] == s[lo+1] {
					lo++
				}
				lo++
				hi--
			}
		}
	}
	return out
}

// triplesOutput sorts each list, then the lists, and prints [[a, b, c], ...].
func triplesOutput(lists [][]int) string {
	for _, l := range lists {
		slices.Sort(l)
	}
	slices.SortFunc(lists, func(a, b []int) int { return slices.Compare(a, b) })
	return gridOutput(lists)
}

func randomHeights(rng *rand.Rand, n int) []int {
	h := make([]int, n)
	top := []int{10, 1000, 100_000}[rng.Intn(3)]
	switch rng.Intn(3) {
	case 0: // noise
		for i := range h {
			h[i] = rng.Intn(top + 1)
		}
	case 1: // a valley: high at the ends
		for i := range h {
			d := min(i, n-1-i)
			h[i] = max(0, top-d*top/max(1, n/2)-rng.Intn(top/10+1))
		}
	default: // a mountain range with dips
		v := rng.Intn(top + 1)
		for i := range h {
			v = max(0, min(top, v+rng.Intn(top/5+1)-top/10))
			h[i] = v
		}
	}
	return h
}

// container: wall heights; the answer is the most water two walls can hold.
func genContainer(rng *rand.Rand, size int) Case {
	h := randomHeights(rng, max(2, size))
	return Case{Input: arrayInput(h), Expected: strconv.Itoa(container(h))}
}

func container(h []int) int {
	best := 0
	for i, j := 0, len(h)-1; i < j; {
		best = max(best, min(h[i], h[j])*(j-i))
		if h[i] < h[j] {
			i++
		} else {
			j--
		}
	}
	return best
}

// rain-trap: column heights; the answer is how much rain stays between them.
func genRainTrap(rng *rand.Rand, size int) Case {
	h := randomHeights(rng, max(1, size))
	return Case{Input: arrayInput(h), Expected: strconv.Itoa(rainTrap(h))}
}

func rainTrap(h []int) int {
	n := len(h)
	left, right := make([]int, n), make([]int, n)
	for i := range h {
		left[i] = h[i]
		if i > 0 {
			left[i] = max(left[i], left[i-1])
		}
	}
	for i := n - 1; i >= 0; i-- {
		right[i] = h[i]
		if i < n-1 {
			right[i] = max(right[i], right[i+1])
		}
	}
	total := 0
	for i := range h {
		total += min(left[i], right[i]) - h[i]
	}
	return total
}
