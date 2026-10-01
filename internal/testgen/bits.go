package testgen

import (
	"math"
	"math/bits"
	"math/rand"
	"strconv"
	"strings"
)

// Generators for the "Bits and math" group.

func init() {
	register("count-set-bits", genCountSetBits)
	register("power-of-two", genPowerOfTwo)
	register("bit-edits", genBitEdits)
	register("lone-number", genLoneNumber)
	register("missing-number", genMissingNumber)
	register("fast-power", genFastPower)
}

// randomInt32 favours the values where bit tricks go wrong: 0, -1, the
// extremes, powers of two and their neighbours.
func randomInt32(rng *rand.Rand) int32 {
	switch rng.Intn(6) {
	case 0:
		return []int32{0, -1, 1, math.MaxInt32, math.MinInt32}[rng.Intn(5)]
	case 1:
		return int32(1)<<rng.Intn(31) + int32(rng.Intn(3)-1)
	case 2:
		return int32(rng.Intn(1000))
	default:
		return int32(rng.Uint32())
	}
}

// count-set-bits: a 32-bit signed integer; the answer is the number of 1
// bits in its two's complement form.
func genCountSetBits(rng *rand.Rand, size int) Case {
	n := randomInt32(rng)
	return Case{Input: strconv.Itoa(int(n)), Expected: strconv.Itoa(bits.OnesCount32(uint32(n)))}
}

// power-of-two: a 32-bit signed integer; the answer is whether it is 2^k for
// some k ≥ 0. Half the cases are exact powers.
func genPowerOfTwo(rng *rand.Rand, size int) Case {
	n := randomInt32(rng)
	if rng.Intn(2) == 0 {
		n = int32(1) << rng.Intn(31)
	}
	return Case{Input: strconv.Itoa(int(n)), Expected: boolText(n > 0 && n&(n-1) == 0)}
}

// bit-edits: a value, then operations "get i", "set i", "clear i" and
// "update i b" applied in order. get gives the bit; the others change the
// value and give it. The answer lists every result.
func genBitEdits(rng *rand.Rand, size int) Case {
	n := int32(rng.Intn(math.MaxInt32))
	if rng.Intn(4) == 0 {
		n = int32(rng.Intn(64))
	}
	q := 1 + rng.Intn(max(1, size))
	lines := []string{strconv.Itoa(int(n)) + " " + strconv.Itoa(q)}
	results := make([]int, q)
	for k := range results {
		i := rng.Intn(31)
		if rng.Intn(2) == 0 {
			i = rng.Intn(6) // low bits, so edits often hit the same bit
		}
		switch op := rng.Intn(4); op {
		case 0:
			lines = append(lines, "get "+strconv.Itoa(i))
			results[k] = int(n >> i & 1)
		case 1:
			lines = append(lines, "set "+strconv.Itoa(i))
			n |= 1 << i
			results[k] = int(n)
		case 2:
			lines = append(lines, "clear "+strconv.Itoa(i))
			n &^= 1 << i
			results[k] = int(n)
		default:
			b := rng.Intn(2)
			lines = append(lines, "update "+strconv.Itoa(i)+" "+strconv.Itoa(b))
			n = n&^(1<<i) | int32(b)<<i
			results[k] = int(n)
		}
	}
	return Case{Input: strings.Join(lines, "\n"), Expected: listOutput(results)}
}

// lone-number: every value appears twice except one; the answer is that one.
func genLoneNumber(rng *rand.Rand, size int) Case {
	pairs := rng.Intn(max(1, size/2) + 1)
	top := []int{10, 1_000_000_000}[rng.Intn(2)]
	vals := distinct(rng, min(pairs+1, 2*top+1), -top, top)
	lone := vals[0]
	nums := []int{lone}
	for _, v := range vals[1:] {
		nums = append(nums, v, v)
	}
	rng.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })
	return Case{Input: arrayInput(nums), Expected: strconv.Itoa(lone)}
}

// missing-number: n different values from 0 to n; the answer is the one
// value in that range that is not there.
func genMissingNumber(rng *rand.Rand, size int) Case {
	n := max(1, size)
	missing := rng.Intn(n + 1)
	if rng.Intn(5) == 0 {
		missing = []int{0, n}[rng.Intn(2)]
	}
	nums := make([]int, 0, n)
	for v := 0; v <= n; v++ {
		if v != missing {
			nums = append(nums, v)
		}
	}
	rng.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })
	return Case{Input: arrayInput(nums), Expected: strconv.Itoa(missing)}
}

// fast-power: a and e; the answer is a^e modulo 10^9 + 7, with 0^0 = 1.
func genFastPower(rng *rand.Rand, size int) Case {
	a := uint64(rng.Intn(1_000_000_001))
	e := uint64(rng.Int63n(1_000_000_000_000_000_001))
	switch rng.Intn(5) {
	case 0:
		a = uint64(rng.Intn(4)) // 0, 1, 2, 3
	case 1:
		e = uint64(rng.Intn(4))
	case 2:
		a = 1_000_000_000 - uint64(rng.Intn(3)) // the largest bases
	}
	result, base := uint64(1), a%mod
	for k := e; k > 0; k >>= 1 {
		if k&1 == 1 {
			result = result * base % mod
		}
		base = base * base % mod
	}
	return Case{Input: strconv.FormatUint(a, 10) + " " + strconv.FormatUint(e, 10), Expected: strconv.FormatUint(result, 10)}
}
