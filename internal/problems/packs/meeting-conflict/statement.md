You are given meetings as `[start, end]` times. A meeting runs from its start up to, but not including, its end, so one meeting may start exactly when another ends. Return `true` if one person can attend all of them, meaning no two overlap, and `false` otherwise.

Sort the meetings by start time. Then only neighbours can clash.

## Constraints

- 0 ≤ number of meetings ≤ 10⁴
- 0 ≤ start < end ≤ 10⁶
