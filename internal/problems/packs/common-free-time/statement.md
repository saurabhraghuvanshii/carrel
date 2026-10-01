You are given the busy times of `k` people. Each person's list holds intervals `[start, end]` sorted by start, with gaps between them. Return the free times shared by everyone: the gaps of positive length, between the earliest start and the latest end, when nobody is busy. Return them sorted. Two busy intervals that touch leave no gap.

Merging all the busy intervals gives the answer as the spaces between them. Each list is already sorted, so a heap can walk the lists together, as in "Merge k sorted lists".

## Constraints

- 1 ≤ k ≤ 50
- 1 ≤ total number of intervals ≤ 10⁴
- 0 ≤ start < end ≤ 10⁸
