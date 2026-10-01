Read a whole number from the start of a text, following these rules in order:

1. Skip any spaces at the start.
2. If the next character is `+` or `-`, it sets the sign. Only one sign is allowed.
3. Read digits until the first character that is not a digit, or the end. Leading zeros are fine.
4. If no digits were read, the answer is 0.
5. If the number is below −2³¹, return −2³¹; if it is above 2³¹ − 1, return 2³¹ − 1.

Anything after the digits is ignored. Check for overflow before it happens: before multiplying by 10 and adding a digit, compare with the limit (or keep the value in a 64-bit integer and stop once it passes the limit).

## Constraints

- 0 ≤ text length ≤ 200
- Letters, digits, spaces, `+`, `-` and `.`
