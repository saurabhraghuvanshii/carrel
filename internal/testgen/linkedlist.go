package testgen

import (
	"math/rand"
	"slices"
	"strconv"
)

// Generators for the "Linked list" group. Lists use the array format; cycle
// problems add a line with the position the last node points back to, or -1.

func init() {
	register("merge-lists", genMergeLists)
	register("list-middle", genListMiddle)
	register("mirror-list", genMirrorList)
	register("list-cycle", genListCycle)
	register("cycle-start", genCycleStart)
	register("remove-nth", genRemoveNth)
	register("reorder-list", genReorderList)
	register("add-digits", genAddDigits)
	register("reverse-groups", genReverseGroups)
}

func listSize(rng *rand.Rand, size int, allowEmpty bool) int {
	switch {
	case allowEmpty && rng.Intn(12) == 0:
		return 0
	case rng.Intn(12) == 0:
		return 1
	}
	return max(1, size)
}

// merge-lists: two sorted lists; the answer is them merged.
func genMergeLists(rng *rand.Rand, size int) Case {
	total := max(0, size)
	m := rng.Intn(total + 1)
	span := []int{3, 1000}[rng.Intn(2)]
	a, b := sortedRandom(rng, m, span), sortedRandom(rng, total-m, span)
	merged := append(slices.Clone(a), b...)
	slices.Sort(merged)
	return Case{Input: arrayInput(a) + "\n" + arrayInput(b), Expected: listOutput(merged)}
}

// list-middle: the answer is the list from the middle node on (the second
// middle when the length is even).
func genListMiddle(rng *rand.Rand, size int) Case {
	vals := randomInts(rng, listSize(rng, size, false), -1000, 1000)
	return Case{Input: arrayInput(vals), Expected: listOutput(vals[len(vals)/2:])}
}

// mirror-list: half the lists read the same both ways, the rest differ in
// one place.
func genMirrorList(rng *rand.Rand, size int) Case {
	n := listSize(rng, size, false)
	vals := randomInts(rng, n, 0, []int{2, 9, 1000}[rng.Intn(3)])
	for i := 0; i < n/2; i++ {
		vals[n-1-i] = vals[i]
	}
	if n > 1 && rng.Intn(2) == 0 {
		vals[rng.Intn(n)] += 1 + rng.Intn(3)
	}
	return Case{Input: arrayInput(vals), Expected: boolText(isMirror(vals))}
}

func isMirror(vals []int) bool {
	for i, j := 0, len(vals)-1; i < j; i, j = i+1, j-1 {
		if vals[i] != vals[j] {
			return false
		}
	}
	return true
}

// cycleCase builds a list and a loop position, -1 about a third of the time.
func cycleCase(rng *rand.Rand, size int) ([]int, int) {
	n := listSize(rng, size, true)
	vals := randomInts(rng, n, -1000, 1000)
	pos := -1
	if n > 0 && rng.Intn(3) != 0 {
		pos = rng.Intn(n)
	}
	return vals, pos
}

// list-cycle: the answer is whether the list loops.
func genListCycle(rng *rand.Rand, size int) Case {
	vals, pos := cycleCase(rng, size)
	return Case{Input: arrayInput(vals) + "\n" + strconv.Itoa(pos), Expected: boolText(pos >= 0)}
}

// cycle-start: the answer is the position where the loop starts, or -1.
func genCycleStart(rng *rand.Rand, size int) Case {
	vals, pos := cycleCase(rng, size)
	return Case{Input: arrayInput(vals) + "\n" + strconv.Itoa(pos), Expected: strconv.Itoa(pos)}
}

// remove-nth: a list and n; the answer is the list without its n-th node
// from the end.
func genRemoveNth(rng *rand.Rand, size int) Case {
	vals := randomInts(rng, listSize(rng, size, false), -1000, 1000)
	n := 1 + rng.Intn(len(vals))
	if rng.Intn(4) == 0 {
		n = []int{1, len(vals)}[rng.Intn(2)]
	}
	at := len(vals) - n
	rest := append(slices.Clone(vals[:at]), vals[at+1:]...)
	return Case{Input: arrayInput(vals) + "\n" + strconv.Itoa(n), Expected: listOutput(rest)}
}

// reorder-list: the answer is first, last, second, second to last, ...
func genReorderList(rng *rand.Rand, size int) Case {
	vals := randomInts(rng, listSize(rng, size, false), -1000, 1000)
	var out []int
	for i, j := 0, len(vals)-1; i <= j; i, j = i+1, j-1 {
		out = append(out, vals[i])
		if i != j {
			out = append(out, vals[j])
		}
	}
	return Case{Input: arrayInput(vals), Expected: listOutput(out)}
}

// digitList returns a number's digits, lowest first, with no leading zero.
func digitList(rng *rand.Rand, n int) []int {
	d := randomInts(rng, n, 0, 9)
	if n > 1 && d[n-1] == 0 {
		d[n-1] = 1 + rng.Intn(9)
	}
	if rng.Intn(5) == 0 {
		for i := range d {
			d[i] = 9 // carries all the way
		}
	}
	return d
}

// add-digits: two numbers as digit lists, lowest digit first; the answer is
// their sum the same way.
func genAddDigits(rng *rand.Rand, size int) Case {
	a := digitList(rng, 1+rng.Intn(max(1, size)))
	b := digitList(rng, 1+rng.Intn(max(1, size)))
	return Case{Input: arrayInput(a) + "\n" + arrayInput(b), Expected: listOutput(addDigits(a, b))}
}

func addDigits(a, b []int) []int {
	var out []int
	carry := 0
	for i := 0; i < len(a) || i < len(b) || carry > 0; i++ {
		s := carry
		if i < len(a) {
			s += a[i]
		}
		if i < len(b) {
			s += b[i]
		}
		out = append(out, s%10)
		carry = s / 10
	}
	return out
}

// reverse-groups: a list and k; the answer reverses each full group of k
// nodes and leaves a shorter last group as it is.
func genReverseGroups(rng *rand.Rand, size int) Case {
	vals := randomInts(rng, listSize(rng, size, false), -1000, 1000)
	k := 1 + rng.Intn(min(len(vals)+2, 12))
	out := slices.Clone(vals)
	for i := 0; i+k <= len(out); i += k {
		slices.Reverse(out[i : i+k])
	}
	return Case{Input: arrayInput(vals) + "\n" + strconv.Itoa(k), Expected: listOutput(out)}
}
