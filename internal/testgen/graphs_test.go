package testgen

import (
	"math/bits"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func parseGraph(t *testing.T, lines []string) (int, []edge, []string) {
	t.Helper()
	nm := strings.Fields(lines[0])
	n, _ := strconv.Atoi(nm[0])
	m, _ := strconv.Atoi(nm[1])
	edges := make([]edge, m)
	for i := range edges {
		f := strings.Fields(lines[i+1])
		edges[i].u, _ = strconv.Atoi(f[0])
		edges[i].v, _ = strconv.Atoi(f[1])
		if len(f) > 2 {
			edges[i].w, _ = strconv.Atoi(f[2])
		}
		if edges[i].u < 0 || edges[i].u >= n || edges[i].v < 0 || edges[i].v >= n {
			t.Fatalf("edge %v out of range", edges[i])
		}
	}
	return n, edges, lines[m+1:]
}

// reach marks nodes reachable from s by depth-first search.
func reach(n int, edges []edge, s int, directed bool) []bool {
	adj := make([][]int, n)
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], e.v)
		if !directed {
			adj[e.v] = append(adj[e.v], e.u)
		}
	}
	seen := make([]bool, n)
	stack := []int{s}
	seen[s] = true
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, v := range adj[u] {
			if !seen[v] {
				seen[v] = true
				stack = append(stack, v)
			}
		}
	}
	return seen
}

func TestPathExistsAndProvinces(t *testing.T) {
	for i, c := range cases(t, "path-exists", 50) {
		n, edges, rest := parseGraph(t, strings.Split(c.Input, "\n"))
		sd := strings.Fields(rest[0])
		s, _ := strconv.Atoi(sd[0])
		d, _ := strconv.Atoi(sd[1])
		if boolText(reach(n, edges, s, false)[d]) != c.Expected {
			t.Fatalf("path case %d differs", i)
		}
	}
	for i, c := range cases(t, "provinces", 50) {
		n, edges, _ := parseGraph(t, strings.Split(c.Input, "\n"))
		groups := 0
		seen := make([]bool, n)
		for s := 0; s < n; s++ {
			if !seen[s] {
				groups++
				for v, ok := range reach(n, edges, s, false) {
					seen[v] = seen[v] || ok
				}
			}
			if n > 1500 {
				break
			}
		}
		if n <= 1500 && strconv.Itoa(groups) != c.Expected {
			t.Fatalf("provinces case %d: want %d, generator says %s", i, groups, c.Expected)
		}
	}
}

func TestIslands(t *testing.T) {
	for i, c := range cases(t, "islands", 50) {
		lines := strings.Split(c.Input, "\n")
		rows := lines[1:]
		r, cols := len(rows), len(rows[0])
		// Union-find over land cells instead of a flood fill.
		sets := newSets(r * cols)
		land := 0
		for y := 0; y < r; y++ {
			for x := 0; x < cols; x++ {
				if rows[y][x] != '1' {
					continue
				}
				land++
				if x+1 < cols && rows[y][x+1] == '1' && sets.union(y*cols+x, y*cols+x+1) {
					land--
				}
				if y+1 < r && rows[y+1][x] == '1' && sets.union(y*cols+x, (y+1)*cols+x) {
					land--
				}
			}
		}
		if strconv.Itoa(land) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, land, c.Expected)
		}
	}
}

func parseNumberGrid(lines []string) [][]int {
	var g [][]int
	for _, line := range lines[1:] {
		var row []int
		for _, f := range strings.Fields(line) {
			v, _ := strconv.Atoi(f)
			row = append(row, v)
		}
		g = append(g, row)
	}
	return g
}

func TestSpreading(t *testing.T) {
	for i, c := range cases(t, "spreading", 50) {
		g := parseNumberGrid(strings.Split(c.Input, "\n"))
		// Simulate minute by minute.
		minutes := 0
		for {
			healthy := 0
			var next [][2]int
			for y := range g {
				for x := range g[y] {
					if g[y][x] != 1 {
						continue
					}
					healthy++
					for _, d := range steps4 {
						yy, xx := y+d[0], x+d[1]
						if yy >= 0 && yy < len(g) && xx >= 0 && xx < len(g[0]) && g[yy][xx] == 2 {
							next = append(next, [2]int{y, x})
							break
						}
					}
				}
			}
			if healthy == 0 {
				break
			}
			if len(next) == 0 {
				minutes = -1
				break
			}
			for _, p := range next {
				g[p[0]][p[1]] = 2
			}
			minutes++
		}
		if strconv.Itoa(minutes) != c.Expected {
			t.Fatalf("case %d: simulation gives %d, generator says %s", i, minutes, c.Expected)
		}
	}
}

