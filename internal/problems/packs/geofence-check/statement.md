A delivery app sends the customer a message each time the courier enters the zone around their home. The zone is a rectangle `[x1, y1, x2, y2]` with corners `(x1, y1)` and `(x2, y2)`; a point on the edge counts as inside. The courier's phone reports a list of positions in order.

Return how many times the courier enters the zone: the number of reported positions that are inside when the position before was outside. If the first position is inside, that counts as entering too.

Work out inside or outside for each position and compare it with the one before. A run of positions inside the zone is one visit, however long it is.

## Constraints

- 1 ≤ n ≤ 10⁴
- −10⁶ ≤ x, y ≤ 10⁶
- x1 ≤ x2 and y1 ≤ y2

## Follow-up

Phone positions jump around by a few metres, so a courier standing at the edge seems to enter and leave again and again. How would you stop the customer getting a message each time?
