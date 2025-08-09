package algorithm

import (
	"testing"
)

func TestIsSymmetric(t *testing.T) {
	root := TreeNode{Val: 2, Left: &TreeNode{Val: 3}, Right: &TreeNode{Val: 3}}
	t.Log(isSymmetric(&root))
	// output:
	// true

	root = TreeNode{Val: 2, Left: &TreeNode{Val: 3}, Right: &TreeNode{Val: 5}}
	t.Log(isSymmetric(&root))
	// output:
	// false

	root = TreeNode{Val: 2,
		Left:  &TreeNode{Val: 3, Left: &TreeNode{Val: 4}, Right: &TreeNode{Val: 5}},
		Right: &TreeNode{Val: 3, Left: nil, Right: &TreeNode{Val: 4}}}
	t.Log(isSymmetric(&root))
	// output:
	// false

}

func TestTraverseTree(t *testing.T) {
	root := TreeNode{Val: 2, Left: &TreeNode{Val: 3}, Right: &TreeNode{Val: 3}}
	traverseTree(&root)

	root = TreeNode{Val: 2,
		Left:  &TreeNode{Val: 3, Left: &TreeNode{Val: 4}, Right: &TreeNode{Val: 5}},
		Right: &TreeNode{Val: 3, Left: nil, Right: &TreeNode{Val: 4}}}
	traverseTree(&root)
}