func TestGridPath(t *testing.T) {
	for i, c := range cases(t, "grid-path", 50) {
		g := parseNumberGrid(strings.Split(c.Input, "\n"))
		r, cols := len(g), len(g[0])
		if r*cols > 400 {
			continue
		}
		// Relax every cell until nothing changes.
		const inf = 1 << 30
		dist := make([]int, r*cols)
		for k := range dist {
			dist[k] = inf
		}
		if g[0][0] == 0 {
			dist[0] = 1
		}
		for changed := true; changed; {
			changed = false
			for y := 0; y < r; y++ {
				for x := 0; x < cols; x++ {
					if g[y][x] != 0 || dist[y*cols+x] == inf {
						continue
					}
					for _, d := range steps4 {
						yy, xx := y+d[0], x+d[1]
						if yy >= 0 && yy < r && xx >= 0 && xx < cols && g[yy][xx] == 0 && dist[y*cols+x]+1 < dist[yy*cols+xx] {
							dist[yy*cols+xx] = dist[y*cols+x] + 1
							changed = true
						}
					}
				}
			}
		}
		want := dist[r*cols-1]
		if want == inf || g[r-1][cols-1] != 0 {
			want = -1
		}
		if strconv.Itoa(want) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, want, c.Expected)
		}
	}
}

func TestTwoColour(t *testing.T) {
	seen := map[string]bool{}
	for i, c := range cases(t, "two-colour", 60) {
		n, edges, _ := parseGraph(t, strings.Split(c.Input, "\n"))
		if n > 16 {
			seen[c.Expected] = true
			continue
		}
		ok := false
		for mask := 0; mask < 1<<n && !ok; mask++ {
			good := true
			for _, e := range edges {
				good = good && (mask>>e.u&1 != mask>>e.v&1)
			}
			ok = good
		}
		if boolText(ok) != c.Expected {
			t.Fatalf("case %d: trying every colouring gives %v, generator says %s", i, ok, c.Expected)
		}
		seen[c.Expected] = true
	}
	if !seen["true"] || !seen["false"] {
		t.Fatal("both answers must appear")
	}
}

func TestCourseOrder(t *testing.T) {
	seen := map[string]bool{}
	for i, c := range cases(t, "course-order", 60) {
		n, edges, _ := parseGraph(t, strings.Split(c.Input, "\n"))
		// A loop exists if some node reaches itself.
		loop := false
		for s := 0; s < n && !loop && n <= 800; s++ {
			for _, e := range edges {
				if e.u == s && reach(n, edges, e.v, true)[s] {
					loop = true
					break
				}
			}
		}
		if n > 800 {
			continue
		}
		want := "valid order"
		if loop {
			want = "[]"
		}
		if want != c.Expected {
			t.Fatalf("case %d: want %s, generator says %s", i, want, c.Expected)
		}
		seen[want] = true
	}
	if !seen["[]"] || !seen["valid order"] {
		t.Fatal("both answers must appear")
	}
}

func TestExtraLink(t *testing.T) {
	for i, c := range cases(t, "extra-link", 40) {
		n, edges, _ := parseGraph(t, strings.Split(c.Input, "\n"))
		if len(edges) != n || n > 600 {
			continue
		}
		// The last edge whose removal leaves a connected tree.
		want := edge{-1, -1, 0}
		for k := range edges {
			rest := append(slices.Clone(edges[:k]), edges[k+1:]...)
			if !slices.Contains(reach(n, rest, 0, false), false) {
				want = edges[k]
			}
		}
		if listOutput([]int{want.u, want.v}) != c.Expected {
			t.Fatalf("case %d: want %v, generator says %s", i, want, c.Expected)
		}
	}
}

func TestCloneGraph(t *testing.T) {
	for i, c := range cases(t, "clone-graph", 40) {
		n, edges, _ := parseGraph(t, strings.Split(c.Input, "\n"))
		if slices.Contains(reach(n, edges, 0, false), false) {
			t.Fatalf("case %d: graph is not connected", i)
		}
		adj := make([][]int, n)
		for j := range adj {
			adj[j] = []int{}
		}
		for _, e := range edges {
			if e.u == e.v {
				t.Fatalf("case %d: self loop", i)
			}
			adj[e.u] = append(adj[e.u], e.v)
			adj[e.v] = append(adj[e.v], e.u)
		}
		if gridOutput(adj) != c.Expected {
			t.Fatalf("case %d: adjacency differs", i)
		}
	}
}

