Given a 32-bit signed integer `n`, return how many of its 32 bits are 1. Negative numbers are stored in two's complement, so −1 has all 32 bits set.

Test the lowest bit with `n & 1` and shift right, 32 times. In Java, use `>>>` so the sign bit is not copied in. A shorter loop: `n & (n − 1)` clears the lowest 1 bit, so count how many times you can do that before `n` is 0.

## Constraints

- −2³¹ ≤ n ≤ 2³¹ − 1
