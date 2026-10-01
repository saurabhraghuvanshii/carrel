An image service caches thumbnails, and popular images should stay even if nobody asked for them in the last minute. The cache holds at most `capacity` entries. Every `get` of a stored key and every `put` counts as a use of that key. When the cache is full and a new key comes in, it removes the key with the fewest uses; if several keys tie, it removes the one whose last use was longest ago.

Write a class `LFUCache`:

- `LFUCache(capacity)` makes an empty cache.
- `get(key)` returns the value, or −1 if the key is not stored. A miss is not a use.
- `put(key, value)` stores the value, replacing any old one, and removes an entry first if a new key would go over capacity.

Both calls should take constant time on average. Keep a map from key to its value and use count, and for every use count a list of keys in order of last use (a linked hash set in Java, or a list plus a map to its positions in C++). Also track the smallest use count present: a new key always has count 1.

Each test is a list of calls. The result line has one entry per call: the returned value, and `null` for the constructor and for calls that return nothing.

## Constraints

- 1 ≤ capacity ≤ 1000
- 0 ≤ key ≤ 10⁴
- 0 ≤ value ≤ 10⁵
- Up to 5000 calls in one test

## Follow-up

A key that was very popular last month keeps its high count and can never be removed, even if nobody wants it now. How would you let old uses count for less over time?
