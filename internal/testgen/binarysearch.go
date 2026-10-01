package testgen

import (
	"math"
	"math/rand"
	"slices"
	"strconv"
)

// Generators for the "Binary search" group.

func init() {
	register("find-sorted", genFindSorted)
	register("int-sqrt", genIntSqrt)
	register("first-last", genFirstLast)
	register("min-rotated", genMinRotated)
	register("search-rotated", genSearchRotated)
	register("mountain-top", genMountainTop)
	register("eating-speed", genEatingSpeed)
	register("median-two", genMedianTwo)
}

// sortedDistinct returns n different values in ascending order.
func sortedDistinct(rng *rand.Rand, n int) []int {
	nums := distinct(rng, n, -1_000_000_000, 1_000_000_000)
	slices.Sort(nums)
	return nums
}

// pickTarget returns a value from nums half the time, otherwise one that is
// probably missing.
func pickTarget(rng *rand.Rand, nums []int) int {
	if len(nums) > 0 && rng.Intn(2) == 0 {
		return nums[rng.Intn(len(nums))]
	}
	if len(nums) > 0 && rng.Intn(2) == 0 {
		return nums[rng.Intn(len(nums))] + 1
	}
	return rng.Intn(2_000_000_001) - 1_000_000_000
}

// find-sorted: different values in ascending order and a target; the answer
// is its position or -1.
func genFindSorted(rng *rand.Rand, size int) Case {
	nums := sortedDistinct(rng, max(1, size))
	target := pickTarget(rng, nums)
	return Case{Input: arrayInput(nums) + "\n" + strconv.Itoa(target), Expected: strconv.Itoa(slicesIndex(nums, target))}
}

func slicesIndex(nums []int, v int) int {
	if i, ok := slices.BinarySearch(nums, v); ok {
		return i
	}
	return -1
}

// int-sqrt: one whole number x; the answer is the floor of its square root.
// Values cluster around perfect squares and reach the 32-bit limit.
func genIntSqrt(rng *rand.Rand, size int) Case {
	var x int
	switch rng.Intn(5) {
	case 0:
		x = rng.Intn(100)
	case 1:
		r := rng.Intn(46341)
		x = r*r + rng.Intn(3) - 1
	case 2:
		x = math.MaxInt32 - rng.Intn(1000)
	default:
		x = rng.Intn(math.MaxInt32)
	}
	x = max(0, min(x, math.MaxInt32))
	return Case{Input: strconv.Itoa(x), Expected: strconv.Itoa(intSqrt(x))}
}

func intSqrt(x int) int {
	r := int(math.Sqrt(float64(x)))
	for r*r > x {
		r--
	}
	for (r+1)*(r+1) <= x {
		r++
	}
	return r
}

// first-last: ascending values with long repeats and a target; the answer
// is [first, last] position or [-1, -1].
func genFirstLast(rng *rand.Rand, size int) Case {
	n := size
	if rng.Intn(12) == 0 {
		n = 0
	}
	nums := make([]int, n)
	v := rng.Intn(2001) - 1000
	for i := range nums {
		if rng.Intn(3) == 0 {
			v += 1 + rng.Intn(3)
		}
		nums[i] = v
	}
	target := rng.Intn(2001) - 1000
	if n > 0 && rng.Intn(3) != 0 {
		target = nums[rng.Intn(n)]
	}
	first, last := -1, -1
	for i, x := range nums {
		if x == target {
			if first < 0 {
				first = i
			}
			last = i
		}
	}
	return Case{Input: arrayInput(nums) + "\n" + strconv.Itoa(target), Expected: listOutput([]int{first, last})}
}

// rotated returns sorted different values turned so they start at a random
// position.
func rotated(rng *rand.Rand, n int) []int {
	nums := sortedDistinct(rng, n)
	k := rng.Intn(n)
	return append(slices.Clone(nums[k:]), nums[:k]...)
}

// min-rotated: the answer is the smallest value of a rotated sorted list.
func genMinRotated(rng *rand.Rand, size int) Case {
	nums := rotated(rng, max(1, size))
	return Case{Input: arrayInput(nums), Expected: strconv.Itoa(slices.Min(nums))}
}

