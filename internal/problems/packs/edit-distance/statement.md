You can change a word one letter at a time: insert a letter anywhere, delete a letter, or replace a letter with another. Return the fewest changes that turn the first word into the second.

Let `d[i][j]` be the answer for the first `i` letters of one word and the first `j` of the other. If those last letters match, nothing needs to change there. Otherwise one change is needed, and the rest comes from the cheapest of the three neighbouring cells. Turning a prefix into an empty word costs its length.

## Constraints

- 1 ≤ length of each word ≤ 500
- Only `a` to `z`
