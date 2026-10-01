A warehouse ships orders in the order they were placed. A van takes a batch of orders that are next to each other in that order, and the weights in one batch may add up to at most the van's capacity. The warehouse can send at most `k` vans today, all of the same size. Return the smallest capacity that lets every order ship.

If a capacity works, every larger capacity works too, so you can binary search on the answer. To test one capacity, fill vans greedily: keep adding orders to the current van until the next one does not fit, then start a new van, and count the vans. The answer lies between the heaviest order and the sum of all weights.

## Constraints

- 1 ≤ k ≤ n ≤ 10⁴
- 1 ≤ weight ≤ 10⁴

## Follow-up

Orders may now ship in any order, as long as each van stays within its capacity. Why does the greedy test stop working, and what would you do in practice?
