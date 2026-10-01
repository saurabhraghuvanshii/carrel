A printer has jobs of different sizes and `h` hours. Running at speed `k`, it prints `k` pages an hour, but it works on one job at a time: if a job has fewer than `k` pages left, it finishes that job and waits for the next hour. Return the smallest whole speed `k` that finishes every job within `h` hours.

The time needed only goes down as `k` goes up, so binary search over `k`. Add up the hours in a 64-bit integer: they can pass 2³¹.

## Constraints

- 1 ≤ number of jobs ≤ 10⁴
- 1 ≤ job size ≤ 10⁹
- number of jobs ≤ h ≤ 10⁹
