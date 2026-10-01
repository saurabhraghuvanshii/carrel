package testgen

import (
	"math/rand"
	"slices"
	"strconv"
	"strings"
)

// Generators for the "Trees" group. A tree is one line in level order with
// "null" for a missing child; answers print trees as [3, 9, null, 7] and an
// empty tree as [].

func init() {
	register("same-tree", genSameTree)
	register("mirror-tree", genMirrorTree)
	register("subtree", genSubtree)
	register("tree-diameter", genTreeDiameter)
	register("level-order", genLevelOrder)
	register("right-view", genRightView)
	register("sorted-tree", genSortedTree)
	register("valid-bst", genValidBST)
	register("kth-bst", genKthBST)
	register("tree-lca", genTreeLCA)
	register("path-sum", genPathSum)
	register("tree-codec", genTreeCodec)
}

func treeOutput(root *treeNode) string {
	if root == nil {
		return "[]"
	}
	return "[" + strings.ReplaceAll(levelOrder(root), " ", ", ") + "]"
}

func treeSize(rng *rand.Rand, size int, allowEmpty bool) int {
	switch {
	case allowEmpty && rng.Intn(12) == 0:
		return 0
	case rng.Intn(12) == 0:
		return 1
	}
	return max(1, size)
}

// valuedTree is a random shape with values drawn from [lo, hi].
func valuedTree(rng *rand.Rand, n, lo, hi int) *treeNode {
	root := randomTree(rng, n, rng.Intn(5) == 0)
	for _, t := range nodesOf(root) {
		t.val = lo + rng.Intn(hi-lo+1)
	}
	return root
}

// nodesOf lists the nodes in level order.
func nodesOf(root *treeNode) []*treeNode {
	if root == nil {
		return nil
	}
	out := []*treeNode{root}
	for i := 0; i < len(out); i++ {
		for _, c := range []*treeNode{out[i].left, out[i].right} {
			if c != nil {
				out = append(out, c)
			}
		}
	}
	return out
}

func cloneTree(t *treeNode) *treeNode {
	if t == nil {
		return nil
	}
	c := &treeNode{val: t.val}
	// Iterative so a 10 000-node chain does not recurse deeply.
	type pair struct{ from, to *treeNode }
	stack := []pair{{t, c}}
	for len(stack) > 0 {
		p := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if p.from.left != nil {
			p.to.left = &treeNode{val: p.from.left.val}
			stack = append(stack, pair{p.from.left, p.to.left})
		}
		if p.from.right != nil {
			p.to.right = &treeNode{val: p.from.right.val}
			stack = append(stack, pair{p.from.right, p.to.right})
		}
	}
	return c
}

// changeOnce alters a copy of t in one random way: a value, an extra leaf,
// or a removed leaf.
func changeOnce(rng *rand.Rand, t *treeNode) *treeNode {
	c := cloneTree(t)
	nodes := nodesOf(c)
	if len(nodes) == 0 {
		return &treeNode{val: rng.Intn(3)}
	}
	n := nodes[rng.Intn(len(nodes))]
	switch rng.Intn(3) {
	case 0:
		n.val += 1 + rng.Intn(2)
	case 1:
		if n.left == nil {
			n.left = &treeNode{val: n.val}
		} else if n.right == nil {
			n.right = &treeNode{val: n.val}
		} else {
			n.val--
		}
	default:
		for _, p := range nodes {
			if p.left != nil && p.left.left == nil && p.left.right == nil {
				p.left = nil
				return c
			}
		}
		n.val++
	}
	return c
}

// same-tree: two trees on two lines; half are equal copies.
func genSameTree(rng *rand.Rand, size int) Case {
	a := valuedTree(rng, treeSize(rng, size, true), 0, []int{2, 100}[rng.Intn(2)])
	b := cloneTree(a)
	switch r := rng.Intn(4); {
	case r == 0 && b != nil && b.left != nil:
		// Turn the root to the right: the in-order values stay the same but
		// the shape changes, which catches comparing values only.
		l := b.left
		b.left, l.right = l.right, b
		b = l
	case r < 2:
		b = changeOnce(rng, a)
	}
	return Case{Input: levelOrder(a) + "\n" + levelOrder(b), Expected: boolText(levelOrder(a) == levelOrder(b))}
}

func mirror(t *treeNode) *treeNode {
	c := cloneTree(t)
	for _, n := range nodesOf(c) {
		n.left, n.right = n.right, n.left
	}
	return c
}

