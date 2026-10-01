You are given a list of whole numbers sorted from smallest to largest, which may repeat, and a target. Return the first and last positions of the target as `[first, last]`, or `[-1, -1]` if it is missing.

A plain binary search finds some copy of the target. Run it twice, once leaning left and once leaning right, to find the edges in O(log n).

## Constraints

- 0 ≤ length ≤ 10⁴
- −10⁴ ≤ value, target ≤ 10⁴
- The list is sorted from smallest to largest
