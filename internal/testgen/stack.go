package testgen

import (
	"math/rand"
	"strconv"
	"strings"
)

// Generators for the "Stack" group.

func init() {
	register("brackets", genBrackets)
	register("adjacent-twins", genAdjacentTwins)
	register("min-stack", genMinStack)
	register("postfix", genPostfix)
	register("redundant-brackets", genRedundantBrackets)
	register("warmer-day", genWarmerDay)
	register("price-span", genPriceSpan)
	register("decode-text", genDecodeText)
	register("robot-rail", genRobotRail)
	register("histogram", genHistogram)
}

const openers, closers = "([{", ")]}"

// brackets: a word of ()[]{}; half are balanced, the rest are broken by one
// change, deletion or swap.
func genBrackets(rng *rand.Rand, size int) Case {
	pairs := max(1, size/2)
	var b, stack []byte
	for opened := 0; opened < pairs || len(stack) > 0; {
		if opened < pairs && (len(stack) == 0 || rng.Intn(2) == 0) {
			k := rng.Intn(3)
			b = append(b, openers[k])
			stack = append(stack, closers[k])
			opened++
		} else {
			b = append(b, stack[len(stack)-1])
			stack = stack[:len(stack)-1]
		}
	}
	if rng.Intn(2) == 0 {
		i := rng.Intn(len(b))
		switch rng.Intn(3) {
		case 0:
			b[i] = "()[]{}"[rng.Intn(6)]
		case 1:
			if len(b) > 1 {
				b = append(b[:i], b[i+1:]...)
			}
		default:
			if i+1 < len(b) {
				b[i], b[i+1] = b[i+1], b[i]
			}
		}
	}
	return Case{Input: string(b), Expected: boolText(balanced(string(b)))}
}

func balanced(s string) bool {
	var stack []byte
	for i := 0; i < len(s); i++ {
		if k := strings.IndexByte(openers, s[i]); k >= 0 {
			stack = append(stack, closers[k])
		} else if len(stack) == 0 || stack[len(stack)-1] != s[i] {
			return false
		} else {
			stack = stack[:len(stack)-1]
		}
	}
	return len(stack) == 0
}

// adjacent-twins: a lowercase word; the answer is what is left after
// removing equal neighbours again and again, in quotes.
func genAdjacentTwins(rng *rand.Rand, size int) Case {
	n := max(1, size)
	letters := "abcdefghijklmnopqrstuvwxyz"[:2+rng.Intn(4)]
	var word string
	if rng.Intn(4) == 0 {
		half := randomWord(rng, max(1, n/2), letters)
		b := []byte(half)
		for i := len(half) - 1; i >= 0; i-- {
			b = append(b, half[i])
		}
		word = string(b) // cancels completely
	} else {
		word = randomWord(rng, n, letters)
	}
	return Case{Input: word, Expected: `"` + removeTwins(word) + `"`}
}

func removeTwins(word string) string {
	var stack []byte
	for i := 0; i < len(word); i++ {
		if len(stack) > 0 && stack[len(stack)-1] == word[i] {
			stack = stack[:len(stack)-1]
		} else {
			stack = append(stack, word[i])
		}
	}
	return string(stack)
}

// min-stack: a call sequence of new, push x, pop, top and min. pop, top and
// min only happen on a non-empty stack.
func genMinStack(rng *rand.Rand, size int) Case {
	k := min(max(size, 5), 5000)
	calls := []string{"new"}
	var values, mins []int
	var out []string
	out = append(out, "null")
	for len(calls) < k {
		r := rng.Intn(10)
		switch {
		case len(values) == 0 || r < 4:
			v := rng.Intn(200_001) - 100_000
			if len(mins) > 0 && rng.Intn(4) == 0 {
				v = mins[len(mins)-1] // equal to the current minimum
			}
			calls = append(calls, "push "+strconv.Itoa(v))
			values = append(values, v)
			if len(mins) == 0 || v <= mins[len(mins)-1] {
				mins = append(mins, v)
			}
			out = append(out, "null")
		case r < 6:
			calls = append(calls, "pop")
			if values[len(values)-1] == mins[len(mins)-1] {
				mins = mins[:len(mins)-1]
			}
			values = values[:len(values)-1]
			out = append(out, "null")
		case r < 8:
			calls = append(calls, "top")
			out = append(out, strconv.Itoa(values[len(values)-1]))
		default:
			calls = append(calls, "min")
			out = append(out, strconv.Itoa(mins[len(mins)-1]))
		}
	}
	return Case{Input: strconv.Itoa(len(calls)) + "\n" + strings.Join(calls, "\n"), Expected: strings.Join(out, " ")}
}

