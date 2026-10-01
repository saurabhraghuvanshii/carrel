A ride app sets a price multiplier from recent demand. For every minute it knows how many rides were requested and how many drivers were free. Looking at the last `k` minutes together (the current minute and the `k − 1` before it), let `R` be the total requests and `D` the total free drivers:

- If `R ≤ D`, the multiplier is 1.
- Otherwise it is `R / D` rounded up, but never more than 5. With `D = 0` it is 5.

Return the multiplier for every minute that has `k` full minutes of history, from the oldest to the newest.

Keep both window totals as you move one minute forward: add the new minute and subtract the one that falls out, instead of adding up `k` minutes again.

## Constraints

- 1 ≤ k ≤ n ≤ 10⁵
- 0 ≤ requests, drivers in a minute ≤ 1000

## Follow-up

Minutes now arrive one at a time as a live stream, for thousands of city zones at once. What would you keep in memory for each zone, and how would you handle a minute that arrives late?
