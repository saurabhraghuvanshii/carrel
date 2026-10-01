Return aᵉ modulo 10⁹ + 7. Here 0⁰ counts as 1.

The exponent can be as large as 10¹⁸, so multiplying `a` by itself `e` times is far too slow. Look at the bits of `e` instead. Keep squaring the base (a, a², a⁴, a⁸, …), and multiply it into the answer whenever the matching bit of `e` is 1. Take the modulo after every multiplication, and use 64-bit integers so the product of two values below 10⁹ + 7 fits.

## Constraints

- 0 ≤ a ≤ 10⁹
- 0 ≤ e ≤ 10¹⁸
