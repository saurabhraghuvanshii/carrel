package testgen

import (
	"math/rand"
	"slices"
	"strconv"
	"strings"
)

// Generators for the real-interview group "Logs and analytics".

func init() {
	register("top-pages", genTopPages)
	register("funnel", genFunnel)
	register("sessions", genSessions)
	register("error-spike", genErrorSpike)
	register("unique-visitors", genUniqueVisitors)
	register("merge-logs", genMergeLogs)
}

// top-pages: "n k", then n page ids. The answer is the k most viewed pages,
// most viewed first, the smaller id first on a tie; fewer if there are not k.
func genTopPages(rng *rand.Rand, size int) Case {
	n := max(1, size)
	pages := 1 + rng.Intn([]int{3, 20, 1000}[rng.Intn(3)])
	k := 1 + rng.Intn(5)
	views := randomInts(rng, n, 1, pages)
	count := map[int]int{}
	for _, v := range views {
		count[v]++
	}
	ids := make([]int, 0, len(count))
	for id := range count {
		ids = append(ids, id)
	}
	slices.SortFunc(ids, func(a, b int) int {
		if count[a] != count[b] {
			return count[b] - count[a]
		}
		return a - b
	})
	return Case{Input: ints(n, k) + "\n" + joinInts(views, " "), Expected: listOutput(ids[:min(k, len(ids))])}
}

// funnel: "n steps", then events "user step" in time order. A user reaches
// step j with an event for j after reaching j - 1; other events are ignored.
// The answer counts the users who reached each step.
func genFunnel(rng *rand.Rand, size int) Case {
	n := max(1, size)
	steps := 2 + rng.Intn(5)
	users := 1 + rng.Intn(max(1, n/3))
	reached := map[int]int{}
	rows := make([][]int, n)
	for i := range rows {
		u := rng.Intn(users)
		step := 1 + rng.Intn(steps)
		if rng.Intn(5) < 3 && reached[u] < steps {
			step = reached[u] + 1
		}
		if step == reached[u]+1 {
			reached[u]++
		}
		rows[i] = []int{u, step}
	}
	counts := make([]int, steps)
	for _, r := range reached {
		for j := 0; j < r; j++ {
			counts[j]++
		}
	}
	return Case{Input: ints(n, steps) + rowsText(rows), Expected: listOutput(counts)}
}

// sessions: "n timeout", then events "user time" in any order. A user's
// events sorted by time split into sessions wherever the gap is more than
// timeout. The answer is [number of sessions, longest session length].
func genSessions(rng *rand.Rand, size int) Case {
	n := max(1, size)
	timeout := []int{0, 5, 30}[rng.Intn(3)]
	users := 1 + rng.Intn(5)
	reach := max(10, n*[]int{1, 5, 40}[rng.Intn(3)])
	rows := make([][]int, n)
	byUser := map[int][]int{}
	for i := range rows {
		rows[i] = []int{rng.Intn(users), rng.Intn(reach)}
		byUser[rows[i][0]] = append(byUser[rows[i][0]], rows[i][1])
	}
	sessions, longest := 0, 0
	for _, times := range byUser {
		slices.Sort(times)
		start := times[0]
		sessions++
		for i := 1; i < len(times); i++ {
			if times[i]-times[i-1] > timeout {
				sessions++
				start = times[i]
			}
			longest = max(longest, times[i]-start)
		}
	}
	return Case{Input: ints(n, timeout) + rowsText(rows), Expected: listOutput([]int{sessions, longest})}
}

// error-spike: "n window", then n error times in order. The answer is the
// most errors in any stretch (t - window, t].
func genErrorSpike(rng *rand.Rand, size int) Case {
	n := max(1, size)
	window := 1 + rng.Intn([]int{1, 10, 100}[rng.Intn(3)])
	times := make([]int, n)
	t := rng.Intn(5)
	for i := range times {
		t += []int{0, 0, 1, rng.Intn(2*window + 1)}[rng.Intn(4)]
		times[i] = t
	}
	best, lo := 0, 0
	for hi := range times {
		for times[lo] <= times[hi]-window {
			lo++
		}
		best = max(best, hi-lo+1)
	}
	return Case{Input: ints(n, window) + "\n" + joinInts(times, " "), Expected: strconv.Itoa(best)}
}

// unique-visitors: calls new window; visit user time; unique time. unique is
// the number of different users with a visit in (time - window, time].
func genUniqueVisitors(rng *rand.Rand, size int) Case {
	window := 1 + rng.Intn(20)
	users := 1 + rng.Intn(8)
	calls := []string{"new " + strconv.Itoa(window)}
	results := []string{"null"}
	last := map[int]int{}
	t := rng.Intn(5)
	for len(calls) < callCount(size) {
		t += rng.Intn(4) * rng.Intn(3)
		if rng.Intn(5) < 3 {
			u := rng.Intn(users)
			last[u] = t
			calls, results = append(calls, "visit "+ints(u, t)), append(results, "null")
			continue
		}
		n := 0
		for _, at := range last {
			if at > t-window {
				n++
			}
		}
		calls, results = append(calls, "unique "+strconv.Itoa(t)), append(results, strconv.Itoa(n))
	}
	return callCase(calls, results)
}

// merge-logs: "k limit", then k clock offsets, then for every server "m" and
// its m times in order. An entry's true time is its time plus its server's
// offset. The answer is the server of each of the first limit entries by
// true time, the lower server first on a tie.
func genMergeLogs(rng *rand.Rand, size int) Case {
	k := 1 + rng.Intn(min(8, max(1, size)))
	total := max(1, min(size, 5000))
	offsets := make([]int, k)
	if rng.Intn(3) != 0 {
		offsets = randomInts(rng, k, -50, 50)
	}
	sizes := make([]int, k)
	for i := 0; i < total; i++ {
		sizes[rng.Intn(k)]++
	}
	top := []int{5, 100, 1_000_000}[rng.Intn(3)]
	logs := make([][]int, k)
	var b strings.Builder
	for s := range logs {
		logs[s] = randomInts(rng, sizes[s], 0, top)
		slices.Sort(logs[s])
		b.WriteString("\n" + arrayInput(logs[s]))
	}
	limit := total
	if rng.Intn(2) == 0 {
		limit = 1 + rng.Intn(total)
	}
	// Take the smallest head each time; the scan finds the lowest server first.
	heads := make([]int, k)
	order := make([]int, 0, limit)
	for len(order) < limit {
		best := -1
		for s := range logs {
			if heads[s] < len(logs[s]) && (best < 0 || logs[s][heads[s]]+offsets[s] < logs[best][heads[best]]+offsets[best]) {
				best = s
			}
		}
		heads[best]++
		order = append(order, best)
	}
	return Case{Input: ints(k, limit) + "\n" + joinInts(offsets, " ") + b.String(), Expected: listOutput(order)}
}
