A grocery chain has stores along one long road, and each customer lives somewhere on that road. Positions are distances in metres from the start of the road. For every customer, return the distance to the store nearest to them.

Sort the stores once. For each customer, binary search for the first store at or after their position; the nearest store is either that one or the one just before it.

## Constraints

- 1 ≤ n, m ≤ 10⁴
- 0 ≤ position ≤ 10⁹
- Stores may share a position, and the stores are not given in order

## Follow-up

Real stores sit on a map, not on one road. Sorting by one coordinate no longer finds the nearest. How would you organise the stores so that a lookup still checks only a few of them?
