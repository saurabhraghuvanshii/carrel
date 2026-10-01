A stretch of a word is a part of it without gaps, at least one letter long. Return how many different stretches the word has; a stretch that appears in several places counts once.

Every stretch is the beginning of some suffix. Insert every suffix into a trie: each path from the root spells one stretch, so the answer is the number of nodes, not counting the root.

## Constraints

- 1 ≤ word length ≤ 500
- Only `a` to `z`
