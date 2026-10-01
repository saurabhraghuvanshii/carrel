You are given a list of different whole numbers sorted from smallest to largest, and a target. Return the position of the target, counting from 0, or −1 if it is not in the list.

Look at the middle value: it tells you which half the target must be in. Halving the list each step takes O(log n) time.

## Constraints

- 1 ≤ length ≤ 10⁴
- −10⁹ ≤ value, target ≤ 10⁹
- The values are different and sorted from smallest to largest
