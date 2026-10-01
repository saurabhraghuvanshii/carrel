You are given different words, and no word is the beginning of another. For every word, find its shortest beginning that no other word starts with. Return them in the same order as the words.

Insert every word into a trie and count, at each node, how many words pass through it. A word's answer ends at the first node on its path with a count of 1.

## Constraints

- 1 ≤ number of words ≤ 1000
- 1 ≤ word length ≤ 13
- Only `a` to `z`
