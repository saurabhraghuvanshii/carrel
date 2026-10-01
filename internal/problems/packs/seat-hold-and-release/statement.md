A ticket site holds seats for a customer while they type in their card details. One row has seats 0 to `seats − 1`. A hold made at time `t` lasts until `t + ttl`: from that moment on, if it was not bought or released, its seats are free again. Write a class `SeatHolder`:

- `SeatHolder(seats, ttl)` starts with every seat free.
- `hold(user, count, time)` finds the lowest-numbered block of `count` free seats side by side, holds it for the user and returns its first seat. It returns −1, and holds nothing, if there is no such block or if the user still has a hold.
- `buy(user, time)` turns the user's hold into a sale and returns `true`; the seats are sold for good. It returns `false` if the user has no hold at that time.
- `release(user, time)` gives the user's held seats back and returns `true`, or returns `false` if the user has no hold at that time.
- `freeSeats(time)` returns the number of free seats.

Times never go down from one call to the next. Expired holds must be cleared before any call answers.

Because every hold lasts the same `ttl`, holds expire in the order they were made: a queue is enough to find the expired ones. A hold in the queue may already have been bought or released, so check that it is still the user's current hold before freeing its seats. Scanning the row for a free block is fine at this size.

Each test is a list of calls. The result line has one entry per call: the returned value, and `null` for the constructor and for calls that return nothing.

## Constraints

- 1 ≤ seats ≤ 1000
- 1 ≤ ttl ≤ 10⁶
- 1 ≤ count ≤ 1000
- 0 ≤ user ≤ 10⁵
- 0 ≤ time ≤ 10⁹
- Up to 5000 calls in one test

## Follow-up

A stadium has eighty thousand seats in many rows, and thousands of people try to hold seats in the same second. How would you find a free block without scanning, and how would you stop two servers from holding the same seat?
