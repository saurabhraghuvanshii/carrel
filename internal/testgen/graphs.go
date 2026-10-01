package testgen

import (
	"container/heap"
	"math/rand"
	"slices"
	"strconv"
	"strings"
)

// Generators for the "Graphs" group. A graph is "n m" then m edges "u v"
// (or "u v w" when weighted); nodes are 0 to n-1.

func init() {
	register("path-exists", genPathExists)
	register("provinces", genProvinces)
	register("islands", genIslands)
	register("spreading", genSpreading)
	register("grid-path", genGridPath)
	register("two-colour", genTwoColour)
	register("course-order", genCourseOrder)
	register("extra-link", genExtraLink)
	register("clone-graph", genCloneGraph)
	register("min-connect", genMinConnect)
	register("network-delay", genNetworkDelay)
	register("cheapest-stops", genCheapestStops)
	register("word-ladder", genWordLadder)
}

type edge struct{ u, v, w int }

func graphInput(n int, edges []edge, weighted bool) string {
	var b strings.Builder
	b.WriteString(strconv.Itoa(n) + " " + strconv.Itoa(len(edges)))
	for _, e := range edges {
		b.WriteString("\n" + strconv.Itoa(e.u) + " " + strconv.Itoa(e.v))
		if weighted {
			b.WriteString(" " + strconv.Itoa(e.w))
		}
	}
	return b.String()
}

// randomEdges returns m undirected edges between different nodes; repeats
// are possible.
func randomEdges(rng *rand.Rand, n, m int) []edge {
	var edges []edge
	for len(edges) < m && n > 1 {
		u, v := rng.Intn(n), rng.Intn(n)
		if u != v {
			edges = append(edges, edge{u, v, 1 + rng.Intn(1000)})
		}
	}
	return edges
}

// parentOf is a union-find over n nodes.
type parentOf []int

func newSets(n int) parentOf {
	p := make(parentOf, n)
	for i := range p {
		p[i] = i
	}
	return p
}

func (p parentOf) find(x int) int {
	for p[x] != x {
		p[x] = p[p[x]]
		x = p[x]
	}
	return x
}

func (p parentOf) union(a, b int) bool {
	a, b = p.find(a), p.find(b)
	if a == b {
		return false
	}
	p[a] = b
	return true
}

// path-exists: an undirected graph and two nodes; the answer is whether a
// path joins them.
func genPathExists(rng *rand.Rand, size int) Case {
	n := max(1, size)
	edges := randomEdges(rng, n, rng.Intn(n+n/2+1))
	s, d := rng.Intn(n), rng.Intn(n)
	sets := newSets(n)
	for _, e := range edges {
		sets.union(e.u, e.v)
	}
	return Case{Input: graphInput(n, edges, false) + "\n" + strconv.Itoa(s) + " " + strconv.Itoa(d), Expected: boolText(sets.find(s) == sets.find(d))}
}

// provinces: an undirected graph; the answer is how many connected groups.
func genProvinces(rng *rand.Rand, size int) Case {
	n := max(1, size)
	edges := randomEdges(rng, n, rng.Intn(n+1))
	sets := newSets(n)
	groups := n
	for _, e := range edges {
		if sets.union(e.u, e.v) {
			groups--
		}
	}
	return Case{Input: graphInput(n, edges, false), Expected: strconv.Itoa(groups)}
}

func gridSize(rng *rand.Rand, size int) (int, int) {
	side := max(1, min(60, size/4+1))
	return 1 + rng.Intn(side), 1 + rng.Intn(side)
}

func cellGrid(rng *rand.Rand, r, c int, weights []int) [][]int {
	g := make([][]int, r)
	total := 0
	for _, w := range weights {
		total += w
	}
	for i := range g {
		g[i] = make([]int, c)
		for j := range g[i] {
			x := rng.Intn(total)
			for v, w := range weights {
				if x < w {
					g[i][j] = v
					break
				}
				x -= w
			}
		}
	}
	return g
}

