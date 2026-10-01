You are given a grid of letters and a word. Return `true` if the word can be spelled by moving from cell to neighbouring cell (up, down, left or right), using each cell at most once, and `false` otherwise.

Start from every cell that matches the first letter. Mark a cell as used while you explore from it and unmark it when you step back.

## Constraints

- 1 ≤ rows, cols ≤ 6
- 1 ≤ length of the word ≤ 15
- Only `a` to `z`
