A building has floors 0 to `floors − 1` and one lift. The lift starts at floor 0, heading up. It keeps going in its direction as long as someone is waiting that way, and only then turns round, so nobody waits forever. Write a class `Elevator`:

- `Elevator(floors)` puts the lift at floor 0, heading up, with no calls.
- `call(floor)` adds a stop at that floor. A call for the floor the lift is at right now is ignored, and calling a floor twice is the same as calling it once.
- `next()` moves the lift to its next stop and returns that floor. The next stop is the nearest called floor in the current direction. If there is none in that direction, the lift turns round and takes the nearest called floor the other way. The stop is then removed. With no calls at all it returns −1 and the lift does not move or turn.
- `travelled()` returns the total number of floors the lift has moved.

Keep the called floors in an ordered set (`TreeSet` in Java, `set` in C++). The nearest floor above the lift and the nearest floor below it are each one lookup.

Each test is a list of calls. The result line has one entry per call: the returned value, and `null` for the constructor and for calls that return nothing.

## Constraints

- 2 ≤ floors ≤ 10⁶
- 0 ≤ floor < floors
- Up to 5000 calls in one test

## Follow-up

The building adds three more lifts. When someone calls a floor, which lift should take it? Think about lifts that are already heading that way, and about not sending every lift to the same floor.
