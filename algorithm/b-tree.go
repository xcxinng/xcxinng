package algorithm

import "fmt"

// Binary Search Tree特点是每个节点一个key(数据)，2个指针
// M-Way Search Tree 每个节点(m-1)个key，m个指针，其中key按顺序排列
type BinarySearchTreeNode struct {
	Value int
	Left  *BinarySearchTreeNode
	Right *BinarySearchTreeNode
}

type BSTService struct {
	bn *BinarySearchTreeNode
}

// 什么条件下应该触发自平衡，如何检测？
func (ts *BSTService) Insert(value int) error {
	if ts.bn == nil {
		ts.bn = &BinarySearchTreeNode{Value: value}
		return nil
	}

	node := BinarySearchTreeNode{Value: value}
	cursor := ts.bn
	for cursor != nil {
		if value < cursor.Value {
			if cursor.Left == nil {
				cursor.Left = &node
				break
			} else {
				cursor = cursor.Left
			}
		} else {
			if cursor.Right == nil {
				cursor.Right = &node
				break
			} else {
				cursor = cursor.Right
			}
		}
	}
	return nil
}

func (ts *BSTService) Traverse() {
	traverse2(ts.bn)
}

// 左序遍历
func traverse2(root *BinarySearchTreeNode) {
	if root == nil {
		return
	}
	traverse2(root.Left)
	fmt.Println(root.Value)
	traverse2(root.Right)
}

// Node represents a node in the M-Way Search Tree
type Node struct {
	values   []int   // Values stored in the node
	children []*Node // Children of the node
	isLeaf   bool    // Indicates if the node is a leaf
}

// NewNode creates a new node with the given values
func NewNode(values []int) *Node {
	return &Node{
		values:   values,
		children: make([]*Node, len(values)+1),
		isLeaf:   true,
	}
}

// MWaySearchTree represents the M-Way Search Tree
type MWaySearchTree struct {
	root *Node
	m    int // Degree of the tree
}

// NewMWaySearchTree creates a new M-Way Search Tree with the given degree
func NewMWaySearchTree(m int) *MWaySearchTree {
	return &MWaySearchTree{
		root: nil,
		m:    m,
	}
}

// insertValue inserts a value into the tree and returns the new root
func (tree *MWaySearchTree) InsertValue(value int) *Node {
	if tree.root == nil {
		// 如果树为空，创建一个新节点作为根节点
		tree.root = NewNode([]int{value})
		return tree.root
	}

	// 找到插入位置的辅助函数
	var findInsertionPlace func(node *Node, value int) (*Node, int)
	findInsertionPlace = func(node *Node, value int) (*Node, int) {
		if node.isLeaf {
			return node, 0
		}
		for i := 0; i < len(node.values); i++ {
			if value < node.values[i] {
				return findInsertionPlace(node.children[i], value)
			}
		}
		// 如果值大于或等于最后一个值，则在下一个子树中查找
		return findInsertionPlace(node.children[len(node.values)], value)
	}

	currentNode, index := findInsertionPlace(tree.root, value)

	// 插入值
	if len(currentNode.values) < tree.m-1 {
		// 如果当前节点未满，直接插入
		currentNode.values = append(currentNode.values[:index], append([]int{value}, currentNode.values[index:]...)...)
		for i := len(currentNode.values); i > index+1; i-- {
			currentNode.children[i] = currentNode.children[i-1]
		}
		currentNode.children[index+1] = nil
	} else {
		// 如果当前节点已满，分割节点
		newValues := make([]int, 0, tree.m-1)
		newValues = append(newValues, currentNode.values[:len(currentNode.values)/2]...)
		newNode := NewNode(currentNode.values[len(currentNode.values)/2:])
		for i := len(currentNode.values)/2 + 1; i <= len(currentNode.values); i++ {
			newNode.children[i-len(currentNode.values)/2] = currentNode.children[i]
		}
		currentNode.values = newValues
		currentNode.children = currentNode.children[:len(currentNode.values)+1]
		if value > currentNode.values[len(currentNode.values)-1] {
			newNode.children[0] = currentNode.children[len(currentNode.values)]
			currentNode.children[len(currentNode.values)] = newNode
		} else {
			currentNode.children[len(currentNode.values)] = newNode
		}
	}

	return tree.root
}