var steps4 = [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

// islands: rows of 0 (water) and 1 (land); the answer is how many islands,
// joined up, down, left and right.
func genIslands(rng *rand.Rand, size int) Case {
	r, c := gridSize(rng, size)
	g := cellGrid(rng, r, c, []int{1 + rng.Intn(4), 1 + rng.Intn(4)})
	seen := make([][]bool, r)
	for i := range seen {
		seen[i] = make([]bool, c)
	}
	count := 0
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			if g[i][j] != 1 || seen[i][j] {
				continue
			}
			count++
			stack := [][2]int{{i, j}}
			seen[i][j] = true
			for len(stack) > 0 {
				p := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				for _, d := range steps4 {
					y, x := p[0]+d[0], p[1]+d[1]
					if y >= 0 && y < r && x >= 0 && x < c && g[y][x] == 1 && !seen[y][x] {
						seen[y][x] = true
						stack = append(stack, [2]int{y, x})
					}
				}
			}
		}
	}
	rows := make([]string, r)
	for i := range g {
		var b strings.Builder
		for _, v := range g[i] {
			b.WriteByte(byte('0' + v))
		}
		rows[i] = b.String()
	}
	return Case{Input: strconv.Itoa(r) + " " + strconv.Itoa(c) + "\n" + strings.Join(rows, "\n"), Expected: strconv.Itoa(count)}
}

func numberGridInput(g [][]int) string {
	var b strings.Builder
	b.WriteString(strconv.Itoa(len(g)) + " " + strconv.Itoa(len(g[0])))
	for _, row := range g {
		b.WriteString("\n" + joinInts(row, " "))
	}
	return b.String()
}

// spreading: 0 empty, 1 healthy, 2 infected; each minute infection spreads
// to side neighbours. The answer is the minutes until nobody is healthy, or
// -1 if someone never gets infected.
func genSpreading(rng *rand.Rand, size int) Case {
	r, c := gridSize(rng, size)
	g := cellGrid(rng, r, c, []int{1 + rng.Intn(3), 2 + rng.Intn(6), 1})
	dist := make([][]int, r)
	var queue [][2]int
	for i := range g {
		dist[i] = make([]int, c)
		for j := range g[i] {
			dist[i][j] = -1
			if g[i][j] == 2 {
				dist[i][j] = 0
				queue = append(queue, [2]int{i, j})
			}
		}
	}
	for k := 0; k < len(queue); k++ {
		p := queue[k]
		for _, d := range steps4 {
			y, x := p[0]+d[0], p[1]+d[1]
			if y >= 0 && y < r && x >= 0 && x < c && g[y][x] == 1 && dist[y][x] < 0 {
				dist[y][x] = dist[p[0]][p[1]] + 1
				queue = append(queue, [2]int{y, x})
			}
		}
	}
	minutes := 0
	for i := range g {
		for j := range g[i] {
			if g[i][j] == 1 {
				if dist[i][j] < 0 {
					minutes = -1
				} else if minutes >= 0 {
					minutes = max(minutes, dist[i][j])
				}
			}
		}
	}
	return Case{Input: numberGridInput(g), Expected: strconv.Itoa(minutes)}
}

// grid-path: 0 open, 1 wall; the answer is the number of cells on the
// shortest side-step path from the top-left to the bottom-right, or -1.
func genGridPath(rng *rand.Rand, size int) Case {
	r, c := gridSize(rng, size)
	g := cellGrid(rng, r, c, []int{3 + rng.Intn(5), 1 + rng.Intn(3)})
	if rng.Intn(4) != 0 {
		g[0][0], g[r-1][c-1] = 0, 0
	}
	return Case{Input: numberGridInput(g), Expected: strconv.Itoa(gridDistance(g))}
}

func gridDistance(g [][]int) int {
	r, c := len(g), len(g[0])
	if g[0][0] != 0 || g[r-1][c-1] != 0 {
		return -1
	}
	dist := make([][]int, r)
	for i := range dist {
		dist[i] = make([]int, c)
	}
	dist[0][0] = 1
	queue := [][2]int{{0, 0}}
	for k := 0; k < len(queue); k++ {
		p := queue[k]
		for _, d := range steps4 {
			y, x := p[0]+d[0], p[1]+d[1]
			if y >= 0 && y < r && x >= 0 && x < c && g[y][x] == 0 && dist[y][x] == 0 {
				dist[y][x] = dist[p[0]][p[1]] + 1
				queue = append(queue, [2]int{y, x})
			}
		}
	}
	if dist[r-1][c-1] == 0 {
		return -1
	}
	return dist[r-1][c-1]
}

