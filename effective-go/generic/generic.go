package main

// A FuncTree must be created with NewTreeFunc.
type FuncTree[E any] struct {
	root *funcNode[E]
	cmp  func(E, E) int
}

func NewFuncTree[E any](cmp func(E, E) int) *FuncTree[E] {
	return &FuncTree[E]{cmp: cmp}
}

func (t *FuncTree[E]) Insert(element E) {
	t.root = t.root.insert(t.cmp, element)
}

type funcNode[E any] struct {
	value E
	left  *funcNode[E]
	right *funcNode[E]
}

func (n *funcNode[E]) insert(cmp func(E, E) int, element E) *funcNode[E] {
	if n == nil {
		return &funcNode[E]{value: element}
	}
	sign := cmp(element, n.value)
	switch {
	case sign < 0:
		n.left = n.left.insert(cmp, element)
	case sign > 0:
		n.right = n.right.insert(cmp, element)
	}
	return n
}
