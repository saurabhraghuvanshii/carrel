You have coins of a few different values and as many of each as you want. Return the fewest coins that add up to `amount` exactly, or −1 if no mix of coins does. An amount of 0 needs no coins.

Always taking the largest coin that fits does not work for every set of coins. Instead work out the best answer for every amount from 0 up to `amount`: the best for `v` is one more than the best for `v − coin`, over every coin that fits.

## Constraints

- 1 ≤ number of coins ≤ 12
- 1 ≤ coin ≤ 400, and the coins are all different
- 0 ≤ amount ≤ 10⁴
