package testgen

import (
	"container/heap"
	"math/rand"
	"slices"
	"strconv"
	"strings"
)

// Generators for the real-interview group "Ride dispatch".

func init() {
	register("surge-window", genSurgeWindow)
	register("shift-gaps", genShiftGaps)
	register("pickup-batches", genPickupBatches)
	register("eta-lights", genEtaLights)
	register("pool-match", genPoolMatch)
}

// surgeFor is the multiplier for one window: 1 when supply covers demand,
// otherwise demand over supply rounded up, at most 5.
func surgeFor(requests, drivers int) int {
	if requests <= drivers {
		return 1
	}
	if drivers == 0 {
		return 5
	}
	return min(5, (requests+drivers-1)/drivers)
}

// surge-window: requests and free drivers per minute, and a window k; the
// answer is the multiplier for every full window, oldest first.
func genSurgeWindow(rng *rand.Rand, size int) Case {
	n := max(1, size)
	k := 1 + rng.Intn(min(n, []int{3, 10, n}[rng.Intn(3)]))
	reqTop, drvTop := []int{5, 100, 1000}[rng.Intn(3)], []int{2, 50, 1000}[rng.Intn(3)]
	req, drv := randomInts(rng, n, 0, reqTop), randomInts(rng, n, 0, drvTop)
	out := make([]int, 0, n-k+1)
	r, d := 0, 0
	for i := 0; i < n; i++ {
		r, d = r+req[i], d+drv[i]
		if i >= k {
			r, d = r-req[i-k], d-drv[i-k]
		}
		if i >= k-1 {
			out = append(out, surgeFor(r, d))
		}
	}
	return Case{
		Input:    strconv.Itoa(n) + " " + strconv.Itoa(k) + "\n" + joinInts(req, " ") + "\n" + joinInts(drv, " "),
		Expected: listOutput(out),
	}
}

// shift-gaps: shifts [start, end) in a day [0, dayEnd); the answer lists the
// longest stretches when fewer than m drivers are on shift, in time order.
func genShiftGaps(rng *rand.Rand, size int) Case {
	n := max(1, size)
	dayEnd := []int{20, 1000, 1_000_000_000}[rng.Intn(3)]
	m := 1 + rng.Intn([]int{1, 3, 6}[rng.Intn(3)])
	shifts := make([]span, n)
	for i := range shifts {
		s := rng.Intn(dayEnd)
		shifts[i] = span{s, s + 1 + rng.Intn(max(1, min(dayEnd-s, dayEnd/max(1, n/[]int{1, 4}[rng.Intn(2)]))))}
		shifts[i].e = min(shifts[i].e, dayEnd)
		if rng.Intn(5) == 0 && i > 0 { // start exactly when another shift ends
			prev := shifts[rng.Intn(i)]
			if prev.e < dayEnd {
				shifts[i] = span{prev.e, prev.e + 1 + rng.Intn(dayEnd-prev.e)}
			}
		}
	}
	return Case{
		Input:    strconv.Itoa(n) + " " + strconv.Itoa(m) + " " + strconv.Itoa(dayEnd) + "\n" + strings.TrimPrefix(intervalsInput(shifts), strconv.Itoa(n)+"\n"),
		Expected: intervalsOutput(shiftGaps(shifts, m, dayEnd)),
	}
}

func shiftGaps(shifts []span, m, dayEnd int) []span {
	type ev struct{ t, d int }
	evs := []ev{{0, 0}, {dayEnd, 0}}
	for _, s := range shifts {
		evs = append(evs, ev{s.s, 1}, ev{s.e, -1})
	}
	slices.SortFunc(evs, func(a, b ev) int { return a.t - b.t })
	var out []span
	on := 0
	for i := 0; i < len(evs); {
		t := evs[i].t
		for i < len(evs) && evs[i].t == t {
			on += evs[i].d
			i++
		}
		if i < len(evs) && on < m && t < dayEnd {
			next := evs[i].t
			if len(out) > 0 && out[len(out)-1].e == t {
				out[len(out)-1].e = next
			} else if next > t {
				out = append(out, span{t, next})
			}
		}
	}
	return out
}

func floorDiv(a, b int) int {
	q := a / b
	if a%b != 0 && a < 0 {
		q--
	}
	return q
}

