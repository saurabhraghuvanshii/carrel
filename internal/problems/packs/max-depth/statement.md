You get the root of a binary tree. Every node has at most two children, a left one and a right one. The depth of the tree is the number of nodes on the longest path that starts at the root and goes down to a node with no children. Return the depth. An empty tree has depth 0.

The examples write a tree in level order, row by row from the top, with `null` where a child is missing.

## Constraints

- 0 ≤ number of nodes ≤ 10⁴
- −100 ≤ value ≤ 100
- The tree can be one long path, so a recursive solution may go 10⁴ calls deep