// mirror-tree: the answer is the tree flipped left to right.
func genMirrorTree(rng *rand.Rand, size int) Case {
	t := valuedTree(rng, treeSize(rng, size, true), -1000, 1000)
	return Case{Input: levelOrder(t), Expected: treeOutput(mirror(t))}
}

// subtree: a tree and a smaller tree with small values; half the time the
// second is a real subtree of the first.
func genSubtree(rng *rand.Rand, size int) Case {
	root := valuedTree(rng, treeSize(rng, size, false), 0, 2)
	nodes := nodesOf(root)
	sub := cloneTree(nodes[rng.Intn(len(nodes))])
	if rng.Intn(2) == 0 {
		sub = changeOnce(rng, sub)
	}
	return Case{Input: levelOrder(root) + "\n" + levelOrder(sub), Expected: boolText(hasSubtree(root, sub))}
}

// hasSubtree compares the text of every subtree with the same root value and
// size as sub; only those can match.
func hasSubtree(root, sub *treeNode) bool {
	want, wantSize := levelOrder(sub), len(nodesOf(sub))
	size := map[*treeNode]int{nil: 0}
	nodes := nodesOf(root)
	for i := len(nodes) - 1; i >= 0; i-- {
		n := nodes[i]
		size[n] = 1 + size[n.left] + size[n.right]
	}
	for _, n := range nodes {
		if n.val == sub.val && size[n] == wantSize && levelOrder(n) == want {
			return true
		}
	}
	return false
}

// heights returns the height of every node, children before parents.
func heights(root *treeNode) map[*treeNode]int {
	h := map[*treeNode]int{nil: 0}
	nodes := nodesOf(root)
	for i := len(nodes) - 1; i >= 0; i-- {
		n := nodes[i]
		h[n] = 1 + max(h[n.left], h[n.right])
	}
	return h
}

// tree-diameter: the answer is the number of nodes on the longest path
// between any two nodes.
func genTreeDiameter(rng *rand.Rand, size int) Case {
	t := valuedTree(rng, treeSize(rng, size, true), -100, 100)
	if t != nil && rng.Intn(3) == 0 {
		// A lopsided root: the longest path usually stays inside the subtree.
		t = &treeNode{val: rng.Intn(201) - 100, left: t}
	}
	h := heights(t)
	best := 0
	for _, n := range nodesOf(t) {
		best = max(best, h[n.left]+h[n.right]+1)
	}
	return Case{Input: levelOrder(t), Expected: strconv.Itoa(best)}
}

func levels(root *treeNode) [][]int {
	var out [][]int
	level := []*treeNode{root}
	if root == nil {
		level = nil
	}
	for len(level) > 0 {
		var vals []int
		var next []*treeNode
		for _, n := range level {
			vals = append(vals, n.val)
			for _, c := range []*treeNode{n.left, n.right} {
				if c != nil {
					next = append(next, c)
				}
			}
		}
		out = append(out, vals)
		level = next
	}
	return out
}

// level-order: the answer is the values row by row, [[root], [next row], ...].
func genLevelOrder(rng *rand.Rand, size int) Case {
	t := valuedTree(rng, treeSize(rng, size, true), -1000, 1000)
	return Case{Input: levelOrder(t), Expected: gridOutput(levels(t))}
}

// right-view: the answer is the last value of each row.
func genRightView(rng *rand.Rand, size int) Case {
	t := valuedTree(rng, treeSize(rng, size, true), -1000, 1000)
	var view []int
	for _, row := range levels(t) {
		view = append(view, row[len(row)-1])
	}
	return Case{Input: levelOrder(t), Expected: listOutput(view)}
}

// balancedTree builds a search tree from sorted values, taking the middle (the
// left one of two) as the root each time.
func balancedTree(vals []int) *treeNode {
	if len(vals) == 0 {
		return nil
	}
	mid := (len(vals) - 1) / 2
	return &treeNode{val: vals[mid], left: balancedTree(vals[:mid]), right: balancedTree(vals[mid+1:])}
}

// sorted-tree: different values in ascending order; the answer is the
// balanced search tree built by taking the middle each time.
func genSortedTree(rng *rand.Rand, size int) Case {
	vals := sortedDistinct(rng, treeSize(rng, size, false))
	return Case{Input: arrayInput(vals), Expected: treeOutput(balancedTree(vals))}
}

