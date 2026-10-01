An analytics job receives click events as `[user, time]`. The events are not in time order, because they come from many servers. Sort each user's events by time: a new session starts at the user's first event, and again whenever the gap since that user's previous event is more than `timeout`. A gap of exactly `timeout` stays in the same session.

Return two numbers: the total number of sessions over all users, and the length of the longest session, measured from its first event to its last (a session with one event has length 0).

Group the times by user in a hash map, sort each list, and walk it once.

## Constraints

- 1 ≤ n ≤ 10⁴
- 0 ≤ timeout ≤ 3600
- 0 ≤ user ≤ 10⁵
- 0 ≤ time ≤ 10⁹

## Follow-up

Events now arrive as a live stream, mostly in order but sometimes a few seconds late. When can you safely say a session is finished, and what do you do with an event that arrives after that?
