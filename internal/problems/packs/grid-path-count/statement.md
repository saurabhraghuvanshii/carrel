A grid holds `0` (open) and `1` (wall). Starting on the top-left cell and moving only right or down through open cells, count the paths that reach the bottom-right cell. If either corner is a wall, the answer is 0.

The paths into a cell come from the cell above it and the cell to its left, so its count is the sum of theirs. A wall has a count of 0.

## Constraints

- 1 ≤ rows, cols ≤ 16
- Every cell is 0 or 1
