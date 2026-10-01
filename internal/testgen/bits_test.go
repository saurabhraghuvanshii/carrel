package testgen

import (
	"math/big"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func TestCountSetBits(t *testing.T) {
	for i, c := range cases(t, "count-set-bits", 80) {
		n, _ := strconv.ParseInt(c.Input, 10, 32)
		word := strconv.FormatUint(uint64(uint32(n)), 2)
		if want := strconv.Itoa(strings.Count(word, "1")); want != c.Expected {
			t.Fatalf("case %d (%d): want %s, generator says %s", i, n, want, c.Expected)
		}
	}
}

func TestPowerOfTwo(t *testing.T) {
	seen := map[string]bool{}
	for i, c := range cases(t, "power-of-two", 80) {
		n, _ := strconv.Atoi(c.Input)
		ok := false
		for p := 1; p <= n && !ok; p *= 2 {
			ok = p == n
		}
		if boolText(ok) != c.Expected {
			t.Fatalf("case %d (%d): want %v, generator says %s", i, n, ok, c.Expected)
		}
		seen[c.Expected] = true
	}
	if !seen["true"] || !seen["false"] {
		t.Fatal("both answers must appear")
	}
}

func TestBitEdits(t *testing.T) {
	for i, c := range cases(t, "bit-edits", 60) {
		lines := strings.Split(c.Input, "\n")
		first := strings.Fields(lines[0])
		n, _ := strconv.Atoi(first[0])
		// Keep the value as 31 binary digits, lowest bit last.
		word := []byte(strconv.FormatInt(int64(n), 2))
		word = append([]byte(strings.Repeat("0", 31-len(word))), word...)
		value := func() int {
			v, _ := strconv.ParseInt(string(word), 2, 64)
			return int(v)
		}
		var got []int
		for _, l := range lines[1:] {
			f := strings.Fields(l)
			bit, _ := strconv.Atoi(f[1])
			at := 30 - bit
			switch f[0] {
			case "get":
				got = append(got, int(word[at]-'0'))
				continue
			case "set":
				word[at] = '1'
			case "clear":
				word[at] = '0'
			case "update":
				word[at] = f[2][0]
			}
			got = append(got, value())
		}
		if listOutput(got) != c.Expected {
			t.Fatalf("case %d differs", i)
		}
	}
}

func TestLoneNumber(t *testing.T) {
	for i, c := range cases(t, "lone-number", 60) {
		nums, _ := parseArray(t, c.Input)
		count := map[int]int{}
		for _, v := range nums {
			count[v]++
		}
		var once []int
		for v, k := range count {
			if k == 1 {
				once = append(once, v)
			} else if k != 2 {
				t.Fatalf("case %d: %d appears %d times", i, v, k)
			}
		}
		if len(once) != 1 || strconv.Itoa(once[0]) != c.Expected {
			t.Fatalf("case %d: values seen once %v, generator says %s", i, once, c.Expected)
		}
	}
}

func TestMissingNumber(t *testing.T) {
	for i, c := range cases(t, "missing-number", 60) {
		nums, _ := parseArray(t, c.Input)
		s := slices.Clone(nums)
		slices.Sort(s)
		want := len(s)
		for k, v := range s {
			if v != k {
				want = k
				break
			}
		}
		if strconv.Itoa(want) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, want, c.Expected)
		}
	}
}

func TestFastPower(t *testing.T) {
	for i, c := range cases(t, "fast-power", 80) {
		f := strings.Fields(c.Input)
		a, _ := new(big.Int).SetString(f[0], 10)
		e, _ := new(big.Int).SetString(f[1], 10)
		want := new(big.Int).Exp(a, e, big.NewInt(mod))
		if e.Sign() == 0 {
			want.SetInt64(1)
		}
		if want.String() != c.Expected {
			t.Fatalf("case %d (%s): want %s, generator says %s", i, c.Input, want, c.Expected)
		}
	}
}
