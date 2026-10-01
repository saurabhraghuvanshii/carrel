package testgen

import (
	"strconv"
	"strings"
	"testing"
)

func TestGenerateIsDeterministic(t *testing.T) {
	for _, name := range []string{"pair-sum", "nearest-driver"} {
		a, err := Generate(name, 42, 20)
		if err != nil {
			t.Fatal(err)
		}
		b, _ := Generate(name, 42, 20)
		for i := range a {
			if a[i] != b[i] {
				t.Fatalf("%s: case %d differs for the same seed", name, i)
			}
		}
	}
}

func TestPairSumHasExactlyOnePair(t *testing.T) {
	cases, err := Generate("pair-sum", 7, 40)
	if err != nil {
		t.Fatal(err)
	}
	for i, c := range cases {
		lines := strings.Split(c.Input, "\n")
		fields := strings.Fields(lines[1])
		target, _ := strconv.Atoi(lines[2])
		nums := make([]int, len(fields))
		for k, f := range fields {
			nums[k], _ = strconv.Atoi(f)
		}
		if countPairs(nums, target) != 1 {
			t.Fatalf("case %d: expected exactly one pair", i)
		}
		if len(nums) > 400 {
			continue // skip the slow brute force on big inputs
		}
		found := ""
		for p := 0; p < len(nums); p++ {
			for q := p + 1; q < len(nums); q++ {
				if nums[p]+nums[q] == target {
					found = strconv.Itoa(p) + " " + strconv.Itoa(q)
				}
			}
		}
		if found != c.Expected {
			t.Fatalf("case %d: brute force says %q, generator says %q", i, found, c.Expected)
		}
	}
}

func TestNearestDriverMatchesBruteForce(t *testing.T) {
	cases, err := Generate("nearest-driver", 11, 40)
	if err != nil {
		t.Fatal(err)
	}
	for i, c := range cases {
		lines := strings.Split(c.Input, "\n")
		n, _ := strconv.Atoi(lines[0])
		best, bestDist := -1, -1
		var rx, ry int
		rider := strings.Fields(lines[n+1])
		rx, _ = strconv.Atoi(rider[0])
		ry, _ = strconv.Atoi(rider[1])
		for k := 1; k <= n; k++ {
			f := strings.Fields(lines[k])
			id, _ := strconv.Atoi(f[0])
			x, _ := strconv.Atoi(f[1])
			y, _ := strconv.Atoi(f[2])
			free, _ := strconv.Atoi(f[3])
			if free == 0 {
				continue
			}
			d := abs(x-rx) + abs(y-ry)
			if bestDist == -1 || d < bestDist || (d == bestDist && id < best) {
				best, bestDist = id, d
			}
		}
		if strconv.Itoa(best) != c.Expected {
			t.Fatalf("case %d: expected %d, generator says %s", i, best, c.Expected)
		}
	}
}
