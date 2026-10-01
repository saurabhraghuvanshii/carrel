You are given the first node of a linked list. Normally the last node links to nothing, but here it may link back to an earlier node, making the list loop forever. Return `true` if the list loops and `false` otherwise.

The tests describe a loop by the position the last node points back to, or −1 for no loop; your code only gets the first node. Use a fixed amount of extra memory: a fast pointer that moves two steps at a time meets a slow one if, and only if, there is a loop.

## Constraints

- 0 ≤ length ≤ 10⁴
- −1000 ≤ value ≤ 1000
