You are given a start word, an end word and a list of words, all the same length. A ladder goes from the start word to the end word, changing exactly one letter at each step, and every word after the start must be in the list. Return the number of words in the shortest ladder, counting the start and the end, or 0 if there is none. The start and end words are different, and the start word need not be in the list.

Each word is a node, and words one letter apart are joined. Search breadth-first from the start word. To find neighbours quickly, try every letter at every position and look the result up in a set.

## Constraints

- 3 ≤ word length ≤ 5
- 1 ≤ number of words ≤ 3000
- Only `a` to `z`
