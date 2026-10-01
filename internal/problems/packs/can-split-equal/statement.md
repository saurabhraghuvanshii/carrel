You are given a list of positive whole numbers. Return `true` if you can put every number into one of two groups so that both groups have the same sum, and `false` otherwise.

An odd total can never split. With an even total, the question becomes: does some selection of the numbers add up to exactly half? Keep a table of which sums are reachable, adding one number at a time. Go through the sums from high to low so each number is used at most once.

## Constraints

- 1 ≤ length ≤ 200
- 1 ≤ value ≤ 100
