Fixed windows let a user send `2 × limit` requests in a short burst around a window edge. A sliding window avoids that: a request at `time` is allowed only if the user had fewer than `limit` allowed requests in the last `window` units, that is in `(time − window, time]`. A refused request does not count.

Write a class `SlidingLimiter`:

- `SlidingLimiter(limit, window)` sets the limit and the window length.
- `allow(user, time)` answers whether this request is allowed, and records it if it is.

Times never go down from one call to the next. Users do not affect each other.

Keep a queue of allowed times for each user. Before deciding, drop times from the front that are `time − window` or older; the queue then holds exactly the requests that count.

Each test is a list of calls. The result line has one entry per call: the returned value, and `null` for the constructor and for calls that return nothing.

## Constraints

- 1 ≤ limit ≤ 100
- 1 ≤ window ≤ 10⁶
- 0 ≤ user ≤ 10⁵
- 0 ≤ time ≤ 10⁹
- Up to 5000 calls in one test

## Follow-up

With a limit of a million requests per hour, a queue per user holds a million times. How could you answer with a few counters per user instead, and how wrong could that answer be?
