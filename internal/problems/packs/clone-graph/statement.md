You are given one node of a connected undirected graph. Each node has a value and a list of neighbours. Return a deep copy of the whole graph: new nodes with the same values and the same neighbour lists, in the same order, pointing only at new nodes.

Keep a map from each original node to its copy. When you meet a neighbour that already has a copy, link to that copy instead of making another one.

The tests build the graph from an edge list, where node `i` has value `i`, give you node 0, and print each copied node's neighbour values. Using an original node in the copy is reported.

## Constraints

- 1 ≤ number of nodes ≤ 2000
- The graph is connected and has no repeated edges or self-loops
