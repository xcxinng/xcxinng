// Copyright 2023 Dmitry Dikun
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Types and methods in this package are not thread-safe.

package algorithm

type Key interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~string
}

type KeyValue[K Key, V any] struct {
	Key   K
	Value any
}

// What's the point of this thing?
type collision[V any] []V

type Iterator[K Key, V any] interface {
	Next() (KeyValue[K, V], bool)
}

const MinDegree = 3

type BPTree[K Key, V any] struct {
	root *node[K, V]
	size int
}

// NewBPTree returns a new BPTree. Degree measures the capacity of nodes, i.e. maximum
// allowed number of direct child nodes for internal nodes, and maximum key-value pairs for leaf
// nodes.
// maxDegree should be greater or equal MinOrder, otherwise BPTree will be initialized with MinDegree.
func NewBPTree[K Key, V any](maxDegree int) *BPTree[K, V] {
	if maxDegree < MinDegree {
		maxDegree = MinDegree
	}
	return &BPTree[K, V]{
		root: newLeafNode[K, V](maxDegree),
	}
}

// Clear tree.
func (t *BPTree[K, V]) Clear() {
	if t.root.isLeaf() {
		t.root = newLeafNode[K, V](cap(t.root.keys))
	} else {
		t.root = newLeafNode[K, V](cap(t.root.children))
	}
	t.size = 0
}

// Size returns a number of key-value pairs currently stored in a tree.
func (t *BPTree[K, V]) Size() int {
	return t.size
}

// Find returns a (value, true) for a given key, or (nil, false) if not found.
func (t *BPTree[K, V]) Find(key K) (V, bool) {
	if v, ok := t.find(key); ok {
		if v, ok := v.(collision[V]); ok {
			return v[0], true
		}
		return v.(V), true
	}
	var zero V
	return zero, false
}

// FindAll returns a ([]value, true) for a given key, or (nil, false) if not found.
func (t *BPTree[K, V]) FindAll(key K) ([]V, bool) {
	if v, ok := t.find(key); ok {
		if v, ok := v.(collision[V]); ok {
			return v, true
		}
		return []V{v.(V)}, true
	}
	return nil, false
}

func (t *BPTree[K, V]) find(key K) (any, bool) {
	n := t.root
NodesLoop:
	for n.isInternal() {
		for i, c := range n.children {
			if i == len(n.keys) || key < n.keys[i] {
				n = c
				continue NodesLoop
			}
		}
	}
	for i, k := range n.keys {
		if k == key {
			return n.values[i], true
		}
	}
	return nil, false
}

// Insert puts a key-value pair to the tree. If given key is present in tree, it's value will be replaced.
func (t *BPTree[K, V]) Insert(key K, val V) {
	t.insert(key, val, true)
}

// Append puts a key-value pair to the tree. If given key is present in tree, val will be appended to it's values.
func (t *BPTree[K, V]) Append(key K, val V) {
	t.insert(key, val, false)
}

func (t *BPTree[K, V]) insert(key K, val V, replace bool) {
	n := t.root
	ok, key2, n2 := n.insert(key, val, replace)
	if n2 != nil {
		if n.isLeaf() {
			t.root = newInternalNode[K, V](cap(n.keys))
		} else {
			t.root = newInternalNode[K, V](cap(n.children))
		}
		t.root.keys = t.root.keys[:1]
		t.root.keys[0] = key2
		t.root.children = t.root.children[:2]
		t.root.children[0] = n
		t.root.children[1] = n2
	}
	if ok {
		t.size++
	}
}

// Delete removes a key-value pair and returns it's (value, true) if success, or (nil, false) if not found.
// If multiply values are found, last added will be removed.
func (t *BPTree[K, V]) Delete(key K) (val V, ok bool) {
	if v, ok := t.delete(key, false, -1); ok {
		return v.(V), true
	}
	return
}

// DeleteOne is like Delete, but removes concrete value if multiply are.
func (t *BPTree[K, V]) DeleteOne(key K, idx int) (val V, ok bool) {
	if v, ok := t.delete(key, false, idx); ok {
		return v.(V), true
	}
	return
}

// DeleteAll is like Delete, but removes all values id multiply are.
func (t *BPTree[K, V]) DeleteAll(key K) (vals []V, ok bool) {
	if v, ok := t.delete(key, true, 0); ok {
		return v.(collision[V]), true
	}
	return nil, false
}

