package tree

import (
	"fmt"
	"sort"
)

type edge struct {
	label string
	n     *node
}

func newNode(depth int) *node {
	return &node{
		edges: []*edge{},
		depth: depth,
	}
}

type node struct {
	Value        interface{}
	isCollection bool
	edges        []*edge
	depth        int
}

func (n *node) Add(ks []string, v interface{}) {
	if len(ks) == 0 {
		n.flatten(v)
		return
	}

	for _, e := range n.edges {
		if e.label == ks[0] {
			e.n.Add(ks[1:], v)
			return
		}
	}

	child := newNode(n.depth + 1)
	n.edges = append(n.edges, &edge{label: ks[0], n: child})
	child.Add(ks[1:], v)
}

func (n *node) Del(ks ...string) {
	lenKs := len(ks)

	if lenKs == 0 || n.IsLeaf() {
		return
	}

	if ks[0] == wildcard {
		if lenKs > 1 {
			for _, e := range n.edges {
				e.n.Del(ks[1:]...)
			}
			return
		}

		n.edges = []*edge{}
		return
	}

	for i, e := range n.edges {
		if e.label == ks[0] {
			if lenKs == 1 {
				if i == 0 {
					n.edges = n.edges[1:]
					return
				}
				if i == len(n.edges)-1 {
					n.edges = n.edges[:i]
					return
				}
				n.edges = append(n.edges[:i], n.edges[i+1:]...)
				return
			}
			e.n.Del(ks[1:]...)
			return
		}
	}
}

func (n *node) Get(ks ...string) interface{} {
	lenKs := len(ks)
	if n.IsLeaf() && lenKs > 0 {
		return nil
	}

	if lenKs == 0 {
		return n.expand()
	}

	if lenKs == 1 {
		for _, e := range n.edges {
			if e.label == ks[0] {
				return e.n.Get()
			}
		}
		return nil
	}

	for _, e := range n.edges {
		if e.label == ks[0] {
			return e.n.Get(ks[1:]...)
		}
	}
	return nil
}

// Depth returns the node's depth.
func (n *node) Depth() int {
	return n.depth
}

func (n *node) SetDepth(d int) {
	n.depth = d
	for _, e := range n.edges {
		e.n.SetDepth(d + 1)
	}
}

// IsLeaf returns whether the node is a leaf.
func (n *node) IsLeaf() bool {
	return len(n.edges) == 0
}

func (n *node) expand() interface{} {
	children := len(n.edges)
	if children == 0 {
		return n.Value
	}

	if n.isCollection {
		res := make([]interface{}, children)
		for i, e := range n.edges {
			res[i] = e.n.Get()
		}

		return res
	}

	res := map[string]interface{}{}
	for _, e := range n.edges {
		res[e.label] = e.n.Get()
	}
	return res
}

func (n *node) flatten(i interface{}) {
	switch v := i.(type) {
	case map[string]interface{}:
		n.isCollection = false
		for k, e := range v {
			n.Add([]string{k}, e)
		}
	case []interface{}:
		n.isCollection = true
		for i, e := range v {
			n.Add([]string{fmt.Sprintf("%d", i)}, e)
		}
	default:
		n.isCollection = false
		n.Value = v
	}
}

func (n *node) sort() {
	sort.Sort(n)
	for _, e := range n.edges {
		e.n.sort()
	}
}

func (n *node) Len() int {
	return len(n.edges)
}

func (n *node) Less(i, j int) bool {
	if n.isCollection {
		return i < j
	}
	return n.edges[i].label < n.edges[j].label
}

func (n *node) Swap(i, j int) {
	n.edges[i], n.edges[j] = n.edges[j], n.edges[i]
}