You are given a connected undirected graph with `n` nodes, where each edge `[u, v, w]` costs `w`. Choose edges so that every node is connected to every other, at the smallest total cost, and return that cost.

Grow a connected set from node 0: always add the cheapest edge that reaches a new node (a min-heap of edges does this). Or sort all edges and add each one that joins two separate groups.

## Constraints

- 1 ≤ n ≤ 10⁴
- 0 ≤ number of edges ≤ 3 × 10⁴
- 1 ≤ w ≤ 1000
- The graph is connected; edges may repeat with different costs
