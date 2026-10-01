You are given two lists of whole numbers, each sorted from smallest to largest; one of them may be empty. Return the median of all the values together: the middle value, or the average of the two middle values when the count is even. The answer is shown with one decimal place.

Merging the lists, as in "Merge two sorted lists in place", takes O(m + n). Aim for O(log(min(m, n))): binary search for where to cut the shorter list.

## Constraints

- 0 ≤ m, n ≤ 10⁴, and m + n ≥ 1
- −10⁶ ≤ value ≤ 10⁶
