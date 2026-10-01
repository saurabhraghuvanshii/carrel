An on-call engineer wants to know how bad the worst burst of errors was. You are given the time, in seconds, of every error in the log, in order; several errors can share a second. A stretch of `window` seconds ending at time `t` covers the times after `t − window` up to and including `t`. Return the largest number of errors that fall inside one such stretch.

The best stretch can always end on an error, so try each error as the end. Move a second pointer forward past the errors that are `window` seconds or more older; the errors between the two pointers are the ones inside. Both pointers only move forward.

## Constraints

- 1 ≤ n ≤ 10⁵
- 1 ≤ window ≤ 3600
- 0 ≤ time ≤ 10⁹, and times never go down

## Follow-up

The alert must fire while the errors are still coming in, as soon as some stretch reaches a threshold. What do you keep in memory, and how do you stop the alert from firing again every second during one long incident?
