package testgen

import (
	"math/rand"
	"slices"
	"strconv"
	"strings"
)

// Generators for the real-interview group "Search and text".

func init() {
	register("suggestion-counts", genSuggestionCounts)
	register("near-matches", genNearMatches)
	register("trending-tags", genTrendingTags)
	register("autocomplete", genAutocomplete)
	register("merge-contacts", genMergeContacts)
}

// suggestion-counts: "n q", the n stored searches, the q prefixes. The
// answer is, for every prefix, how many stored searches start with it.
func genSuggestionCounts(rng *rand.Rand, size int) Case {
	letters := "abc"[:2+rng.Intn(2)]
	n, q := max(1, size), 1+rng.Intn(max(1, size))
	root := newTrieNode()
	words := make([]string, n)
	for i := range words {
		words[i] = randomWord(rng, 1+rng.Intn(6), letters)
		root.insert(words[i])
	}
	prefixes := make([]string, q)
	counts := make([]int, q)
	for i := range prefixes {
		prefixes[i] = randomWord(rng, 1+rng.Intn(4), letters)
		if rng.Intn(2) == 0 {
			w := words[rng.Intn(n)]
			prefixes[i] = w[:1+rng.Intn(len(w))]
		}
		if node := root.find(prefixes[i]); node != nil {
			counts[i] = node.pass
		}
	}
	return Case{Input: ints(n, q) + "\n" + strings.Join(words, " ") + "\n" + strings.Join(prefixes, " "), Expected: listOutput(counts)}
}

// near-matches: "n q", n different words, q queries. The answer is, for
// every query, how many words have its length and differ from it in at most
// one place.
func genNearMatches(rng *rand.Rand, size int) Case {
	letters := "abcd"[:2+rng.Intn(3)]
	set := map[string]bool{}
	var words []string
	for tries := 0; tries < 3*max(1, size) && len(words) < max(1, size); tries++ {
		w := randomWord(rng, 1+rng.Intn(5), letters)
		if !set[w] {
			set[w] = true
			words = append(words, w)
		}
	}
	// masked[p] counts the words that match p, where one letter is a '*'.
	masked := map[string]int{}
	mask := func(w string, i int) string { return w[:i] + "*" + w[i+1:] }
	for _, w := range words {
		for i := range w {
			masked[mask(w, i)]++
		}
	}
	q := 1 + rng.Intn(max(1, size))
	queries := make([]string, q)
	counts := make([]int, q)
	for k := range queries {
		w := []byte(words[rng.Intn(len(words))])
		if rng.Intn(3) != 0 {
			w[rng.Intn(len(w))] = letters[rng.Intn(len(letters))]
		}
		if rng.Intn(5) == 0 {
			w = []byte(randomWord(rng, 1+rng.Intn(5), letters))
		}
		queries[k] = string(w)
		for i := range w {
			counts[k] += masked[mask(queries[k], i)]
		}
		if set[queries[k]] { // an exact match was counted once per position
			counts[k] -= len(w) - 1
		}
	}
	return Case{Input: ints(len(words), q) + "\n" + strings.Join(words, " ") + "\n" + strings.Join(queries, " "), Expected: listOutput(counts)}
}

// trending-tags: "n window k", then posts "time tag" in time order. With T
// the last time, a tag's growth is its posts in (T - window, T] minus its
// posts in (T - 2*window, T - window]. The answer is up to k tags with
// growth above zero, largest growth first, the smaller tag on a tie.
func genTrendingTags(rng *rand.Rand, size int) Case {
	n := max(1, size)
	window, k := 1+rng.Intn(20), 1+rng.Intn(4)
	tags := 1 + rng.Intn([]int{3, 10}[rng.Intn(2)])
	rows := make([][]int, n)
	t := 0
	for i := range rows {
		t += rng.Intn(3) * rng.Intn(3)
		rows[i] = []int{t, 1 + rng.Intn(tags)}
	}
	growth := map[int]int{}
	for _, r := range rows {
		switch {
		case r[0] > t-window:
			growth[r[1]]++
		case r[0] > t-2*window:
			growth[r[1]]--
		}
	}
	var up []int
	for tag, g := range growth {
		if g > 0 {
			up = append(up, tag)
		}
	}
	slices.SortFunc(up, func(a, b int) int {
		if growth[a] != growth[b] {
			return growth[b] - growth[a]
		}
		return a - b
	})
	return Case{Input: ints(n, window, k) + rowsText(rows), Expected: listOutput(up[:min(k, len(up))])}
}

// autocomplete: calls new; add word score; suggest prefix. add raises the
// word's total score. suggest gives up to three words that start with the
// prefix, highest total first, alphabetical on a tie.
func genAutocomplete(rng *rand.Rand, size int) Case {
	letters := "abc"[:2+rng.Intn(2)]
	pool := make([]string, 4+rng.Intn(10))
	for i := range pool {
		pool[i] = randomWord(rng, 1+rng.Intn(5), letters)
	}
	calls := []string{"new"}
	results := []string{"null"}
	total := map[string]int{}
	for len(calls) < callCount(size) {
		lastCall := len(calls) == callCount(size)-1 // end on a question
		if rng.Intn(9) < 5 && !(lastCall && len(total) > 0) {
			w, score := pool[rng.Intn(len(pool))], 1+rng.Intn(5)
			total[w] += score
			calls, results = append(calls, "add "+w+" "+strconv.Itoa(score)), append(results, "null")
			continue
		}
		w := pool[rng.Intn(len(pool))]
		prefix := w[:1+rng.Intn(len(w))]
		if lastCall {
			prefix = prefix[:1]
		} else if rng.Intn(6) == 0 {
			prefix = randomWord(rng, 1+rng.Intn(3), letters)
		}
		var match []string
		for word := range total {
			if strings.HasPrefix(word, prefix) {
				match = append(match, word)
			}
		}
		slices.SortFunc(match, func(a, b string) int {
			if total[a] != total[b] {
				return total[b] - total[a]
			}
			return strings.Compare(a, b)
		})
		calls, results = append(calls, "suggest "+prefix), append(results, quoted(match[:min(3, len(match))]))
	}
	return callCase(calls, results)
}

// merge-contacts: "n", then n lines "k email ... email". Contacts that
// share an email are the same person, and so on through chains. The answer
// is, for every contact, the lowest contact number of its person.
func genMergeContacts(rng *rand.Rand, size int) Case {
	n := max(1, min(size, 3000))
	pool := max(2, n*[]int{1, 3, 8}[rng.Intn(3)])
	sets := newSets(n)
	first := map[int]int{}
	var b strings.Builder
	b.WriteString(strconv.Itoa(n))
	for i := 0; i < n; i++ {
		ids := distinct(rng, 1+rng.Intn(min(3, pool)), 0, pool-1)
		b.WriteString("\n" + strconv.Itoa(len(ids)))
		for _, id := range ids {
			b.WriteString(" u" + strconv.Itoa(id) + "@m.co")
			if j, ok := first[id]; ok {
				sets.union(i, j)
			} else {
				first[id] = i
			}
		}
	}
	lowest := map[int]int{}
	owners := make([]int, n)
	for i := range owners {
		root := sets.find(i)
		if _, ok := lowest[root]; !ok {
			lowest[root] = i
		}
		owners[i] = lowest[root]
	}
	return Case{Input: b.String(), Expected: listOutput(owners)}
}
