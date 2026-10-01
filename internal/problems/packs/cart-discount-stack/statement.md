A shop's checkout page lets a customer enter coupons one after another, and take the latest one back. Prices are in cents. Each entry is `[type, value]`:

- `[1, p]` is a coupon for `p` percent off.
- `[2, f]` is a coupon for `f` cents off.
- `[0, 0]` removes the most recent coupon that is still in force. If there is none, it does nothing.

When the customer pays, the coupons still in force apply in the order they were entered. A percent coupon takes `p` percent of the running total, rounded down to a whole cent, off the total. A flat coupon subtracts `f`, but the total never goes below 0. Return the final total.

A stack fits the "take the latest one back" rule: push each coupon and pop on a removal. Then go through what is left from the bottom up. The total can pass the 32-bit range.

## Constraints

- 1 ≤ number of items ≤ 10⁴
- 1 ≤ price ≤ 10⁶
- 0 ≤ number of entries ≤ 20
- 1 ≤ p ≤ 100
- 1 ≤ f ≤ 10⁹

## Follow-up

The shop wants to show the customer the best order to apply the same coupons in. Is there a rule that always gives the lowest total, and does it still hold once totals cannot go below zero?
