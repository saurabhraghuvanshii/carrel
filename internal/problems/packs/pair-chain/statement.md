You are given pairs `[a, b]` with `a < b`. A pair `[c, d]` can follow `[a, b]` in a chain when `c > b`; equal values do not chain. Using each pair at most once and in any order, return the length of the longest chain.

Sort the pairs by their second value. Walk through them and take every pair that can follow the last one taken: finishing as early as possible leaves the most room for the rest.

## Constraints

- 1 ≤ number of pairs ≤ 10⁴
- −10⁶ ≤ a < b ≤ 10⁶