const exprLimit = 100_000_000

// postfixTree builds a random expression with the given number of numbers
// and returns its tokens and value. Every step stays within ±10⁸ and never
// divides by zero.
func postfixTree(rng *rand.Rand, leaves int) ([]string, int) {
	if leaves == 1 {
		v := rng.Intn(201) - 100
		return []string{strconv.Itoa(v)}, v
	}
	left := 1 + rng.Intn(leaves-1)
	lt, a := postfixTree(rng, left)
	rt, b := postfixTree(rng, leaves-left)
	ops := []byte{"+-*/"[rng.Intn(4)], '+', '-', '/'}
	for _, op := range ops {
		var v int
		switch op {
		case '+':
			v = a + b
		case '-':
			v = a - b
		case '*':
			v = a * b
		case '/':
			if b == 0 {
				continue
			}
			v = a / b
		}
		if v >= -exprLimit && v <= exprLimit {
			return append(append(lt, rt...), string(op)), v
		}
	}
	// b == 0 and a+b, a-b both out of range cannot happen: then a+b == a.
	panic("unreachable")
}

// postfix: tokens of a postfix expression; the answer is its value.
func genPostfix(rng *rand.Rand, size int) Case {
	leaves := max(1, min(size, 3000)/2+1)
	tokens, v := postfixTree(rng, leaves)
	return Case{Input: strconv.Itoa(len(tokens)) + "\n" + strings.Join(tokens, " "), Expected: strconv.Itoa(v)}
}

// redundantExpr writes a random valid expression of single letters, with
// some groups wrapped in brackets, some wrapped twice, and some single
// letters wrapped (which does not count as pointless).
func redundantExpr(rng *rand.Rand, leaves int, twice float64) string {
	var s string
	if leaves == 1 {
		s = string(rune('a' + rng.Intn(26)))
		if rng.Intn(6) == 0 {
			s = "(" + s + ")"
		}
	} else {
		left := 1 + rng.Intn(leaves-1)
		s = redundantExpr(rng, left, twice) + string("+-*/"[rng.Intn(4)]) + redundantExpr(rng, leaves-left, twice)
		if rng.Intn(2) == 0 {
			s = "(" + s + ")"
		}
	}
	if rng.Float64() < twice {
		s = "((" + s + "))"
	}
	return s
}

// redundant-brackets: the answer is whether some bracket pair holds nothing
// but other bracket groups.
func genRedundantBrackets(rng *rand.Rand, size int) Case {
	leaves := max(1, min(size, 2000)/2+1)
	twice := 0.0
	if rng.Intn(2) == 0 {
		twice = 1.0 / float64(leaves+1)
	}
	expr := redundantExpr(rng, leaves, twice)
	return Case{Input: expr, Expected: boolText(hasPointless(expr))}
}

// hasPointless is the stack method: at each ')' count what is popped before
// its '('. A group leaves nothing behind once closed.
func hasPointless(expr string) bool {
	var stack []byte
	for i := 0; i < len(expr); i++ {
		if expr[i] != ')' {
			stack = append(stack, expr[i])
			continue
		}
		count := 0
		for stack[len(stack)-1] != '(' {
			stack = stack[:len(stack)-1]
			count++
		}
		stack = stack[:len(stack)-1]
		if count == 0 {
			return true
		}
	}
	return false
}

// warmer-day: daily temperatures; the answer is, for each day, how many days
// until a warmer one, or 0.
func genWarmerDay(rng *rand.Rand, size int) Case {
	n := max(1, size)
	temps := make([]int, n)
	t := rng.Intn(111) - 50
	for i := range temps {
		t = max(-50, min(60, t+rng.Intn(9)-4))
		temps[i] = t
	}
	return Case{Input: arrayInput(temps), Expected: listOutput(warmerDays(temps))}
}

func warmerDays(temps []int) []int {
	out := make([]int, len(temps))
	var stack []int
	for i, t := range temps {
		for len(stack) > 0 && temps[stack[len(stack)-1]] < t {
			j := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			out[j] = i - j
		}
		stack = append(stack, i)
	}
	return out
}

// price-span: daily prices with flat stretches; the answer is, for each day,
// how many days in a row ending today had a price at most today's.
func genPriceSpan(rng *rand.Rand, size int) Case {
	n := max(1, size)
	prices := make([]int, n)
	p := 1 + rng.Intn(100_000)
	for i := range prices {
		if rng.Intn(3) != 0 {
			p = max(1, min(100_000, p+rng.Intn(21)-10))
		}
		prices[i] = p
	}
	return Case{Input: arrayInput(prices), Expected: listOutput(priceSpans(prices))}
}

