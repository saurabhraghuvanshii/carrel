You are given a whole number `x` that is zero or more. Return the largest whole number `r` with `r × r ≤ x`, the square root rounded down. Do not use a library square root.

Binary search over the answer: if `mid × mid` is too large, the root is smaller. Watch out: `mid × mid` can be larger than a 32-bit integer.

## Constraints

- 0 ≤ x ≤ 2³¹ − 1
