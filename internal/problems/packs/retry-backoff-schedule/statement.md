A payment client calls a bank that is down, so every attempt fails and is tried again. The first attempt starts at time 0. The wait before the next attempt starts at `base` and doubles after every attempt, but is never more than `cap`: the waits are `min(cap, base)`, `min(cap, 2 × base)`, `min(cap, 4 × base)`, and so on. An attempt takes no time.

Return how many attempts start at or before `deadline`, counting the first one.

While the wait is still below `cap` it doubles, so that part takes only about 30 steps. Once the wait has reached `cap` every later wait is the same, and one division tells you how many more attempts fit. Stepping through them one at a time is far too slow when `cap` is small and the deadline is large. Use 64-bit integers.

## Constraints

- 1 ≤ base ≤ 10⁹
- 1 ≤ cap ≤ 10⁹
- 0 ≤ deadline ≤ 10¹⁸

## Follow-up

A thousand clients fail at the same moment and all retry on this same schedule, so the bank gets hit by waves. What would you add to the wait to spread them out, and what does that do to the worst case?
