package testgen

import (
	"strconv"
	"strings"
	"testing"
)

func TestNearestStore(t *testing.T) {
	for i, c := range cases(t, "nearest-store", 60) {
		lines := strings.Split(c.Input, "\n")
		stores, customers := intsOf(lines[1]), intsOf(lines[2])
		if len(stores)*len(customers) > 4_000_000 {
			continue
		}
		out := make([]int, len(customers))
		for k, at := range customers {
			out[k] = 1 << 62
			for _, s := range stores {
				out[k] = min(out[k], max(s-at, at-s))
			}
		}
		if listOutput(out) != c.Expected {
			t.Fatalf("case %d differs", i)
		}
	}
}

func TestGeofence(t *testing.T) {
	seen := map[string]bool{}
	for i, c := range cases(t, "geofence", 60) {
		lines := strings.Split(c.Input, "\n")
		z := intsOf(lines[0])
		n, _ := strconv.Atoi(lines[1])
		// Write the path as O and I letters and count "OI".
		trail := []byte{'O'}
		for _, p := range rowsOf(lines, 2, n) {
			ch := byte('O')
			if p[0] >= z[0] && p[0] <= z[2] && p[1] >= z[1] && p[1] <= z[3] {
				ch = 'I'
			}
			trail = append(trail, ch)
		}
		if want := strconv.Itoa(strings.Count(string(trail), "OI")); want != c.Expected {
			t.Fatalf("case %d: want %s, generator says %s", i, want, c.Expected)
		}
		seen[c.Expected] = true
	}
	if len(seen) < 4 {
		t.Fatal("too little variety in the answers")
	}
}

func TestStreetWalk(t *testing.T) {
	for i, c := range cases(t, "street-walk", 80) {
		lines := strings.Split(c.Input, "\n")
		start := intsOf(lines[0])[1]
		stops := intsOf(lines[1])
		if len(stops) > 7 {
			continue
		}
		best := -1
		permutations(len(stops), func(p []int) {
			at, walked := start, 0
			for _, k := range p {
				walked += max(stops[k]-at, at-stops[k])
				at = stops[k]
			}
			if best < 0 || walked < best {
				best = walked
			}
		})
		if strconv.Itoa(best) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, best, c.Expected)
		}
	}
}

func TestZoneLookup(t *testing.T) {
	hits, misses := 0, 0
	for i, c := range cases(t, "zone-lookup", 60) {
		lines := strings.Split(c.Input, "\n")
		head := intsOf(lines[0])
		rows := rowsOf(lines, 1, head[0])
		codes := intsOf(lines[1+head[0]])
		if len(rows)*len(codes) > 4_000_000 {
			continue
		}
		for a := range rows {
			for b := a + 1; b < len(rows) && len(rows) <= 300; b++ {
				if rows[a][0] <= rows[b][1] && rows[b][0] <= rows[a][1] {
					t.Fatalf("case %d: ranges overlap", i)
				}
			}
		}
		out := make([]int, len(codes))
		for k, code := range codes {
			out[k] = -1
			for _, r := range rows {
				if r[0] <= code && code <= r[1] {
					out[k] = r[2]
					hits++
				}
			}
			if out[k] < 0 {
				misses++
			}
		}
		if listOutput(out) != c.Expected {
			t.Fatalf("case %d differs", i)
		}
	}
	if hits == 0 || misses == 0 {
		t.Fatal("need codes inside and outside the ranges")
	}
}

func TestCourierAssignment(t *testing.T) {
	spare := false
	for i, c := range cases(t, "courier-assignment", 80) {
		lines := strings.Split(c.Input, "\n")
		orders, couriers := intsOf(lines[1]), intsOf(lines[2])
		spare = spare || len(couriers) > len(orders)
		if len(couriers) > 7 {
			continue
		}
		// Try every ordering of the couriers; order k gets the k-th one.
		best := -1
		permutations(len(couriers), func(p []int) {
			total := 0
			for k, o := range orders {
				total += max(o-couriers[p[k]], couriers[p[k]]-o)
			}
			if best < 0 || total < best {
				best = total
			}
		})
		if strconv.Itoa(best) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, best, c.Expected)
		}
	}
	if !spare {
		t.Fatal("no case has spare couriers")
	}
}
