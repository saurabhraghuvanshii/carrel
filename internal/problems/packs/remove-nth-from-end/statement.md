You are given the first node of a linked list and a number `n`. Remove the `n`-th node counting from the end (the last node is 1st from the end) and return the first node of the changed list.

Try it in one pass: start a lead pointer `n` nodes ahead, then move both until the lead reaches the end.

## Constraints

- 1 ≤ length ≤ 10⁴
- 1 ≤ n ≤ length
- −1000 ≤ value ≤ 1000
