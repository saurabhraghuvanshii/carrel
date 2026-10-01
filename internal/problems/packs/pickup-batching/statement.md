A shared-ride service groups pickup requests into batches. The map is cut into square zones 100 units wide: the point `(x, y)` is in zone `(⌊x / 100⌋, ⌊y / 100⌋)`, rounding down, so `x = −50` is in column −1.

Requests arrive in time order as `[time, x, y]`. Each zone has at most one open batch. A request joins its zone's open batch if the batch has fewer than `limit` requests and the request's time is at most `window` after the batch's first request. Otherwise the request starts a new batch for its zone, which replaces the old one. Return the number of batches.

Use a hash map from zone to its open batch (first time and size). Watch the division: in Java and C++, `−50 / 100` is 0, not −1. Use `Math.floorDiv` in Java, or adjust the result for negative values in C++.

## Constraints

- 1 ≤ n ≤ 10⁴
- 1 ≤ limit ≤ 10
- 0 ≤ window ≤ 1000
- 0 ≤ time ≤ 10⁹, times never go down
- −10⁶ ≤ x, y ≤ 10⁶

## Follow-up

A batch should also be sent off once its window runs out, even if no new request comes to its zone. How would you find the batches that are due without scanning every zone, and what changes when requests for one city are spread over several servers?