// search-rotated: a rotated sorted list of different values and a target;
// the answer is its position or -1. Most found targets come from the front
// part, before the drop, where a plain binary search looks the wrong way.
func genSearchRotated(rng *rand.Rand, size int) Case {
	nums := rotated(rng, max(1, size))
	target := pickTarget(rng, nums)
	if drop := dropAt(nums); drop > 0 && rng.Intn(3) != 0 {
		target = nums[rng.Intn(drop)]
	}
	return Case{Input: arrayInput(nums) + "\n" + strconv.Itoa(target), Expected: strconv.Itoa(slices.Index(nums, target))}
}

// dropAt returns the position of the smallest value, where the list drops.
func dropAt(nums []int) int {
	return slices.Index(nums, slices.Min(nums))
}

// mountain-top: values that strictly rise and then strictly fall (either part
// may be empty); the answer is the position of the top.
func genMountainTop(rng *rand.Rand, size int) Case {
	n := max(1, size)
	top := rng.Intn(n)
	if rng.Intn(6) == 0 {
		top = []int{0, n - 1}[rng.Intn(2)]
	}
	vals := sortedDistinct(rng, n)
	// The largest value goes at top; the rest split into a rising left part
	// and a falling right part.
	rest := vals[:n-1]
	rng.Shuffle(len(rest), func(i, j int) { rest[i], rest[j] = rest[j], rest[i] })
	left := slices.Clone(rest[:top])
	right := slices.Clone(rest[top:])
	slices.Sort(left)
	slices.Sort(right)
	slices.Reverse(right)
	nums := append(append(left, vals[n-1]), right...)
	return Case{Input: arrayInput(nums), Expected: strconv.Itoa(top)}
}

// eating-speed: job sizes and an hour limit h ≥ n; the answer is the
// smallest speed that finishes every job within h hours, one job per hour
// at most.
func genEatingSpeed(rng *rand.Rand, size int) Case {
	n := max(1, size)
	top := []int{10, 10_000, 1_000_000_000}[rng.Intn(3)]
	jobs := randomInts(rng, n, 1, top)
	total := 0
	for _, j := range jobs {
		total += j
	}
	h := n + rng.Intn(min(total, 1_000_000_000-n)+1)
	if rng.Intn(4) == 0 {
		h = n
	}
	return Case{Input: arrayInput(jobs) + "\n" + strconv.Itoa(h), Expected: strconv.Itoa(eatingSpeed(jobs, h))}
}

func hoursAt(jobs []int, k int) int {
	hours := 0
	for _, j := range jobs {
		hours += (j + k - 1) / k
	}
	return hours
}

func eatingSpeed(jobs []int, h int) int {
	lo, hi := 1, slices.Max(jobs)
	for lo < hi {
		mid := lo + (hi-lo)/2
		if hoursAt(jobs, mid) <= h {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}

// median-two: two sorted lists, together not empty; the answer is the median
// of all values with one decimal place.
func genMedianTwo(rng *rand.Rand, size int) Case {
	total := max(1, size)
	m := rng.Intn(total + 1)
	span := []int{3, 1000, 1_000_000}[rng.Intn(3)]
	a, b := sortedRandom(rng, m, span), sortedRandom(rng, total-m, span)
	return Case{Input: arrayInput(a) + "\n" + arrayInput(b), Expected: medianText(append(slices.Clone(a), b...))}
}

func medianText(all []int) string {
	slices.Sort(all)
	n := len(all)
	twice := 2 * all[n/2]
	if n%2 == 0 {
		twice = all[n/2-1] + all[n/2]
	}
	return halfText(twice)
}

// halfText prints twice/2 with one decimal place, as "-5.5" or "3.0".
func halfText(twice int) string {
	sign := ""
	if twice < 0 {
		sign, twice = "-", -twice
	}
	s := sign + strconv.Itoa(twice/2) + ".0"
	if twice%2 != 0 {
		s = sign + strconv.Itoa(twice/2) + ".5"
	}
	return s
}
