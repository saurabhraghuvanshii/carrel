An office has `rooms` meeting rooms, numbered from 0. Write a class `RoomBooker`:

- `RoomBooker(rooms)` starts with every room free.
- `book(start, end)` books the lowest-numbered room that is free for the whole of `[start, end)` and returns its number, or returns −1 if no room is free for all of it. A booking that ends at 10 does not block one that starts at 10.
- `cancel(room, start)` removes the booking in `room` that starts at `start`. It returns `true` if there was one and `false` otherwise.

Requests come in any order of time: someone may book next week before someone else books tomorrow.

For each room keep its bookings in an ordered map from start to end (`TreeMap` in Java, `map` in C++). A new booking fits when the booking just before it has ended by `start` and the booking just after it starts at `end` or later, so each room needs only two lookups.

Each test is a list of calls. The result line has one entry per call: the returned value, and `null` for the constructor and for calls that return nothing.

## Constraints

- 1 ≤ rooms ≤ 10
- 0 ≤ start < end ≤ 10⁹
- 0 ≤ room < rooms
- Up to 5000 calls in one test

## Follow-up

Rooms now have sizes, and a booking asks for a number of seats. The office wants the smallest room that fits, not the lowest number. How does the search change, and what if no single room is large enough?
