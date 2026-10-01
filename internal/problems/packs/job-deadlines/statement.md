Each job takes one unit of time and pays its profit if it finishes by its deadline. Time starts at 0, so a job with deadline `d` must be done in one of the slots 1 to `d`. You can do one job per slot. Return the largest total profit.

Take the jobs from the richest down. Put each one in the latest free slot at or before its deadline, so the early slots stay free for jobs that need them. Skip a job when no such slot is free.

## Constraints

- 1 ≤ number of jobs ≤ 10⁴
- 1 ≤ deadline ≤ number of jobs
- 1 ≤ profit ≤ 1000
