You are given the first node of a linked list and a number `k`. Reverse the nodes in each group of `k` neighbours, from the front. If fewer than `k` nodes are left at the end, leave them as they are. Return the first node of the changed list.

Reversing one group is "Reverse a linked list" on a piece of the list. The hard part is linking each reversed group to the one before it and the one after it.

## Constraints

- 1 ≤ length ≤ 10⁴
- 1 ≤ k ≤ 13
- −1000 ≤ value ≤ 1000
