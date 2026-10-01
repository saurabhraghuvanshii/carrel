You are given a list of whole numbers and a window size `k`. Slide a window of `k` neighbouring values from the left end to the right end, one step at a time, and return the largest value in the window at each step, in order.

## Constraints

- 1 ≤ k ≤ length ≤ 10⁴
- −10⁴ ≤ value ≤ 10⁴

## Follow-up

Keep a queue of positions whose values get smaller from front to back. Why does each position enter and leave the queue only once, giving O(n) time?
