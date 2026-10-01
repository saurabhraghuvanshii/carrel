You are given meetings as `[start, end]` times, each running up to but not including its end. Return the smallest number of rooms needed so that no two meetings share a room at the same time. A meeting may use a room the moment another meeting there ends.

Sort by start and keep a min-heap of end times for the rooms in use. Or sort all starts and all ends separately and sweep through them.

## Constraints

- 1 ≤ number of meetings ≤ 10⁴
- 0 ≤ start < end ≤ 10⁶
