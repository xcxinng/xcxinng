// Package algorithm records algorithms from "shopee code league" on leetcode.
package algorithm

/* Symmetric Tree

Given the root of a binary tree, check whether it is a mirror of itself (i.e., symmetric around its center).

Example 1:
Input: root = [1,2,2,3,4,4,3]
Output: true

Example 2:
Input: root = [1,2,2,null,3,null,3]
Output: false

Constraints:
The number of nodes in the tree is in the range [1, 1000].
-100 <= Node.val <= 100
*/

// isSymmetric judge a root(in this context, it represents a binary tree) is
// symmetric or not.
func isSymmetric(root *TreeNode) bool {
	return check(root, root)
}

// check works in recursion way, It judges whether p is symmetric with p.
//
// check moves p and q simultaneously to traverse the corresponding tree.
// Each moving direction is on the opposite, respectively, e.g. every
// time p move to the right subtree, p move to its left subtree.
//
// In each recursion, check will check value of the current node of
// p or q, if both are equal, it will continue,  otherwise, they are not
// symmetric trees.
//
// Q1: What's a symmetric tree?
// A1: A tree is symmetric if it's left subtree is symmetric with its
// right subtree.
//
// Q2: How do you understand a tree is symmetric with another tree?
// A2: If they meet the requirements below, we would say they are
// symmetric with each other:
//
//  1. values of root node must be the same
//  2. each value of its right subtree is mirror equal to its left subtree.
//
func check(p, q *TreeNode) bool {
	// both values are nil
	if p == nil && q == nil {
		return true
	}

	// either of (p,q) is nil
	if p == nil || q == nil {
		return false
	}

	// check values and move them to the opposite direction,respectively.
	return p.Val == q.Val && check(p.Left, q.Right) && check(p.Right, q.Left)
}

// isSymmetric2 solve the same problem as isSymmetric.
// Unlike isSymmetric, isSymmetric2 works in iteration way.
//
// It uses a queue initialized with 2 elements(both are root nodes),
// and in each iteration, pops 2 elements out of the queue, and compares
// values. And in the following, put respectively left and right sub-node
// in an opposite order into the queue.
func isSymmetric2(root *TreeNode) bool {
	u, v := root, root
	var q []*TreeNode
	q = append(q, u)
	q = append(q, v)
	for len(q) > 0 {
		u, v = q[0], q[1]
		q = q[2:]
		if u == nil && v == nil {
			continue
		}
		if u == nil || v == nil {
			return false
		}
		if u.Val != v.Val {
			return false
		}

		// For a symmetric tree, u.Left == v.Right && u.Right == v.Left.
		// Append them in opposite order into q.
		q = append(q, u.Left)
		q = append(q, v.Right)

		q = append(q, u.Right)
		q = append(q, v.Left)
	}
	return true
}

/* Binary Tree Level Order Traversal

Given the root of a binary tree, return the level order traversal of
its nodes' values. (i.e., from left to right, level by level).

Example 1:
Input: root = [3,9,20,null,null,15,7]
Output: [[3],[9,20],[15,7]]

Example 2:
Input: root = [1]
Output: [[1]]

Example 3:
Input: root = []
Output: []

Constraints:
The number of nodes in the tree is in the range [0, 2000].
-1000 <= Node.val <= 1000
*/
func traverseTree(root *TreeNode) {
	if root == nil {
		return
	}
	println(root.Val)
	traverseTree(root.Left)
	traverseTree(root.Right)
}