// two-colour: an undirected graph built from two sides, sometimes with one
// edge inside a side; the answer is whether two colours suffice.
func genTwoColour(rng *rand.Rand, size int) Case {
	n := max(1, size)
	side := make([]int, n)
	for i := range side {
		side[i] = rng.Intn(2)
	}
	var edges []edge
	for _, e := range randomEdges(rng, n, rng.Intn(n+n/2+1)) {
		if side[e.u] != side[e.v] {
			edges = append(edges, e)
		}
	}
	if n > 3 && rng.Intn(2) == 0 {
		// A triangle away from node 0: three colours needed there.
		a, b, c := 1+rng.Intn(n-1), 1+rng.Intn(n-1), 1+rng.Intn(n-1)
		if a != b && b != c && a != c {
			edges = append(edges, edge{a, b, 1}, edge{b, c, 1}, edge{c, a, 1})
		}
	}
	rng.Shuffle(len(edges), func(i, j int) { edges[i], edges[j] = edges[j], edges[i] })
	return Case{Input: graphInput(n, edges, false), Expected: boolText(twoColourable(n, edges))}
}

func twoColourable(n int, edges []edge) bool {
	adj := make([][]int, n)
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
	}
	colour := make([]int, n)
	for i := range colour {
		colour[i] = -1
	}
	for s := 0; s < n; s++ {
		if colour[s] >= 0 {
			continue
		}
		colour[s] = 0
		queue := []int{s}
		for k := 0; k < len(queue); k++ {
			u := queue[k]
			for _, v := range adj[u] {
				if colour[v] < 0 {
					colour[v] = 1 - colour[u]
					queue = append(queue, v)
				} else if colour[v] == colour[u] {
					return false
				}
			}
		}
	}
	return true
}

// course-order: "u v" means u must come before v. Edges follow a random
// ranking, sometimes with one edge against it. The answer is "valid order"
// when an order exists, else "[]"; the driver checks the learner's order.
func genCourseOrder(rng *rand.Rand, size int) Case {
	n := max(1, size)
	rank := rng.Perm(n)
	var edges []edge
	for _, e := range randomEdges(rng, n, rng.Intn(2*n+1)) {
		if rank[e.u] < rank[e.v] {
			edges = append(edges, e)
		} else {
			edges = append(edges, edge{e.v, e.u, 1})
		}
	}
	if n > 1 && rng.Intn(3) == 0 {
		u, v := rng.Intn(n), rng.Intn(n)
		if u != v {
			if rank[u] < rank[v] {
				u, v = v, u
			}
			edges = append(edges, edge{u, v, 1}) // against the ranking
		}
	}
	want := "[]"
	if courseOrderExists(n, edges) {
		want = "valid order"
	}
	return Case{Input: graphInput(n, edges, false), Expected: want}
}

func courseOrderExists(n int, edges []edge) bool {
	indeg := make([]int, n)
	adj := make([][]int, n)
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], e.v)
		indeg[e.v]++
	}
	var queue []int
	for i, d := range indeg {
		if d == 0 {
			queue = append(queue, i)
		}
	}
	for k := 0; k < len(queue); k++ {
		for _, v := range adj[queue[k]] {
			if indeg[v]--; indeg[v] == 0 {
				queue = append(queue, v)
			}
		}
	}
	return len(queue) == n
}

// extra-link: a tree on n nodes plus one extra edge, in random order; the
// answer is the edge to remove, the last one in the input if several work.
func genExtraLink(rng *rand.Rand, size int) Case {
	n := max(3, size)
	adjacent := map[[2]int]bool{}
	var edges []edge
	for v := 1; v < n; v++ {
		u := rng.Intn(v)
		if rng.Intn(3) == 0 && v > 1 {
			u = v - 1 - rng.Intn(min(v, 3)) // some long thin paths
		}
		edges = append(edges, edge{u, v, 1})
		adjacent[[2]int{u, v}], adjacent[[2]int{v, u}] = true, true
	}
	for {
		u, v := rng.Intn(n), rng.Intn(n)
		if u != v && !adjacent[[2]int{u, v}] {
			edges = append(edges, edge{u, v, 1})
			break
		}
	}
	rng.Shuffle(len(edges), func(i, j int) { edges[i], edges[j] = edges[j], edges[i] })
	for i := range edges {
		if rng.Intn(2) == 0 {
			edges[i].u, edges[i].v = edges[i].v, edges[i].u
		}
	}
	sets := newSets(n)
	var extra edge
	for _, e := range edges {
		if !sets.union(e.u, e.v) {
			extra = e
		}
	}
	return Case{Input: graphInput(n, edges, false), Expected: listOutput([]int{extra.u, extra.v})}
}

