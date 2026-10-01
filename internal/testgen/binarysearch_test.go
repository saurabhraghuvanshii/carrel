package testgen

import (
	"math"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func TestFindSorted(t *testing.T) {
	found := 0
	for i, c := range cases(t, "find-sorted", 60) {
		nums, rest := parseArray(t, c.Input)
		target, _ := strconv.Atoi(rest[0])
		want := -1
		for j, v := range nums {
			if v == target {
				want = j
			}
		}
		if !slices.IsSorted(nums) || strconv.Itoa(want) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, want, c.Expected)
		}
		if want >= 0 {
			found++
		}
	}
	if found == 0 || found == 60 {
		t.Fatal("need both found and missing targets")
	}
}

func TestIntSqrt(t *testing.T) {
	for i, c := range cases(t, "int-sqrt", 80) {
		x, _ := strconv.Atoi(c.Input)
		r, _ := strconv.Atoi(c.Expected)
		if x < 0 || x > math.MaxInt32 || r*r > x || (r+1)*(r+1) <= x {
			t.Fatalf("case %d: %d is not the floor square root of %d", i, r, x)
		}
	}
}

func TestFirstLast(t *testing.T) {
	for i, c := range cases(t, "first-last", 60) {
		nums, rest := parseArray(t, c.Input)
		target, _ := strconv.Atoi(rest[0])
		first, last := slices.Index(nums, target), -1
		for j := len(nums) - 1; j >= 0; j-- {
			if nums[j] == target {
				last = j
				break
			}
		}
		if listOutput([]int{first, last}) != c.Expected {
			t.Fatalf("case %d: want [%d, %d], generator says %s", i, first, last, c.Expected)
		}
	}
}

// isRotatedSorted reports whether nums is a rotation of a strictly rising list.
func isRotatedSorted(nums []int) bool {
	drops := 0
	for i := range nums {
		if nums[i] >= nums[(i+1)%len(nums)] && len(nums) > 1 {
			drops++
		}
	}
	return drops <= 1
}

func TestMinRotated(t *testing.T) {
	for i, c := range cases(t, "min-rotated", 50) {
		nums, _ := parseArray(t, c.Input)
		m := nums[0]
		for _, v := range nums {
			m = min(m, v)
		}
		if !isRotatedSorted(nums) || strconv.Itoa(m) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, m, c.Expected)
		}
	}
}

func TestSearchRotated(t *testing.T) {
	for i, c := range cases(t, "search-rotated", 60) {
		nums, rest := parseArray(t, c.Input)
		target, _ := strconv.Atoi(rest[0])
		want := -1
		for j, v := range nums {
			if v == target {
				want = j
			}
		}
		if !isRotatedSorted(nums) || strconv.Itoa(want) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, want, c.Expected)
		}
	}
}

func TestMountainTop(t *testing.T) {
	for i, c := range cases(t, "mountain-top", 60) {
		nums, _ := parseArray(t, c.Input)
		top := 0
		for j := range nums {
			if nums[j] > nums[top] {
				top = j
			}
		}
		for j := 1; j < len(nums); j++ {
			if (j <= top && nums[j] <= nums[j-1]) || (j > top && nums[j] >= nums[j-1]) {
				t.Fatalf("case %d: not a strict mountain at %d", i, j)
			}
		}
		if strconv.Itoa(top) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, top, c.Expected)
		}
	}
}

func TestEatingSpeed(t *testing.T) {
	for i, c := range cases(t, "eating-speed", 60) {
		jobs, rest := parseArray(t, c.Input)
		h, _ := strconv.Atoi(rest[0])
		k, _ := strconv.Atoi(c.Expected)
		if h < len(jobs) || h > 1_000_000_000 {
			t.Fatalf("case %d: h=%d out of range", i, h)
		}
		if hoursAt(jobs, k) > h || (k > 1 && hoursAt(jobs, k-1) <= h) {
			t.Fatalf("case %d: %d is not the smallest speed", i, k)
		}
	}
}

func TestMedianTwo(t *testing.T) {
	for i, c := range cases(t, "median-two", 60) {
		a, rest := parseArray(t, c.Input)
		b, _ := parseArray(t, strings.Join(rest, "\n"))
		all := append(slices.Clone(a), b...)
		slices.Sort(all)
		n := len(all)
		var m float64
		if n%2 == 1 {
			m = float64(all[n/2])
		} else {
			m = float64(all[n/2-1]+all[n/2]) / 2
		}
		if want := strconv.FormatFloat(m, 'f', 1, 64); want != c.Expected {
			t.Fatalf("case %d: want %s, generator says %s", i, want, c.Expected)
		}
	}
}
