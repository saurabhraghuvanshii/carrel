package testgen

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

func init() {
	register("pair-sum", genPairSum)
	register("nearest-driver", genNearestDriver)
}

// pair-sum: input is n, the numbers, then the target.
// Expected is the two indices in ascending order. The generator plants a pair
// and keeps the case only if that pair is the only one.
func genPairSum(rng *rand.Rand, size int) Case {
	if size < 2 {
		size = 2
	}
	span := size * size * 8
	if span < 100 {
		span = 100
	}
	if span > 2_000_000_000 {
		span = 2_000_000_000
	}
	for {
		nums := make([]int, size)
		for i := range nums {
			nums[i] = rng.Intn(span) - span/2
		}
		i := rng.Intn(size)
		j := rng.Intn(size - 1)
		if j >= i {
			j++
		}
		if i > j {
			i, j = j, i
		}
		target := nums[i] + nums[j]
		if countPairs(nums, target) != 1 {
			continue
		}
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(size))
		sb.WriteString("\n")
		for k, v := range nums {
			if k > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteString("\n")
		sb.WriteString(strconv.Itoa(target))
		return Case{Input: sb.String(), Expected: fmt.Sprintf("%d %d", i, j)}
	}
}

// countPairs counts index pairs p<q with nums[p]+nums[q]==target, stopping at 2.
func countPairs(nums []int, target int) int {
	seen := map[int]int{}
	count := 0
	for _, v := range nums {
		count += seen[target-v]
		if count > 1 {
			return count
		}
		seen[v]++
	}
	return count
}

// nearest-driver: input is n, then n lines "id x y free", then "rx ry".
// Expected is the id of the closest free driver (rows plus columns), ties go to
// the lower id, and -1 when nobody is free.
func genNearestDriver(rng *rand.Rand, size int) Case {
	n := size
	ids := rng.Perm(n)
	grid := 20 + n
	if grid > 1000 {
		grid = 1000
	}
	freeChance := 0.7
	if rng.Intn(8) == 0 {
		freeChance = 0 // sometimes nobody is free
	}
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(n))
	sb.WriteString("\n")
	drivers := make([][4]int, n)
	for k := 0; k < n; k++ {
		free := 0
		if rng.Float64() < freeChance {
			free = 1
		}
		drivers[k] = [4]int{ids[k] + 1, rng.Intn(grid), rng.Intn(grid), free}
		fmt.Fprintf(&sb, "%d %d %d %d\n", drivers[k][0], drivers[k][1], drivers[k][2], drivers[k][3])
	}
	rx, ry := rng.Intn(grid), rng.Intn(grid)
	fmt.Fprintf(&sb, "%d %d", rx, ry)
	return Case{Input: sb.String(), Expected: strconv.Itoa(nearestDriver(drivers, rx, ry))}
}

func nearestDriver(drivers [][4]int, rx, ry int) int {
	best, bestDist := -1, 1<<62
	for _, d := range drivers {
		if d[3] == 0 {
			continue
		}
		dist := abs(d[1]-rx) + abs(d[2]-ry)
		if dist < bestDist || (dist == bestDist && d[0] < best) {
			best, bestDist = d[0], dist
		}
	}
	return best
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
