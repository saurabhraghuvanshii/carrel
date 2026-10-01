package testgen

import (
	"math/rand"
	"slices"
	"strconv"
	"strings"
)

// Generators for the "Tries" group.

func init() {
	register("prefix-tree", genPrefixTree)
	register("unique-prefix", genUniquePrefix)
	register("buildable-word", genBuildableWord)
	register("distinct-substrings", genDistinctSubstrings)
}

type trieNode struct {
	next map[byte]*trieNode
	word bool
	pass int // words that go through this node
}

func newTrieNode() *trieNode { return &trieNode{next: map[byte]*trieNode{}} }

func (t *trieNode) insert(w string) int {
	made := 0
	cur := t
	for i := 0; i < len(w); i++ {
		n, ok := cur.next[w[i]]
		if !ok {
			n = newTrieNode()
			cur.next[w[i]] = n
			made++
		}
		n.pass++
		cur = n
	}
	cur.word = true
	return made
}

func (t *trieNode) find(w string) *trieNode {
	cur := t
	for i := 0; i < len(w) && cur != nil; i++ {
		cur = cur.next[w[i]]
	}
	return cur
}

// prefix-tree: operations "insert w", "search w" and "prefix p" on an empty
// tree. The answer lists the result of every search and prefix, in order.
func genPrefixTree(rng *rand.Rand, size int) Case {
	q := 1 + rng.Intn(max(1, size))
	letters := "abcdefghijklmnopqrstuvwxyz"[:2+rng.Intn(3)]
	root := newTrieNode()
	var added []string
	lines := []string{strconv.Itoa(q)}
	var results []string
	for k := 0; k < q; k++ {
		op := rng.Intn(3)
		if k == q-1 && op == 0 {
			op = 1 + rng.Intn(2) // end on a question so the answer is never empty
		}
		w := randomWord(rng, 1+rng.Intn(6), letters)
		if op > 0 && len(added) > 0 && rng.Intn(2) == 0 {
			w = added[rng.Intn(len(added))]
			if op == 2 || rng.Intn(3) == 0 {
				w = w[:1+rng.Intn(len(w))] // a prefix of a word, asked either way
			}
		}
		switch op {
		case 0:
			root.insert(w)
			added = append(added, w)
			lines = append(lines, "insert "+w)
		case 1:
			n := root.find(w)
			results = append(results, boolText(n != nil && n.word))
			lines = append(lines, "search "+w)
		default:
			results = append(results, boolText(root.find(w) != nil))
			lines = append(lines, "prefix "+w)
		}
	}
	return Case{Input: strings.Join(lines, "\n"), Expected: "[" + strings.Join(results, ", ") + "]"}
}

// prefixFree drops words that start another word or another word starts, and
// repeats, keeping the first of each.
func prefixFree(words []string) []string {
	var out []string
	for _, w := range words {
		ok := true
		for _, o := range out {
			ok = ok && !strings.HasPrefix(w, o) && !strings.HasPrefix(o, w)
		}
		if ok {
			out = append(out, w)
		}
	}
	return out
}

// unique-prefix: different words, none the start of another; the answer is,
// for every word in order, its shortest beginning that no other word has.
func genUniquePrefix(rng *rand.Rand, size int) Case {
	letters := "abcdefghijklmnopqrstuvwxyz"[:2+rng.Intn(25)]
	n := 1 + rng.Intn(min(1000, max(1, size)))
	stem := randomWord(rng, rng.Intn(4), letters) // shared starts make it harder
	words := make([]string, n)
	for i := range words {
		words[i] = randomWord(rng, 1+rng.Intn(10), letters)
		if rng.Intn(2) == 0 {
			words[i] = stem + words[i]
		}
	}
	words = prefixFree(words)
	root := newTrieNode()
	for _, w := range words {
		root.insert(w)
	}
	out := make([]string, len(words))
	for i, w := range words {
		cur := root
		for k := 0; k < len(w); k++ {
			cur = cur.next[w[k]]
			if cur.pass == 1 {
				out[i] = w[:k+1]
				break
			}
		}
	}
	return Case{Input: strconv.Itoa(len(words)) + "\n" + strings.Join(words, " "), Expected: quoted(out)}
}

// buildable-word: the answer is the longest word whose every beginning is
// also in the list, the alphabetically first of equal lengths, or "" if none.
func genBuildableWord(rng *rand.Rand, size int) Case {
	letters := "abcdefghijklmnopqrstuvwxyz"[:2+rng.Intn(5)]
	set := map[string]bool{}
	for k := 1 + rng.Intn(min(300, max(1, size))); k > 0; k-- {
		w := randomWord(rng, 1+rng.Intn(8), letters)
		set[w] = true
		if rng.Intn(2) == 0 { // add most beginnings, so long chains appear
			for p := 1; p < len(w); p++ {
				if rng.Intn(6) != 0 {
					set[w[:p]] = true
				}
			}
		}
	}
	words := make([]string, 0, len(set))
	for w := range set {
		words = append(words, w)
	}
	slices.Sort(words)
	rng.Shuffle(len(words), func(i, j int) { words[i], words[j] = words[j], words[i] })
	root := newTrieNode()
	for _, w := range words {
		root.insert(w)
	}
	best := ""
	var walk func(n *trieNode, path []byte)
	walk = func(n *trieNode, path []byte) {
		if len(path) > len(best) {
			best = string(path)
		}
		for c := byte('a'); c <= 'z'; c++ {
			if next, ok := n.next[c]; ok && next.word {
				walk(next, append(path, c))
			}
		}
	}
	walk(root, nil)
	return Case{Input: strconv.Itoa(len(words)) + "\n" + strings.Join(words, " "), Expected: `"` + best + `"`}
}

// distinct-substrings: the answer is how many different non-empty stretches
// the word has. Every new node of a trie of all suffixes is one of them.
func genDistinctSubstrings(rng *rand.Rand, size int) Case {
	w := randomWord(rng, 1+rng.Intn(min(500, max(1, size))), "abcdefghijklmnopqrstuvwxyz"[:1+rng.Intn([]int{2, 4, 26}[rng.Intn(3)])])
	root := newTrieNode()
	count := 0
	for i := range w {
		count += root.insert(w[i:])
	}
	return Case{Input: w, Expected: strconv.Itoa(count)}
}
