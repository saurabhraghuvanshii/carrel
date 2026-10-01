A session service caches login tokens that expire. A value stored at time `t` with a time to live `ttl` can be read at times before `t + ttl`, and is gone from `t + ttl` on. Storing a key again replaces both its value and its expiry.

Write a class `TTLCache`:

- `TTLCache()` makes an empty cache.
- `put(key, value, ttl, time)` stores the value.
- `get(key, time)` returns the value if it has not expired, or −1.
- `count(time)` returns how many keys have a value that has not expired.

Times never go down from one call to the next.

`get` only needs a hash map. For `count`, keep a min-heap of `(expiry, key)` and, before answering, pop every entry that has expired. An entry is stale if the key was stored again with a different expiry: check the map before removing the key.

Each test is a list of calls. The result line has one entry per call: the returned value, and `null` for the constructor and for calls that return nothing.

## Constraints

- 0 ≤ key ≤ 10⁴
- 0 ≤ value ≤ 10⁵
- 1 ≤ ttl ≤ 10⁶
- 0 ≤ time ≤ 10⁹
- Up to 5000 calls in one test

## Follow-up

The cache now holds ten million keys and most are never read again. Popping expired keys only during `count` lets memory grow. How would you clean up in the background without slowing down `get` and `put`?
