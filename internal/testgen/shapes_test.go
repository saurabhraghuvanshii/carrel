package testgen

import (
	"slices"
	"strconv"
	"strings"
	"testing"
)

// simpleLRU keeps keys in a slice, most recent last.
func simpleLRU(input string) string {
	lines := strings.Split(input, "\n")[1:]
	var keys []int
	vals := map[int]int{}
	capacity := 0
	var out []string
	touch := func(k int) {
		if i := slices.Index(keys, k); i >= 0 {
			keys = slices.Delete(keys, i, i+1)
		}
		keys = append(keys, k)
	}
	for _, line := range lines {
		f := strings.Fields(line)
		a, _ := strconv.Atoi(f[1])
		switch f[0] {
		case "new":
			capacity, keys, vals = a, nil, map[int]int{}
			out = append(out, "null")
		case "get":
			if v, ok := vals[a]; ok {
				touch(a)
				out = append(out, strconv.Itoa(v))
			} else {
				out = append(out, "-1")
			}
		case "put":
			b, _ := strconv.Atoi(f[2])
			vals[a] = b
			touch(a)
			if len(keys) > capacity {
				delete(vals, keys[0])
				keys = keys[1:]
			}
			out = append(out, "null")
		}
	}
	return strings.Join(out, " ")
}

func TestLRUMatchesSimpleVersion(t *testing.T) {
	cases, err := Generate("lru-cache", 5, 50)
	if err != nil {
		t.Fatal(err)
	}
	evictions := false
	for i, c := range cases {
		lines := strings.Split(c.Input, "\n")
		k, _ := strconv.Atoi(lines[0])
		if k != len(lines)-1 || k < 5 || k > 5000 || !strings.HasPrefix(lines[1], "new ") {
			t.Fatalf("case %d: bad shape", i)
		}
		if got := simpleLRU(c.Input); got != c.Expected {
			t.Fatalf("case %d: simple LRU says %q, generator says %q", i, got, c.Expected)
		}
		if strings.Count(c.Expected, " ")+1 != k {
			t.Fatalf("case %d: want %d results", i, k)
		}
		if strings.Contains(c.Expected, " -1") {
			evictions = true
		}
	}
	if !evictions {
		t.Fatal("no get ever missed, so evictions are not being tested")
	}
}

func TestReverseListMatchesBruteForce(t *testing.T) {
	cases, err := Generate("reverse-list", 9, 60)
	if err != nil {
		t.Fatal(err)
	}
	seenEmpty, seenOne := false, false
	for i, c := range cases {
		lines := strings.Split(c.Input, "\n")
		n, _ := strconv.Atoi(lines[0])
		fields := strings.Fields(lines[1])
		if len(lines) != 2 || len(fields) != n {
			t.Fatalf("case %d: bad shape %q", i, c.Input)
		}
		slices.Reverse(fields)
		if want := "[" + strings.Join(fields, ", ") + "]"; want != c.Expected {
			t.Fatalf("case %d: want %q, generator says %q", i, want, c.Expected)
		}
		seenEmpty = seenEmpty || n == 0
		seenOne = seenOne || n == 1
	}
	if !seenEmpty || !seenOne {
		t.Fatal("empty and single-node lists must both appear")
	}
}

// parseTree reads level order with "null" for a missing child.
func parseTree(line string) *treeNode {
	tok := strings.Fields(line)
	if len(tok) == 0 || tok[0] == "null" {
		return nil
	}
	mk := func(s string) *treeNode { v, _ := strconv.Atoi(s); return &treeNode{val: v} }
	root := mk(tok[0])
	queue := []*treeNode{root}
	for i := 1; i < len(tok); {
		t := queue[0]
		queue = queue[1:]
		for _, side := range []**treeNode{&t.left, &t.right} {
			if i < len(tok) && tok[i] != "null" {
				*side = mk(tok[i])
				queue = append(queue, *side)
			}
			i++
		}
	}
	return root
}

func countNodes(t *treeNode) int {
	if t == nil {
		return 0
	}
	return 1 + countNodes(t.left) + countNodes(t.right)
}

func recursiveDepth(t *treeNode) int {
	if t == nil {
		return 0
	}
	return 1 + max(recursiveDepth(t.left), recursiveDepth(t.right))
}

func TestMaxDepthMatchesBruteForce(t *testing.T) {
	cases, err := Generate("max-depth", 3, 80)
	if err != nil {
		t.Fatal(err)
	}
	seenEmpty, seenOne, seenDeep := false, false, false
	for i, c := range cases {
		if strings.Contains(c.Input, "\n") || strings.HasSuffix(c.Input, " null") {
			t.Fatalf("case %d: tree must be one line without trailing nulls", i)
		}
		root := parseTree(c.Input)
		if levelOrder(root) != c.Input {
			t.Fatalf("case %d: the tree does not read back the same", i)
		}
		want := recursiveDepth(root)
		if strconv.Itoa(want) != c.Expected {
			t.Fatalf("case %d: brute force says %d, generator says %s", i, want, c.Expected)
		}
		n := countNodes(root)
		seenEmpty = seenEmpty || n == 0
		seenOne = seenOne || n == 1
		seenDeep = seenDeep || (n > 100 && want > n/2)
	}
	if !seenEmpty || !seenOne || !seenDeep {
		t.Fatalf("missing shapes: empty %v, single %v, deep %v", seenEmpty, seenOne, seenDeep)
	}
}

func TestNewGeneratorsAreDeterministic(t *testing.T) {
	for _, name := range []string{"lru-cache", "reverse-list", "max-depth"} {
		a, _ := Generate(name, 77, 30)
		b, _ := Generate(name, 77, 30)
		if !slices.Equal(a, b) {
			t.Fatalf("%s: same seed gave different cases", name)
		}
	}
}
