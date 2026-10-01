Balloons hang over a straight line; each covers `[start, end]`, including both ends. An arrow shot at position `x` bursts every balloon with `start ≤ x ≤ end`. Return the fewest arrows that burst every balloon.

Sort by end. Shoot at the end of the first balloon still whole: that arrow bursts as many later balloons as any arrow could.

## Constraints

- 1 ≤ number of balloons ≤ 10⁴
- 0 ≤ start ≤ end ≤ 10⁶
