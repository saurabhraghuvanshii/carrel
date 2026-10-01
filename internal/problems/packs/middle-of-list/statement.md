You are given the first node of a linked list. Return its middle node. If the list has an even number of nodes, return the second of the two middle ones. The answer is shown as the list from that node to the end.

Move one pointer one step at a time and another two steps at a time. When the fast one reaches the end, the slow one is in the middle.

## Constraints

- 1 ≤ length ≤ 10⁴
- −1000 ≤ value ≤ 1000
