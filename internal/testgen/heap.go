package testgen

import (
	"container/heap"
	"math/rand"
	"slices"
	"strconv"
	"strings"
)

// Generators for the "Heap" group.

func init() {
	register("kth-largest", genKthLargest)
	register("smash-stones", genSmashStones)
	register("rope-cost", genRopeCost)
	register("closest-points", genClosestPoints)
	register("cooldown", genCooldown)
	register("merge-k", genMergeK)
	register("running-median", genRunningMedian)
}

// maxHeap is a heap of ints with the largest on top.
type maxHeap []int

func (h maxHeap) Len() int           { return len(h) }
func (h maxHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h maxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *maxHeap) Pop() any {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

// kth-largest: values with repeats and k; the answer is the k-th largest,
// counting repeats.
func genKthLargest(rng *rand.Rand, size int) Case {
	n := max(1, size)
	span := []int{10, 10_000}[rng.Intn(2)] // a small span means many repeats
	nums := randomInts(rng, n, -span, span)
	k := 1 + rng.Intn(n)
	s := slices.Clone(nums)
	slices.Sort(s)
	return Case{Input: arrayInput(nums) + "\n" + strconv.Itoa(k), Expected: strconv.Itoa(s[n-k])}
}

// smash-stones: stone weights; the two heaviest collide again and again.
func genSmashStones(rng *rand.Rand, size int) Case {
	stones := randomInts(rng, max(1, size), 1, []int{10, 1000}[rng.Intn(2)])
	h := maxHeap(slices.Clone(stones))
	heap.Init(&h)
	for h.Len() > 1 {
		y, x := heap.Pop(&h).(int), heap.Pop(&h).(int)
		if y != x {
			heap.Push(&h, y-x)
		}
	}
	last := 0
	if h.Len() == 1 {
		last = h[0]
	}
	return Case{Input: arrayInput(stones), Expected: strconv.Itoa(last)}
}

// rope-cost: rope lengths; joining two costs their sum. The answer is the
// cheapest total, always joining the two shortest.
func genRopeCost(rng *rand.Rand, size int) Case {
	ropes := randomInts(rng, max(1, size), 1, []int{10, 1000}[rng.Intn(2)])
	h := maxHeap{}
	for _, r := range ropes {
		heap.Push(&h, -r)
	}
	cost := 0
	for h.Len() > 1 {
		a, b := -heap.Pop(&h).(int), -heap.Pop(&h).(int)
		cost += a + b
		heap.Push(&h, -(a + b))
	}
	return Case{Input: arrayInput(ropes), Expected: strconv.Itoa(cost)}
}

// closest-points: n points on lines "x y", then k. The answer lists the
// positions of the k points nearest the origin, ascending; the k-th and the
// next distance always differ, so the answer is unique.
func genClosestPoints(rng *rand.Rand, size int) Case {
	n := max(1, size)
	span := []int{3, 100, 10_000}[rng.Intn(3)]
	for {
		xs, ys := randomInts(rng, n, -span, span), randomInts(rng, n, -span, span)
		d := make([]int, n)
		for i := range d {
			d[i] = xs[i]*xs[i] + ys[i]*ys[i]
		}
		k := 1 + rng.Intn(n)
		sorted := slices.Clone(d)
		slices.Sort(sorted)
		if k < n && sorted[k-1] == sorted[k] {
			continue
		}
		var b strings.Builder
		b.WriteString(strconv.Itoa(n))
		var picked []int
		for i := range d {
			b.WriteString("\n" + strconv.Itoa(xs[i]) + " " + strconv.Itoa(ys[i]))
			if d[i] <= sorted[k-1] {
				picked = append(picked, i)
			}
		}
		return Case{Input: b.String() + "\n" + strconv.Itoa(k), Expected: listOutput(picked)}
	}
}

// cooldown: a word of task letters and n; the same task needs n other time
// units between runs. The answer is the fewest time units, idle included.
func genCooldown(rng *rand.Rand, size int) Case {
	letters := 1 + rng.Intn(26)
	if rng.Intn(2) == 0 {
		letters = 1 + rng.Intn(4) // few letters, so idle time matters
	}
	tasks := randomWord(rng, max(1, size), "ABCDEFGHIJKLMNOPQRSTUVWXYZ"[:letters])
	n := rng.Intn(min(10, max(1, size)) + 1)
	var count [26]int
	for i := 0; i < len(tasks); i++ {
		count[tasks[i]-'A']++
	}
	most := slices.Max(count[:])
	withMost := 0
	for _, c := range count {
		if c == most {
			withMost++
		}
	}
	return Case{Input: tasks + "\n" + strconv.Itoa(n), Expected: strconv.Itoa(max(len(tasks), (most-1)*(n+1)+withMost))}
}

// merge-k: k sorted lists, any of them empty; the answer is them merged.
func genMergeK(rng *rand.Rand, size int) Case {
	total := max(0, size)
	k := rng.Intn(min(total, 1000) + 1)
	if rng.Intn(10) == 0 {
		k = 1
	}
	var b strings.Builder
	b.WriteString(strconv.Itoa(k))
	var all []int
	left := total
	for i := 0; i < k; i++ {
		n := 0
		if i == k-1 {
			n = left
		} else if left > 0 {
			n = rng.Intn(min(left, 2*total/max(1, k)) + 1)
		}
		left -= n
		list := sortedRandom(rng, n, []int{5, 10_000}[rng.Intn(2)])
		all = append(all, list...)
		b.WriteString("\n" + arrayInput(list))
	}
	slices.Sort(all)
	return Case{Input: b.String(), Expected: listOutput(all)}
}

// running-median: a call sequence of new, add x and median; median only
// comes after at least one add. Medians print with one decimal place.
func genRunningMedian(rng *rand.Rand, size int) Case {
	k := min(max(size, 3), 5000)
	calls, out := []string{"new"}, []string{"null"}
	var sorted []int
	for len(calls) < k {
		if len(sorted) == 0 || rng.Intn(5) < 3 {
			v := rng.Intn(200_001) - 100_000
			if rng.Intn(4) == 0 && len(sorted) > 0 {
				v = sorted[rng.Intn(len(sorted))] // a repeat
			}
			i, _ := slices.BinarySearch(sorted, v)
			sorted = slices.Insert(sorted, i, v)
			calls, out = append(calls, "add "+strconv.Itoa(v)), append(out, "null")
		} else {
			n := len(sorted)
			twice := 2 * sorted[n/2]
			if n%2 == 0 {
				twice = sorted[n/2-1] + sorted[n/2]
			}
			calls, out = append(calls, "median"), append(out, halfText(twice))
		}
	}
	return Case{Input: strconv.Itoa(len(calls)) + "\n" + strings.Join(calls, "\n"), Expected: strings.Join(out, " ")}
}
