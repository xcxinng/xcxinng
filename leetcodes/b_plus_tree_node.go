package algorithm

import "math"

// Node represents an B+ Tree node, including root,internal and leaf node.
type node[K Key, V any] struct {
	// Using link list to link keys together would be better for insert and delete,
	// but slice is the first choice for simplicity.
	keys   []K
	values []any

	children []*node[K, V]
	left     *node[K, V] // sibling that less than current node's min key
	right    *node[K, V] //  sibling node that greater than current node's max key

	// CANNOT UNDERSTAND
	// node's least capacity (for split) ??
	bmin int
}

func newInternalNode[K Key, V any](size int) *node[K, V] {
	return &node[K, V]{
		keys:     make([]K, 0, size-1),
		children: make([]*node[K, V], 0, size),
		bmin:     int(math.Ceil(float64(size) / 2)),
	}
}

func newLeafNode[K Key, V any](size int) *node[K, V] {
	return &node[K, V]{
		keys:   make([]K, 0, size),
		values: make([]any, 0, size),
		bmin:   int(math.Ceil(float64(size) / 2)),
	}
}

// Internal node has at least one children, but wait, what about root node, what's the differences between them?
// Internal node has no data.
func (n *node[K, V]) isInternal() bool {
	return n.children != nil
}

// Leaf node stores the real data.
func (n *node[K, V]) isLeaf() bool {
	return n.values != nil
}

// Insert a k-v pair into B+ Tree in a recursive way.
func (n *node[K, V]) insert(key K, val V, replace bool) (ok bool, key2 K, n2 *node[K, V]) {
	if n.isLeaf() {
		return n.insertToLeaf(key, val, replace)
	}

	// 遍历当前内部节点的所有子节点
	for i, childNode := range n.children {
		// i == len(n.keys) 检查是否到达了最后一个子节点
		// n.keys 是当前节点的键数组，其长度加 1 等于子节点的数量。
		// 如果 i 等于 n.keys 的长度，这意味着已经检查完所有键，并且当前键应该插入到最后一个子节点中。

		// key < n.keys[i] 检查当前键是否小于当前子节点的键
		if i == len(n.keys) || key < n.keys[i] {
			ok, key2, n2 = childNode.insert(key, val, replace)
			break
		}
	}

	if n2 != nil {
		key2, n2 = n.insertToInternal(key2, n2)
	}
	return
}

// Insert an k-v pair into a leaf node.
func (n *node[K, V]) insertToLeaf(key K, val V, replace bool) (ok bool, key2 K, n2 *node[K, V]) {
	var pos int
	for i, k := range n.keys {
		if k > key {
			break
		}

		// key already exists
		if k == key {
			if replace {
				n.values[i] = val
				return false, key2, n2
			} else {
				if c, ok := n.values[i].(collision[V]); !ok {
					c = collision[V]{n.values[i].(V), val}
					n.values[i] = c
				} else {
					n.values[i] = append(c, val)
				}
				return true, key2, n2
			}
		}

		if k < key {
			pos = i + 1
			continue
		}
	}

	// leaf node is not full
	if len(n.keys) < cap(n.keys) {
		n.keys = n.keys[:len(n.keys)+1]
		n.values = n.values[:len(n.values)+1]
		copy(n.keys[pos+1:], n.keys[pos:len(n.keys)-1])
		copy(n.values[pos+1:], n.values[pos:len(n.values)-1])
		n.keys[pos] = key
		n.values[pos] = val
		return true, key2, n2
	}

	n2 = newLeafNode[K, V](cap(n.keys))
	n2.right = n.right
	if n.right != nil {
		n.right.left = n2
	}
	n.right = n2
	n2.left = n
	n2.keys = n2.keys[:cap(n.keys)+1-n.bmin]
	n2.values = n2.values[:cap(n.values)+1-n.bmin]
	if pos < n.bmin {
		copy(n2.keys, n.keys[n.bmin-1:])
		copy(n2.values, n.values[n.bmin-1:])
		n.keys = n.keys[:n.bmin]
		n.values = n.values[:n.bmin]
		copy(n.keys[pos+1:], n.keys[pos:n.bmin-1])
		copy(n.values[pos+1:], n.values[pos:n.bmin-1])
		n.keys[pos] = key
		n.values[pos] = val
	} else {
		pos2 := pos - n.bmin
		copy(n2.keys, n.keys[n.bmin:pos])
		copy(n2.values, n.values[n.bmin:pos])
		n2.keys[pos2] = key
		n2.values[pos2] = val
		copy(n2.keys[pos2+1:], n.keys[pos:])
		copy(n2.values[pos2+1:], n.values[pos:])
		n.keys = n.keys[:n.bmin]
		n.values = n.values[:n.bmin]
	}
	trimValueSlice(n.values)
	return true, n2.keys[0], n2
}

