A debugging tool shows one timeline built from the logs of `k` servers. Each server's log is a list of times in the order they were written, so it never goes down. The servers' clocks do not agree: the true time of an entry is its logged time plus that server's `offset`, which can be negative.

Merge all entries by true time. When two entries have the same true time, the one from the lower-numbered server comes first; entries of one server keep their own order. Return, for each of the first `limit` entries of the merged timeline, the number of the server it came from.

Put the first entry of every server in a min-heap ordered by (true time, server). Pop the smallest, write down its server, and push that server's next entry. Stop after `limit` pops: the rest of the logs is never touched.

## Constraints

- 1 ≤ k ≤ 100
- 0 ≤ entries per server, and 1 ≤ total entries ≤ 5000
- 1 ≤ limit ≤ total entries
- 0 ≤ time ≤ 10⁹
- −1000 ≤ offset ≤ 1000

## Follow-up

The logs are files far too large for memory, and the tool only shows one page of the timeline at a time. How would you show page 500 without merging the 499 pages before it every time?
