Given a word, replace every run of the same letter by a single copy of that letter, and return the result. Only letters that sit side by side count as a run, so a letter can still appear more than once in the answer.

Walk the word and keep a letter only when it differs from the letter just before it. Build the answer with a `StringBuilder` (Java) or by appending to a `string` (C++).

## Constraints

- 1 ≤ word length ≤ 10⁵
- Only `a` to `z`
