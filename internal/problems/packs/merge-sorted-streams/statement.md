You are given `k` lists of whole numbers, each sorted from smallest to largest; some may be empty. Return one sorted list with all their values.

Merging the lists two at a time, as in "Merge two sorted linked lists", repeats a lot of work when `k` is large. Keep a min-heap holding the next value of each list: take the smallest, then add the next value from the same list. That is O(N log k) for N values in total.

## Constraints

- 0 ≤ k ≤ 1000
- 0 ≤ total number of values ≤ 10⁴
- −10⁴ ≤ value ≤ 10⁴
