package testgen

import (
	"math/rand"
	"strconv"
)

// Generators for the real-interview group "Scheduling and booking".

func init() {
	register("calendar-conflict", genCalendarConflict)
	register("room-booking", genRoomBooking)
	register("job-deps", genJobDeps)
	register("parking-lot", genParkingLot)
	register("elevator", genElevator)
	register("seat-hold", genSeatHold)
}

// calendar-conflict: meetings [start, end), then a new meeting on the last
// line; the answer is how many meetings the new one overlaps.
func genCalendarConflict(rng *rand.Rand, size int) Case {
	ivs := randomSpans(rng, max(1, size))
	add := randomSpans(rng, max(1, size))[0]
	switch rng.Intn(3) {
	case 0: // start exactly where a meeting ends
		length := add.e - add.s
		add.s = ivs[rng.Intn(len(ivs))].e
		add.e = add.s + length
	case 1: // reach from one meeting to another
		a, b := ivs[rng.Intn(len(ivs))], ivs[rng.Intn(len(ivs))]
		add = span{max(0, min(a.s, b.s)-rng.Intn(3)), max(a.e, b.e) + rng.Intn(3)}
	}
	n := 0
	for _, iv := range ivs {
		if iv.s < add.e && add.s < iv.e {
			n++
		}
	}
	return Case{Input: intervalsInput(ivs) + "\n" + ints(add.s, add.e), Expected: strconv.Itoa(n)}
}

// room-booking: calls new rooms; book start end; cancel room start. book
// gives the lowest room that is free for all of [start, end), or -1. cancel
// removes that room's booking starting at start and says whether it existed.
func genRoomBooking(rng *rand.Rand, size int) Case {
	rooms := 1 + rng.Intn(4)
	reach := []int{12, 60}[rng.Intn(2)]
	calls := []string{"new " + strconv.Itoa(rooms)}
	results := []string{"null"}
	booked := make([][]span, rooms)
	for len(calls) < callCount(size) {
		if rng.Intn(4) == 0 {
			room, start := rng.Intn(rooms), rng.Intn(reach)
			if l := booked[room]; len(l) > 0 && rng.Intn(4) != 0 {
				start = l[rng.Intn(len(l))].s
			}
			found := false
			for i, b := range booked[room] {
				if b.s == start {
					booked[room] = append(booked[room][:i:i], booked[room][i+1:]...)
					found = true
					break
				}
			}
			calls, results = append(calls, "cancel "+ints(room, start)), append(results, boolText(found))
			continue
		}
		s := rng.Intn(reach)
		want := span{s, s + 1 + rng.Intn(8)}
		got := -1
		for room := 0; room < rooms && got < 0; room++ {
			free := true
			for _, b := range booked[room] {
				free = free && (b.e <= want.s || want.e <= b.s)
			}
			if free {
				got = room
				booked[room] = append(booked[room], want)
			}
		}
		calls, results = append(calls, "book "+ints(want.s, want.e)), append(results, strconv.Itoa(got))
	}
	return callCase(calls, results)
}

// job-deps: "n m", then n durations, then m lines "a b": job a must finish
// before job b starts. With any number of workers, the answer is the
// earliest time every job is done, or -1 if the jobs wait on each other in
// a circle.
func genJobDeps(rng *rand.Rand, size int) Case {
	n := max(1, min(size, 10_000))
	durations := randomInts(rng, n, 1, []int{5, 10_000}[rng.Intn(2)])
	m := rng.Intn(2*n + 1)
	order := rng.Perm(n)
	circle := rng.Intn(5) == 0 // any pair, so a circle is likely
	if n == 1 && !circle {
		m = 0
	}
	deps := make([][]int, m)
	for i := range deps {
		a, b := rng.Intn(n), rng.Intn(n)
		if !circle {
			for a == b {
				b = rng.Intn(n)
			}
			a, b = min(a, b), max(a, b)
		}
		deps[i] = []int{order[a], order[b]}
	}
	// Kahn's order, carrying finish times.
	after := make([][]int, n)
	waiting := make([]int, n)
	for _, d := range deps {
		after[d[0]] = append(after[d[0]], d[1])
		waiting[d[1]]++
	}
	start := make([]int, n)
	var ready []int
	for j := 0; j < n; j++ {
		if waiting[j] == 0 {
			ready = append(ready, j)
		}
	}
	done, finish := 0, 0
	for len(ready) > 0 {
		j := ready[0]
		ready = ready[1:]
		done++
		end := start[j] + durations[j]
		finish = max(finish, end)
		for _, next := range after[j] {
			start[next] = max(start[next], end)
			if waiting[next]--; waiting[next] == 0 {
				ready = append(ready, next)
			}
		}
	}
	if done < n {
		finish = -1
	}
	return Case{Input: ints(n, len(deps)) + "\n" + joinInts(durations, " ") + rowsText(deps), Expected: strconv.Itoa(finish)}
}

