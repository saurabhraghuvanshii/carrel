You are given an undirected graph with `n` nodes and a list of edges. Return `true` if you can paint every node one of two colours so that every edge joins two different colours, and `false` otherwise.

Colour a node, give its neighbours the other colour, and keep going. If an edge ever joins two nodes of the same colour, it cannot be done. The graph may come in several separate pieces.

## Constraints

- 1 ≤ n ≤ 10⁴
- 0 ≤ number of edges ≤ 2 × 10⁴
