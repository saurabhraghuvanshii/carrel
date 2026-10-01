package testgen

import (
	"math/rand"
	"slices"
	"strconv"
	"strings"
)

// Generators for the real-interview group "Orders and payments".

func init() {
	register("cart-discounts", genCartDiscounts)
	register("refund-window", genRefundWindow)
	register("retry-backoff", genRetryBackoff)
	register("duplicate-payments", genDuplicatePayments)
	register("order-batching", genOrderBatching)
	register("inventory-reserve", genInventoryReserve)
}

// rowsText prints every row on its own line, each starting with a newline,
// so it can follow a header line directly.
func rowsText(rows [][]int) string {
	var b strings.Builder
	for _, r := range rows {
		b.WriteString("\n" + joinInts(r, " "))
	}
	return b.String()
}

func boolsOutput(vals []bool) string {
	parts := make([]string, len(vals))
	for i, v := range vals {
		parts[i] = boolText(v)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

// cart-discounts: item prices, then coupons "type value" applied in order:
// 1 p takes p percent off the running total (the discount rounds down), 2 f
// takes f off (never below 0), and 0 0 removes the latest coupon still in
// force. The answer is the final total.
func genCartDiscounts(rng *rand.Rand, size int) Case {
	prices := randomInts(rng, max(1, min(size, 10_000)), 1, []int{20, 1000, 1_000_000}[rng.Intn(3)])
	sum := 0
	for _, p := range prices {
		sum += p
	}
	m := rng.Intn(min(21, size+1))
	var coupons, active [][]int
	for i := 0; i < m; i++ {
		switch k := rng.Intn(5); {
		case k == 0:
			coupons = append(coupons, []int{0, 0})
			if len(active) > 0 {
				active = active[:len(active)-1]
			}
		case k <= 2:
			c := []int{1, 1 + rng.Intn([]int{30, 100}[rng.Intn(2)])}
			coupons, active = append(coupons, c), append(active, c)
		default:
			top := max(1, sum/4)
			if rng.Intn(6) == 0 {
				top = sum // sometimes enough to reach zero
			}
			c := []int{2, 1 + rng.Intn(min(top, 1_000_000_000))}
			coupons, active = append(coupons, c), append(active, c)
		}
	}
	total := sum
	for _, c := range active {
		if c[0] == 1 {
			total -= total * c[1] / 100
		} else {
			total = max(0, total-c[1])
		}
	}
	return Case{Input: arrayInput(prices) + "\n" + strconv.Itoa(m) + rowsText(coupons), Expected: strconv.Itoa(total)}
}

// refund-window: purchases "order day" and refund requests "order day" in
// the order they arrive. A request is approved when the order exists, the
// request is 0 to window days after the purchase, and the order was not
// refunded before. The answer has one true or false per request.
func genRefundWindow(rng *rand.Rand, size int) Case {
	n := max(1, min(size, 5000))
	m := 1 + rng.Intn(max(1, min(size, 5000)))
	window := []int{0, 7, 30}[rng.Intn(3)]
	ids := distinct(rng, n, 1, max(3*n, 10))
	bought := map[int]int{}
	purchases := make([][]int, n)
	for i, id := range ids {
		bought[id] = rng.Intn(100)
		purchases[i] = []int{id, bought[id]}
	}
	refunded := map[int]bool{}
	requests := make([][]int, m)
	results := make([]bool, m)
	for i := range requests {
		id := ids[rng.Intn(n)]
		if rng.Intn(6) == 0 {
			id = 1 + rng.Intn(3*n+10) // maybe unknown
		}
		day, known := bought[id]
		at := max(0, day+rng.Intn(2*window+4)-1)
		requests[i] = []int{id, at}
		results[i] = known && at >= day && at-day <= window && !refunded[id]
		if results[i] {
			refunded[id] = true
		}
	}
	return Case{
		Input:    ints(n, m, window) + rowsText(purchases) + rowsText(requests),
		Expected: boolsOutput(results),
	}
}

// retry-backoff: "base cap deadline". Attempt 1 starts at time 0; the wait
// before attempt j + 1 is min(cap, base * 2^(j-1)). The answer is how many
// attempts start at or before the deadline.
func genRetryBackoff(rng *rand.Rand, size int) Case {
	base := 1 + rng.Intn([]int{5, 1000, 1_000_000_000}[rng.Intn(3)])
	limit := 1 + rng.Intn([]int{5, 1000, 1_000_000_000}[rng.Intn(3)])
	deadline := rng.Int63n([]int64{50, 1_000_000, 1_000_000_000_000_000_000}[rng.Intn(3)] + 1)
	t, count, wait := int64(0), int64(1), int64(base)
	for wait < int64(limit) && t+wait <= deadline {
		t += wait
		count++
		wait *= 2
	}
	if wait >= int64(limit) {
		count += (deadline - t) / int64(limit)
	}
	return Case{
		Input:    strconv.Itoa(base) + " " + strconv.Itoa(limit) + " " + strconv.FormatInt(deadline, 10),
		Expected: strconv.FormatInt(count, 10),
	}
}

// duplicate-payments: "n window", then payments "time user merchant amount"
// in time order. A payment is a duplicate when the latest earlier payment
// with the same user, merchant and amount is at most window before it. The
// answer lists the positions of the duplicates.
func genDuplicatePayments(rng *rand.Rand, size int) Case {
	n := max(1, min(size, 10_000))
	window := []int{0, 5, 60}[rng.Intn(3)]
	users, merchants, amounts := 1+rng.Intn(3), 1+rng.Intn(2), 1+rng.Intn(3)
	last := map[[3]int]int{}
	rows := make([][]int, n)
	var flagged []int
	t := 0
	for i := range rows {
		t += rng.Intn(4) * rng.Intn(3)
		key := [3]int{rng.Intn(users), rng.Intn(merchants), 100 * (1 + rng.Intn(amounts))}
		rows[i] = []int{t, key[0], key[1], key[2]}
		if before, ok := last[key]; ok && t-before <= window {
			flagged = append(flagged, i)
		}
		last[key] = t
	}
	return Case{Input: ints(n, window) + rowsText(rows), Expected: listOutput(flagged)}
}

// batchesNeeded counts the batches when orders ship in order and no batch
// weighs more than capacity.
func batchesNeeded(weights []int, capacity int) int {
	batches, load := 1, 0
	for _, w := range weights {
		if load+w > capacity {
			batches++
			load = 0
		}
		load += w
	}
	return batches
}

// order-batching: "n k", then n weights. Orders ship in order, in at most k
// batches; the answer is the smallest batch capacity that allows it.
func genOrderBatching(rng *rand.Rand, size int) Case {
	n := max(1, min(size, 10_000))
	weights := randomInts(rng, n, 1, []int{5, 100, 10_000}[rng.Intn(3)])
	k := 1 + rng.Intn(min(n, []int{2, 10, n}[rng.Intn(3)]))
	lo, hi := slices.Max(weights), 0
	for _, w := range weights {
		hi += w
	}
	for lo < hi {
		mid := (lo + hi) / 2
		if batchesNeeded(weights, mid) <= k {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return Case{Input: ints(n, k) + "\n" + joinInts(weights, " "), Expected: strconv.Itoa(lo)}
}

// inventory-reserve: calls new; add item qty; reserve order item qty;
// cancel order; confirm order; available item. A reservation holds stock
// until it is cancelled (stock returns) or confirmed (stock is gone).
func genInventoryReserve(rng *rand.Rand, size int) Case {
	calls := []string{"new"}
	results := []string{"null"}
	items, orders := 1+rng.Intn(3), 2+rng.Intn(6)
	free := map[int]int{}
	held := map[int][2]int{}
	for len(calls) < callCount(size) {
		item, order := rng.Intn(items), rng.Intn(orders)
		switch k := rng.Intn(10); {
		case k < 2:
			qty := 1 + rng.Intn(10)
			free[item] += qty
			calls, results = append(calls, "add "+ints(item, qty)), append(results, "null")
		case k < 5:
			qty := 1 + rng.Intn(5)
			_, busy := held[order]
			ok := !busy && free[item] >= qty
			if ok {
				free[item] -= qty
				held[order] = [2]int{item, qty}
			}
			calls, results = append(calls, "reserve "+ints(order, item, qty)), append(results, boolText(ok))
		case k < 8:
			h, ok := held[order]
			name := "confirm "
			if k == 5 || (k == 6 && rng.Intn(2) == 0) {
				name = "cancel "
				if ok {
					free[h[0]] += h[1]
				}
			}
			delete(held, order)
			calls, results = append(calls, name+strconv.Itoa(order)), append(results, boolText(ok))
		default:
			calls, results = append(calls, "available "+strconv.Itoa(item)), append(results, strconv.Itoa(free[item]))
		}
	}
	return callCase(calls, results)
}
