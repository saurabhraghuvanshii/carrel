A build system runs `n` jobs, numbered from 0. Job `i` takes `durations[i]` minutes. A dependency `[a, b]` means job `a` must finish before job `b` can start. There are enough machines to run any number of jobs at the same time, so a job starts the moment its last dependency finishes; a job with no dependencies starts at time 0.

Return the time when the last job finishes. If the dependencies form a circle, so that some jobs can never start, return −1.

Process the jobs in topological order: start with the jobs that wait for nothing, and release a job when all the jobs it waits for are done. A job's start time is the largest finish time among its dependencies. If fewer than `n` jobs were ever released, there is a circle.

## Constraints

- 1 ≤ n ≤ 10⁴
- 0 ≤ m ≤ 2 × 10⁴
- 1 ≤ duration ≤ 10⁴
- Dependencies may repeat, and a job may be listed as depending on itself

## Follow-up

There are now only `k` machines, so at most `k` jobs run at once. Why is the largest-finish-time rule no longer enough, and which waiting job would you start first?
