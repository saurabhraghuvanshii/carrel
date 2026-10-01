Orders and couriers are spread along one main road, each at a position. There are at least as many couriers as orders. Give every order its own courier (a courier takes at most one order; spare couriers stay idle) so that the total distance the couriers travel to their orders is as small as possible. Return that total.

With equal numbers, sorting both lists and pairing them in order is best. With spare couriers you must also choose who stays idle, and taking the nearest courier for each order in turn can go wrong. Sort both lists and let `best[i][j]` be the smallest total for the first `i` orders using only the first `j` couriers: courier `j` either stays idle, or takes order `i`. The total can pass the 32-bit range.

## Constraints

- 1 ≤ n ≤ m ≤ 1000
- 0 ≤ position ≤ 10⁹

## Follow-up

Orders arrive one by one and each must get a courier at once, before the next order is known. What rule would you use, and how far from the best total can it end up?
