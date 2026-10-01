You are given a short pattern and a text, both lowercase words. Return `true` if some stretch of the text uses exactly the same letters as the pattern, each the same number of times, in any order. Otherwise return `false`.

This is "Same letters" for every window of the pattern's length. Keep the letter counts of the window up to date as it slides instead of counting each window again.

## Constraints

- 1 ≤ length of the pattern ≤ 30
- 1 ≤ length of the text ≤ 10⁴
- Only `a` to `z`
