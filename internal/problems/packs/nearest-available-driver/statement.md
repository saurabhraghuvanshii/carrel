A rider requests a pickup at a point on a city grid. Each driver has an id, a position, and a flag saying whether they are free. Return the id of the closest free driver by walking distance (rows plus columns). If two drivers tie, pick the lower id. If nobody is free, return −1.

Each driver is given as `{id, x, y, free}`, where `free` is 1 for a free driver and 0 for a busy one.

## Constraints

- Up to 10⁴ drivers
- Ids are unique
- Coordinates are non-negative whole numbers

## Follow-up

Drivers move every few seconds. What would you change so each request stays fast?
