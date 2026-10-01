You are given a list of whole numbers and a number `k`. Return the `k`-th largest value. Repeated values count separately, so in `[9, 9, 7]` the 2nd largest is 9.

Sorting works in O(n log n). Keep a min-heap of the `k` largest values seen so far instead: its top is the answer, in O(n log k).

## Constraints

- 1 ≤ k ≤ length ≤ 10⁴
- −10⁴ ≤ value ≤ 10⁴
