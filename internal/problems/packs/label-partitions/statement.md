Cut a word into as many pieces as you can so that each letter appears in one piece only. Return the lengths of the pieces, from left to right.

First record where each letter last appears. Then walk the word, stretching the current piece's end to the last place of every letter you meet. When you reach that end, the piece is complete.

## Constraints

- 1 ≤ word length ≤ 10⁴
- Only `a` to `z`
