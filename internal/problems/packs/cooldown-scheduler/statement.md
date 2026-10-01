A machine runs tasks, one per time unit. Each task is a capital letter, and the same letter can repeat. After running a task, the machine must wait at least `n` time units before running that same letter again; it may run other tasks or stay idle meanwhile. Tasks can run in any order. Return the fewest time units needed to run them all.

You can simulate with a max-heap of remaining counts. There is also a counting shortcut: the most frequent letter decides how many idle gaps there must be.

## Constraints

- 1 ≤ number of tasks ≤ 10⁴
- 0 ≤ n ≤ 10
- Only `A` to `Z`
