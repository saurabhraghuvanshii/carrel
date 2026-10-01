A shuttle drives along a fixed route with stops 0, 1, …, `stops`. It has `seats` seats. Each ride request `[from, to)` wants one seat from stop `from` until it gets off at stop `to`; someone getting off at a stop frees the seat for someone getting on there. Choose which requests to accept so that no more than `seats` riders are on board between any two stops, and return the largest number you can accept.

Go through the requests in order of drop-off stop, earliest first, and accept each one that still fits. Riders who get off early leave the most room for the rest. To check whether a request fits, keep the number on board for each stretch between two stops.

## Constraints

- 1 ≤ stops ≤ 1000
- 1 ≤ seats ≤ 50
- 1 ≤ n ≤ 2000
- 0 ≤ from < to ≤ stops

## Follow-up

Requests now arrive one at a time while people wait for an answer, so you must accept or refuse each one at once, without seeing the later ones. How good can any rule be, and what would you do differently with several shuttles on the same route?