// parking-lot: calls new small large; park car size; leave car; freeSpots.
// Spots 0 to small-1 are small, the rest large. A small car (size 1) takes
// the lowest free small spot, or the lowest free large one if no small spot
// is free; a large car (size 2) takes the lowest free large spot.
func genParkingLot(rng *rand.Rand, size int) Case {
	small, large := rng.Intn(5), rng.Intn(5)
	if small+large == 0 {
		large = 1
	}
	calls := []string{"new " + ints(small, large)}
	results := []string{"null"}
	taken := make([]bool, small+large)
	spotOf := map[int]int{}
	lowest := func(from, to int) int {
		for s := from; s < to; s++ {
			if !taken[s] {
				return s
			}
		}
		return -1
	}
	cars := 2 + rng.Intn(small+large+3)
	for len(calls) < callCount(size) {
		car := rng.Intn(cars)
		switch k := rng.Intn(10); {
		case k < 5:
			carSize := 1 + rng.Intn(2)
			spot := -1
			if _, parked := spotOf[car]; !parked {
				if carSize == 1 {
					spot = lowest(0, small)
				}
				if spot < 0 {
					spot = lowest(small, small+large)
				}
			}
			if spot >= 0 {
				taken[spot], spotOf[car] = true, spot
			}
			calls, results = append(calls, "park "+ints(car, carSize)), append(results, strconv.Itoa(spot))
		case k < 8:
			spot, parked := spotOf[car]
			if parked {
				taken[spot] = false
				delete(spotOf, car)
			} else {
				spot = -1
			}
			calls, results = append(calls, "leave "+strconv.Itoa(car)), append(results, strconv.Itoa(spot))
		default:
			calls, results = append(calls, "freeSpots"), append(results, strconv.Itoa(small+large-len(spotOf)))
		}
	}
	return callCase(calls, results)
}

// elevator: calls new floors; call floor; next; travelled. The lift starts
// at floor 0 going up. next moves to the nearest called floor in the current
// direction, turning round only when there is none that way, and gives that
// floor (-1 with no calls). travelled is the floors moved so far.
func genElevator(rng *rand.Rand, size int) Case {
	floors := 2 + rng.Intn([]int{6, 30, 1000}[rng.Intn(3)])
	calls := []string{"new " + strconv.Itoa(floors)}
	results := []string{"null"}
	pending := make([]bool, floors)
	at, dir, moved := 0, 1, 0
	seek := func(d int) int {
		for f := at + d; f >= 0 && f < floors; f += d {
			if pending[f] {
				return f
			}
		}
		return -1
	}
	for len(calls) < callCount(size) {
		switch k := rng.Intn(10); {
		case k < 5:
			f := rng.Intn(floors)
			if rng.Intn(8) == 0 {
				f = at
			}
			if f != at {
				pending[f] = true
			}
			calls, results = append(calls, "call "+strconv.Itoa(f)), append(results, "null")
		case k < 9:
			f := seek(dir)
			if f < 0 {
				if f = seek(-dir); f >= 0 {
					dir = -dir
				}
			}
			if f >= 0 {
				moved += max(f-at, at-f)
				at = f
				pending[f] = false
			}
			calls, results = append(calls, "next"), append(results, strconv.Itoa(f))
		default:
			calls, results = append(calls, "travelled"), append(results, strconv.Itoa(moved))
		}
	}
	return callCase(calls, results)
}

// seat-hold: calls new seats ttl; hold user count time; buy user time;
// release user time; freeSeats time. hold takes the lowest block of count
// free seats side by side and gives its first seat, or -1 (also when the
// user already has a hold). A hold made at time t lasts until t + ttl and
// then frees its seats, unless it was bought or released first.
func genSeatHold(rng *rand.Rand, size int) Case {
	seats, ttl := 1+rng.Intn(30), 1+rng.Intn(15)
	users := 2 + rng.Intn(5)
	calls := []string{"new " + ints(seats, ttl)}
	results := []string{"null"}
	const free, held, sold = 0, 1, 2
	state := make([]int, seats)
	type hold struct{ first, count, until int }
	holds := map[int]hold{}
	t := 0
	for len(calls) < callCount(size) {
		t += rng.Intn(4) * rng.Intn(3)
		for u, h := range holds {
			if h.until <= t {
				for s := h.first; s < h.first+h.count; s++ {
					state[s] = free
				}
				delete(holds, u)
			}
		}
		u := rng.Intn(users)
		h, active := holds[u]
		switch k := rng.Intn(10); {
		case k < 4:
			count := 1 + rng.Intn(min(4, seats))
			first := -1
			if !active {
				run := 0
				for s := 0; s < seats && first < 0; s++ {
					if state[s] == free {
						run++
					} else {
						run = 0
					}
					if run == count {
						first = s - count + 1
					}
				}
			}
			if first >= 0 {
				for s := first; s < first+count; s++ {
					state[s] = held
				}
				holds[u] = hold{first, count, t + ttl}
			}
			calls, results = append(calls, "hold "+ints(u, count, t)), append(results, strconv.Itoa(first))
		case k < 8:
			name, to := "buy ", sold
			if k >= 6 {
				name, to = "release ", free
			}
			if active {
				for s := h.first; s < h.first+h.count; s++ {
					state[s] = to
				}
				delete(holds, u)
			}
			calls, results = append(calls, name+ints(u, t)), append(results, boolText(active))
		default:
			n := 0
			for _, s := range state {
				if s == free {
					n++
				}
			}
			calls, results = append(calls, "freeSeats "+strconv.Itoa(t)), append(results, strconv.Itoa(n))
		}
	}
	return callCase(calls, results)
}
