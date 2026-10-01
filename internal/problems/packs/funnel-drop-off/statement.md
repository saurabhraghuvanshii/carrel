A checkout has `steps` steps, numbered from 1: cart, address, payment, and so on. The event log has one `[user, step]` entry, in time order, each time a user loads a step. A user has reached step 1 once they have an event for it. They have reached step `j` once they have an event for step `j` that comes after they reached step `j − 1`. Events that skip ahead or repeat an earlier step change nothing.

Return a list with one number per step: how many users reached it. The product team reads the drop-off from one step to the next.

Keep, for each user, the highest step reached so far. An event moves the user forward only when it is for exactly the next step.

## Constraints

- 1 ≤ n ≤ 10⁴
- 2 ≤ steps ≤ 6
- 0 ≤ user ≤ 10⁵
- 1 ≤ step ≤ steps

## Follow-up

The team now wants the funnel to count only users who finish all steps within 30 minutes of starting. What extra state does each user need, and what happens when a user starts over?
