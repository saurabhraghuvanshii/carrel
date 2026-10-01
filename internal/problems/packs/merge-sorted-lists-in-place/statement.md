You are given two lists of whole numbers, each sorted from smallest to largest. The first list has room at its end: it holds `m` values followed by `n` empty slots, where `n` is the length of the second list. The empty slots hold 0.

Merge the second list into the first so the first list holds all `m + n` values in sorted order. Change the first list; return nothing.

## Constraints

- 0 ≤ m, n ≤ 10⁴, and m + n ≥ 1
- −10⁹ ≤ value ≤ 10⁹
- Both lists are sorted from smallest to largest

## Follow-up

Filling the first list from the back means you never overwrite a value you still need. Why?
