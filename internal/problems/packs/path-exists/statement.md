You are given an undirected graph with `n` nodes numbered 0 to `n − 1` and a list of edges, plus two nodes `source` and `target`. Return `true` if you can walk from `source` to `target` along edges, and `false` otherwise. A node can always reach itself.

Turn the edge list into an adjacency list, then explore from `source` with a depth-first or breadth-first search, marking nodes as visited so you never loop.

## Constraints

- 1 ≤ n ≤ 10⁴
- 0 ≤ number of edges ≤ 2 × 10⁴
- Edges may repeat
