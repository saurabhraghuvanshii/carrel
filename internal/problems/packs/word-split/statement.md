You are given a text with no spaces and a list of words. Return `true` if the text can be cut into pieces that are all on the list, and `false` otherwise. A word may be used any number of times.

Trying every first word and recursing on the rest takes exponential time on texts like `aaaa…ab`. Remember, for every position, whether the text from there to the end can be cut. Then each position is worked out once.

## Constraints

- 1 ≤ text length ≤ 300
- 1 ≤ number of words ≤ 20, all different
- 1 ≤ word length ≤ 10
- Only `a` to `z`
