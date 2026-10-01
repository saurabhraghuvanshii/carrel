You are given a list of whole numbers that strictly rises and then strictly falls; either part may be empty. Return the position of the largest value in O(log n) time.

Compare the middle value with the one after it: if the list is still rising there, the top is to the right.

## Constraints

- 1 ≤ length ≤ 10⁴
- −10⁹ ≤ value ≤ 10⁹
- The values strictly rise to the top and strictly fall after it