// clone-graph: a connected undirected graph; the answer lists every node's
// neighbours in the order the edges add them, as the copy must have.
func genCloneGraph(rng *rand.Rand, size int) Case {
	n := max(1, min(size, 2000))
	var edges []edge
	for v := 1; v < n; v++ {
		edges = append(edges, edge{rng.Intn(v), v, 1})
	}
	edges = append(edges, randomEdges(rng, n, rng.Intn(n))...)
	seen := map[[2]int]bool{}
	var simple []edge
	for _, e := range edges {
		k := [2]int{min(e.u, e.v), max(e.u, e.v)}
		if !seen[k] {
			seen[k] = true
			simple = append(simple, e)
		}
	}
	rng.Shuffle(len(simple), func(i, j int) { simple[i], simple[j] = simple[j], simple[i] })
	adj := make([][]int, n)
	for _, e := range simple {
		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
	}
	for i := range adj {
		if adj[i] == nil {
			adj[i] = []int{}
		}
	}
	return Case{Input: graphInput(n, simple, false), Expected: gridOutput(adj)}
}

// connectedWeighted makes a connected graph with weights 1 to 1000.
func connectedWeighted(rng *rand.Rand, n int) []edge {
	var edges []edge
	for v := 1; v < n; v++ {
		edges = append(edges, edge{rng.Intn(v), v, 1 + rng.Intn(1000)})
	}
	edges = append(edges, randomEdges(rng, n, rng.Intn(2*n))...)
	rng.Shuffle(len(edges), func(i, j int) { edges[i], edges[j] = edges[j], edges[i] })
	return edges
}

// min-connect: a connected weighted graph; the answer is the smallest total
// weight of edges that keep every node connected.
func genMinConnect(rng *rand.Rand, size int) Case {
	n := max(1, size)
	edges := connectedWeighted(rng, n)
	s := slices.Clone(edges)
	slices.SortFunc(s, func(a, b edge) int { return a.w - b.w })
	sets := newSets(n)
	total := 0
	for _, e := range s {
		if sets.union(e.u, e.v) {
			total += e.w
		}
	}
	return Case{Input: graphInput(n, edges, true), Expected: strconv.Itoa(total)}
}

type distItem struct{ node, dist int }
type distHeap []distItem

func (h distHeap) Len() int           { return len(h) }
func (h distHeap) Less(i, j int) bool { return h[i].dist < h[j].dist }
func (h distHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *distHeap) Push(x any)        { *h = append(*h, x.(distItem)) }
func (h *distHeap) Pop() any {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

func shortestFrom(n int, edges []edge, s int) []int {
	adj := make([][]edge, n)
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], e)
	}
	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}
	h := &distHeap{{s, 0}}
	for h.Len() > 0 {
		it := heap.Pop(h).(distItem)
		if dist[it.node] >= 0 {
			continue
		}
		dist[it.node] = it.dist
		for _, e := range adj[it.node] {
			if dist[e.v] < 0 {
				heap.Push(h, distItem{e.v, it.dist + e.w})
			}
		}
	}
	return dist
}

// network-delay: a directed weighted graph and a source; the answer is how
// long until every node hears the signal, or -1 if some never does.
func genNetworkDelay(rng *rand.Rand, size int) Case {
	n := max(1, size)
	var edges []edge
	if rng.Intn(3) != 0 {
		order := rng.Perm(n)
		for i := 1; i < n; i++ { // a spanning arborescence from order[0]
			edges = append(edges, edge{order[rng.Intn(i)], order[i], 1 + rng.Intn(100)})
		}
	}
	for _, e := range randomEdges(rng, n, rng.Intn(2*n+1)) {
		e.w = 1 + rng.Intn(100)
		edges = append(edges, e)
	}
	rng.Shuffle(len(edges), func(i, j int) { edges[i], edges[j] = edges[j], edges[i] })
	s := rng.Intn(n)
	best := 0
	for _, d := range shortestFrom(n, edges, s) {
		if d < 0 {
			best = -1
			break
		}
		best = max(best, d)
	}
	return Case{Input: graphInput(n, edges, true) + "\n" + strconv.Itoa(s), Expected: strconv.Itoa(best)}
}

