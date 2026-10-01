package testgen

import (
	"math/big"
	"slices"
	"strconv"
	"strings"
	"testing"
)

type node struct {
	val  int
	next *node
}

// buildList links the values, and the last node to position pos if pos ≥ 0.
func buildList(vals []int, pos int) (*node, []*node) {
	nodes := make([]*node, len(vals))
	for i, v := range vals {
		nodes[i] = &node{val: v}
		if i > 0 {
			nodes[i-1].next = nodes[i]
		}
	}
	if len(nodes) == 0 {
		return nil, nil
	}
	if pos >= 0 {
		nodes[len(nodes)-1].next = nodes[pos]
	}
	return nodes[0], nodes
}

func values(head *node) []int {
	var out []int
	for p := head; p != nil; p = p.next {
		out = append(out, p.val)
	}
	return out
}

func TestMergeLists(t *testing.T) {
	for i, c := range cases(t, "merge-lists", 50) {
		a, rest := parseArray(t, c.Input)
		b, _ := parseArray(t, strings.Join(rest, "\n"))
		ha, _ := buildList(a, -1)
		hb, _ := buildList(b, -1)
		dummy := &node{}
		tail := dummy
		for ha != nil && hb != nil {
			if ha.val <= hb.val {
				tail.next, ha = ha, ha.next
			} else {
				tail.next, hb = hb, hb.next
			}
			tail = tail.next
		}
		if ha != nil {
			tail.next = ha
		} else {
			tail.next = hb
		}
		if listOutput(values(dummy.next)) != c.Expected {
			t.Fatalf("case %d: merge differs", i)
		}
	}
}

func TestListMiddle(t *testing.T) {
	for i, c := range cases(t, "list-middle", 50) {
		vals, _ := parseArray(t, c.Input)
		slow, _ := buildList(vals, -1)
		fast := slow
		for fast != nil && fast.next != nil {
			slow, fast = slow.next, fast.next.next
		}
		if listOutput(values(slow)) != c.Expected {
			t.Fatalf("case %d: middle differs", i)
		}
	}
}

func TestMirrorList(t *testing.T) {
	seen := map[string]bool{}
	for i, c := range cases(t, "mirror-list", 60) {
		vals, _ := parseArray(t, c.Input)
		rev := slices.Clone(vals)
		slices.Reverse(rev)
		if want := boolText(slices.Equal(vals, rev)); want != c.Expected {
			t.Fatalf("case %d: want %s, generator says %s", i, want, c.Expected)
		}
		seen[c.Expected] = true
	}
	if !seen["true"] || !seen["false"] {
		t.Fatal("both answers must appear")
	}
}

// loopStart walks the list with a visited set and returns where it first
// comes back, or -1.
func loopStart(vals []int, pos int) int {
	head, nodes := buildList(vals, pos)
	seen := map[*node]bool{}
	for p := head; p != nil; p = p.next {
		if seen[p] {
			return slices.Index(nodes, p)
		}
		seen[p] = true
	}
	return -1
}

func TestListCycleAndStart(t *testing.T) {
	loops := 0
	for _, name := range []string{"list-cycle", "cycle-start"} {
		for i, c := range cases(t, name, 60) {
			vals, rest := parseArray(t, c.Input)
			pos, _ := strconv.Atoi(rest[0])
			start := loopStart(vals, pos)
			want := strconv.Itoa(start)
			if name == "list-cycle" {
				want = boolText(start >= 0)
			}
			if want != c.Expected {
				t.Fatalf("%s case %d: want %s, generator says %s", name, i, want, c.Expected)
			}
			if start >= 0 {
				loops++
			}
		}
	}
	if loops == 0 {
		t.Fatal("no case has a loop")
	}
}

func TestRemoveNth(t *testing.T) {
	for i, c := range cases(t, "remove-nth", 60) {
		vals, rest := parseArray(t, c.Input)
		n, _ := strconv.Atoi(rest[0])
		if n < 1 || n > len(vals) {
			t.Fatalf("case %d: n=%d out of range", i, n)
		}
		// Two pointers n apart, as a learner would.
		head, _ := buildList(vals, -1)
		dummy := &node{next: head}
		lead, lag := dummy, dummy
		for j := 0; j <= n; j++ {
			lead = lead.next
		}
		for lead != nil {
			lead, lag = lead.next, lag.next
		}
		lag.next = lag.next.next
		if listOutput(values(dummy.next)) != c.Expected {
			t.Fatalf("case %d: removal differs", i)
		}
	}
}

func TestReorderList(t *testing.T) {
	for i, c := range cases(t, "reorder-list", 50) {
		vals, _ := parseArray(t, c.Input)
		var out []int
		used := make([]bool, len(vals))
		for len(out) < len(vals) {
			for j := range vals {
				if !used[j] {
					used[j] = true
					out = append(out, vals[j])
					break
				}
			}
			for j := len(vals) - 1; j >= 0 && len(out) < len(vals); j-- {
				if !used[j] {
					used[j] = true
					out = append(out, vals[j])
					break
				}
			}
		}
		if listOutput(out) != c.Expected {
			t.Fatalf("case %d: reorder differs", i)
		}
	}
}

func digitsNumber(d []int) *big.Int {
	n := new(big.Int)
	for i := len(d) - 1; i >= 0; i-- {
		n.Mul(n, big.NewInt(10))
		n.Add(n, big.NewInt(int64(d[i])))
	}
	return n
}

func TestAddDigits(t *testing.T) {
	for i, c := range cases(t, "add-digits", 50) {
		a, rest := parseArray(t, c.Input)
		b, _ := parseArray(t, strings.Join(rest, "\n"))
		for _, d := range [][]int{a, b} {
			if len(d) > 1 && d[len(d)-1] == 0 {
				t.Fatalf("case %d: leading zero", i)
			}
		}
		sum := new(big.Int).Add(digitsNumber(a), digitsNumber(b)).String()
		var want []int
		for j := len(sum) - 1; j >= 0; j-- {
			want = append(want, int(sum[j]-'0'))
		}
		if listOutput(want) != c.Expected {
			t.Fatalf("case %d: big-number sum differs", i)
		}
	}
}

func TestReverseGroups(t *testing.T) {
	for i, c := range cases(t, "reverse-groups", 60) {
		vals, rest := parseArray(t, c.Input)
		k, _ := strconv.Atoi(rest[0])
		var out []int
		for start := 0; start < len(vals); start += k {
			end := min(start+k, len(vals))
			group := slices.Clone(vals[start:end])
			if end-start == k {
				for l, r := 0, len(group)-1; l < r; l, r = l+1, r-1 {
					group[l], group[r] = group[r], group[l]
				}
			}
			out = append(out, group...)
		}
		if listOutput(out) != c.Expected {
			t.Fatalf("case %d: group reversal differs", i)
		}
	}
}