func priceSpans(prices []int) []int {
	out := make([]int, len(prices))
	var stack []int
	for i, p := range prices {
		for len(stack) > 0 && prices[stack[len(stack)-1]] <= p {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			out[i] = i + 1
		} else {
			out[i] = i - stack[len(stack)-1]
		}
		stack = append(stack, i)
	}
	return out
}

const decodedLimit = 20_000

// encodeText writes letters and k[...] groups whose decoded length stays
// within budget. It returns the encoded and the decoded text.
func encodeText(rng *rand.Rand, budget, depth int) (string, string) {
	var enc, dec strings.Builder
	for parts := 1 + rng.Intn(3); parts > 0 && budget > 0; parts-- {
		if depth < 5 && budget >= 2 && rng.Intn(2) == 0 {
			k := 1 + rng.Intn(min(12, budget))
			if rng.Intn(8) == 0 {
				k = 10 + rng.Intn(21)
			}
			inner := max(1, budget/k)
			e, d := encodeText(rng, inner, depth+1)
			if len(d)*k > budget {
				k = max(1, budget/len(d))
			}
			enc.WriteString(strconv.Itoa(k) + "[" + e + "]")
			dec.WriteString(strings.Repeat(d, k))
			budget -= len(d) * k
		} else {
			w := randomWord(rng, 1+rng.Intn(min(4, budget)), "abcdefghijklmnopqrstuvwxyz")
			enc.WriteString(w)
			dec.WriteString(w)
			budget -= len(w)
		}
	}
	if enc.Len() == 0 {
		return "a", "a"
	}
	return enc.String(), dec.String()
}

// decode-text: k[...] means the bracket's text k times; the answer is the
// decoded text.
func genDecodeText(rng *rand.Rand, size int) Case {
	budget := min(decodedLimit, max(1, size*2))
	enc, dec := encodeText(rng, budget, 0)
	return Case{Input: enc, Expected: dec}
}

// robot-rail: nonzero strengths, the sign is the direction; the answer is
// the robots left after every collision, in order.
func genRobotRail(rng *rand.Rand, size int) Case {
	n := max(1, size)
	top := []int{3, 10, 1000}[rng.Intn(3)]
	robots := make([]int, n)
	for i := range robots {
		robots[i] = 1 + rng.Intn(top)
		if rng.Intn(2) == 0 {
			robots[i] = -robots[i]
		}
	}
	return Case{Input: arrayInput(robots), Expected: listOutput(robotRail(robots))}
}

func robotRail(robots []int) []int {
	stack := []int{}
	for _, r := range robots {
		alive := true
		for alive && r < 0 && len(stack) > 0 && stack[len(stack)-1] > 0 {
			top := stack[len(stack)-1]
			switch {
			case top < -r:
				stack = stack[:len(stack)-1]
			case top == -r:
				stack = stack[:len(stack)-1]
				alive = false
			default:
				alive = false
			}
		}
		if alive {
			stack = append(stack, r)
		}
	}
	return stack
}

// histogram: bar heights; the answer is the largest rectangle under them.
func genHistogram(rng *rand.Rand, size int) Case {
	n := max(1, size)
	top := []int{5, 100, 10_000}[rng.Intn(3)]
	h := make([]int, n)
	switch rng.Intn(3) {
	case 0:
		for i := range h {
			h[i] = rng.Intn(top + 1)
		}
	case 1: // rising then falling
		peak := rng.Intn(n)
		for i := range h {
			d := peak - i
			if d < 0 {
				d = -d
			}
			h[i] = max(0, top-d*top/max(1, n))
		}
	default: // plateaus
		v := rng.Intn(top + 1)
		for i := range h {
			if rng.Intn(4) == 0 {
				v = rng.Intn(top + 1)
			}
			h[i] = v
		}
	}
	return Case{Input: arrayInput(h), Expected: strconv.Itoa(largestRectangle(h))}
}

func largestRectangle(h []int) int {
	best := 0
	var stack []int
	for i := 0; i <= len(h); i++ {
		cur := 0
		if i < len(h) {
			cur = h[i]
		}
		for len(stack) > 0 && h[stack[len(stack)-1]] >= cur {
			height := h[stack[len(stack)-1]]
			stack = stack[:len(stack)-1]
			left := -1
			if len(stack) > 0 {
				left = stack[len(stack)-1]
			}
			best = max(best, height*(i-left-1))
		}
		stack = append(stack, i)
	}
	return best
}
