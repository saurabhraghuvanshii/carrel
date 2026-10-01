A staircase has `n` steps. Each move climbs either one step or two steps. Return the number of different move sequences that end exactly on the top step.

To stand on step `i`, your last move came from step `i − 1` or from step `i − 2`. So the count for `i` is the sum of the counts for those two steps. Build the counts upward from the bottom instead of recomputing them.

## Constraints

- 1 ≤ n ≤ 45
- The answer fits in a 32-bit signed integer
