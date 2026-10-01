A courier parks the van at position `start` on a long straight street and delivers parcels on foot. Each parcel goes to a position on the same street. The courier may visit the stops in any order and finishes at the last stop, without walking back to the van. Return the shortest total distance walked.

Only the stop furthest to the left and the stop furthest to the right matter; everything between them is passed on the way. If stops lie on both sides of the van, one side is walked twice (there and back) and the other once, so walk the shorter side first.

## Constraints

- 1 ≤ n ≤ 10⁴
- 0 ≤ position, start ≤ 10⁹
- Several parcels may go to the same position

## Follow-up

Now the courier must come back to the van when it is full, because they can carry only `k` parcels at a time. Which stops would you serve together on each trip?
