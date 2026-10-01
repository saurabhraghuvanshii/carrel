Children stand in a row, each with a rating. Give every child at least one sweet, and give a child more sweets than each neighbour with a lower rating. Neighbours with equal ratings have no rule between them. Return the fewest sweets in total.

Do two passes. Left to right, a child rated above the left neighbour gets one more than that neighbour. Right to left, a child rated above the right neighbour needs at least one more than that neighbour; keep the larger of the two numbers.

## Constraints

- 1 ≤ number of children ≤ 10⁴
- 0 ≤ rating ≤ 20000
