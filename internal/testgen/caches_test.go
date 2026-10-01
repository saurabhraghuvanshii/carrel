package testgen

import (
	"strconv"
	"strings"
	"testing"
)

// replay feeds every call after "new" to step and compares the results.
func replay(t *testing.T, name string, n int, step func(newArgs []int) func(call string, args []int) string) {
	t.Helper()
	seen := map[string]bool{}
	for i, c := range cases(t, name, n) {
		lines := strings.Split(c.Input, "\n")
		run := step(intsOf(strings.TrimPrefix(lines[1], "new")))
		got := []string{"null"}
		for _, l := range lines[2:] {
			f := strings.Fields(l)
			got = append(got, run(f[0], intsOf(strings.Join(f[1:], " "))))
		}
		if strings.Join(got, " ") != c.Expected {
			t.Fatalf("case %d differs", i)
		}
		for _, r := range got {
			seen[r] = true
		}
	}
	if seen["true"] != seen["false"] {
		t.Fatal("both answers must appear")
	}
}

func TestFixedWindow(t *testing.T) {
	replay(t, "fixed-window", 50, func(a []int) func(string, []int) string {
		limit, window := a[0], a[1]
		var allowed [][2]int
		return func(_ string, x []int) string {
			// Count this user's allowed requests that share the window.
			n := 0
			for _, r := range allowed {
				if r[0] == x[0] && r[1]/window == x[1]/window {
					n++
				}
			}
			if n >= limit {
				return "false"
			}
			allowed = append(allowed, [2]int{x[0], x[1]})
			return "true"
		}
	})
}

func TestSlidingLimiter(t *testing.T) {
	replay(t, "sliding-limiter", 50, func(a []int) func(string, []int) string {
		limit, window := a[0], a[1]
		var allowed [][2]int
		return func(_ string, x []int) string {
			n := 0
			for _, r := range allowed {
				if r[0] == x[0] && x[1]-window < r[1] && r[1] <= x[1] {
					n++
				}
			}
			if n >= limit {
				return "false"
			}
			allowed = append(allowed, [2]int{x[0], x[1]})
			return "true"
		}
	})
}

func TestTokenBucket(t *testing.T) {
	replay(t, "token-bucket", 50, func(a []int) func(string, []int) string {
		capacity, refill := a[0], a[1]
		tokens, clock := capacity, 0
		return func(_ string, x []int) string {
			// Step the clock one unit at a time.
			for clock < x[0] {
				clock++
				if clock%refill == 0 && tokens < capacity {
					tokens++
				}
			}
			if tokens < x[1] {
				return "false"
			}
			tokens -= x[1]
			return "true"
		}
	})
}

func TestTTLCache(t *testing.T) {
	replay(t, "ttl-cache", 50, func([]int) func(string, []int) string {
		type put struct{ key, value, until int }
		var puts []put
		latest := func(key, at int) (put, bool) {
			for i := len(puts) - 1; i >= 0; i-- {
				if puts[i].key == key {
					return puts[i], puts[i].until > at
				}
			}
			return put{}, false
		}
		return func(call string, x []int) string {
			switch call {
			case "put":
				puts = append(puts, put{x[0], x[1], x[3] + x[2]})
				return "null"
			case "get":
				if p, ok := latest(x[0], x[1]); ok {
					return strconv.Itoa(p.value)
				}
				return "-1"
			}
			n := 0
			for key := 0; key < 100; key++ {
				if _, ok := latest(key, x[0]); ok {
					n++
				}
			}
			return strconv.Itoa(n)
		}
	})
}

func TestLFUCache(t *testing.T) {
	replay(t, "lfu-cache", 50, func(a []int) func(string, []int) string {
		capacity := a[0]
		type entry struct{ key, value, uses, last int }
		var cache []*entry
		clock := 0
		find := func(key int) *entry {
			for _, e := range cache {
				if e.key == key {
					return e
				}
			}
			return nil
		}
		return func(call string, x []int) string {
			clock++
			e := find(x[0])
			if call == "get" {
				if e == nil {
					return "-1"
				}
				e.uses, e.last = e.uses+1, clock
				return strconv.Itoa(e.value)
			}
			if e != nil {
				e.value, e.uses, e.last = x[1], e.uses+1, clock
				return "null"
			}
			if len(cache) == capacity {
				worst := 0
				for i, c := range cache {
					w := cache[worst]
					if c.uses < w.uses || (c.uses == w.uses && c.last < w.last) {
						worst = i
					}
				}
				cache = append(cache[:worst], cache[worst+1:]...)
			}
			cache = append(cache, &entry{x[0], x[1], 1, clock})
			return "null"
		}
	})
}
