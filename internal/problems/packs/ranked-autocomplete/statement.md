A search box suggests the most popular completions of what the user has typed. Write a class `Autocomplete`:

- `Autocomplete()` starts with no words.
- `add(word, score)` adds `score` to the word's total. The first `add` of a word creates it.
- `suggest(prefix)` returns up to three words that start with `prefix`, the highest total first; on a tie the word that comes first alphabetically wins. A word starts with itself. It returns an empty list if nothing matches.

`suggest` runs on every keystroke, so it should not look at every stored word. Build a trie and keep, at every node, the best three words below it. An `add` changes one word's total, so it only has to repair the lists on that word's own path: put the word in each list if it is missing, sort the (at most four) entries, and keep three. That works because totals only ever go up.

Each test is a list of calls. The result line has one entry per call: the returned value, and `null` for the constructor and for calls that return nothing. A `suggest` result is printed as a list such as `["car", "care"]`.

## Constraints

- 1 ≤ word and prefix length ≤ 10
- Only `a` to `z`
- 1 ≤ score ≤ 1000
- Up to 5000 calls in one test

## Follow-up

Scores should fade: a search that was popular last year should lose to one that is popular this week. Why does that break the "totals only go up" shortcut, and what would you do instead?
