// Package shopee records algorithms from "shopee code league" on leetcode.
package shopee

import (
	"sort"
	"strings"
)

/*
Given an integer array nums, return all the triplets [nums[i], nums[j],nums[k]]
such that i != j, i != k, and j != k, and nums[i] + nums[j] + nums[k] == 0.

Notice that the solution set must not contain duplicate triplets.

Example 1:
Input: nums = [-1,0,1,2,-1,-4]
Output: [[-1,-1,2],[-1,0,1]]
Explanation:
nums[0] + nums[1] + nums[1] = (-1) + 0 + 1 = 0.
nums[1] + nums[2] + nums[4] = 0 + 1 + (-1) = 0.
nums[0] + nums[3] + nums[4] = (-1) + 2 + (-1) = 0.
The distinct triplets are [-1,0,1] and [-1,-1,2].
Notice that the order of the output and the order of the triplets does not matter.

Example 2:
Input: nums = [0,1,1]
Output: []
Explanation: The only possible triplet does not sum up to 0.

Example 3:
Input: nums = [0,0,0]
Output: [[0,0,0]]
Explanation: The only possible triplet sums up to 0.

*/

// method_1: sort + "two pointer"
//
// Solution steps (Based on the traditional triple for-loop solution):
//
//  1. As the algorithm required, tuples should be unique, so tuple
//  (a,b,c) must meet the requirement of: a<=b<=c, which means we
//  have to sort nums.
//
//  2. In the meanwhile, in each iteration, should examine if
//  there are the same num after the current num, if it does
//  that it should be skipped as well.
//
//
func threeSum(nums []int) (result [][]int) {
	if len(nums) < 3 {
		return
	}

	sort.Ints(nums)
	result = make([][]int, 0)
	for i := 0; i < len(nums)-1; i++ {
		if nums[i] > 0 {
			return
		}
		// skip nums that is the same as the before.
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		left := i + 1
		right := len(nums) - 1
		for left < right {
			sum := nums[i] + nums[left] + nums[right]
			switch {
			case sum == 0:
				result = append(result, []int{nums[i], nums[left], nums[right]})
				for left < right && nums[left] == nums[left+1] {
					// ignore the duplicated neighboring nums in the left
					left++
				}
				for left < right && nums[right] == nums[right-1] {
					// ignore the duplicated neighboring nums in the right
					right--
				}
				left++
				right--
			case sum > 0: // means that the num in the right is too large
				right--
			case sum < 0: // means that the num in the left is too small
				left++
			}
		}
	}
	return
}

/*
Given a string s containing just the characters '(', ')', '{', '}', '[' and ']',
determine if the input string is valid.

An input string is valid only if:

- Open brackets must be closed by the same type of brackets.
- Open brackets must be closed in the correct order.

Example 1:
Input: s = "()"
Output: true

Example 2:
Input: s = "()[]{}"
Output: true

Example 3:
Input: s = "(]"
Output: false

*/
// Solution 1 using stack structure to solve efficiently.
func isValid(s string) bool {
	if len(s)%2 != 0 || len(s) == 0 {
		return false
	}

	var mapping = map[string]string{"(": ")", "{": "}", "[": "]"}
	stack := make([]string, 0, len(s)/2)
	for _, c := range s {
		if _, exist := mapping[string(c)]; exist { // is an open bucket
			stack = append(stack, string(c))
			continue
		}
		if len(stack) == 0 { // invalid string e.g. ")("
			return false
		}

		// In most cases, we use slice to replace with stack in golang,
		// cuz golang does not have a proper support for stack.
		//
		// 2 lines below works like: stack.Pop()
		left := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if mapping[left] != string(c) {
			return false
		}
	}
	return len(stack) == 0
}

// Solution 2 woks in an inefficient way.
func isValid2(s string) bool {
	for strings.Contains(s, "{}") || strings.Contains(s, "[]") || strings.Contains(s, "()") {
		s = strings.Replace(s, "{}", "", 1)
		s = strings.Replace(s, "()", "", 1)
		s = strings.Replace(s, "[]", "", 1)
	}
	return s == ""
}

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

// TreeNode represents a binary tree structure.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

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

/* Rotate Image

You are given an n x n 2D matrix representing an image, rotate the image by 90 degrees (clockwise).
You have to rotate the image in-place, which means you have to modify the input 2D matrix directly.
DO NOT allocate another 2D matrix and do the rotation.

Example 1:
Input: matrix = [[1,2,3],[4,5,6],[7,8,9]]
Output: [[7,4,1],[8,5,2],[9,6,3]]

Example 2:
Input: matrix = [[5,1,9,11],[2,4,8,10],[13,3,6,7],[15,14,12,16]]
Output: [[15,13,2,5],[14,3,4,1],[12,6,8,9],[16,7,10,11]]

Constraints:
n == matrix.length == matrix[i].length
1 <= n <= 20
-1000 <= matrix[i][j] <= 1000
*/

// rotate will rotate a matrix by 90 degrees by a extra matrix,
// it has not meet the requirement yet.
//
// TODO: rotates matrix in-place.
func rotate(matrix [][]int) {
	if len(matrix) == 0 {
		return
	}

	// init tmp
	mLen := len(matrix)
	var tmp = make([][]int, mLen)
	subLen := len(matrix[0])
	for i := 0; i < mLen; i++ {
		tmp[i] = make([]int, subLen)
	}

	for i := len(matrix) - 1; i >= 0; i-- {
		actualIndex := len(matrix) - 1 - i
		for j := 0; j < len(matrix[i]); j++ {
			tmp[j][actualIndex] = matrix[i][j]
		}
	}
	copy(matrix, tmp)
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
