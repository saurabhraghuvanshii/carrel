package testgen

import (
	"math/rand"
	"strconv"
	"strings"
)

func init() { register("reverse-list", genReverseList) }

// reverse-list: a linked list in the array format, "n" then the values
// ("0" and an empty line when empty). Expected is the reversed list as
// "[a, b, c]", or "[]".
func genReverseList(rng *rand.Rand, size int) Case {
	n := size
	switch rng.Intn(12) {
	case 0:
		n = 0
	case 1:
		n = 1
	}
	vals := make([]int, n)
	for i := range vals {
		vals[i] = rng.Intn(2001) - 1000
	}
	rev := make([]int, n)
	for i, v := range vals {
		rev[n-1-i] = v
	}
	return Case{Input: strconv.Itoa(n) + "\n" + joinInts(vals, " "), Expected: "[" + joinInts(rev, ", ") + "]"}
}

func joinInts(nums []int, sep string) string {
	parts := make([]string, len(nums))
	for i, v := range nums {
		parts[i] = strconv.Itoa(v)
	}
	return strings.Join(parts, sep)
}
