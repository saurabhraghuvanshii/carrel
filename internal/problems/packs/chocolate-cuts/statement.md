A chocolate bar has `rows × cols` squares. You must break it into single squares, cutting along the lines between them. Every cut goes straight across one piece. Cutting along horizontal line `i` costs `horizontal[i]` each time, and vertical line `j` costs `vertical[j]` each time; a line must be cut once in every piece it passes through. Return the cheapest total.

A cut is repeated once for every piece the other way's earlier cuts have made. So make the expensive cuts first, while there are few pieces: always take the most expensive line left. The total can exceed a 32-bit integer.

## Constraints

- 2 ≤ rows, cols ≤ 1000
- 1 ≤ cost ≤ 10⁴
