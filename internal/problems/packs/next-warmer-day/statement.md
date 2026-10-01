You are given daily temperatures. For each day, return how many days you must wait for a warmer day, or 0 if none comes. An equal temperature is not warmer.

Keep a stack of days that are still waiting. Each new day answers every waiting day that is colder.

## Constraints

- 1 ≤ number of days ≤ 10⁴
- −50 ≤ temperature ≤ 60
