package testgen

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
)

// Generators for the "Strings" group. Text answers print in double quotes so
// an empty answer or spaces at the ends stay visible.

func init() {
	register("squash-repeats", genSquashRepeats)
	register("run-length", genRunLength)
	register("walk-distance", genWalkDistance)
	register("common-start", genCommonStart)
	register("reverse-words", genReverseWords)
	register("first-unique", genFirstUnique)
	register("mirror-stretch", genMirrorStretch)
	register("text-to-number", genTextToNumber)
}

// runsWord makes a word of letter runs, each 1 to maxRun long.
func runsWord(rng *rand.Rand, n, maxRun int, letters string) string {
	var b strings.Builder
	for b.Len() < n {
		b.WriteString(strings.Repeat(letters[rng.Intn(len(letters)):][:1], 1+rng.Intn(maxRun)))
	}
	return b.String()[:n]
}

// squash-repeats: the answer keeps one letter of every run of equal letters.
func genSquashRepeats(rng *rand.Rand, size int) Case {
	w := runsWord(rng, max(1, size), []int{1, 3, 12}[rng.Intn(3)], "abcdefghijklmnopqrstuvwxyz"[:1+rng.Intn(5)])
	var b strings.Builder
	for i := 0; i < len(w); i++ {
		if i == 0 || w[i] != w[i-1] {
			b.WriteByte(w[i])
		}
	}
	return Case{Input: w, Expected: `"` + b.String() + `"`}
}

// run-length: each run becomes its letter, followed by its length when the
// run is longer than one.
func genRunLength(rng *rand.Rand, size int) Case {
	w := runsWord(rng, max(1, size), []int{2, 4, 30}[rng.Intn(3)], "abcdefghijklmnopqrstuvwxyz"[:1+rng.Intn(5)])
	var b strings.Builder
	for i := 0; i < len(w); {
		j := i
		for j < len(w) && w[j] == w[i] {
			j++
		}
		b.WriteByte(w[i])
		if j-i > 1 {
			b.WriteString(strconv.Itoa(j - i))
		}
		i = j
	}
	return Case{Input: w, Expected: `"` + b.String() + `"`}
}

// walk-distance: steps N, S, E and W of one unit; the answer is the squared
// straight-line distance from start to finish. Some walks mostly go one way,
// so the square passes 32 bits.
func genWalkDistance(rng *rand.Rand, size int) Case {
	n := max(1, size*[]int{1, 10}[rng.Intn(2)])
	n = min(n, 100_000)
	steps := []byte(randomWord(rng, n, "NSEW"))
	if rng.Intn(3) == 0 {
		lean := "NSEW"[rng.Intn(4)]
		for i := range steps {
			if rng.Intn(5) != 0 {
				steps[i] = lean
			}
		}
	}
	x, y := 0, 0
	for _, s := range steps {
		switch s {
		case 'N':
			y++
		case 'S':
			y--
		case 'E':
			x++
		default:
			x--
		}
	}
	return Case{Input: string(steps), Expected: strconv.Itoa(x*x + y*y)}
}

// common-start: n words; the answer is the longest beginning they all share.
func genCommonStart(rng *rand.Rand, size int) Case {
	n := 1 + rng.Intn(min(200, max(1, size)))
	letters := "abc"[:1+rng.Intn(3)]
	stem := randomWord(rng, rng.Intn(min(50, max(1, size))+1), letters)
	words := make([]string, n)
	for i := range words {
		cut := len(stem)
		if rng.Intn(4) == 0 {
			cut = rng.Intn(len(stem) + 1)
		}
		words[i] = stem[:cut] + randomWord(rng, rng.Intn(5), letters)
		if words[i] == "" {
			words[i] = letters[:1]
		}
	}
	shared := words[0]
	for _, w := range words[1:] {
		k := 0
		for k < len(shared) && k < len(w) && shared[k] == w[k] {
			k++
		}
		shared = shared[:k]
	}
	return Case{Input: strconv.Itoa(n) + "\n" + strings.Join(words, " "), Expected: `"` + shared + `"`}
}

