You are given intervals `[start, end]` that include both ends. Merge every group of intervals that overlap or touch, so `[3, 5]` and `[5, 8]` become `[3, 8]`. Return the merged intervals sorted by start.

After sorting by start, each interval either extends the last merged one or begins a new one.

## Constraints

- 1 ≤ number of intervals ≤ 10⁴
- 0 ≤ start < end ≤ 10⁶
