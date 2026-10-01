Write four small methods on a non-negative integer `n`. Bit 0 is the lowest bit.

- `getBit(n, i)` returns bit `i` of `n` (0 or 1).
- `setBit(n, i)` returns `n` with bit `i` set to 1.
- `clearBit(n, i)` returns `n` with bit `i` set to 0.
- `updateBit(n, i, b)` returns `n` with bit `i` set to `b`.

The checker starts from `n` and runs the operations in order. Each change replaces `n` with the value you return. It prints the result of every operation.

Each method needs one mask, `1 << i`. OR with it sets the bit, AND with its complement `~(1 << i)` clears it, and AND with it reads it. Watch the brackets: `~1 << i` is a different number.

## Constraints

- 0 ≤ n < 2³¹, and n stays in that range
- 0 ≤ i ≤ 30, b is 0 or 1
- 1 ≤ number of operations ≤ 10⁴
