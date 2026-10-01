A word can be built if every one of its beginnings is also in the list: for "tap", the list must hold "t", "ta" and "tap". Return the longest word that can be built. If several have the same length, return the one that comes first in alphabetical order. If no word can be built, return the empty string.

Insert all the words into a trie, then search it depth-first, stepping only onto nodes that end a word. Trying the letters from `a` to `z` means the first longest word you reach is also the first alphabetically.

## Constraints

- 1 ≤ number of words ≤ 2500, all different
- 1 ≤ word length ≤ 8
- Only `a` to `z`
