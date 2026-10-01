A shop holds stock for a customer while they pay, so two people cannot buy the last item. Write a class `Inventory`:

- `Inventory()` starts with no stock.
- `add(item, qty)` puts `qty` more units of `item` on the shelf.
- `reserve(order, item, qty)` holds `qty` units for `order` and returns `true`. It returns `false` and changes nothing if the order already has a hold, or if fewer than `qty` units are available.
- `cancel(order)` ends the order's hold and puts its units back. It returns `false` if the order has no hold.
- `confirm(order)` ends the order's hold because the customer paid: the units are gone for good. It returns `false` if the order has no hold.
- `available(item)` returns the units of `item` that are neither held nor sold.

After a cancel or a confirm, the same order number may reserve again.

Two hash maps are enough: item to available units, and order to the item and quantity it holds.

Each test is a list of calls. The result line has one entry per call: the returned value, and `null` for the constructor and for calls that return nothing.

## Constraints

- 0 ≤ item, order ≤ 10⁵
- 1 ≤ qty ≤ 1000
- Up to 5000 calls in one test

## Follow-up

Customers sometimes close the page and never pay, so their holds stay forever. Add an expiry time to each hold. What must `reserve` and `available` do before they answer, and how do you find expired holds quickly?
