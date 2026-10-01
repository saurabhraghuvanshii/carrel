A garage has `small` small spots, numbered 0 to `small − 1`, followed by `large` large spots, numbered from `small`. The lowest numbers are nearest the entrance. Write a class `ParkingLot`:

- `ParkingLot(small, large)` starts with every spot free.
- `park(car, size)` parks the car and returns its spot. A small car (`size` 1) takes the lowest-numbered free small spot; if no small spot is free it takes the lowest-numbered free large spot. A large car (`size` 2) takes the lowest-numbered free large spot. It returns −1, and nothing changes, if there is no spot for the car or the car is already parked.
- `leave(car)` frees the car's spot and returns its number, or returns −1 if the car is not parked.
- `freeSpots()` returns the number of free spots.

Keep the free spots of each kind in a min-heap, so the lowest free one is always on top, and a hash map from car to its spot.

Each test is a list of calls. The result line has one entry per call: the returned value, and `null` for the constructor and for calls that return nothing.

## Constraints

- 0 ≤ small, large ≤ 1000, and small + large ≥ 1
- 0 ≤ car ≤ 10⁵
- size is 1 or 2
- Up to 5000 calls in one test

## Follow-up

The garage gets several floors and wants to send each car to the free spot nearest the lift it will use. What does "lowest number" turn into, and can the heaps still answer it?
