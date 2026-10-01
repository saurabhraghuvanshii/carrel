Letters arrive one at a time, in the order of the word. After each letter arrives, write down the earliest letter so far that has appeared exactly once, or `#` if there is none. Return the written letters as one string, the same length as the word.

Keep a count for each letter and a queue of letters in arrival order. After each new letter, drop letters from the front of the queue while their count is above 1; the front is then the answer for this step.

## Constraints

- 1 ≤ word length ≤ 10⁵
- Only `a` to `z`
