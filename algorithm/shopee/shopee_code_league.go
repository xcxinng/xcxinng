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
// Source code is copied from [symmetric binary tree].
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
// [symmetric binary tree]: https://leetcode.cn/problems/symmetric-tree/solution/dui-cheng-er-cha-shu-by-leetcode-solution/
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