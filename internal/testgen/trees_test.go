package testgen

import (
	"slices"
	"strconv"
	"strings"
	"testing"
)

func sameShape(a, b *treeNode) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.val == b.val && sameShape(a.left, b.left) && sameShape(a.right, b.right)
}

func twoTrees(c Case) (*treeNode, *treeNode) {
	lines := strings.Split(c.Input, "\n")
	return parseTree(lines[0]), parseTree(lines[1])
}

func TestSameTree(t *testing.T) {
	seen := map[string]bool{}
	for i, c := range cases(t, "same-tree", 60) {
		a, b := twoTrees(c)
		if want := boolText(sameShape(a, b)); want != c.Expected {
			t.Fatalf("case %d: want %s, generator says %s", i, want, c.Expected)
		}
		seen[c.Expected] = true
	}
	if !seen["true"] || !seen["false"] {
		t.Fatal("both answers must appear")
	}
}

func flip(t *treeNode) *treeNode {
	if t == nil {
		return nil
	}
	return &treeNode{val: t.val, left: flip(t.right), right: flip(t.left)}
}

func TestMirrorTree(t *testing.T) {
	for i, c := range cases(t, "mirror-tree", 50) {
		if want := treeOutput(flip(parseTree(c.Input))); want != c.Expected {
			t.Fatalf("case %d: mirror differs", i)
		}
	}
}

func anySame(root, sub *treeNode) bool {
	if root == nil {
		return false
	}
	return sameShape(root, sub) || anySame(root.left, sub) || anySame(root.right, sub)
}

func TestSubtree(t *testing.T) {
	seen := map[string]bool{}
	for i, c := range cases(t, "subtree", 60) {
		root, sub := twoTrees(c)
		if sub == nil {
			t.Fatalf("case %d: empty subtree", i)
		}
		if want := boolText(anySame(root, sub)); want != c.Expected {
			t.Fatalf("case %d: want %s, generator says %s", i, want, c.Expected)
		}
		seen[c.Expected] = true
	}
	if !seen["true"] || !seen["false"] {
		t.Fatal("both answers must appear")
	}
}

func TestTreeDiameter(t *testing.T) {
	for i, c := range cases(t, "tree-diameter", 50) {
		root := parseTree(c.Input)
		if countNodes(root) > 1500 {
			continue
		}
		best := 0
		var walk func(n *treeNode)
		walk = func(n *treeNode) {
			if n == nil {
				return
			}
			best = max(best, recursiveDepth(n.left)+recursiveDepth(n.right)+1)
			walk(n.left)
			walk(n.right)
		}
		walk(root)
		if strconv.Itoa(best) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, best, c.Expected)
		}
	}
}

// byDepth collects values depth first; rows come out in left-to-right order.
func byDepth(n *treeNode, d int, rows *[][]int) {
	if n == nil {
		return
	}
	if d == len(*rows) {
		*rows = append(*rows, nil)
	}
	(*rows)[d] = append((*rows)[d], n.val)
	byDepth(n.left, d+1, rows)
	byDepth(n.right, d+1, rows)
}

func TestLevelOrderAndRightView(t *testing.T) {
	for i, c := range cases(t, "level-order", 50) {
		var rows [][]int
		byDepth(parseTree(c.Input), 0, &rows)
		if gridOutput(rows) != c.Expected {
			t.Fatalf("level-order case %d differs", i)
		}
	}
	for i, c := range cases(t, "right-view", 50) {
		var view []int
		var walk func(n *treeNode, d int)
		walk = func(n *treeNode, d int) {
			if n == nil {
				return
			}
			if d == len(view) {
				view = append(view, n.val)
			}
			walk(n.right, d+1)
			walk(n.left, d+1)
		}
		walk(parseTree(c.Input), 0)
		if listOutput(view) != c.Expected {
			t.Fatalf("right-view case %d differs", i)
		}
	}
}

func TestSortedTree(t *testing.T) {
	for i, c := range cases(t, "sorted-tree", 50) {
		vals, _ := parseArray(t, c.Input)
		out := strings.ReplaceAll(strings.Trim(c.Expected, "[]"), ", ", " ")
		var check func(n *treeNode, vals []int) bool
		check = func(n *treeNode, vals []int) bool {
			if len(vals) == 0 {
				return n == nil
			}
			mid := (len(vals) - 1) / 2
			return n != nil && n.val == vals[mid] && check(n.left, vals[:mid]) && check(n.right, vals[mid+1:])
		}
		if !check(parseTree(out), vals) {
			t.Fatalf("case %d: tree does not follow the middle rule", i)
		}
	}
}

func withinBounds(n *treeNode, lo, hi int) bool {
	if n == nil {
		return true
	}
	return lo < n.val && n.val < hi && withinBounds(n.left, lo, n.val) && withinBounds(n.right, n.val, hi)
}

func TestValidBST(t *testing.T) {
	seen := map[string]bool{}
	for i, c := range cases(t, "valid-bst", 60) {
		root := parseTree(c.Input)
		if want := boolText(withinBounds(root, -1<<40, 1<<40)); want != c.Expected {
			t.Fatalf("case %d: want %s, generator says %s", i, want, c.Expected)
		}
		seen[c.Expected] = true
	}
	if !seen["true"] || !seen["false"] {
		t.Fatal("both answers must appear")
	}
}

func TestKthBST(t *testing.T) {
	for i, c := range cases(t, "kth-bst", 50) {
		lines := strings.Split(c.Input, "\n")
		root := parseTree(lines[0])
		k, _ := strconv.Atoi(lines[1])
		var vals []int
		for _, n := range nodesOf(root) {
			vals = append(vals, n.val)
		}
		slices.Sort(vals)
		if !withinBounds(root, -1<<40, 1<<40) || k < 1 || k > len(vals) || strconv.Itoa(vals[k-1]) != c.Expected {
			t.Fatalf("case %d: k-th smallest differs", i)
		}
	}
}

func lcaOf(n *treeNode, p, q int) *treeNode {
	if n == nil || n.val == p || n.val == q {
		return n
	}
	l, r := lcaOf(n.left, p, q), lcaOf(n.right, p, q)
	switch {
	case l == nil:
		return r
	case r == nil:
		return l
	}
	return n
}

func TestTreeLCA(t *testing.T) {
	for i, c := range cases(t, "tree-lca", 60) {
		lines := strings.Split(c.Input, "\n")
		root := parseTree(lines[0])
		pq := strings.Fields(lines[1])
		p, _ := strconv.Atoi(pq[0])
		q, _ := strconv.Atoi(pq[1])
		if got := lcaOf(root, p, q); got == nil || strconv.Itoa(got.val) != c.Expected || p == q {
			t.Fatalf("case %d: lowest common ancestor differs", i)
		}
	}
}

func TestPathSum(t *testing.T) {
	for i, c := range cases(t, "path-sum", 50) {
		root := parseTree(c.Input)
		best := root.val
		var gain func(n *treeNode) int
		gain = func(n *treeNode) int {
			if n == nil {
				return 0
			}
			l, r := max(0, gain(n.left)), max(0, gain(n.right))
			best = max(best, n.val+l+r)
			return n.val + max(l, r)
		}
		gain(root)
		if strconv.Itoa(best) != c.Expected {
			t.Fatalf("case %d: want %d, generator says %s", i, best, c.Expected)
		}
	}
}

func TestTreeCodec(t *testing.T) {
	for i, c := range cases(t, "tree-codec", 50) {
		if treeOutput(parseTree(c.Input)) != c.Expected {
			t.Fatalf("case %d: tree does not read back the same", i)
		}
	}
}
