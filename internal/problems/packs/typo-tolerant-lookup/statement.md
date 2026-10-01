A product search should still find things when the user mistypes one letter. A query matches a word from the catalogue when both have the same length and they differ in at most one position. So "cut" matches "cat" and "cot", and a query also matches the word that is exactly equal to it.

Given the catalogue words (all different) and a list of queries, return one number per query: how many catalogue words it matches.

For every word, and every position in it, store the word with that letter replaced by `*` in a hash map of counts: "cat" adds `*at`, `c*t` and `ca*`. A query then looks up its own masked forms and adds the counts. One catch: a word equal to the query shows up under every one of the query's masked forms, so it is counted once per letter instead of once. Correct for that.

## Constraints

- 1 ≤ n, m ≤ 10⁴
- 1 ≤ length of a word or query ≤ 8
- Only `a` to `z`

## Follow-up

Users also drop a letter or type one too many. How would you extend the same idea to those typos, and how much more would the map hold?
