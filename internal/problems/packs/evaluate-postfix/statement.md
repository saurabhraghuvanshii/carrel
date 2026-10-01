You are given an arithmetic expression in postfix order: each operator comes after its two operands, so `3 4 +` means 3 + 4 and `5 1 2 + *` means 5 × (1 + 2). The tokens are whole numbers and the operators `+`, `-`, `*` and `/`. Return the value of the expression.

Division drops the fraction, rounding toward zero, so `7 -2 /` is −3. A token is an operator only if it is exactly one of those four characters; `-8` is a number.

## Constraints

- 1 ≤ number of tokens ≤ 3000
- −100 ≤ number ≤ 100
- The expression is valid, never divides by zero, and every step stays within ±10⁸