func (n *node[K, V]) insertToInternal(key K, child *node[K, V]) (key2 K, n2 *node[K, V]) {
	var pos int
	for i, k := range n.keys {
		if k < key {
			pos = i + 1
			continue
		}
		break
	}

	childPosition := pos + 1
	if len(n.children) < cap(n.children) {
		n.keys = n.keys[:len(n.keys)+1]
		n.children = n.children[:len(n.children)+1]
		copy(n.keys[pos+1:], n.keys[pos:len(n.keys)-1])
		copy(n.children[childPosition+1:], n.children[childPosition:len(n.children)-1])
		n.keys[pos] = key
		n.children[childPosition] = child
		return
	}

	n2 = newInternalNode[K, V](cap(n.children))
	n2.right = n.right
	if n.right != nil {
		n.right.left = n2
	}

	n.right = n2
	n2.left = n
	n2.keys = n2.keys[:cap(n.keys)+1-n.bmin]
	n2.children = n2.children[:cap(n.children)+1-n.bmin]

	if pos < n.bmin-1 {
		key2 = n.keys[n.bmin-2]
		copy(n2.keys, n.keys[n.bmin-1:])
		copy(n2.children, n.children[n.bmin-1:])
		n.keys = n.keys[:n.bmin-1]
		n.children = n.children[:n.bmin]
		copy(n.keys[pos+1:], n.keys[pos:n.bmin-2])
		copy(n.children[childPosition+1:], n.children[childPosition:n.bmin-1])
		n.keys[pos] = key
		n.children[childPosition] = child
	} else if pos == n.bmin-1 {
		key2 = key
		copy(n2.keys, n.keys[n.bmin-1:])
		copy(n2.children[1:], n.children[n.bmin:])
		n2.children[0] = child
		n.keys = n.keys[:n.bmin-1]
		n.children = n.children[:n.bmin]

	} else { // pos > n.bmin-1
		key2 = n.keys[n.bmin-1]
		pos2, cpos2 := pos-n.bmin, childPosition-n.bmin
		copy(n2.keys, n.keys[n.bmin:pos])
		copy(n2.children, n.children[n.bmin:childPosition])
		n2.keys[pos2] = key
		n2.children[cpos2] = child
		copy(n2.keys[pos2+1:], n.keys[pos:])
		copy(n2.children[cpos2+1:], n.children[childPosition:])
		n.keys = n.keys[:n.bmin-1]
		n.children = n.children[:n.bmin]
	}
	trimNodeSlice(n.children)
	return
}

func (n *node[K, V]) delete(key K, all bool, idx int) (val any, ok bool) {
	if n.isLeaf() {
		return n.deleteFromLeaf(key, all, idx)
	}
	var i int
	var c *node[K, V]
	for i, c = range n.children {
		if i == len(n.keys) || key < n.keys[i] {
			val, ok = c.delete(key, all, idx)
			break
		}
	}
	if ok {
		if c.isLeaf() {
			if len(c.values) < n.bmin {
				n.balanceLeaf(i)
			}
		} else {
			if len(c.children) < n.bmin {
				n.balanceInternal(i)
			}
		}
	}
	return
}