// cheapest-stops: directed flights with prices, a start, an end and k; the
// answer is the cheapest price using at most k stops in between, or -1.
func genCheapestStops(rng *rand.Rand, size int) Case {
	n := max(2, min(size, 40))
	var edges []edge
	for _, e := range randomEdges(rng, n, rng.Intn(4*n+1)) {
		e.w = 1 + rng.Intn(1000)
		edges = append(edges, e)
	}
	src, dst := rng.Intn(n), rng.Intn(n)
	for dst == src {
		dst = rng.Intn(n)
	}
	k := rng.Intn(min(n, 3)) // small limits make the stop rule matter
	if rng.Intn(4) == 0 {
		k = rng.Intn(n)
	}
	if n >= k+3 && rng.Intn(2) == 0 {
		// A cheap route with one stop too many next to a pricier one that fits.
		var others []int
		for _, v := range rng.Perm(n) {
			if v != src && v != dst {
				others = append(others, v)
			}
		}
		at := src
		for _, v := range others[:k+1] {
			edges = append(edges, edge{at, v, 1})
			at = v
		}
		edges = append(edges, edge{at, dst, 1}, edge{src, dst, 900 + rng.Intn(100)})
		rng.Shuffle(len(edges), func(i, j int) { edges[i], edges[j] = edges[j], edges[i] })
	}
	return Case{
		Input:    graphInput(n, edges, true) + "\n" + strconv.Itoa(src) + " " + strconv.Itoa(dst) + " " + strconv.Itoa(k),
		Expected: strconv.Itoa(cheapestWithStops(n, edges, src, dst, k)),
	}
}

// cheapestWithStops relaxes every flight k+1 times, from the previous round.
func cheapestWithStops(n int, edges []edge, src, dst, k int) int {
	const inf = 1 << 60
	cost := make([]int, n)
	for i := range cost {
		cost[i] = inf
	}
	cost[src] = 0
	for round := 0; round <= k; round++ {
		next := slices.Clone(cost)
		for _, e := range edges {
			if cost[e.u] < inf && cost[e.u]+e.w < next[e.v] {
				next[e.v] = cost[e.u] + e.w
			}
		}
		cost = next
	}
	if cost[dst] == inf {
		return -1
	}
	return cost[dst]
}

// word-ladder: a start word, an end word and a word list; changing one
// letter at a time through listed words, the answer is how many words the
// shortest chain has (start and end included), or 0.
func genWordLadder(rng *rand.Rand, size int) Case {
	length := 3 + rng.Intn(3)
	letters := "abcdefghijklmnopqrstuvwxyz"[:3+rng.Intn(4)]
	n := max(1, min(size, 3000))
	set := map[string]bool{}
	var list []string
	for len(list) < n && len(set) < pow(len(letters), length) {
		w := randomWord(rng, length, letters)
		if !set[w] {
			set[w] = true
			list = append(list, w)
		}
	}
	end := list[rng.Intn(len(list))]
	if rng.Intn(5) == 0 {
		end = randomWord(rng, length, letters) // may be missing from the list
	}
	begin := randomWord(rng, length, letters)
	for begin == end {
		begin = randomWord(rng, length, letters)
	}
	if set[end] && rng.Intn(3) != 0 {
		// Plant a ladder: walk from begin to end one letter at a time.
		w := []byte(begin)
		for string(w) != end {
			i := rng.Intn(length)
			if w[i] == end[i] {
				continue
			}
			w[i] = end[i]
			if !set[string(w)] {
				set[string(w)] = true
				list = append(list, string(w))
			}
		}
		rng.Shuffle(len(list), func(i, j int) { list[i], list[j] = list[j], list[i] })
	}
	return Case{
		Input:    begin + " " + end + "\n" + strconv.Itoa(len(list)) + "\n" + strings.Join(list, " "),
		Expected: strconv.Itoa(ladderLength(begin, end, list)),
	}
}

func pow(b, e int) int {
	r := 1
	for ; e > 0; e-- {
		r *= b
	}
	return r
}

func ladderLength(begin, end string, list []string) int {
	words := map[string]bool{}
	for _, w := range list {
		words[w] = true
	}
	if !words[end] {
		return 0
	}
	dist := map[string]int{begin: 1}
	queue := []string{begin}
	for k := 0; k < len(queue); k++ {
		w := queue[k]
		if w == end {
			return dist[w]
		}
		b := []byte(w)
		for i := range b {
			keep := b[i]
			for ch := byte('a'); ch <= 'z'; ch++ {
				b[i] = ch
				next := string(b)
				if words[next] && dist[next] == 0 {
					dist[next] = dist[w] + 1
					queue = append(queue, next)
				}
			}
			b[i] = keep
		}
	}
	return 0
}
