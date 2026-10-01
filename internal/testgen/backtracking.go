package testgen

import (
	"math/rand"
	"slices"
	"strconv"
	"strings"
)

// Generators for the "Backtracking" group. Lists of answers print sorted, so
// a learner may return them in any order.

func init() {
	register("no-adjacent-ones", genNoAdjacentOnes)
	register("all-subsets", genAllSubsets)
	register("all-orderings", genAllOrderings)
	register("combo-sum", genComboSum)
	register("phone-letters", genPhoneLetters)
	register("bracket-strings", genBracketStrings)
	register("palindrome-splits", genPalindromeSplits)
	register("grid-word", genGridWord)
	register("queens-holes", genQueensHoles)
}

func quoted(words []string) string {
	q := make([]string, len(words))
	for i, w := range words {
		q[i] = `"` + w + `"`
	}
	return "[" + strings.Join(q, ", ") + "]"
}

// stringsOutput sorts the words and prints ["a", "b"].
func stringsOutput(words []string) string {
	s := slices.Clone(words)
	slices.Sort(s)
	return quoted(s)
}

// splitsOutput sorts lists of words and prints [["a", "b"], ["ab"]].
func splitsOutput(lists [][]string) string {
	s := slices.Clone(lists)
	slices.SortFunc(s, func(a, b []string) int { return slices.Compare(a, b) })
	parts := make([]string, len(s))
	for i, l := range s {
		parts[i] = quoted(l)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

// no-adjacent-ones: n; the answer is every binary string of length n with
// no two 1s side by side.
func genNoAdjacentOnes(rng *rand.Rand, size int) Case {
	n := 1 + rng.Intn(min(16, max(1, size)))
	var out []string
	var walk func(s string)
	walk = func(s string) {
		if len(s) == n {
			out = append(out, s)
			return
		}
		walk(s + "0")
		if !strings.HasSuffix(s, "1") {
			walk(s + "1")
		}
	}
	walk("")
	return Case{Input: strconv.Itoa(n), Expected: stringsOutput(out)}
}

func distinctLetters(rng *rand.Rand, n int) string {
	return string([]byte(shuffled(rng, "abcdefghijklmnopqrstuvwxyz"))[:n])
}

// all-subsets: a word of different letters; the answer is every subset, its
// letters in word order.
func genAllSubsets(rng *rand.Rand, size int) Case {
	word := distinctLetters(rng, 1+rng.Intn(min(10, max(1, size))))
	var out []string
	for mask := 0; mask < 1<<len(word); mask++ {
		var b []byte
		for i := 0; i < len(word); i++ {
			if mask>>i&1 == 1 {
				b = append(b, word[i])
			}
		}
		out = append(out, string(b))
	}
	return Case{Input: word, Expected: stringsOutput(out)}
}

// all-orderings: a word of different letters; the answer is every ordering.
func genAllOrderings(rng *rand.Rand, size int) Case {
	word := distinctLetters(rng, 1+rng.Intn(min(7, max(1, size))))
	var out []string
	var walk func(rest, so string)
	walk = func(rest, so string) {
		if rest == "" {
			out = append(out, so)
			return
		}
		for i := 0; i < len(rest); i++ {
			walk(rest[:i]+rest[i+1:], so+rest[i:i+1])
		}
	}
	walk(word, "")
	return Case{Input: word, Expected: stringsOutput(out)}
}

// combos lists every way to reach target with the candidates, each used any
// number of times, each combination ascending.
func combos(cands []int, target int) [][]int {
	s := slices.Clone(cands)
	slices.Sort(s)
	var out [][]int
	var walk func(from, left int, cur []int)
	walk = func(from, left int, cur []int) {
		if left == 0 {
			out = append(out, slices.Clone(cur))
			return
		}
		for i := from; i < len(s) && s[i] <= left; i++ {
			walk(i, left-s[i], append(cur, s[i]))
		}
	}
	walk(0, target, nil)
	return out
}

// combo-sum: different positive candidates and a target; the answer is every
// combination (reuse allowed) that adds up to it.
func genComboSum(rng *rand.Rand, size int) Case {
	for {
		n := 1 + rng.Intn(min(12, max(1, size)))
		top := []int{12, 40}[rng.Intn(2)] // small candidates make reuse matter
		lo := 1 + rng.Intn(3)
		cands := distinct(rng, min(n, top-lo+1), lo, top)
		target := 1 + rng.Intn(30)
		res := combos(cands, target)
		if len(res) > 3000 || (len(res) == 0 && rng.Intn(5) != 0) {
			continue
		}
		return Case{Input: arrayInput(cands) + "\n" + strconv.Itoa(target), Expected: gridOutput(res)}
	}
}

var keypad = map[byte]string{'2': "abc", '3': "def", '4': "ghi", '5': "jkl", '6': "mno", '7': "pqrs", '8': "tuv", '9': "wxyz"}

// phone-letters: digits 2 to 9; the answer is every word the keypad letters
// can spell.
func genPhoneLetters(rng *rand.Rand, size int) Case {
	digits := randomWord(rng, 1+rng.Intn(min(6, max(1, size))), "23456789")
	out := []string{""}
	for i := 0; i < len(digits); i++ {
		var next []string
		for _, w := range out {
			for _, ch := range keypad[digits[i]] {
				next = append(next, w+string(ch))
			}
		}
		out = next
	}
	return Case{Input: digits, Expected: stringsOutput(out)}
}

// bracket-strings: n; the answer is every balanced string of n pairs.
func genBracketStrings(rng *rand.Rand, size int) Case {
	n := 1 + rng.Intn(min(9, max(1, size)))
	var out []string
	var walk func(s string, open, close int)
	walk = func(s string, open, close int) {
		if len(s) == 2*n {
			out = append(out, s)
			return
		}
		if open < n {
			walk(s+"(", open+1, close)
		}
		if close < open {
			walk(s+")", open, close+1)
		}
	}
	walk("", 0, 0)
	return Case{Input: strconv.Itoa(n), Expected: stringsOutput(out)}
}

func isPal(s string) bool {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		if s[i] != s[j] {
			return false
		}
	}
	return true
}

// palindrome-splits: a word; the answer is every way to cut it into pieces
// that each read the same both ways.
func genPalindromeSplits(rng *rand.Rand, size int) Case {
	word := randomWord(rng, 1+rng.Intn(min(12, max(1, size))), "abcd"[:1+rng.Intn(4)])
	var out [][]string
	var walk func(rest string, cur []string)
	walk = func(rest string, cur []string) {
		if rest == "" {
			out = append(out, slices.Clone(cur))
			return
		}
		for i := 1; i <= len(rest); i++ {
			if isPal(rest[:i]) {
				walk(rest[i:], append(cur, rest[:i]))
			}
		}
	}
	walk(word, nil)
	return Case{Input: word, Expected: splitsOutput(out)}
}

// grid-word: a letter grid ("rows cols", then the rows as words) and a word;
// the answer is whether a path of side-by-side cells, each used once, spells
// it. Half the words are planted along a real path.
func genGridWord(rng *rand.Rand, size int) Case {
	r, c := 1+rng.Intn(min(6, max(1, size))), 1+rng.Intn(min(6, max(1, size)))
	letters := "abcdef"[:2+rng.Intn(5)]
	grid := make([][]byte, r)
	for i := range grid {
		grid[i] = []byte(randomWord(rng, c, letters))
	}
	word := randomWord(rng, 1+rng.Intn(min(12, r*c+1)), letters)
	if r*c > 1 && r*c <= 26 && rng.Intn(4) == 0 {
		// Every letter different, and a walk that steps forward and then back
		// onto the cell it came from: spelling it needs a cell twice.
		letters := shuffled(rng, "abcdefghijklmnopqrstuvwxyz")
		for i := range grid {
			grid[i] = []byte(letters[i*c : i*c+c])
		}
		y, x := rng.Intn(r), rng.Intn(c)
		path := []byte{grid[y][x]}
		for len(path) < 2 {
			d := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}[rng.Intn(4)]
			if ny, nx := y+d[0], x+d[1]; ny >= 0 && ny < r && nx >= 0 && nx < c {
				path = append(path, grid[ny][nx], grid[y][x]) // there and back
			}
		}
		word = string(path)
	} else if rng.Intn(2) == 0 {
		// Walk a random path without revisiting and read it off.
		y, x := rng.Intn(r), rng.Intn(c)
		used := map[[2]int]bool{{y, x}: true}
		path := []byte{grid[y][x]}
		for len(path) < 1+rng.Intn(min(12, r*c)) {
			var next [][2]int
			for _, d := range [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}} {
				ny, nx := y+d[0], x+d[1]
				if ny >= 0 && ny < r && nx >= 0 && nx < c && !used[[2]int{ny, nx}] {
					next = append(next, [2]int{ny, nx})
				}
			}
			if len(next) == 0 {
				break
			}
			p := next[rng.Intn(len(next))]
			y, x = p[0], p[1]
			used[p] = true
			path = append(path, grid[y][x])
		}
		word = string(path)
	}
	rows := make([]string, r)
	for i := range grid {
		rows[i] = string(grid[i])
	}
	return Case{
		Input:    strconv.Itoa(r) + " " + strconv.Itoa(c) + "\n" + strings.Join(rows, "\n") + "\n" + word,
		Expected: boolText(gridHasWord(rows, word)),
	}
}