func (n *node[K, V]) deleteFromLeaf(key K, all bool, idx int) (val any, ok bool) {
	for i, k := range n.keys {
		if k == key {
			if all {
				if c, ok := n.values[i].(collision[V]); !ok {
					val = collision[V]{n.values[i].(V)}
				} else {
					val = c
				}
			} else {
				if c, ok := n.values[i].(collision[V]); !ok {
					if idx > 0 {
						return nil, false
					}
					val = n.values[i]
				} else {
					if idx >= len(c) {
						return nil, false
					}
					var zero V
					if idx < 0 {
						val = c[len(c)-1]
						c[len(c)-1] = zero
						n.values[i] = c[:len(c)-1]
					} else {
						val = c[idx]
						copy(c[idx:], c[idx+1:])
						c[len(c)-1] = zero
						n.values[i] = c[:len(c)-1]
					}
					if len(n.values[i].(collision[V])) != 0 {
						return val, true
					}
				}
			}
			ok = true
			copy(n.keys[i:len(n.keys)-1], n.keys[i+1:len(n.keys)])
			copy(n.values[i:len(n.values)-1], n.values[i+1:len(n.values)])
			n.keys = n.keys[:len(n.keys)-1]
			n.values[len(n.values)-1] = nil
			n.values = n.values[:len(n.values)-1]
			return
		}
	}
	return
}

func (n *node[K, V]) balanceLeaf(i int) {
	c := n.children[i]
	if i != 0 && len(n.children[i-1].values) > n.bmin {
		n.keys[i-1] = c.takeFromLeftSiblingLeaf(n.children[i-1])
		return
	}
	if i != len(n.children)-1 && len(n.children[i+1].values) > n.bmin {
		n.keys[i] = c.takeFromRightSiblingLeaf(n.children[i+1])
		return
	}
	if i != 0 && (i == len(n.children)-1 || len(n.children[i-1].values) < len(n.children[i+1].values)) {
		mergeLeafs(n.children[i-1], c)
		n.deleteChild(i)
	} else {
		mergeLeafs(c, n.children[i+1])
		n.deleteChild(i + 1)
	}
}

func (n *node[K, V]) takeFromLeftSiblingLeaf(n2 *node[K, V]) K {
	n.keys = n.keys[:len(n.keys)+1]
	copy(n.keys[1:], n.keys[:len(n.keys)-1])
	n.keys[0] = n2.keys[len(n2.keys)-1]
	n2.keys = n2.keys[:len(n2.keys)-1]
	n.values = n.values[:len(n.values)+1]
	copy(n.values[1:], n.values[:len(n.values)-1])
	n.values[0] = n2.values[len(n2.values)-1]
	n2.values[len(n2.values)-1] = nil
	n2.values = n2.values[:len(n2.values)-1]
	return n.keys[0]
}

func (n *node[K, V]) takeFromRightSiblingLeaf(n2 *node[K, V]) K {
	n.keys = n.keys[:len(n.keys)+1]
	n.keys[len(n.keys)-1] = n2.keys[0]
	copy(n2.keys[:len(n2.keys)-1], n2.keys[1:len(n2.keys)])
	n2.keys = n2.keys[:len(n2.keys)-1]
	n.values = n.values[:len(n.values)+1]
	n.values[len(n.values)-1] = n2.values[0]
	copy(n2.values[:len(n2.values)-1], n2.values[1:len(n2.values)])
	n2.values[len(n2.values)-1] = nil
	n2.values = n2.values[:len(n2.values)-1]
	return n2.keys[0]
}

