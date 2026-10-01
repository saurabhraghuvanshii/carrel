You are given points on a grid, each as whole numbers `x` and `y`, and a number `k`. Return the positions in the list (counting from 0) of the `k` points closest to the origin (0, 0), by straight-line distance. The tests are chosen so there is no tie at the cut-off. Return the positions in any order.

Compare squared distances, `x × x + y × y`: they order points the same way and need no square root.

## Constraints

- 1 ≤ k ≤ number of points ≤ 10⁴
- −10⁴ ≤ x, y ≤ 10⁴
