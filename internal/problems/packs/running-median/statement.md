Build a class `MedianFinder` that receives numbers one at a time and can report the median of everything received so far: the middle value, or the average of the two middle values when the count is even.

- `MedianFinder()` starts empty.
- `add(x)` receives the number `x`.
- `median()` returns the current median. It is only called after at least one `add`.

Each test is a list of calls. The result line has one entry per call: the median with one decimal place, and `null` for the others.

Keep the smaller half in a max-heap and the larger half in a min-heap, with sizes that differ by at most one. The median is then on top.

## Constraints

- −10⁵ ≤ x ≤ 10⁵
- Up to 5000 calls in one test
