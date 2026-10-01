package testgen

import (
	"math/rand"
	"slices"
	"strconv"
	"strings"
)

// Generators for the "Intervals" group. A list of intervals is "n" then n
// lines "start end"; answers print as [[1, 3], [5, 8]].

func init() {
	register("meeting-conflict", genMeetingConflict)
	register("merge-intervals", genMergeIntervals)
	register("insert-interval", genInsertInterval)
	register("rooms", genRooms)
	register("remove-overlaps", genRemoveOverlaps)
	register("arrows", genArrows)
	register("free-time", genFreeTime)
}

type span struct{ s, e int }

func intervalsInput(ivs []span) string {
	var b strings.Builder
	b.WriteString(strconv.Itoa(len(ivs)))
	for _, iv := range ivs {
		b.WriteString("\n" + strconv.Itoa(iv.s) + " " + strconv.Itoa(iv.e))
	}
	return b.String()
}

func intervalsOutput(ivs []span) string {
	rows := make([][]int, len(ivs))
	for i, iv := range ivs {
		rows[i] = []int{iv.s, iv.e}
	}
	return gridOutput(rows)
}

// randomSpans makes n intervals with start < end. A tight range gives many
// overlaps; a loose one gives few.
func randomSpans(rng *rand.Rand, n int) []span {
	room := max(10, n*[]int{1, 5, 50}[rng.Intn(3)])
	length := max(1, room/max(1, n)*[]int{1, 3, 10}[rng.Intn(3)])
	ivs := make([]span, n)
	for i := range ivs {
		s := rng.Intn(room)
		ivs[i] = span{s, s + 1 + rng.Intn(length)}
	}
	return ivs
}

// disjointSpans walks right, leaving gaps of minGap to minGap+gap. With
// minGap 0, neighbours may touch.
func disjointSpans(rng *rand.Rand, n, minGap, gap int) []span {
	ivs := make([]span, n)
	t := rng.Intn(10)
	for i := range ivs {
		t += minGap + rng.Intn(gap+1)
		e := t + 1 + rng.Intn(10)
		ivs[i] = span{t, e}
		t = e
	}
	return ivs
}

func byStart(ivs []span) []span {
	s := slices.Clone(ivs)
	slices.SortFunc(s, func(a, b span) int {
		if a.s != b.s {
			return a.s - b.s
		}
		return a.e - b.e
	})
	return s
}

// meeting-conflict: meetings [start, end); the answer is whether none
// overlap. Half the cases are disjoint, shuffled.
func genMeetingConflict(rng *rand.Rand, size int) Case {
	n := max(1, size)
	var ivs []span
	if rng.Intn(2) == 0 {
		ivs = disjointSpans(rng, n, 0, rng.Intn(3))
		rng.Shuffle(n, func(i, j int) { ivs[i], ivs[j] = ivs[j], ivs[i] })
	} else {
		ivs = randomSpans(rng, n)
	}
	ok := true
	s := byStart(ivs)
	for i := 1; i < len(s); i++ {
		ok = ok && s[i].s >= s[i-1].e
	}
	return Case{Input: intervalsInput(ivs), Expected: boolText(ok)}
}

// mergeClosed merges intervals that overlap or touch.
func mergeClosed(ivs []span) []span {
	var out []span
	for _, iv := range byStart(ivs) {
		if len(out) > 0 && iv.s <= out[len(out)-1].e {
			out[len(out)-1].e = max(out[len(out)-1].e, iv.e)
		} else {
			out = append(out, iv)
		}
	}
	return out
}

// merge-intervals: the answer merges every overlapping or touching pair,
// sorted by start.
func genMergeIntervals(rng *rand.Rand, size int) Case {
	ivs := randomSpans(rng, max(1, size))
	return Case{Input: intervalsInput(ivs), Expected: intervalsOutput(mergeClosed(ivs))}
}

