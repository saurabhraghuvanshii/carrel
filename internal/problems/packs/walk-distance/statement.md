A walk is a string of steps: `N` moves one unit north, `S` south, `E` east and `W` west. Return the square of the straight-line distance between where the walk starts and where it ends. The square is a whole number, so there is no rounding to worry about.

Track the position `(x, y)`. At the end the answer is x² + y². The walk can be long enough that this passes the 32-bit range.

## Constraints

- 1 ≤ number of steps ≤ 10⁵
- Only `N`, `S`, `E` and `W`
