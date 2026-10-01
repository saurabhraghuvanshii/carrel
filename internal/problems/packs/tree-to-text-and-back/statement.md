Write a class `Codec` with two methods:

- `serialize(root)` turns a binary tree into text, in any format you choose.
- `deserialize(text)` turns that text back into a new tree with the same shape and values, and returns its root.

The tests serialize a tree, deserialize the text with a fresh `Codec`, and show the tree that comes back. It must match the original. Values can be negative and can repeat, and the tree can be empty.

## Constraints

- 0 ≤ number of nodes ≤ 10⁴
- −1000 ≤ value ≤ 1000

## Follow-up

Level order with a marker for missing children works. So does pre-order with markers. Which is easier to read back without recursion?
