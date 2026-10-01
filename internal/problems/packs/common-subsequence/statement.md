A subsequence of a word keeps some of its letters in their original order, possibly with gaps. Given two words, return the length of the longest subsequence they share.

Compare prefixes. If the last letters of two prefixes match, the answer for them is one more than for both prefixes without that letter. If not, drop the last letter of one word or the other and keep the better result.

## Constraints

- 1 ≤ length of each word ≤ 1000
- Only `a` to `z`
