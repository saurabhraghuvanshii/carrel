A payments API limits how often each user may call it. Time is cut into fixed windows of `window` units: `[0, window)`, `[window, 2 × window)`, and so on. In each window a user may make at most `limit` requests that are allowed; a refused request does not count.

Write a class `RateLimiter`:

- `RateLimiter(limit, window)` sets the limit and the window length.
- `allow(user, time)` answers whether this request is allowed, and records it if it is.

Times never go down from one call to the next. Users do not affect each other.

Keep a hash map from user to the window they were last seen in and how many requests were allowed there. The window of a time is `time / window`, rounded down.

Each test is a list of calls. The result line has one entry per call: the returned value, and `null` for the constructor and for calls that return nothing.

## Constraints

- 1 ≤ limit ≤ 100
- 1 ≤ window ≤ 10⁶
- 0 ≤ user ≤ 10⁵
- 0 ≤ time ≤ 10⁹
- Up to 5000 calls in one test

## Follow-up

The API now runs on twenty servers behind a load balancer, and a user's requests can reach any of them. How would you keep the limit for each user, and what would you accept being slightly wrong?
