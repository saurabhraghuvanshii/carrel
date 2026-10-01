A driver leaves junction 0 at time 0 and heads for junction `n − 1`. Roads are one-way: road `(u, v, t)` takes `t` minutes from `u` to `v`. Every junction has a traffic light: a car may leave junction `v` only at a time that is a multiple of `period[v]`, so it may have to wait (a period of 1 means no waiting). Return the earliest time the driver can arrive at junction `n − 1`, or −1 if it cannot be reached.

This is Dijkstra's algorithm with one change. When you take junction `u` off the heap with arrival time `d`, the car leaves at `d` rounded up to a multiple of `period[u]`, and every road out of `u` starts from that time. Arriving earlier never makes you leave later, which is why the usual algorithm still works.

## Constraints

- 2 ≤ n ≤ 10⁴
- 0 ≤ m ≤ 4 × 10⁴
- 1 ≤ t ≤ 1000
- 1 ≤ period[v] ≤ 100
- Roads may repeat or loop back to the same junction

## Follow-up

Thousands of drivers ask for an arrival time every second, on a city graph with a million junctions. What would you precompute, and how would you speed up a single search?
