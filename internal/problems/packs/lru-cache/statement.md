A product page service keeps the items people looked at recently in a small cache in memory. The cache holds at most `capacity` entries. When it is full and a new key comes in, it throws out the entry that was used longest ago. Reading a key counts as using it, and so does writing it.

Write a class `LRUCache` with these calls:

- `LRUCache(capacity)` makes an empty cache that holds at most `capacity` entries.
- `get(key)` returns the value stored for `key`, or −1 if the key is not in the cache.
- `put(key, value)` stores the value. If the key is already there, its value is replaced. If the cache now holds more than `capacity` entries, the least recently used one is removed.

Both calls should take constant time on average.

Each test is a list of calls. The result line has one entry per call: the value for `get`, and `null` for the constructor and `put`.

## Constraints

- 1 ≤ capacity ≤ 1000
- 0 ≤ key ≤ 10⁴
- 0 ≤ value ≤ 10⁵
- Up to 5000 calls in one test

## Follow-up

Many requests now reach the cache at the same moment from different threads. What can go wrong, and how would you keep it correct without making every call wait for every other call?
