Write a class `Trie` for words of lowercase letters, with three methods:

- `insert(word)` adds the word.
- `search(word)` returns `true` if that exact word was added.
- `startsWith(prefix)` returns `true` if some added word begins with `prefix` (a word begins with itself).

The checker makes one new `Trie` for each test, runs the operations in order and prints the results of every `search` and `startsWith`.

Give each node an array of 26 children and a flag that marks the end of a word. All three methods walk down from the root one letter at a time, so each takes time proportional to the word's length, however many words are stored.

## Constraints

- 1 ≤ number of operations ≤ 10⁴
- 1 ≤ word length ≤ 6
- Only `a` to `z`
