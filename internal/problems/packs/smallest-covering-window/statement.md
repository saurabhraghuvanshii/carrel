You are given a text and a pattern, both made of letters; upper and lower case count as different letters. Find the shortest stretch of the text that contains every letter of the pattern at least as many times as the pattern does.

If several stretches are equally short, return the one that starts first. If no stretch works, return an empty string. The answer is shown in quotes.

## Constraints

- 1 ≤ length of the text ≤ 10⁴
- 1 ≤ length of the pattern ≤ 10⁴
- Only `a` to `z` and `A` to `Z`

## Follow-up

Grow the window until it covers the pattern, then shrink it from the left while it still does. How do you know it still covers without recounting?