func (t *BPTree[K, V]) delete(key K, all bool, idx int) (val any, ok bool) {
	val, ok = t.root.delete(key, all, idx)
	if ok {
		if t.root.isInternal() && len(t.root.children) == 1 {
			t.root = t.root.children[0]
		}
		if all {
			c, _ := val.(collision[V])
			t.size -= len(c)
			return c, true
		} else {
			t.size--
		}
	}
	return
}

type iterator[K Key, V any] struct {
	from *K
	to   *K
	n    *node[K, V]
	i    int
	c    collision[V]
	ckey K
	ci   int
}

func (i *iterator[K, V]) Next() (KeyValue[K, V], bool) {
SEARCH:
	for i.n != nil {
		if i.c != nil {
			if i.ci < len(i.c) {
				kv := KeyValue[K, V]{Key: i.ckey, Value: i.c[i.ci]}
				i.ci++
				return kv, true
			}
			i.c = nil
		}
		for ; i.i < len(i.n.keys); i.i++ {
			k := i.n.keys[i.i]
			if i.from != nil && k < *i.from {
				continue
			}
			if i.to != nil && k >= *i.to {
				i.n = nil
				break SEARCH
			}
			if c, ok := i.n.values[i.i].(collision[V]); ok {
				i.c = c
				i.ckey = i.n.keys[i.i]
				kv := KeyValue[K, V]{Key: i.ckey, Value: c[0]}
				i.ci = 1
				i.i++
				return kv, true
			}
			kv := KeyValue[K, V]{Key: i.n.keys[i.i], Value: i.n.values[i.i]}
			i.i++
			return kv, true
		}
		i.n = i.n.right
		i.i = 0
	}
	return KeyValue[K, V]{}, false
}

// Iterator returns an Iterator for key-value pairs from interval [*from; *to). Nil given as a parameter will
// be interpreted as begin or end whole tree key diapason.
func (t *BPTree[K, V]) Iterator(from *K, to *K) Iterator[K, V] {
	if from != nil && to != nil && *from >= *to {
		return &iterator[K, V]{}
	}
	n := t.root
NodesLoop:
	for n.isInternal() {
		for i, c := range n.children {
			if from == nil || i == len(n.keys) || *from < n.keys[i] {
				n = c
				continue NodesLoop
			}
		}
	}
	return &iterator[K, V]{
		from: from,
		to:   to,
		n:    n,
	}
}

// Range returns a slice of key-value pairs from interval [*from; *to). Nil given as a parameter will
// be interpreted as begin or end whole tree key diapason. If there are no keys found, returns nil.
func (t *BPTree[K, V]) Range(from *K, to *K) []KeyValue[K, V] {
	i := t.Iterator(from, to)
	var result []KeyValue[K, V]
	for kv, ok := i.Next(); ok; kv, ok = i.Next() {
		result = append(result, kv)
	}
	return result
}

// Entries returns a slice of all key-value pairs stored in tree. If tree is empty, returns nil.
func (t *BPTree[K, V]) Entries() []KeyValue[K, V] {
	return t.Range(nil, nil)
}

// First returns (key-value, true) for the minimal key in tree, or (zero, false) if tree is empty.
func (t *BPTree[K, V]) First() (KeyValue[K, V], bool) {
	if t.size == 0 {
		return KeyValue[K, V]{}, false
	}
	n := t.root
	for n.isInternal() {
		n = n.children[0]
	}
	v := n.values[0]
	if c, ok := v.(collision[V]); ok {
		v = c[0]
	}
	return KeyValue[K, V]{Key: n.keys[0], Value: v}, true
}

// Last returns (key-value, true) for the maximal key in tree, or (zero, false) if tree is empty.
func (t *BPTree[K, V]) Last() (KeyValue[K, V], bool) {
	if t.size == 0 {
		return KeyValue[K, V]{}, false
	}
	n := t.root
	for n.isInternal() {
		n = n.children[len(n.children)-1]
	}
	v := n.values[len(n.values)-1]
	if c, ok := v.(collision[V]); ok {
		v = c[len(c)-1]
	}
	return KeyValue[K, V]{Key: n.keys[len(n.keys)-1], Value: v}, true
}
