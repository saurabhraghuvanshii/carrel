Drivers work shifts `[start, end)`: on shift from `start`, off at `end`. The city needs at least `m` drivers on shift at every moment of the day `[0, dayEnd)`. Return every stretch of the day with fewer than `m` drivers on shift, as `[from, to)` pairs in time order. Make each stretch as long as possible, so two returned stretches never touch.

Turn each shift into two events, +1 at its start and −1 at its end, and sort them by time. Between one event time and the next the number on shift does not change, so walk the events and record the stretches where it is below `m`. A shift that ends at the moment another starts leaves no gap.

## Constraints

- 1 ≤ n ≤ 10⁴
- 1 ≤ m ≤ 10
- 0 ≤ start < end ≤ dayEnd ≤ 10⁹

## Follow-up

Shifts are added and cancelled all through the day, and a manager wants the understaffed list after every change. How would you avoid starting the sweep from scratch each time?