func (n *node[K, V]) balanceInternal(i int) {
	c := n.children[i]
	if i != 0 && len(n.children[i-1].children) > n.bmin {
		n.keys[i-1] = c.takeFromLeftSiblingInternal(n.children[i-1], n.keys[i-1])
		return
	}
	if i != len(n.children)-1 && len(n.children[i+1].children) > n.bmin {
		n.keys[i] = c.takeFromRightSiblingInternal(n.children[i+1], n.keys[i])
		return
	}
	if i != 0 && (i == len(n.children)-1 || len(n.children[i-1].children) < len(n.children[i+1].children)) {
		mergeInternal(n.children[i-1], c, n.keys[i-1])
		n.deleteChild(i)
	} else {
		mergeInternal(c, n.children[i+1], n.keys[i])
		n.deleteChild(i + 1)
	}
}

func (n *node[K, V]) takeFromLeftSiblingInternal(n2 *node[K, V], key K) K {
	n.keys = n.keys[:len(n.keys)+1]
	copy(n.keys[1:], n.keys[:len(n.keys)-1])
	mkey := n2.keys[len(n2.keys)-1]
	n.keys[0] = key
	n2.keys = n2.keys[:len(n2.keys)-1]
	n.children = n.children[:len(n.children)+1]
	copy(n.children[1:], n.children[:len(n.children)-1])
	n.children[0] = n2.children[len(n2.children)-1]
	n2.children[len(n2.children)-1] = nil
	n2.children = n2.children[:len(n2.children)-1]
	return mkey
}

func (n *node[K, V]) takeFromRightSiblingInternal(n2 *node[K, V], key K) K {
	n.keys = n.keys[:len(n.keys)+1]
	n.keys[len(n.keys)-1] = key
	mkey := n2.keys[0]
	copy(n2.keys[:len(n2.keys)-1], n2.keys[1:len(n2.keys)])
	n2.keys = n2.keys[:len(n2.keys)-1]
	n.children = n.children[:len(n.children)+1]
	n.children[len(n.children)-1] = n2.children[0]
	copy(n2.children[:len(n2.children)-1], n2.children[1:len(n2.children)])
	n2.children[len(n2.children)-1] = nil
	n2.children = n2.children[:len(n2.children)-1]
	return mkey
}

func (n *node[K, V]) deleteChild(i int) {
	copy(n.keys[i-1:len(n.keys)-1], n.keys[i:len(n.keys)])
	n.keys = n.keys[:len(n.keys)-1]
	copy(n.children[i:len(n.children)-1], n.children[i+1:len(n.children)])
	n.children[len(n.children)-1] = nil
	n.children = n.children[:len(n.children)-1]
}

func mergeLeafs[K Key, V any](l, r *node[K, V]) {
	l.right = r.right
	if r.right != nil {
		r.right.left = l
	}
	llen, rlen := len(l.keys), len(r.keys)
	l.keys = l.keys[:llen+rlen]
	copy(l.keys[llen:], r.keys)
	l.values = l.values[:llen+rlen]
	copy(l.values[llen:], r.values)
}

func mergeInternal[K Key, V any](l, r *node[K, V], key K) {
	l.right = r.right
	if r.right != nil {
		r.right.left = l
	}
	nlkeys, nlch := len(l.keys), len(l.children)
	l.keys = l.keys[:nlkeys+len(r.keys)+1]
	l.keys[nlkeys] = key
	copy(l.keys[nlkeys+1:], r.keys)
	l.children = l.children[:len(l.keys)+1]
	copy(l.children[nlch:], r.children)
}

func trimNodeSlice[K Key, V any](s []*node[K, V]) {
	s = s[len(s):cap(s)]
	if len(s) == 0 {
		return
	}
	s[0] = nil
	for i := 1; i < len(s); i *= 2 {
		copy(s[i:], s[:i])
	}
}

func trimValueSlice(s []any) {
	s = s[len(s):cap(s)]
	if len(s) == 0 {
		return
	}
	s[0] = nil
	for i := 1; i < len(s); i *= 2 {
		copy(s[i:], s[:i])
	}
}
