You have coins of a few different values and as many of each as you want. Return the number of different ways to pay `amount` exactly. Order does not matter: 2 + 3 and 3 + 2 are the same way. There is one way to pay 0: use no coins.

To avoid counting the same mix twice, bring in the coins one value at a time. After adding a coin value `c`, the ways to pay `v` are the ways without `c` plus the ways to pay `v − c` with `c` allowed.

## Constraints

- 1 ≤ number of coins ≤ 12
- 1 ≤ coin ≤ 400, and the coins are all different
- 0 ≤ amount ≤ 5000
- The answer fits in a 32-bit signed integer
