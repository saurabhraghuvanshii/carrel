A social app shows which tags are trending: not the most used tags, but the ones growing fastest. You are given the posts as `[time, tag]` in time order. Let `T` be the time of the last post. The current window covers the times after `T − window` up to `T`; the previous window is the `window` time units just before it, after `T − 2 × window` up to `T − window`. Older posts do not count.

A tag's growth is its number of posts in the current window minus its number in the previous window. Return up to `k` tags whose growth is above zero, the largest growth first; on a tie the smaller tag comes first. Return an empty list if no tag is growing.

Go through the posts once and keep one number per tag in a hash map: add 1 for a post in the current window and subtract 1 for a post in the previous one. Then sort the tags that ended above zero.

## Constraints

- 1 ≤ n ≤ 10⁴
- 1 ≤ window ≤ 10⁶
- 1 ≤ k ≤ 10
- 0 ≤ time ≤ 10⁹, and times never go down
- 1 ≤ tag ≤ 10⁵

## Follow-up

The list should update live as posts arrive, without going through every post again. What do you keep per tag, and what has to happen when a post slides from the current window into the previous one?
