Given a list of words, return the longest beginning that every word shares. If the first letters already differ, the answer is the empty string.

Compare the words one column at a time: the answer grows while every word has a letter in that column and they all match. Or start with the first word and shorten it until each word starts with it.

## Constraints

- 1 ≤ number of words ≤ 200
- 1 ≤ word length ≤ 60
- Only `a` to `z`
