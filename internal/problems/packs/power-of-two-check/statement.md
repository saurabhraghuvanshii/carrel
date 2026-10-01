Given a 32-bit signed integer `n`, return `true` if `n` equals 2ᵏ for some whole number k ≥ 0, and `false` otherwise. So 1, 2, 4, 8 and so on are powers of two; 0 and negative numbers are not.

A power of two has exactly one 1 bit. Subtracting 1 turns that bit off and every bit below it on, so `n & (n − 1)` is 0 exactly when `n` has a single 1 bit. Check that `n` is positive first.

## Constraints

- −2³¹ ≤ n ≤ 2³¹ − 1
