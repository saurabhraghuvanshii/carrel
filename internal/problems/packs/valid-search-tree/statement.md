You are given the root of a binary tree. Return `true` if it is a search tree: for every node, every value in its left subtree is smaller and every value in its right subtree is larger. Equal values are not allowed. Otherwise return `false`. An empty tree is a search tree.

Comparing each node only with its parent is not enough: a value deep in a left subtree must also be smaller than every ancestor it sits to the left of. Pass the allowed range down.

## Constraints

- 0 ≤ number of nodes ≤ 10⁴
- −2 × 10⁵ ≤ value ≤ 2 × 10⁵
