package testgen

import (
	"slices"
	"strconv"
	"strings"
	"testing"
)

func TestTopPages(t *testing.T) {
	for i, c := range cases(t, "top-pages", 60) {
		lines := strings.Split(c.Input, "\n")
		k := intsOf(lines[0])[1]
		views := intsOf(lines[1])
		if len(views) > 400 {
			continue
		}
		// Pick the best remaining page k times.
		taken := map[int]bool{}
		var got []int
		for len(got) < k {
			best, bestCount := 0, 0
			for _, p := range views {
				if taken[p] {
					continue
				}
				n := 0
				for _, q := range views {
					if q == p {
						n++
					}
				}
				if n > bestCount || (n == bestCount && p < best) {
					best, bestCount = p, n
				}
			}
			if bestCount == 0 {
				break
			}
			taken[best] = true
			got = append(got, best)
		}
		if listOutput(got) != c.Expected {
			t.Fatalf("case %d: want %v, generator says %s", i, got, c.Expected)
		}
	}
}

func TestFunnel(t *testing.T) {
	for i, c := range cases(t, "funnel", 60) {
		lines := strings.Split(c.Input, "\n")
		head := intsOf(lines[0])
		rows := rowsOf(lines, 1, head[0])
		byUser := map[int][]int{}
		for _, r := range rows {
			if r[1] < 1 || r[1] > head[1] {
				t.Fatalf("case %d: step %d out of range", i, r[1])
			}
			byUser[r[0]] = append(byUser[r[0]], r[1])
		}
		// A user reached step j when 1, 2, ..., j appears in order in
		// their events.
		counts := make([]int, head[1])
		for _, evs := range byUser {
			for j := 1; j <= head[1]; j++ {
				want := 1
				for _, s := range evs {
					if s == want && want <= j {
						want++
					}
				}
				if want > j {
					counts[j-1]++
				}
			}
		}
		if listOutput(counts) != c.Expected {
			t.Fatalf("case %d: want %v, generator says %s", i, counts, c.Expected)
		}
	}
}

func TestSessions(t *testing.T) {
	for i, c := range cases(t, "sessions", 60) {
		lines := strings.Split(c.Input, "\n")
		head := intsOf(lines[0])
		if head[0] > 300 {
			continue
		}
		// Each event covers [time, time + timeout]; events of one user that
		// touch belong to one session.
		byUser := map[int][]span{}
		for _, r := range rowsOf(lines, 1, head[0]) {
			byUser[r[0]] = append(byUser[r[0]], span{r[1], r[1] + head[1]})
		}
		sessions, longest := 0, 0
		for _, spans := range byUser {
			for _, m := range mergeRepeatedly(spans) {
				sessions++
				longest = max(longest, m.e-head[1]-m.s)
			}
		}
		if listOutput([]int{sessions, longest}) != c.Expected {
			t.Fatalf("case %d: want [%d, %d], generator says %s", i, sessions, longest, c.Expected)
		}
	}
}

func TestErrorSpike(t *testing.T) {
	for i, c := range cases(t, "error-spike", 60) {
		lines := strings.Split(c.Input, "\n")
		window := intsOf(lines[0])[1]
		times := intsOf(lines[1])
		if !slices.IsSorted(times) {
			t.Fatalf("case %d: times are not in order", i)
		}
		if len(times) > 2000 {
			continue
		}
		best := 0
		for _, end := range times {
			n := 0
			for _, x := range times {
				if end-window < x && x <= end {
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

func TestUniqueVisitors(t *testing.T) {
	replay(t, "unique-visitors", 50, func(a []int) func(string, []int) string {
		window := a[0]
		var visits [][2]int
		return func(call string, x []int) string {
			if call == "visit" {
				visits = append(visits, [2]int{x[0], x[1]})
				return "null"
			}
			users := map[int]bool{}
			for _, v := range visits {
				if x[0]-window < v[1] && v[1] <= x[0] {
					users[v[0]] = true
				}
			}
			return strconv.Itoa(len(users))
		}
	})
}

func TestMergeLogs(t *testing.T) {
	ties := false
	for i, c := range cases(t, "merge-logs", 60) {
		lines := strings.Split(c.Input, "\n")
		head := intsOf(lines[0])
		offsets := intsOf(lines[1])
		type entry struct{ at, server int }
		var all []entry
		for s := 0; s < head[0]; s++ {
			times := intsOf(lines[3+2*s])
			if n, _ := strconv.Atoi(lines[2+2*s]); n != len(times) || !slices.IsSorted(times) {
				t.Fatalf("case %d: server %d is not a sorted list of %d", i, s, n)
			}
			for _, x := range times {
				all = append(all, entry{x + offsets[s], s})
			}
		}
		slices.SortStableFunc(all, func(a, b entry) int {
			if a.at != b.at {
				return a.at - b.at
			}
			return a.server - b.server
		})
		got := make([]int, head[1])
		for j := range got {
			got[j] = all[j].server
			ties = ties || (j > 0 && all[j].at == all[j-1].at && all[j].server != all[j-1].server)
		}
		if listOutput(got) != c.Expected {
			t.Fatalf("case %d differs", i)
		}
	}
	if !ties {
		t.Fatal("no ties between servers")
	}
}
