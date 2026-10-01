You are given intervals `[start, end]` (both ends included) that are sorted by start and already separate: no two overlap or touch. You are also given one more interval. Add it, merging it with any intervals it overlaps or touches, and return the result sorted by start.

The list is already sorted, so you do not need to sort again: copy the intervals that end before the new one, merge the ones it reaches, then copy the rest. That is O(n).

## Constraints

- 0 ≤ number of intervals ≤ 10⁴
- 0 ≤ start < end ≤ 10⁶
