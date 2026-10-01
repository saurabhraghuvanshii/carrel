package testgen

import (
	"container/heap"
	"math/rand"
	"strconv"
	"strings"
)

// Generators for the real-interview group "Caches and rate limits". Each
// case is a list of calls, "new ..." first; the answer has one result per
// call, null for calls that return nothing. Times never go down.

func init() {
	register("fixed-window", genFixedWindow)
	register("sliding-limiter", genSlidingLimiter)
	register("token-bucket", genTokenBucket)
	register("ttl-cache", genTTLCache)
	register("lfu-cache", genLFUCache)
}

func callCase(calls, results []string) Case {
	return Case{Input: strconv.Itoa(len(calls)) + "\n" + strings.Join(calls, "\n"), Expected: strings.Join(results, " ")}
}

func callCount(size int) int { return min(max(size, 4), 5000) }

func ints(vals ...int) string { return joinInts(vals, " ") }

// fixed-window: new limit window; allow user time. Time is cut into windows
// [k*window, (k+1)*window); each user gets at most limit allowed requests
// per window. Refused requests do not count.
func genFixedWindow(rng *rand.Rand, size int) Case {
	limit, window := 1+rng.Intn(5), 1+rng.Intn(20)
	users := 1 + rng.Intn(4)
	calls := []string{"new " + ints(limit, window)}
	results := []string{"null"}
	used := map[[2]int]int{}
	t := rng.Intn(5)
	for len(calls) < callCount(size) {
		t += rng.Intn(4) * rng.Intn(2)
		u := rng.Intn(users)
		key := [2]int{u, t / window}
		ok := used[key] < limit
		if ok {
			used[key]++
		}
		calls = append(calls, "allow "+ints(u, t))
		results = append(results, boolText(ok))
	}
	return callCase(calls, results)
}

// sliding-limiter: new limit window; allow user time. A request is allowed
// when the user had fewer than limit allowed requests in (time - window,
// time].
func genSlidingLimiter(rng *rand.Rand, size int) Case {
	limit, window := 1+rng.Intn(5), 1+rng.Intn(20)
	users := 1 + rng.Intn(4)
	calls := []string{"new " + ints(limit, window)}
	results := []string{"null"}
	recent := map[int][]int{}
	t := rng.Intn(5)
	for len(calls) < callCount(size) {
		t += rng.Intn(4) * rng.Intn(2)
		u := rng.Intn(users)
		q := recent[u]
		for len(q) > 0 && q[0] <= t-window {
			q = q[1:]
		}
		ok := len(q) < limit
		if ok {
			q = append(q, t)
		}
		recent[u] = q
		calls = append(calls, "allow "+ints(u, t))
		results = append(results, boolText(ok))
	}
	return callCase(calls, results)
}

// token-bucket: new capacity refill; take time n. The bucket starts full and
// gains one token at every positive multiple of refill, up to capacity. take
// removes n tokens if there are that many and answers whether it did.
func genTokenBucket(rng *rand.Rand, size int) Case {
	capacity, refill := 1+rng.Intn(10), 1+rng.Intn(10)
	calls := []string{"new " + ints(capacity, refill)}
	results := []string{"null"}
	tokens, last, t := capacity, 0, 0
	for len(calls) < callCount(size) {
		t += rng.Intn(2 * refill)
		n := 1 + rng.Intn(capacity+1)
		if rng.Intn(3) != 0 {
			n = 1 + rng.Intn(2)
		}
		tokens = min(capacity, tokens+t/refill-last/refill)
		last = t
		ok := tokens >= n
		if ok {
			tokens -= n
		}
		calls = append(calls, "take "+ints(t, n))
		results = append(results, boolText(ok))
	}
	return callCase(calls, results)
}

type expiry struct{ at, key int }
type expiryHeap []expiry

func (h expiryHeap) Len() int           { return len(h) }
func (h expiryHeap) Less(i, j int) bool { return h[i].at < h[j].at }
func (h expiryHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *expiryHeap) Push(x any)        { *h = append(*h, x.(expiry)) }
func (h *expiryHeap) Pop() any {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

// ttl-cache: new; put key value ttl time; get key time; count time. A value
// put at time t with ttl d can be read at times before t + d. count is the
// number of keys that can be read at that time.
func genTTLCache(rng *rand.Rand, size int) Case {
	keys := 1 + rng.Intn(20)
	calls := []string{"new"}
	results := []string{"null"}
	type entry struct{ value, until int }
	live := map[int]entry{}
	h := &expiryHeap{}
	t := 0
	for len(calls) < callCount(size) {
		t += rng.Intn(5) * rng.Intn(2)
		for h.Len() > 0 && (*h)[0].at <= t {
			e := heap.Pop(h).(expiry)
			if cur, ok := live[e.key]; ok && cur.until == e.at {
				delete(live, e.key)
			}
		}
		key := rng.Intn(keys)
		switch rng.Intn(5) {
		case 0, 1:
			value, ttl := rng.Intn(1000), 1+rng.Intn(30)
			live[key] = entry{value, t + ttl}
			heap.Push(h, expiry{t + ttl, key})
			calls = append(calls, "put "+ints(key, value, ttl, t))
			results = append(results, "null")
		case 2, 3:
			calls = append(calls, "get "+ints(key, t))
			if e, ok := live[key]; ok {
				results = append(results, strconv.Itoa(e.value))
			} else {
				results = append(results, "-1")
			}
		default:
			calls = append(calls, "count "+strconv.Itoa(t))
			results = append(results, strconv.Itoa(len(live)))
		}
	}
	return callCase(calls, results)
}

// lfu-cache: new capacity; get key; put key value. When full, put removes
// the key used the fewest times, and of those the one used longest ago.
// get and put both count as a use.
func genLFUCache(rng *rand.Rand, size int) Case {
	n := callCount(size)
	capacity := 1 + rng.Intn(max(1, min(n/8, 30)))
	keys := 2*capacity + 2
	calls := []string{"new " + strconv.Itoa(capacity)}
	results := []string{"null"}
	type entry struct{ value, uses int }
	cache := map[int]*entry{}
	// byUse[f] lists keys used f times, oldest use first.
	byUse := map[int][]int{}
	lowest := 0
	remove := func(f, key int) {
		l := byUse[f]
		for i, k := range l {
			if k == key {
				byUse[f] = append(l[:i:i], l[i+1:]...)
				break
			}
		}
	}
	touch := func(key int) {
		e := cache[key]
		remove(e.uses, key)
		if e.uses == lowest && len(byUse[e.uses]) == 0 {
			lowest++
		}
		e.uses++
		byUse[e.uses] = append(byUse[e.uses], key)
	}
	for len(calls) < n {
		key := rng.Intn(keys)
		if rng.Intn(2) == 0 {
			calls = append(calls, "get "+strconv.Itoa(key))
			if e, ok := cache[key]; ok {
				touch(key)
				results = append(results, strconv.Itoa(e.value))
			} else {
				results = append(results, "-1")
			}
			continue
		}
		value := rng.Intn(100_000)
		calls = append(calls, "put "+ints(key, value))
		results = append(results, "null")
		if e, ok := cache[key]; ok {
			e.value = value
			touch(key)
			continue
		}
		if len(cache) == capacity {
			old := byUse[lowest][0]
			remove(lowest, old)
			delete(cache, old)
		}
		cache[key] = &entry{value: value, uses: 1}
		byUse[1] = append(byUse[1], key)
		lowest = 1
	}
	return callCase(calls, results)
}