// pickup-batches: requests "time x y" in time order. A request joins the open
// batch of its zone (floor(x/100), floor(y/100)) when it is at most window
// after that batch's first request and the batch has room; otherwise it
// starts a new batch. The answer is the number of batches.
func genPickupBatches(rng *rand.Rand, size int) Case {
	n := max(1, size)
	limit := 1 + rng.Intn(5)
	window := []int{0, 5, 60}[rng.Intn(3)]
	reach := []int{150, 400, 1_000_000}[rng.Intn(3)]
	lines := []string{strconv.Itoa(n) + " " + strconv.Itoa(limit) + " " + strconv.Itoa(window)}
	type batch struct{ start, size int }
	open := map[[2]int]*batch{}
	count, t := 0, 0
	for i := 0; i < n; i++ {
		t += rng.Intn(10)
		x, y := rng.Intn(2*reach+1)-reach, rng.Intn(2*reach+1)-reach
		lines = append(lines, strconv.Itoa(t)+" "+strconv.Itoa(x)+" "+strconv.Itoa(y))
		zone := [2]int{floorDiv(x, 100), floorDiv(y, 100)}
		if b := open[zone]; b != nil && t-b.start <= window && b.size < limit {
			b.size++
			continue
		}
		open[zone] = &batch{t, 1}
		count++
	}
	return Case{Input: strings.Join(lines, "\n"), Expected: strconv.Itoa(count)}
}

// eta-lights: a road graph with travel times, and a light period at every
// junction: a car may leave junction v only at times that are multiples of
// period[v]. The answer is the earliest arrival at n-1 leaving 0 at time 0,
// or -1.
func genEtaLights(rng *rand.Rand, size int) Case {
	n := 2 + rng.Intn(max(1, min(size, 9999)))
	m := n/2 + rng.Intn(3*n)
	periods := make([]int, n)
	for i := range periods {
		periods[i] = 1 + rng.Intn([]int{1, 5, 100}[rng.Intn(3)])
	}
	lines := []string{strconv.Itoa(n) + " " + strconv.Itoa(m), joinInts(periods, " ")}
	adj := make([][][2]int, n)
	for i := 0; i < m; i++ {
		u, v := rng.Intn(n), rng.Intn(n)
		if rng.Intn(3) != 0 && u+1 < n { // a road forward keeps most cases reachable
			v = u + 1 + rng.Intn(min(3, n-u-1))
		}
		w := 1 + rng.Intn([]int{10, 1000}[rng.Intn(2)])
		adj[u] = append(adj[u], [2]int{v, w})
		lines = append(lines, strconv.Itoa(u)+" "+strconv.Itoa(v)+" "+strconv.Itoa(w))
	}
	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}
	dist[0] = 0
	h := &distHeap{{0, 0}}
	for h.Len() > 0 {
		cur := heap.Pop(h).(distItem)
		if cur.dist > dist[cur.node] {
			continue
		}
		p := periods[cur.node]
		leave := (cur.dist + p - 1) / p * p
		for _, e := range adj[cur.node] {
			if at := leave + e[1]; dist[e[0]] < 0 || at < dist[e[0]] {
				dist[e[0]] = at
				heap.Push(h, distItem{e[0], at})
			}
		}
	}
	return Case{Input: strings.Join(lines, "\n"), Expected: strconv.Itoa(dist[n-1])}
}

// pool-match: a shuttle with seats along stops 0 to stops; each rider rides
// from one stop to a later one. The answer is the most riders it can take
// with never more than seats on board.
func genPoolMatch(rng *rand.Rand, size int) Case {
	n := max(1, min(2000, size))
	stops := []int{5, 30, 1000}[rng.Intn(3)]
	seats := 1 + rng.Intn([]int{1, 3, 50}[rng.Intn(3)])
	rides := make([]span, n)
	for i := range rides {
		a := rng.Intn(stops)
		rides[i] = span{a, a + 1 + rng.Intn(min(stops-a, []int{3, stops}[rng.Intn(2)]))}
	}
	return Case{
		Input:    strconv.Itoa(stops) + " " + strconv.Itoa(seats) + " " + strconv.Itoa(n) + "\n" + strings.TrimPrefix(intervalsInput(rides), strconv.Itoa(n)+"\n"),
		Expected: strconv.Itoa(poolMatch(rides, stops, seats)),
	}
}

// poolMatch takes riders by earliest drop-off, each one that still fits.
func poolMatch(rides []span, stops, seats int) int {
	s := slices.Clone(rides)
	slices.SortFunc(s, func(a, b span) int { return a.e - b.e })
	load := make([]int, stops)
	taken := 0
	for _, r := range s {
		if slices.Max(load[r.s:r.e]) < seats {
			for x := r.s; x < r.e; x++ {
				load[x]++
			}
			taken++
		}
	}
	return taken
}
