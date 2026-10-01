A list holds `n` different values, all from 0 to `n`. Exactly one value in that range is not in the list. Return it.

The values 0 to `n` add up to n(n + 1)/2, so the missing one is that total minus the sum of the list. Or XOR every index and every value together with `n`: everything present cancels out, leaving the missing value.

## Constraints

- 1 ≤ n ≤ 10⁴
- The values are different and 0 ≤ value ≤ n
