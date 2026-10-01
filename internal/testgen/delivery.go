package testgen

import (
	"math/rand"
	"slices"
	"strconv"
)

// Generators for the real-interview group "Maps and delivery".

func init() {
	register("nearest-store", genNearestStore)
	register("geofence", genGeofence)
	register("street-walk", genStreetWalk)
	register("zone-lookup", genZoneLookup)
	register("courier-assignment", genCourierAssignment)
}

// nearest-store: "n m", the n store positions, the m customer positions.
// The answer is each customer's distance to the nearest store.
func genNearestStore(rng *rand.Rand, size int) Case {
	n, m := 1+rng.Intn(max(1, size)), max(1, size)
	top := []int{20, 1000, 1_000_000_000}[rng.Intn(3)]
	stores, customers := randomInts(rng, n, 0, top), randomInts(rng, m, 0, top)
	sorted := sortedCopy(stores)
	out := make([]int, m)
	for i, c := range customers {
		at, _ := slices.BinarySearch(sorted, c)
		out[i] = top + 1
		if at < n {
			out[i] = sorted[at] - c
		}
		if at > 0 {
			out[i] = min(out[i], c-sorted[at-1])
		}
	}
	return Case{Input: ints(n, m) + "\n" + joinInts(stores, " ") + "\n" + joinInts(customers, " "), Expected: listOutput(out)}
}

// geofence: a rectangle "x1 y1 x2 y2", then "n" and n points of a path.
// The answer is how many times the path enters the rectangle (edges are
// inside; starting inside counts as entering).
func genGeofence(rng *rand.Rand, size int) Case {
	x1, y1 := rng.Intn(10), rng.Intn(10)
	x2, y2 := x1+rng.Intn(8), y1+rng.Intn(8)
	n := max(1, size)
	// The path wanders in a band four wide around the rectangle.
	x, y := x1-4+rng.Intn(x2-x1+9), y1-4+rng.Intn(y2-y1+9)
	rows := make([][]int, n)
	entries, wasIn := 0, false
	for i := range rows {
		rows[i] = []int{x, y}
		in := x1 <= x && x <= x2 && y1 <= y && y <= y2
		if in && !wasIn {
			entries++
		}
		wasIn = in
		x, y = x+rng.Intn(7)-3, y+rng.Intn(7)-3
		x, y = max(x1-4, min(x2+4, x)), max(y1-4, min(y2+4, y))
	}
	return Case{Input: ints(x1, y1, x2, y2) + "\n" + strconv.Itoa(n) + rowsText(rows), Expected: strconv.Itoa(entries)}
}

// street-walk: "n start", then n stop positions on a street. The answer is
// the shortest walk from start that visits every stop, in any order,
// without coming back.
func genStreetWalk(rng *rand.Rand, size int) Case {
	n := max(1, size)
	top := []int{20, 1000, 1_000_000_000}[rng.Intn(3)]
	stops := randomInts(rng, n, 0, top)
	start := rng.Intn(top + 1)
	if rng.Intn(4) == 0 {
		start = []int{0, top}[rng.Intn(2)] // often all stops on one side
	}
	left, right := max(0, start-slices.Min(stops)), max(0, slices.Max(stops)-start)
	return Case{Input: ints(n, start) + "\n" + joinInts(stops, " "), Expected: strconv.Itoa(left + right + min(left, right))}
}

// zone-lookup: "n m", then n lines "lo hi zone" (code ranges with both ends
// included, not overlapping, in any order), then the m codes. The answer is
// each code's zone, or -1.
func genZoneLookup(rng *rand.Rand, size int) Case {
	n, m := max(1, size), 1+rng.Intn(max(1, size))
	ranges := disjointSpans(rng, n, 1, rng.Intn(15))
	rows := make([][]int, n)
	for i, r := range ranges {
		if rng.Intn(4) == 0 {
			r.e = r.s // a range of one code
		}
		ranges[i] = r
		rows[i] = []int{r.s, r.e, 1 + rng.Intn(50)}
	}
	sorted := slices.Clone(rows)
	rng.Shuffle(n, func(i, j int) { rows[i], rows[j] = rows[j], rows[i] })
	end := ranges[n-1].e + 5
	codes := make([]int, m)
	out := make([]int, m)
	for i := range codes {
		codes[i] = rng.Intn(end + 1)
		if r := ranges[rng.Intn(n)]; rng.Intn(2) == 0 {
			codes[i] = []int{r.s, r.e, r.s + rng.Intn(r.e-r.s+1), r.e + 1}[rng.Intn(4)]
		}
		// The last range that starts at or before the code.
		at, found := slices.BinarySearchFunc(sorted, codes[i], func(row []int, code int) int { return row[0] - code })
		if !found {
			at--
		}
		out[i] = -1
		if at >= 0 && codes[i] <= sorted[at][1] {
			out[i] = sorted[at][2]
		}
	}
	return Case{Input: ints(n, m) + rowsText(rows) + "\n" + joinInts(codes, " "), Expected: listOutput(out)}
}

// courier-assignment: "n m", the n order positions, the m courier positions
// (m >= n). Every order gets its own courier; the answer is the smallest
// total distance.
func genCourierAssignment(rng *rand.Rand, size int) Case {
	n := max(1, min(size, 1000))
	m := min(1000, n+rng.Intn([]int{1, 4, n + 1}[rng.Intn(3)]))
	top := []int{20, 1000, 1_000_000_000}[rng.Intn(3)]
	orders, couriers := randomInts(rng, n, 0, top), randomInts(rng, m, 0, top)
	o, c := sortedCopy(orders), sortedCopy(couriers)
	// best[j] is the cheapest way to serve the first i orders with the
	// first j couriers.
	const inf = 1 << 62
	prev := make([]int, m+1)
	for i := 1; i <= n; i++ {
		cur := make([]int, m+1)
		for j := range cur {
			cur[j] = inf
		}
		for j := i; j <= m; j++ {
			cur[j] = min(cur[j-1], prev[j-1]+max(o[i-1]-c[j-1], c[j-1]-o[i-1]))
		}
		prev = cur
	}
	return Case{Input: ints(n, m) + "\n" + joinInts(orders, " ") + "\n" + joinInts(couriers, " "), Expected: strconv.Itoa(prev[m])}
}
