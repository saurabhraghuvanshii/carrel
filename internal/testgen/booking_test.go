package testgen

import (
	"strconv"
	"strings"
	"testing"
)

func TestCalendarConflict(t *testing.T) {
	for i, c := range cases(t, "calendar-conflict", 60) {
		ivs, rest := parseSpans(t, strings.Split(c.Input, "\n"))
		add := intsOf(rest[0])
		n := 0
		for _, iv := range ivs {
			// Two meetings overlap when the later start is before the
			// earlier end.
			if max(iv.s, add[0]) < min(iv.e, add[1]) {
				n++
			}
		}
		if strconv.Itoa(n) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, n, c.Expected)
		}
	}
}

func TestRoomBooking(t *testing.T) {
	replay(t, "room-booking", 50, func(a []int) func(string, []int) string {
		// busy[room][u] marks the unit [u, u+1); ends remembers bookings.
		busy := make([]map[int]bool, a[0])
		ends := make([]map[int]int, a[0])
		for r := range busy {
			busy[r], ends[r] = map[int]bool{}, map[int]int{}
		}
		return func(call string, x []int) string {
			if call == "cancel" {
				end, ok := ends[x[0]][x[1]]
				for u := x[1]; u < end; u++ {
					delete(busy[x[0]], u)
				}
				delete(ends[x[0]], x[1])
				return boolText(ok)
			}
			for r := range busy {
				free := true
				for u := x[0]; u < x[1]; u++ {
					free = free && !busy[r][u]
				}
				if free {
					for u := x[0]; u < x[1]; u++ {
						busy[r][u] = true
					}
					ends[r][x[0]] = x[1]
					return strconv.Itoa(r)
				}
			}
			return "-1"
		}
	})
}

func TestJobDeps(t *testing.T) {
	circles, plain := 0, 0
	for i, c := range cases(t, "job-deps", 80) {
		lines := strings.Split(c.Input, "\n")
		head := intsOf(lines[0])
		n := head[0]
		if n > 300 {
			continue
		}
		durations := intsOf(lines[1])
		deps := rowsOf(lines, 2, head[1])
		// Push finish times along the links until nothing changes. With no
		// circle that takes at most n rounds.
		finish := append([]int{}, durations...)
		changed := true
		for round := 0; round <= n && changed; round++ {
			changed = false
			for _, d := range deps {
				if f := finish[d[0]] + durations[d[1]]; f > finish[d[1]] {
					finish[d[1]], changed = f, true
				}
			}
		}
		want := -1
		if !changed {
			for _, f := range finish {
				want = max(want, f)
			}
			plain++
		} else {
			circles++
		}
		if strconv.Itoa(want) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, want, c.Expected)
		}
	}
	if circles == 0 || plain == 0 {
		t.Fatal("need cases with and without a circle")
	}
}

func TestParkingLot(t *testing.T) {
	replay(t, "parking-lot", 50, func(a []int) func(string, []int) string {
		small := a[0]
		spots := make([]int, a[0]+a[1]) // car + 1, or 0 when free
		find := func(car int) int {
			for s, c := range spots {
				if c == car+1 {
					return s
				}
			}
			return -1
		}
		return func(call string, x []int) string {
			switch call {
			case "park":
				if find(x[0]) >= 0 {
					return "-1"
				}
				for pass := 0; pass < 2; pass++ {
					for s := range spots {
						isSmall := s < small
						if spots[s] == 0 && ((pass == 0 && isSmall && x[1] == 1) || (pass == 1 && !isSmall)) {
							spots[s] = x[0] + 1
							return strconv.Itoa(s)
						}
					}
				}
				return "-1"
			case "leave":
				s := find(x[0])
				if s >= 0 {
					spots[s] = 0
				}
				return strconv.Itoa(s)
			}
			n := 0
			for _, c := range spots {
				if c == 0 {
					n++
				}
			}
			return strconv.Itoa(n)
		}
	})
}

func TestElevator(t *testing.T) {
	replay(t, "elevator", 50, func([]int) func(string, []int) string {
		pending := map[int]bool{}
		at, dir, moved := 0, 1, 0
		// nearest gives the closest pending floor on side d, or -1.
		nearest := func(d int) int {
			best := -1
			for f := range pending {
				if (f-at)*d > 0 && (best < 0 || (f-at)*d < (best-at)*d) {
					best = f
				}
			}
			return best
		}
		return func(call string, x []int) string {
			switch call {
			case "call":
				if x[0] != at {
					pending[x[0]] = true
				}
				return "null"
			case "travelled":
				return strconv.Itoa(moved)
			}
			f := nearest(dir)
			if f < 0 {
				if f = nearest(-dir); f < 0 {
					return "-1"
				}
				dir = -dir
			}
			moved += (f - at) * dir
			at = f
			delete(pending, f)
			return strconv.Itoa(f)
		}
	})
}

func TestSeatHold(t *testing.T) {
	replay(t, "seat-hold", 50, func(a []int) func(string, []int) string {
		ttl := a[1]
		type seat struct{ user, until int } // user -1 free, -2 sold
		seats := make([]seat, a[0])
		for s := range seats {
			seats[s].user = -1
		}
		return func(call string, x []int) string {
			now := x[len(x)-1]
			for s := range seats {
				if seats[s].user >= 0 && seats[s].until <= now {
					seats[s].user = -1
				}
			}
			count := func(user int) int {
				n := 0
				for _, s := range seats {
					if s.user == user {
						n++
					}
				}
				return n
			}
			switch call {
			case "freeSeats":
				return strconv.Itoa(count(-1))
			case "buy", "release":
				if count(x[0]) == 0 {
					return "false"
				}
				for s := range seats {
					if seats[s].user == x[0] {
						seats[s].user = map[string]int{"buy": -2, "release": -1}[call]
					}
				}
				return "true"
			}
			if count(x[0]) > 0 {
				return "-1"
			}
			for first := 0; first+x[1] <= len(seats); first++ {
				ok := true
				for s := first; s < first+x[1]; s++ {
					ok = ok && seats[s].user == -1
				}
				if ok {
					for s := first; s < first+x[1]; s++ {
						seats[s] = seat{x[0], now + ttl}
					}
					return strconv.Itoa(first)
				}
			}
			return "-1"
		}
	})
}
