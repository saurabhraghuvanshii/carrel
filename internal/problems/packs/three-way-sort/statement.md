You are given a list in which every value is 0, 1 or 2. Rearrange it in place so that all the 0s come first, then all the 1s, then all the 2s.

Counting each value and rewriting the list takes two passes. Do it in a single pass with a fixed amount of extra memory, using three pointers: the end of the 0s, the current value, and the start of the 2s.

## Constraints

- 1 ≤ length ≤ 10⁴
- Every value is 0, 1 or 2
