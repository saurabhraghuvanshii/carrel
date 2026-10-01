There are `n` people. Each person either stays alone or forms a pair with exactly one other person, and nobody is in two pairs. Return the number of different ways to do this, modulo 10⁹ + 7. Only who is with whom matters, not the order of pairs.

Look at the last person. Either they stay alone, leaving `n − 1` people to arrange, or they pair with one of the other `n − 1` people, leaving `n − 2`. Keep the running counts modulo 10⁹ + 7, and use a 64-bit type for the multiplication.

## Constraints

- 1 ≤ n ≤ 10⁴
