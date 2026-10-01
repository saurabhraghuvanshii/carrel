A signal is sent from node `source` in a network of `n` nodes. Each edge `[u, v, w]` is one-way: the signal takes `w` time to go from `u` to `v`. Return the time when the last node receives the signal, or −1 if some node never does.

This is the shortest time to every node, which Dijkstra's algorithm finds with a min-heap: always settle the closest node not yet settled.

## Constraints

- 1 ≤ n ≤ 10⁴
- 0 ≤ number of edges ≤ 3 × 10⁴
- 1 ≤ w ≤ 100
