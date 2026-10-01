A search tree holding the keys 1 to `n` keeps smaller keys in the left subtree and larger keys in the right subtree of every node. Return the number of different search trees that hold exactly the keys 1 to `n`.

Pick the root `r`. The keys below `r` form the left subtree and the keys above it form the right, and any left shape goes with any right shape. So the count for `n` keys sums, over every root, the count for `r − 1` keys times the count for `n − r` keys. An empty tree counts as one shape.

## Constraints

- 1 ≤ n ≤ 19
