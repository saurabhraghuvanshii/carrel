You are given a list of different whole numbers that was sorted and then turned, as in "Smallest in a turned list", and a target. Return the position of the target, or −1 if it is missing, in O(log n) time.

At every step, one half of the range is in sorted order. Check whether the target falls inside that half.

## Constraints

- 1 ≤ length ≤ 10⁴
- −10⁹ ≤ value, target ≤ 10⁹
- The values are different
