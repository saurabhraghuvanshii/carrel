package testgen

import (
	"strconv"
	"strings"
	"testing"
)

// rowsOf reads count rows of ints starting at lines[from].
func rowsOf(lines []string, from, count int) [][]int {
	rows := make([][]int, count)
	for i := range rows {
		rows[i] = intsOf(lines[from+i])
	}
	return rows
}

func TestCartDiscounts(t *testing.T) {
	undos := 0
	for i, c := range cases(t, "cart-discounts", 80) {
		lines := strings.Split(c.Input, "\n")
		m, _ := strconv.Atoi(lines[2])
		coupons := rowsOf(lines, 3, m)
		// Walk backwards: each undo cancels the nearest earlier coupon that
		// is not cancelled yet.
		keep := make([]bool, m)
		skip := 0
		for k := m - 1; k >= 0; k-- {
			switch {
			case coupons[k][0] == 0:
				skip++
				undos++
			case skip > 0:
				skip--
			default:
				keep[k] = true
			}
		}
		want := cartTotal(intsOf(lines[1]), coupons, keep)
		if strconv.Itoa(want) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, want, c.Expected)
		}
	}
	if undos == 0 {
		t.Fatal("no case has an undo")
	}
}

// cartTotal applies the kept coupons; a percent discount is the largest
// whole d with 100 * d <= total * p.
func cartTotal(prices []int, coupons [][]int, keep []bool) int {
	total := 0
	for _, p := range prices {
		total += p
	}
	for k, cp := range coupons {
		if !keep[k] {
			continue
		}
		if cp[0] == 2 {
			total = max(0, total-cp[1])
			continue
		}
		d := total / 100 * cp[1]
		for 100*(d+1) <= total*cp[1] {
			d++
		}
		total -= d
	}
	return total
}

func TestRefundWindow(t *testing.T) {
	seen := map[bool]bool{}
	for i, c := range cases(t, "refund-window", 60) {
		lines := strings.Split(c.Input, "\n")
		head := intsOf(lines[0])
		purchases, requests := rowsOf(lines, 1, head[0]), rowsOf(lines, 1+head[0], head[1])
		got := make([]bool, len(requests))
		for k, r := range requests {
			for _, p := range purchases {
				if p[0] == r[0] && r[1] >= p[1] && r[1] <= p[1]+head[2] {
					got[k] = true
				}
			}
			for j := 0; j < k; j++ {
				if got[j] && requests[j][0] == r[0] {
					got[k] = false
				}
			}
			seen[got[k]] = true
		}
		if boolsOutput(got) != c.Expected {
			t.Fatalf("case %d differs", i)
		}
	}
	if !seen[true] || !seen[false] {
		t.Fatal("both answers must appear")
	}
}

func TestRetryBackoff(t *testing.T) {
	checked := 0
	for i, c := range cases(t, "retry-backoff", 120) {
		f := strings.Fields(c.Input)
		base, _ := strconv.ParseInt(f[0], 10, 64)
		limit, _ := strconv.ParseInt(f[1], 10, 64)
		deadline, _ := strconv.ParseInt(f[2], 10, 64)
		if deadline/limit > 2_000_000 {
			continue
		}
		// Step through the attempts one at a time.
		count, at, raw := int64(1), int64(0), base
		for {
			at += min(limit, raw)
			if at > deadline {
				break
			}
			count++
			if raw < limit {
				raw *= 2
			}
		}
		if strconv.FormatInt(count, 10) != c.Expected {
			t.Fatalf("case %d (%s): want %d, generator says %s", i, c.Input, count, c.Expected)
		}
		checked++
	}
	if checked < 40 {
		t.Fatalf("only %d cases were small enough to check", checked)
	}
}

func TestDuplicatePayments(t *testing.T) {
	flagged := 0
	for i, c := range cases(t, "duplicate-payments", 60) {
		lines := strings.Split(c.Input, "\n")
		head := intsOf(lines[0])
		rows := rowsOf(lines, 1, head[0])
		if len(rows) > 2000 {
			continue
		}
		var got []int
		for k, r := range rows {
			for j := k - 1; j >= 0; j-- {
				p := rows[j]
				if p[1] == r[1] && p[2] == r[2] && p[3] == r[3] {
					if r[0]-p[0] <= head[1] {
						got = append(got, k)
					}
					break
				}
			}
		}
		if listOutput(got) != c.Expected {
			t.Fatalf("case %d differs", i)
		}
		flagged += len(got)
	}
	if flagged == 0 {
		t.Fatal("no duplicates in any case")
	}
}

func TestOrderBatching(t *testing.T) {
	for i, c := range cases(t, "order-batching", 60) {
		lines := strings.Split(c.Input, "\n")
		k := intsOf(lines[0])[1]
		weights := intsOf(lines[1])
		sum := 0
		for _, w := range weights {
			sum += w
		}
		if sum > 20_000 {
			continue
		}
		// Try every capacity from 1 up; a batch is cut only when it must be.
		want := 0
		for capacity := 1; want == 0; capacity++ {
			batches, load, ok := 1, 0, true
			for _, w := range weights {
				if w > capacity {
					ok = false
					break
				}
				if load+w > capacity {
					batches, load = batches+1, 0
				}
				load += w
			}
			if ok && batches <= k {
				want = capacity
			}
		}
		if strconv.Itoa(want) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, want, c.Expected)
		}
	}
}

func TestInventoryReserve(t *testing.T) {
	replay(t, "inventory-reserve", 50, func([]int) func(string, []int) string {
		// Keep every event and work the stock out from the whole history.
		type hold struct {
			order, item, qty int
			state            string // held, cancelled, confirmed
		}
		var added [][2]int
		var holds []*hold
		active := func(order int) *hold {
			for _, h := range holds {
				if h.order == order && h.state == "held" {
					return h
				}
			}
			return nil
		}
		stock := func(item int) int {
			n := 0
			for _, a := range added {
				if a[0] == item {
					n += a[1]
				}
			}
			for _, h := range holds {
				if h.item == item && h.state != "cancelled" {
					n -= h.qty
				}
			}
			return n
		}
		return func(call string, x []int) string {
			switch call {
			case "add":
				added = append(added, [2]int{x[0], x[1]})
				return "null"
			case "reserve":
				if active(x[0]) != nil || stock(x[1]) < x[2] {
					return "false"
				}
				holds = append(holds, &hold{x[0], x[1], x[2], "held"})
				return "true"
			case "cancel", "confirm":
				h := active(x[0])
				if h == nil {
					return "false"
				}
				h.state = "cancelled"
				if call == "confirm" {
					h.state = "confirmed"
				}
				return "true"
			}
			return strconv.Itoa(stock(x[0]))
		}
	})
}
