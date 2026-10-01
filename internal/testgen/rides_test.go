package testgen

import (
	"fmt"
	"math"
	"math/bits"
	"strconv"
	"strings"
	"testing"
)

func intsOf(line string) []int {
	var out []int
	for _, f := range strings.Fields(line) {
		v, _ := strconv.Atoi(f)
		out = append(out, v)
	}
	return out
}

func TestSurgeWindow(t *testing.T) {
	for i, c := range cases(t, "surge-window", 60) {
		lines := strings.Split(c.Input, "\n")
		nk := intsOf(lines[0])
		req, drv := intsOf(lines[1]), intsOf(lines[2])
		if nk[0] > 3000 {
			continue
		}
		var out []int
		for end := nk[1] - 1; end < nk[0]; end++ {
			r, d := 0, 0
			for j := end - nk[1] + 1; j <= end; j++ {
				r, d = r+req[j], d+drv[j]
			}
			mult := 1
			if r > d {
				mult = 5
				if d > 0 {
					mult = min(5, int(math.Ceil(float64(r)/float64(d))))
				}
			}
			out = append(out, mult)
		}
		if listOutput(out) != c.Expected {
			t.Fatalf("case %d differs", i)
		}
	}
}

func TestShiftGaps(t *testing.T) {
	gaps := 0
	for i, c := range cases(t, "shift-gaps", 80) {
		lines := strings.Split(c.Input, "\n")
		head := intsOf(lines[0])
		n, m, dayEnd := head[0], head[1], head[2]
		if dayEnd > 1000 {
			continue
		}
		// Count drivers on each unit [u, u+1) of the day.
		on := make([]int, dayEnd)
		for _, l := range lines[1 : n+1] {
			s := intsOf(l)
			for u := s[0]; u < s[1]; u++ {
				on[u]++
			}
		}
		var free []span
		for u := 0; u < dayEnd; u++ {
			if on[u] < m {
				if len(free) > 0 && free[len(free)-1].e == u {
					free[len(free)-1].e++
				} else {
					free = append(free, span{u, u + 1})
				}
			}
		}
		if intervalsOutput(free) != c.Expected {
			t.Fatalf("case %d: want %s, generator says %s", i, intervalsOutput(free), c.Expected)
		}
		gaps += len(free)
	}
	if gaps == 0 {
		t.Fatal("no case has a gap")
	}
}

func TestPickupBatches(t *testing.T) {
	negative := false
	for i, c := range cases(t, "pickup-batches", 60) {
		lines := strings.Split(c.Input, "\n")
		head := intsOf(lines[0])
		limit, window := head[1], head[2]
		// Keep every batch per zone as a list of times.
		batches := map[string][][]int{}
		total := 0
		for _, l := range lines[1:] {
			r := intsOf(l)
			zone := fmt.Sprint(math.Floor(float64(r[1])/100), math.Floor(float64(r[2])/100))
			list := batches[zone]
			if k := len(list) - 1; k >= 0 && len(list[k]) < limit && r[0]-list[k][0] <= window {
				list[k] = append(list[k], r[0])
			} else {
				batches[zone] = append(list, []int{r[0]})
				total++
			}
			negative = negative || r[1] < 0
		}
		if strconv.Itoa(total) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, total, c.Expected)
		}
	}
	if !negative {
		t.Fatal("no negative coordinates")
	}
}

func TestEtaLights(t *testing.T) {
	found, none := 0, 0
	for i, c := range cases(t, "eta-lights", 60) {
		lines := strings.Split(c.Input, "\n")
		nm := intsOf(lines[0])
		n := nm[0]
		if n > 300 {
			continue
		}
		periods := intsOf(lines[1])
		// Relax every road until nothing improves.
		dist := make([]int, n)
		for k := range dist {
			dist[k] = math.MaxInt
		}
		dist[0] = 0
		for changed := true; changed; {
			changed = false
			for _, l := range lines[2:] {
				r := intsOf(l)
				if dist[r[0]] == math.MaxInt {
					continue
				}
				p := periods[r[0]]
				leave := dist[r[0]]
				for leave%p != 0 {
					leave++
				}
				if leave+r[2] < dist[r[1]] {
					dist[r[1]] = leave + r[2]
					changed = true
				}
			}
		}
		want := dist[n-1]
		if want == math.MaxInt {
			want = -1
			none++
		} else {
			found++
		}
		if strconv.Itoa(want) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, want, c.Expected)
		}
	}
	if found == 0 || none == 0 {
		t.Fatal("need both reachable and unreachable cases")
	}
}

func TestPoolMatch(t *testing.T) {
	for i, c := range cases(t, "pool-match", 60) {
		lines := strings.Split(c.Input, "\n")
		head := intsOf(lines[0])
		stops, seats, n := head[0], head[1], head[2]
		if n > 12 {
			continue
		}
		rides := make([][]int, n)
		for k := range rides {
			rides[k] = intsOf(lines[k+1])
		}
		best := 0
		for mask := 0; mask < 1<<n; mask++ {
			load := make([]int, stops)
			ok := true
			for k, r := range rides {
				if mask>>k&1 == 1 {
					for x := r[0]; x < r[1]; x++ {
						load[x]++
						ok = ok && load[x] <= seats
					}
				}
			}
			if ok {
				best = max(best, bits.OnesCount(uint(mask)))
			}
		}
		if strconv.Itoa(best) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, best, c.Expected)
		}
	}
}
