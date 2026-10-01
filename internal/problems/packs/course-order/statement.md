There are `n` courses numbered 0 to `n − 1`. Each pair `[u, v]` means course `u` must be taken before course `v`. Return an order in which all `n` courses can be taken, or an empty list if that is impossible because the rules go round in a loop.

Many orders can be right; any of them passes. The tests check your order instead of comparing it with one answer, and show "valid order" when it is right.

Repeatedly take a course that nothing still waiting must come before (in-degree 0). If you run out before taking every course, there is a loop.

## Constraints

- 1 ≤ n ≤ 10⁴
- 0 ≤ number of pairs ≤ 2 × 10⁴
