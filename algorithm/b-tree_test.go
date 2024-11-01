package algorithm

import (
	"fmt"
	"testing"
)

func TestTreeService_Insert(t *testing.T) {
	s := BSTService{}
	s.Insert(10)
	s.Insert(1)
	s.Insert(20)
	s.Insert(5)
	s.Insert(6)
	s.Traverse()
}

func TestMWaySearchTree_InsertValue(t *testing.T) {
	s := NewMWaySearchTree(3)
	s.InsertValue(10)
	s.InsertValue(1)
	s.InsertValue(5)
	s.InsertValue(15)
	s.InsertValue(20)
	s.InsertValue(2)
	printTree(s.root, 4)
}

func printTree(node *Node, level int) {
	if node == nil {
		return
	}

	// 递归地打印左子树
	for i := 0; i < len(node.values); i++ {
		fmt.Printf("%*s%d\n", level*4, "", node.values[i])
		printTree(node.children[i], level+1)
	}

	// 如果是最后一个值的右子树，也打印它
	if len(node.values) == 0 || node.children[len(node.values)] == nil {
		return
	}
	printTree(node.children[len(node.values)], level+1)
}