// reverse-words: text in quotes with words split by one or more spaces, maybe
// with spaces at the ends. The answer is the words in reverse order, one
// space apart.
func genReverseWords(rng *rand.Rand, size int) Case {
	n := 1 + rng.Intn(min(1000, max(1, size)))
	var b strings.Builder
	b.WriteString(strings.Repeat(" ", rng.Intn(3)))
	words := make([]string, n)
	for i := range words {
		words[i] = randomWord(rng, 1+rng.Intn(6), "abcdefghij0123456789")
		if i > 0 {
			b.WriteString(strings.Repeat(" ", 1+rng.Intn(3)*rng.Intn(2)))
		}
		b.WriteString(words[i])
	}
	b.WriteString(strings.Repeat(" ", rng.Intn(3)))
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		words[i], words[j] = words[j], words[i]
	}
	return Case{Input: `"` + b.String() + `"`, Expected: `"` + strings.Join(words, " ") + `"`}
}

// first-unique: for every prefix of the word, the first letter that has
// appeared exactly once so far, or '#'. The answer joins them.
func genFirstUnique(rng *rand.Rand, size int) Case {
	w := randomWord(rng, max(1, size), "abcdefghijklmnopqrstuvwxyz"[:1+rng.Intn(26)])
	count := [26]int{}
	var queue []byte
	out := make([]byte, len(w))
	for i := 0; i < len(w); i++ {
		count[w[i]-'a']++
		queue = append(queue, w[i])
		for len(queue) > 0 && count[queue[0]-'a'] > 1 {
			queue = queue[1:]
		}
		out[i] = '#'
		if len(queue) > 0 {
			out[i] = queue[0]
		}
	}
	return Case{Input: w, Expected: `"` + string(out) + `"`}
}

// mirror-stretch: the answer is the longest stretch that reads the same both
// ways; the leftmost of equal lengths.
func genMirrorStretch(rng *rand.Rand, size int) Case {
	n := max(1, min(1000, size))
	w := []byte(randomWord(rng, n, "abcdefghijklmnopqrstuvwxyz"[:1+rng.Intn(4)]))
	if rng.Intn(2) == 0 && n > 3 { // plant a long mirror
		l := 2 + rng.Intn(n-2)
		at := rng.Intn(n - l + 1)
		for k := 0; k < l/2; k++ {
			w[at+l-1-k] = w[at+k]
		}
	}
	best, from := 0, 0
	for c := 0; c < 2*n-1; c++ {
		lo, hi := c/2, c/2+c%2
		for lo >= 0 && hi < n && w[lo] == w[hi] {
			lo, hi = lo-1, hi+1
		}
		if hi-lo-1 > best || (hi-lo-1 == best && lo+1 < from) {
			best, from = hi-lo-1, lo+1
		}
	}
	return Case{Input: string(w), Expected: `"` + string(w[from:from+best]) + `"`}
}

// text-to-number: text in quotes. Skip leading spaces, read an optional sign
// and then digits until the first non-digit; no digits gives 0. Values past
// the 32-bit range are clamped to it.
func genTextToNumber(rng *rand.Rand, size int) Case {
	var b strings.Builder
	b.WriteString(strings.Repeat(" ", rng.Intn(4)))
	if rng.Intn(8) == 0 {
		b.WriteString(randomWord(rng, 1+rng.Intn(3), "ab. +-"))
	}
	b.WriteString([]string{"", "", "+", "-", "-"}[rng.Intn(5)])
	if rng.Intn(4) == 0 {
		b.WriteString(strings.Repeat("0", 1+rng.Intn(15)))
	}
	digits := []int{1 + rng.Intn(5), 1 + rng.Intn(9), 1 + rng.Intn(9), 9 + rng.Intn(3), 10, 12 + rng.Intn(20)}[rng.Intn(6)]
	if rng.Intn(10) == 0 {
		digits = 0
	}
	b.WriteString(randomWord(rng, digits, "0123456789"))
	if rng.Intn(10) == 0 { // the edges of the range, exactly
		b.Reset()
		b.WriteString([]string{"2147483647", "-2147483648", "2147483648", "-2147483649"}[rng.Intn(4)])
	}
	b.WriteString(randomWord(rng, rng.Intn(5), "abc 1.-+"))
	text := b.String()

	i := 0
	for i < len(text) && text[i] == ' ' {
		i++
	}
	sign := 1
	if i < len(text) && (text[i] == '+' || text[i] == '-') {
		if text[i] == '-' {
			sign = -1
		}
		i++
	}
	value := 0
	for ; i < len(text) && text[i] >= '0' && text[i] <= '9'; i++ {
		value = min(value*10+int(text[i]-'0'), math.MaxInt32+1)
	}
	value = max(math.MinInt32, min(math.MaxInt32, sign*value))
	return Case{Input: `"` + text + `"`, Expected: strconv.Itoa(value)}
}
