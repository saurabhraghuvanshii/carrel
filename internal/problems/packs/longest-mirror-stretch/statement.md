A mirror stretch is a part of the word, without gaps, that reads the same forwards and backwards. Return the longest mirror stretch. If several have the same length, return the one that starts first.

Every mirror stretch has a centre: a letter (odd length) or the gap between two letters (even length). From each of the 2n − 1 centres, grow outwards while the two ends match. That is O(n²) time and O(1) extra memory.

## Constraints

- 1 ≤ word length ≤ 1000
- Only `a` to `z`
