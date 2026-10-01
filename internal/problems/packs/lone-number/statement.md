Every value in the list appears exactly twice, except one value that appears once. Return that value, using constant extra memory.

XOR has two useful rules: `x ^ x = 0` and `x ^ 0 = x`, and the order does not matter. So XOR all the values together: every pair cancels and only the lone value is left.

## Constraints

- 1 ≤ length ≤ 10⁴, and the length is odd
- −10⁹ ≤ value ≤ 10⁹
