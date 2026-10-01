A subsequence keeps some values of a list in their original order, possibly with gaps. Return the length of the longest subsequence in which every value is strictly larger than the one before it.

For each position, work out the longest rising subsequence that ends there: one more than the best among earlier positions with a smaller value.

## Constraints

- 1 ≤ length ≤ 2500
- −10⁴ ≤ value ≤ 10⁴

## Follow-up

Keep, for every length, the smallest value that can end a rising subsequence of that length. These smallest endings are sorted, so binary search can place each new value. That brings the time down to O(n log n).
