package testgen

import (
	"math/rand"
	"strconv"
	"strings"
)

func init() { register("max-depth", genMaxDepth) }

type treeNode struct {
	val         int
	left, right *treeNode
}

// max-depth: a binary tree in level order on one line, "null" for a missing
// child. An empty tree is "null". Expected is the number of nodes on the
// longest path from the root down.
func genMaxDepth(rng *rand.Rand, size int) Case {
	n := size
	switch rng.Intn(12) {
	case 0:
		n = 0
	case 1:
		n = 1
	}
	root := randomTree(rng, n, rng.Intn(4) == 0)
	return Case{Input: levelOrder(root), Expected: strconv.Itoa(depth(root))}
}

// randomTree hangs each new node on a random free child slot. With chain set
// it always hangs it under the newest node, which makes a deep, thin tree.
func randomTree(rng *rand.Rand, n int, chain bool) *treeNode {
	if n == 0 {
		return nil
	}
	newNode := func() *treeNode { return &treeNode{val: rng.Intn(201) - 100} }
	root := newNode()
	type slot struct {
		parent *treeNode
		left   bool
	}
	free := []slot{{root, true}, {root, false}}
	for i := 1; i < n; i++ {
		pick := rng.Intn(len(free))
		if chain {
			pick = len(free) - 1 - rng.Intn(2)
		}
		s := free[pick]
		free[pick] = free[len(free)-1]
		free = free[:len(free)-1]
		child := newNode()
		if s.left {
			s.parent.left = child
		} else {
			s.parent.right = child
		}
		free = append(free, slot{child, true}, slot{child, false})
	}
	return root
}

// levelOrder writes the tree breadth first and drops the trailing nulls.
func levelOrder(root *treeNode) string {
	var parts []string
	queue := []*treeNode{root}
	for len(queue) > 0 {
		t := queue[0]
		queue = queue[1:]
		if t == nil {
			parts = append(parts, "null")
			continue
		}
		parts = append(parts, strconv.Itoa(t.val))
		queue = append(queue, t.left, t.right)
	}
	for len(parts) > 1 && parts[len(parts)-1] == "null" {
		parts = parts[:len(parts)-1]
	}
	return strings.Join(parts, " ")
}

// depth walks level by level, so a chain of 10 000 nodes does not recurse.
func depth(root *treeNode) int {
	d := 0
	level := []*treeNode{root}
	for {
		var next []*treeNode
		for _, t := range level {
			if t != nil {
				next = append(next, t.left, t.right)
			}
		}
		if len(next) == 0 {
			return d
		}
		d++
		level = next
	}
}
