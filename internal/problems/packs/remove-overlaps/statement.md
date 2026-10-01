You are given intervals `[start, end]`, each running up to but not including its end, so touching intervals do not overlap. Return the fewest intervals you must remove so that none of the rest overlap.

Removing as few as possible means keeping as many as possible. Keep the interval that ends first, then the next one that starts after it ends, and so on.

## Constraints

- 1 ≤ number of intervals ≤ 10⁴
- 0 ≤ start < end ≤ 10⁶
