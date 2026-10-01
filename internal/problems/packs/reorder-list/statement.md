You are given the first node of a linked list. Change its links so the nodes come in this order: first, last, second, second to last, third, and so on. Do it in place; return nothing. The answer is shown as the list after the change.

Combine three earlier ideas: find the middle, reverse the second half, then weave the two halves together.

## Constraints

- 1 ≤ length ≤ 10⁴
- −1000 ≤ value ≤ 1000
