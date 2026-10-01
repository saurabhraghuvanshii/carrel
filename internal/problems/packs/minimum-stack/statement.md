Build a class `MinStack` that works like a stack and can also tell you its smallest value at any moment. Every call should take constant time.

- `MinStack()` makes an empty stack.
- `push(x)` puts `x` on top.
- `pop()` removes the top value.
- `top()` returns the top value.
- `minimum()` returns the smallest value in the stack.

Each test is a list of calls. The result line has one entry per call: the value for `top` and `minimum`, and `null` for the others. `pop`, `top` and `minimum` are only called on a stack that is not empty.

## Constraints

- −10⁵ ≤ x ≤ 10⁵
- Up to 5000 calls in one test

## Follow-up

When the smallest value is popped, how does the stack know the next smallest without searching?
