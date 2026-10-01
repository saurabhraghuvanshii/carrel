A network of `n` nodes was a tree: connected, with no loops, using `n − 1` links. Then one more link was added between two different nodes, making a loop. You are given all `n` links. Return a link you can remove to get a tree back. If several links would work, return the one that appears last in the list.

Add the links one by one to a union-find. The first link that joins two nodes already connected closes the loop, and it is the last loop link in the list.

## Constraints

- 3 ≤ n ≤ 10⁴
- The links are a tree plus one more link
