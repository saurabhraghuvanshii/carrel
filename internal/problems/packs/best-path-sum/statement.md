You are given the root of a binary tree whose values may be negative. A path is any chain of one or more nodes joined by parent and child links, visiting each node at most once; it does not have to pass through the root or reach a leaf. Return the largest sum of values along any path.

Like "Longest path in a tree", the best path through a node joins the best downward path from its left child and from its right child. Leave out a side whose best is negative.

## Constraints

- 1 ≤ number of nodes ≤ 10⁴
- −1000 ≤ value ≤ 1000
