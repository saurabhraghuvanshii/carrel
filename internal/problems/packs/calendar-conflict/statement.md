A calendar app warns you before you save a meeting that clashes with others. Each meeting is `[start, end)`: it runs from `start` up to, but not including, `end`. So a meeting that ends at 10 does not clash with one that starts at 10.

Given the meetings already in the calendar (they may overlap each other) and a new meeting `[start, end)`, return how many existing meetings the new one clashes with.

Two meetings clash when each starts before the other ends. Write that as one condition and count the meetings that meet it.

## Constraints

- 1 ≤ n ≤ 10⁴
- 0 ≤ start < end ≤ 10⁹ for every meeting

## Follow-up

The calendar holds ten years of meetings and the check runs on every keystroke while the user drags the new meeting around. How would you store the meetings so that each check looks at only a few of them?
