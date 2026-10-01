A food app remembers what people searched for; the same search can be stored many times. As a user types, the search box shows how many stored searches begin with what they have typed so far.

Given the stored searches and a list of prefixes, return one number per prefix: how many stored searches start with it. A search starts with itself, and repeated searches count each time.

Checking every search against every prefix is fine at this size. The faster way is a trie in which every node counts the searches that pass through it: the answer for a prefix is the count at the node where the prefix ends, or 0 if the path stops early.

## Constraints

- 1 ≤ n, m ≤ 10⁴
- 1 ≤ length of a search or prefix ≤ 10
- Only `a` to `z`

## Follow-up

Now there are a hundred million stored searches and the count must appear within a few milliseconds of every keystroke. What would you build ahead of time, and how would you keep it current as new searches come in?
