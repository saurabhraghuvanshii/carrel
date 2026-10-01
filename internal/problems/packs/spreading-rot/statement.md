A grid holds `0` (empty), `1` (a healthy plant) or `2` (an infected plant). Every minute, each infected plant infects the healthy plants directly up, down, left and right of it. Return the number of minutes until no healthy plant is left, or −1 if some plant never gets infected. If there are no healthy plants to begin with, the answer is 0.

Put every infected plant in the queue at the start, then run a breadth-first search one minute (one layer) at a time.

## Constraints

- 1 ≤ rows, cols ≤ 60
- Every cell is 0, 1 or 2
