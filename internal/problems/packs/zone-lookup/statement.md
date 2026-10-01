A delivery company prices each parcel by zone. Zones are defined by ranges of postal codes: `[lo, hi, zone]` means every code from `lo` to `hi`, both included, belongs to that zone. The ranges do not overlap, but they are not given in order, there can be gaps between them, and several ranges may belong to the same zone.

For each postal code in a list, return its zone, or −1 if no range covers it.

Sort the ranges by `lo`. For a code, binary search for the last range whose `lo` is at or below the code; the code is covered only if that range's `hi` is at or above it.

## Constraints

- 1 ≤ n, m ≤ 10⁴
- 0 ≤ lo ≤ hi ≤ 10⁹
- 1 ≤ zone ≤ 1000
- 0 ≤ code ≤ 10⁹

## Follow-up

The pricing team now adds overrides: a small range inside a big one, which should win. How would you answer lookups when ranges can sit inside each other?
