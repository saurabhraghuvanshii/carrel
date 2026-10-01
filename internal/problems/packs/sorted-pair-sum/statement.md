You are given a list of whole numbers sorted from smallest to largest, and a target. Exactly one pair of positions holds values that add up to the target. Return those two positions, counting from 0.

You already solved this with a map in "Pair with target sum". This time use only a fixed amount of extra memory: the list is sorted, so start one pointer at each end.

## Constraints

- 2 ≤ length ≤ 10⁴
- −5 × 10⁸ ≤ value ≤ 5 × 10⁸
- The list is sorted from smallest to largest
- Exactly one pair adds up to the target
- The pair may come back in any order
