In a binary search tree, every value in a node's left subtree is smaller than the node and every value in its right subtree is larger. You are given different whole numbers sorted from smallest to largest. Build a balanced search tree from them and return its root.

Use the middle value as the root, then build the left half and the right half the same way, as in a binary search. When a part has two middle values, use the left one, so that there is exactly one right answer.

## Constraints

- 1 ≤ length ≤ 10⁴
- −10⁹ ≤ value ≤ 10⁹
- The values are different and sorted