// searchTree inserts different values in random order, sometimes sorted
// order so the tree is one long chain.
func searchTree(rng *rand.Rand, n int) *treeNode {
	vals := distinct(rng, n, -100_000, 100_000)
	if rng.Intn(8) == 0 {
		slices.Sort(vals)
		if n > 2000 {
			vals = vals[:2000]
		}
	}
	var root *treeNode
	for _, v := range vals {
		link := &root
		for *link != nil {
			if v < (*link).val {
				link = &(*link).left
			} else {
				link = &(*link).right
			}
		}
		*link = &treeNode{val: v}
	}
	return root
}

func inorder(root *treeNode) []int {
	var out []int
	var stack []*treeNode
	for n := root; n != nil || len(stack) > 0; {
		for n != nil {
			stack = append(stack, n)
			n = n.left
		}
		n = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		out = append(out, n.val)
		n = n.right
	}
	return out
}

// valid-bst: half the trees are search trees; the rest break the rule in one
// node, often against an ancestor further up rather than the parent.
func genValidBST(rng *rand.Rand, size int) Case {
	t := searchTree(rng, treeSize(rng, size, true))
	if nodes := nodesOf(t); len(nodes) > 1 && rng.Intn(2) == 0 {
		// The grandparent trap: a right grandchild under a left child gets a
		// value above the grandparent. Its parent is still fine with it.
		for _, p := range nodes {
			if c := p.left; c != nil && c.right != nil && c.right.right == nil && rng.Intn(2) == 0 {
				c.right.val = p.val + 1
				goto checked
			}
		}
		n := nodes[1+rng.Intn(len(nodes)-1)]
		switch rng.Intn(3) {
		case 0:
			n.val = t.val // equal to the root
		case 1:
			n.val = -n.val
		default:
			n.val += rng.Intn(2001) - 1000
		}
	}
checked:
	vals := inorder(t)
	valid := true
	for i := 1; i < len(vals); i++ {
		valid = valid && vals[i-1] < vals[i]
	}
	return Case{Input: levelOrder(t), Expected: boolText(valid)}
}

// kth-bst: a search tree and k; the answer is its k-th smallest value.
func genKthBST(rng *rand.Rand, size int) Case {
	t := searchTree(rng, treeSize(rng, size, false))
	vals := inorder(t)
	k := 1 + rng.Intn(len(vals))
	return Case{Input: levelOrder(t) + "\n" + strconv.Itoa(k), Expected: strconv.Itoa(vals[k-1])}
}

// tree-lca: a tree of different values and two of them; the answer is the
// value of their lowest common ancestor.
func genTreeLCA(rng *rand.Rand, size int) Case {
	n := max(2, treeSize(rng, size, false))
	t := randomTree(rng, n, rng.Intn(5) == 0)
	nodes := nodesOf(t)
	vals := distinct(rng, len(nodes), -100_000, 100_000)
	parent := map[*treeNode]*treeNode{}
	for i, nd := range nodes {
		nd.val = vals[i]
		for _, c := range []*treeNode{nd.left, nd.right} {
			if c != nil {
				parent[c] = nd
			}
		}
	}
	p := nodes[rng.Intn(len(nodes))]
	q := nodes[rng.Intn(len(nodes))]
	for q == p {
		q = nodes[rng.Intn(len(nodes))]
	}
	if rng.Intn(4) == 0 && parent[q] != nil {
		p = parent[q] // one is an ancestor of the other
	}
	above := map[*treeNode]bool{}
	for a := p; a != nil; a = parent[a] {
		above[a] = true
	}
	lca := q
	for !above[lca] {
		lca = parent[lca]
	}
	return Case{Input: levelOrder(t) + "\n" + strconv.Itoa(p.val) + " " + strconv.Itoa(q.val), Expected: strconv.Itoa(lca.val)}
}

// path-sum: the answer is the largest sum along any path of one or more
// nodes, going up and down through parents.
func genPathSum(rng *rand.Rand, size int) Case {
	t := valuedTree(rng, treeSize(rng, size, false), -1000, 1000)
	gain := map[*treeNode]int{nil: 0}
	nodes := nodesOf(t)
	best := t.val
	for i := len(nodes) - 1; i >= 0; i-- {
		n := nodes[i]
		l, r := max(0, gain[n.left]), max(0, gain[n.right])
		best = max(best, n.val+l+r)
		gain[n] = n.val + max(l, r)
	}
	return Case{Input: levelOrder(t), Expected: strconv.Itoa(best)}
}

// tree-codec: the learner turns the tree into text and back; the answer is
// the same tree.
func genTreeCodec(rng *rand.Rand, size int) Case {
	t := valuedTree(rng, treeSize(rng, size, true), -1000, 1000)
	return Case{Input: levelOrder(t), Expected: treeOutput(t)}
}
