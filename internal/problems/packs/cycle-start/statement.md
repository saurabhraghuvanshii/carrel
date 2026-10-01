You are given the first node of a linked list that may loop, as in "Does the list loop?". Return the node where the loop begins, or `null` (`nullptr` in C++) if there is no loop. The answer is shown as that node's position, counting from 0, or −1.

After the fast and slow pointers meet, move one of them back to the first node. Step both one node at a time: they meet again exactly where the loop begins.

## Constraints

- 0 ≤ length ≤ 10⁴
- −1000 ≤ value ≤ 1000
- Do not change the list
