You are given the roots of two binary trees, `root` and `sub`. Return `true` if some node of `root`, together with everything below it, is exactly the same tree as `sub`, and `false` otherwise.

Check every node of `root` with your answer to "Same tree". A match must include the whole subtree, all the way down to the leaves.

## Constraints

- 1 ≤ number of nodes in `root` ≤ 10⁴
- 1 ≤ number of nodes in `sub` ≤ number of nodes in `root` + 1
- 0 ≤ value ≤ 10⁴