func gridHasWord(rows []string, word string) bool {
	r, c := len(rows), len(rows[0])
	used := make([][]bool, r)
	for i := range used {
		used[i] = make([]bool, c)
	}
	var walk func(y, x, k int) bool
	walk = func(y, x, k int) bool {
		if y < 0 || y >= r || x < 0 || x >= c || used[y][x] || rows[y][x] != word[k] {
			return false
		}
		if k == len(word)-1 {
			return true
		}
		used[y][x] = true
		found := walk(y+1, x, k+1) || walk(y-1, x, k+1) || walk(y, x+1, k+1) || walk(y, x-1, k+1)
		used[y][x] = false
		return found
	}
	for y := 0; y < r; y++ {
		for x := 0; x < c; x++ {
			if walk(y, x, 0) {
				return true
			}
		}
	}
	return false
}

// queens-holes: an n by n board of '.' and '#'; the answer is how many ways
// n queens fit on '.' squares without attacking each other.
func genQueensHoles(rng *rand.Rand, size int) Case {
	n := 1 + rng.Intn(min(10, max(1, size)))
	holes := []float64{0, 0.1, 0.25}[rng.Intn(3)]
	rows := make([]string, n)
	for i := range rows {
		b := make([]byte, n)
		for j := range b {
			b[j] = '.'
			if rng.Float64() < holes {
				b[j] = '#'
			}
		}
		rows[i] = string(b)
	}
	return Case{Input: strconv.Itoa(n) + "\n" + strings.Join(rows, "\n"), Expected: strconv.Itoa(queensCount(rows))}
}

func queensCount(rows []string) int {
	n := len(rows)
	cols, d1, d2 := make([]bool, n), make([]bool, 2*n), make([]bool, 2*n)
	var place func(r int) int
	place = func(r int) int {
		if r == n {
			return 1
		}
		total := 0
		for c := 0; c < n; c++ {
			if rows[r][c] == '#' || cols[c] || d1[r+c] || d2[r-c+n] {
				continue
			}
			cols[c], d1[r+c], d2[r-c+n] = true, true, true
			total += place(r + 1)
			cols[c], d1[r+c], d2[r-c+n] = false, false, false
		}
		return total
	}
	return place(0)
}
