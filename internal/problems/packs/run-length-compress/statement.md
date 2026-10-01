Shorten a word by writing every run of the same letter as the letter followed by the length of the run. A run of length 1 is written as the letter alone. Return the coded word.

Count each run with an inner loop that moves forward while the letter stays the same, then write the letter and, if the count is above 1, the count.

## Constraints

- 1 ≤ word length ≤ 10⁵
- Only `a` to `z`
