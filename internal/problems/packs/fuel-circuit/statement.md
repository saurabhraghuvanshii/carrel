Fuel stations stand in a circle. Station `i` gives you `gas[i]` units, and driving from station `i` to the next one uses `cost[i]`. You start with an empty tank at some station, fill up there, and drive round once, back to that station. The tank may never go below zero. Return the smallest start index that works, or −1 if none does.

If the total gas is less than the total cost, no start works. Otherwise drive from station 0; whenever the tank would go negative, no station from the start up to here can work, so start again from the next one.

## Constraints

- 1 ≤ number of stations ≤ 10⁴, and both lists have the same length
- 0 ≤ gas[i], cost[i] ≤ 100
