A live dashboard shows how many different people visited a site recently. Write a class `VisitorCounter`:

- `VisitorCounter(window)` sets the length of "recently".
- `visit(user, time)` records a visit.
- `unique(time)` returns how many different users have at least one visit after `time − window` and up to `time`.

Times never go down from one call to the next.

Keep a queue of visits in time order and a hash map from user to the time of their latest visit. Before answering, pop visits that are too old from the front of the queue. A popped visit removes its user from the map only if it is that user's latest visit; otherwise the user came back later and still counts.

Each test is a list of calls. The result line has one entry per call: the returned value, and `null` for the constructor and for calls that return nothing.

## Constraints

- 1 ≤ window ≤ 10⁶
- 0 ≤ user ≤ 10⁵
- 0 ≤ time ≤ 10⁹
- Up to 5000 calls in one test

## Follow-up

The site has a hundred million visitors a day, and an answer that is off by one or two percent is fine. What could you store instead of every user, and how would you still let old visits fall out of the window?
