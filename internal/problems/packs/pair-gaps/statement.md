You are given two lists of the same length. Pair every value in the first list with exactly one value in the second. The gap of a pair is the absolute difference of its two values. Return the smallest possible sum of all the gaps.

Pairing the smallest with the smallest, the second smallest with the second smallest, and so on, is always best: if two pairs cross, swapping their partners never makes the sum larger. The sum can exceed a 32-bit integer.

## Constraints

- 1 ≤ length ≤ 10⁴, and both lists have the same length
- 0 ≤ value ≤ 10⁶
