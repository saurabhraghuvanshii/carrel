An online shop refunds an order only within `window` days of the purchase, and only once. You are given the purchases as `[order, day]` and the refund requests as `[order, day]` in the order they arrive. Approve a request when all of these hold:

- the order exists,
- the request day is the purchase day or later, and at most `window` days after it,
- no earlier request for the same order was approved.

Return one `true` or `false` per request, in order.

Put the purchases in a hash map from order to day so each request is one lookup, and keep a set of orders already refunded.

## Constraints

- 1 ≤ n, m ≤ 5000
- Order numbers are different and between 1 and 10⁹
- 0 ≤ day ≤ 10⁶
- 0 ≤ window ≤ 365

## Follow-up

Requests now come from the website and from the support team at the same moment, and two of them for one order can be handled by different servers. How would you make sure the order is refunded only once?
