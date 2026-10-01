You are given an `n` by `n` board where `.` is a usable square and `#` is a hole. Count the ways to place `n` queens on usable squares so that no two queens attack each other: no two share a row, a column or a diagonal.

There is one queen in each row. Place them row by row and keep a record of which columns and diagonals are taken, so each check is instant.

## Constraints

- 1 ≤ n ≤ 10
- Every square is `.` or `#`
