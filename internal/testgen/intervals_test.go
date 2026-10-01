package testgen

import (
	"math/bits"
	"strconv"
	"strings"
	"testing"
)

func parseSpans(t *testing.T, lines []string) ([]span, []string) {
	t.Helper()
	n, _ := strconv.Atoi(lines[0])
	ivs := make([]span, n)
	for i := range ivs {
		f := strings.Fields(lines[i+1])
		ivs[i].s, _ = strconv.Atoi(f[0])
		ivs[i].e, _ = strconv.Atoi(f[1])
		if ivs[i].s > ivs[i].e {
			t.Fatalf("start after end: %v", ivs[i])
		}
	}
	return ivs, lines[n+1:]
}

func TestMeetingConflict(t *testing.T) {
	seen := map[string]bool{}
	for i, c := range cases(t, "meeting-conflict", 60) {
		ivs, _ := parseSpans(t, strings.Split(c.Input, "\n"))
		if len(ivs) > 1500 {
			continue
		}
		ok := true
		for a := range ivs {
			for b := a + 1; b < len(ivs); b++ {
				if ivs[a].s < ivs[b].e && ivs[b].s < ivs[a].e {
					ok = false
				}
			}
		}
		if boolText(ok) != c.Expected {
			t.Fatalf("case %d: want %v, generator says %s", i, ok, c.Expected)
		}
		seen[c.Expected] = true
	}
	if !seen["true"] || !seen["false"] {
		t.Fatal("both answers must appear")
	}
}

// mergeRepeatedly joins any two closed intervals that share a point, until
// none do.
func mergeRepeatedly(ivs []span) []span {
	out := append([]span{}, ivs...)
	for changed := true; changed; {
		changed = false
		for a := 0; a < len(out) && !changed; a++ {
			for b := a + 1; b < len(out); b++ {
				if out[a].s <= out[b].e && out[b].s <= out[a].e {
					out[a] = span{min(out[a].s, out[b].s), max(out[a].e, out[b].e)}
					out = append(out[:b], out[b+1:]...)
					changed = true
					break
				}
			}
		}
	}
	return byStart(out)
}

func TestMergeAndInsert(t *testing.T) {
	for i, c := range cases(t, "merge-intervals", 50) {
		ivs, _ := parseSpans(t, strings.Split(c.Input, "\n"))
		if len(ivs) > 300 {
			continue
		}
		if intervalsOutput(mergeRepeatedly(ivs)) != c.Expected {
			t.Fatalf("merge case %d differs", i)
		}
	}
	for i, c := range cases(t, "insert-interval", 50) {
		ivs, rest := parseSpans(t, strings.Split(c.Input, "\n"))
		for j := 1; j < len(ivs); j++ {
			if ivs[j].s <= ivs[j-1].e {
				t.Fatalf("insert case %d: the list is not already merged", i)
			}
		}
		f := strings.Fields(rest[0])
		s, _ := strconv.Atoi(f[0])
		e, _ := strconv.Atoi(f[1])
		if len(ivs) > 300 {
			continue
		}
		if intervalsOutput(mergeRepeatedly(append(ivs, span{s, e}))) != c.Expected {
			t.Fatalf("insert case %d differs", i)
		}
	}
}

func TestRooms(t *testing.T) {
	for i, c := range cases(t, "rooms", 50) {
		ivs, _ := parseSpans(t, strings.Split(c.Input, "\n"))
		if len(ivs) > 1500 {
			continue
		}
		best := 0
		for _, at := range ivs {
			n := 0
			for _, iv := range ivs {
				if iv.s <= at.s && at.s < iv.e {
					n++
				}
			}
			best = max(best, n)
		}
		if strconv.Itoa(best) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, best, c.Expected)
		}
	}
}

// TestRemoveOverlaps and TestArrows try every subset on small cases.
func TestRemoveOverlaps(t *testing.T) {
	for i, c := range cases(t, "remove-overlaps", 50) {
		ivs, _ := parseSpans(t, strings.Split(c.Input, "\n"))
		if len(ivs) > 12 {
			continue
		}
		best := 0
		for mask := 0; mask < 1<<len(ivs); mask++ {
			ok := true
			for a := range ivs {
				for b := a + 1; b < len(ivs) && ok; b++ {
					if mask>>a&1 == 1 && mask>>b&1 == 1 && ivs[a].s < ivs[b].e && ivs[b].s < ivs[a].e {
						ok = false
					}
				}
			}
			if ok {
				best = max(best, bits.OnesCount(uint(mask)))
			}
		}
		if strconv.Itoa(len(ivs)-best) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, len(ivs)-best, c.Expected)
		}
	}
}

func TestArrows(t *testing.T) {
	for i, c := range cases(t, "arrows", 50) {
		ivs, _ := parseSpans(t, strings.Split(c.Input, "\n"))
		if len(ivs) > 10 {
			continue
		}
		// Some best set of arrows can always sit on balloon ends.
		var ends []int
		for _, iv := range ivs {
			ends = append(ends, iv.e)
		}
		best := len(ivs)
		for mask := 1; mask < 1<<len(ends); mask++ {
			all := true
			for _, iv := range ivs {
				hit := false
				for j, x := range ends {
					hit = hit || (mask>>j&1 == 1 && iv.s <= x && x <= iv.e)
				}
				all = all && hit
			}
			if all {
				best = min(best, bits.OnesCount(uint(mask)))
			}
		}
		if strconv.Itoa(best) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, best, c.Expected)
		}
	}
}

func TestFreeTime(t *testing.T) {
	gaps := 0
	for i, c := range cases(t, "free-time", 40) {
		lines := strings.Split(c.Input, "\n")
		k, _ := strconv.Atoi(lines[0])
		rest := lines[1:]
		var all []span
		lo, hi := 1<<62, 0
		for j := 0; j < k; j++ {
			ivs, after := parseSpans(t, rest)
			for x := 1; x < len(ivs); x++ {
				if ivs[x].s <= ivs[x-1].e {
					t.Fatalf("case %d: a person's list is not separate and sorted", i)
				}
			}
			for _, iv := range ivs {
				lo, hi = min(lo, iv.s), max(hi, iv.e)
			}
			all = append(all, ivs...)
			rest = after
		}
		if hi-lo > 200_000 {
			continue
		}
		// Mark each unit [u, u+1) busy if some interval covers it.
		busy := make([]bool, hi-lo)
		for _, iv := range all {
			for u := iv.s; u < iv.e; u++ {
				busy[u-lo] = true
			}
		}
		var free []span
		for u := 0; u < len(busy); u++ {
			if !busy[u] {
				start := u
				for u < len(busy) && !busy[u] {
					u++
				}
				free = append(free, span{start + lo, u + lo})
			}
		}
		if intervalsOutput(free) != c.Expected {
			t.Fatalf("case %d: free time differs", i)
		}
		gaps += len(free)
	}
	if gaps == 0 {
		t.Fatal("no case has free time")
	}
}