// insert-interval: sorted separate intervals (gaps of at least 1), then a new
// interval on its own line; the answer inserts and merges it.
func genInsertInterval(rng *rand.Rand, size int) Case {
	n := size
	if rng.Intn(10) == 0 {
		n = 0
	}
	ivs := disjointSpans(rng, n, 1, rng.Intn(20)) // gaps of at least 1: already merged
	end := 20
	if n > 0 {
		end = ivs[n-1].e + 5
	}
	s := rng.Intn(end)
	add := span{s, s + 1 + rng.Intn(max(1, end/[]int{1, 4, 20}[rng.Intn(3)]))}
	return Case{
		Input:    intervalsInput(ivs) + "\n" + strconv.Itoa(add.s) + " " + strconv.Itoa(add.e),
		Expected: intervalsOutput(mergeClosed(append(slices.Clone(ivs), add))),
	}
}

// rooms: meetings [start, end); the answer is the most that run at once.
func genRooms(rng *rand.Rand, size int) Case {
	ivs := randomSpans(rng, max(1, size))
	for i := range ivs {
		if rng.Intn(3) == 0 { // start exactly when another meeting ends
			length := ivs[i].e - ivs[i].s
			ivs[i].s = ivs[rng.Intn(len(ivs))].e
			ivs[i].e = ivs[i].s + length
		}
	}
	type ev struct{ t, d int }
	var evs []ev
	for _, iv := range ivs {
		evs = append(evs, ev{iv.s, 1}, ev{iv.e, -1})
	}
	slices.SortFunc(evs, func(a, b ev) int {
		if a.t != b.t {
			return a.t - b.t
		}
		return a.d - b.d // an end before a start at the same time
	})
	cur, best := 0, 0
	for _, e := range evs {
		cur += e.d
		best = max(best, cur)
	}
	return Case{Input: intervalsInput(ivs), Expected: strconv.Itoa(best)}
}

// remove-overlaps: intervals [start, end); the answer is the fewest to remove
// so the rest do not overlap.
func genRemoveOverlaps(rng *rand.Rand, size int) Case {
	ivs := randomSpans(rng, max(1, size))
	s := slices.Clone(ivs)
	slices.SortFunc(s, func(a, b span) int { return a.e - b.e })
	kept, last := 0, -1<<62
	for _, iv := range s {
		if iv.s >= last {
			kept++
			last = iv.e
		}
	}
	return Case{Input: intervalsInput(ivs), Expected: strconv.Itoa(len(ivs) - kept)}
}

// arrows: balloons [start, end] including both ends; an arrow at x bursts
// every balloon with start ≤ x ≤ end. The answer is the fewest arrows.
func genArrows(rng *rand.Rand, size int) Case {
	ivs := randomSpans(rng, max(1, size))
	for i := range ivs {
		if rng.Intn(6) == 0 {
			ivs[i].e = ivs[i].s // a balloon of zero width
		}
	}
	s := slices.Clone(ivs)
	slices.SortFunc(s, func(a, b span) int { return a.e - b.e })
	arrows, at := 0, -1<<62
	for _, iv := range s {
		if iv.s > at {
			arrows++
			at = iv.e
		}
	}
	return Case{Input: intervalsInput(ivs), Expected: strconv.Itoa(arrows)}
}

// free-time: k people, each with sorted separate busy intervals; the answer
// is the gaps of positive length when nobody is busy, between the first
// start and the last end.
func genFreeTime(rng *rand.Rand, size int) Case {
	total := max(1, size)
	k := 1 + rng.Intn(min(50, total))
	var b strings.Builder
	b.WriteString(strconv.Itoa(k))
	var all []span
	for i := 0; i < k; i++ {
		m := max(1, total/k+rng.Intn(3)-1)
		ivs := disjointSpans(rng, m, 1, rng.Intn(40))
		all = append(all, ivs...)
		b.WriteString("\n" + intervalsInput(ivs))
	}
	busy := mergeClosed(all)
	var free []span
	for i := 1; i < len(busy); i++ {
		free = append(free, span{busy[i-1].e, busy[i].s})
	}
	return Case{Input: b.String(), Expected: intervalsOutput(free)}
}
