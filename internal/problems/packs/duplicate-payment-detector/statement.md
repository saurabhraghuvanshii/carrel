A payments service wants to catch double charges, where a customer taps "Pay" twice. Payments arrive in time order as `[time, user, merchant, amount]`. A payment is a duplicate if the most recent earlier payment with the same user, the same merchant and the same amount happened at most `window` seconds before it. That earlier payment counts whether or not it was flagged itself, so a run of quick repeats is flagged all the way along.

Return the positions (starting from 0) of the duplicates, in increasing order.

Keep a hash map from the triple (user, merchant, amount) to the time it was last seen. In Java, a `Map` keyed by a `List` of the three values or by a combined string works; in C++, a `map` keyed by a `tuple`.

## Constraints

- 1 ≤ n ≤ 10⁴
- 0 ≤ window ≤ 3600
- 0 ≤ time ≤ 10⁹, and times never go down
- 0 ≤ user, merchant ≤ 10⁵
- 1 ≤ amount ≤ 10⁶

## Follow-up

The service handles millions of users, and the map of last payments grows without end. When is it safe to forget an entry, and how would you find those entries cheaply?
