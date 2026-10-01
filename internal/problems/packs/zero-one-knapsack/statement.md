You have `n` items, each with a weight and a value, and a bag that holds a total weight of at most `capacity`. Each item is either packed whole or left out, and there is only one of each. Return the largest total value you can pack.

Packing by best value per unit of weight can leave space wasted. Instead keep, for every capacity from 0 up, the best value so far, and add the items one at a time. Go through the capacities from high to low so no item is packed twice.

## Constraints

- 1 ≤ n ≤ 100
- 1 ≤ weight ≤ 1000
- 1 ≤ value ≤ 1000
- 0 ≤ capacity ≤ 10⁴
