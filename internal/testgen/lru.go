package testgen

import (
	"container/list"
	"math/rand"
	"strconv"
	"strings"
)

func init() { register("lru-cache", genLRU) }

// lru-cache: a call sequence. The first call is "new capacity", then "get key"
// and "put key value". Keys come from a small range so evictions happen.
// Expected is one result per call: the value or -1 for get, null otherwise.
func genLRU(rng *rand.Rand, size int) Case {
	k := min(max(size, 5), 5000)
	capacity := 1 + rng.Intn(max(1, min(k/8, 50)))
	keys := 2*capacity + 2

	calls := []string{"new " + strconv.Itoa(capacity)}
	for len(calls) < k {
		key := rng.Intn(keys)
		if rng.Intn(2) == 0 {
			calls = append(calls, "get "+strconv.Itoa(key))
		} else {
			calls = append(calls, "put "+strconv.Itoa(key)+" "+strconv.Itoa(rng.Intn(100_000)))
		}
	}
	return Case{
		Input:    strconv.Itoa(len(calls)) + "\n" + strings.Join(calls, "\n"),
		Expected: strings.Join(runLRU(calls), " "),
	}
}

// runLRU is the reference: a map into a list kept in recency order.
func runLRU(calls []string) []string {
	type entry struct{ key, value int }
	var (
		capacity int
		order    = list.New()
		byKey    = map[int]*list.Element{}
		out      []string
	)
	for _, call := range calls {
		f := strings.Fields(call)
		arg := func(i int) int { n, _ := strconv.Atoi(f[i]); return n }
		switch f[0] {
		case "new":
			capacity, order, byKey = arg(1), list.New(), map[int]*list.Element{}
			out = append(out, "null")
		case "get":
			if e, ok := byKey[arg(1)]; ok {
				order.MoveToFront(e)
				out = append(out, strconv.Itoa(e.Value.(*entry).value))
			} else {
				out = append(out, "-1")
			}
		case "put":
			key, value := arg(1), arg(2)
			if e, ok := byKey[key]; ok {
				e.Value.(*entry).value = value
				order.MoveToFront(e)
			} else {
				byKey[key] = order.PushFront(&entry{key, value})
				if order.Len() > capacity {
					last := order.Back()
					order.Remove(last)
					delete(byKey, last.Value.(*entry).key)
				}
			}
			out = append(out, "null")
		}
	}
	return out
}
