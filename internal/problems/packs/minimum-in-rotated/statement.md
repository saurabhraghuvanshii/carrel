A list of different whole numbers was sorted from smallest to largest and then turned: some values from the front were moved, in order, to the back. For example `[2, 4, 6, 8]` can become `[6, 8, 2, 4]`. The list may also not have been turned at all. Return the smallest value in O(log n) time.

Compare the middle value with the last one to see which side the turning point is on.

## Constraints

- 1 ≤ length ≤ 10⁴
- −10⁹ ≤ value ≤ 10⁹
- The values are different
