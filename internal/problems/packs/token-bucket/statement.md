A messaging service lets each customer send in bursts but limits the average rate with a token bucket. The bucket holds at most `capacity` tokens and starts full at time 0. One token is added at every time that is a positive multiple of `refill` (`refill`, `2 × refill`, …), unless the bucket is already full. Sending `n` messages needs `n` tokens.

Write a class `TokenBucket`:

- `TokenBucket(capacity, refill)` makes a full bucket.
- `take(time, n)` removes `n` tokens and returns `true` if the bucket has at least `n` at that time; otherwise it removes nothing and returns `false`.

Times never go down from one call to the next.

Do not add tokens one time unit at a time. Remember the time of the last call: the number of new tokens is how many multiples of `refill` lie after it and up to now, which is `time / refill − last / refill` with whole-number division. Then cap at `capacity`.

Each test is a list of calls. The result line has one entry per call: the returned value, and `null` for the constructor and for calls that return nothing.

## Constraints

- 1 ≤ capacity ≤ 1000
- 1 ≤ refill ≤ 10⁶
- 1 ≤ n ≤ 2000
- 0 ≤ time ≤ 10⁹
- Up to 5000 calls in one test

## Follow-up

Customers now buy plans that refill 3 tokens every 7 seconds, not one at a time. How would you change the state you keep so the bucket still never needs a timer?