func TestMinConnect(t *testing.T) {
	for i, c := range cases(t, "min-connect", 50) {
		n, edges, _ := parseGraph(t, strings.Split(c.Input, "\n"))
		if len(edges) > 14 {
			continue
		}
		best := -1
		for mask := 0; mask < 1<<len(edges); mask++ {
			if bits.OnesCount(uint(mask)) != n-1 {
				continue
			}
			var chosen []edge
			total := 0
			for k, e := range edges {
				if mask>>k&1 == 1 {
					chosen = append(chosen, e)
					total += e.w
				}
			}
			if !slices.Contains(reach(n, chosen, 0, false), false) && (best < 0 || total < best) {
				best = total
			}
		}
		if n == 1 {
			best = 0
		}
		if strconv.Itoa(best) != c.Expected {
			t.Fatalf("case %d: every subset gives %d, generator says %s", i, best, c.Expected)
		}
	}
}

func TestNetworkDelay(t *testing.T) {
	for i, c := range cases(t, "network-delay", 50) {
		n, edges, rest := parseGraph(t, strings.Split(c.Input, "\n"))
		s, _ := strconv.Atoi(rest[0])
		if n > 120 {
			continue
		}
		// Floyd-Warshall.
		const inf = 1 << 40
		d := make([][]int, n)
		for a := range d {
			d[a] = make([]int, n)
			for b := range d[a] {
				if a != b {
					d[a][b] = inf
				}
			}
		}
		for _, e := range edges {
			d[e.u][e.v] = min(d[e.u][e.v], e.w)
		}
		for k := 0; k < n; k++ {
			for a := 0; a < n; a++ {
				for b := 0; b < n; b++ {
					d[a][b] = min(d[a][b], d[a][k]+d[k][b])
				}
			}
		}
		want := 0
		for _, v := range d[s] {
			if v >= inf {
				want = -1
				break
			}
			want = max(want, v)
		}
		if strconv.Itoa(want) != c.Expected {
			t.Fatalf("case %d: Floyd-Warshall gives %d, generator says %s", i, want, c.Expected)
		}
	}
}

func TestCheapestStops(t *testing.T) {
	for i, c := range cases(t, "cheapest-stops", 50) {
		n, edges, rest := parseGraph(t, strings.Split(c.Input, "\n"))
		f := strings.Fields(rest[0])
		src, _ := strconv.Atoi(f[0])
		dst, _ := strconv.Atoi(f[1])
		k, _ := strconv.Atoi(f[2])
		if n > 12 || k > 5 {
			continue
		}
		// Try every walk of at most k+1 flights.
		best := -1
		var walk func(u, flights, cost int)
		walk = func(u, flights, cost int) {
			if u == dst {
				if best < 0 || cost < best {
					best = cost
				}
			}
			if flights == k+1 {
				return
			}
			for _, e := range edges {
				if e.u == u {
					walk(e.v, flights+1, cost+e.w)
				}
			}
		}
		walk(src, 0, 0)
		if strconv.Itoa(best) != c.Expected {
			t.Fatalf("case %d: every walk gives %d, generator says %s", i, best, c.Expected)
		}
	}
}

func TestWordLadder(t *testing.T) {
	found := 0
	for i, c := range cases(t, "word-ladder", 40) {
		lines := strings.Split(c.Input, "\n")
		be := strings.Fields(lines[0])
		begin, end := be[0], be[1]
		list := strings.Fields(lines[2])
		if len(list) > 600 || begin == end {
			if begin == end {
				t.Fatalf("case %d: start equals end", i)
			}
			continue
		}
		oneApart := func(a, b string) bool {
			diff := 0
			for j := range a {
				if a[j] != b[j] {
					diff++
				}
			}
			return diff == 1
		}
		want := 0
		if slices.Contains(list, end) {
			dist := map[string]int{begin: 1}
			queue := []string{begin}
			for k := 0; k < len(queue); k++ {
				for _, w := range list {
					if _, ok := dist[w]; !ok && oneApart(queue[k], w) {
						dist[w] = dist[queue[k]] + 1
						queue = append(queue, w)
					}
				}
			}
			want = dist[end]
		}
		if strconv.Itoa(want) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, want, c.Expected)
		}
		if want > 0 {
			found++
		}
	}
	if found == 0 {
		t.Fatal("no case has a ladder")
	}
}
