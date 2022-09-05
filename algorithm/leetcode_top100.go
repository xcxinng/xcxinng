package algorithm

import (
	"fmt"
	"math"
	"math/big"
	"sort"
	"strings"
)

// In Chinese: 两数之和
// Difficulty: simple(easy)
func twoSum(nums []int, target int) []int {
	for i := 0; i < len(nums)-1; i++ {
		for j := len(nums) - 1; j > i; j-- {
			if nums[i]+nums[j] == target {
				return []int{i, j}
			}
		}
	}
	return nil
}

// In Chinese: 两数相加
// Difficulty: medium
func addTwoNumbers(l1, l2 *ListNode) (head *ListNode) {
	var tail *ListNode
	carry := 0
	for l1 != nil || l2 != nil {
		n1, n2 := 0, 0
		if l1 != nil {
			n1 = l1.Val
			l1 = l1.Next
		}
		if l2 != nil {
			n2 = l2.Val
			l2 = l2.Next
		}
		sum := n1 + n2 + carry
		sum, carry = sum%10, sum/10
		if head == nil {
			head = &ListNode{Val: sum}
			tail = head
		} else {
			tail = &ListNode{Val: sum}
			tail = tail.Next
		}
	}
	if carry > 0 {
		tail.Next = &ListNode{Val: carry}
	}
	return
}

// In Chinese: 无重复字符的最长子串
// Difficulty: medium
//
// Description: 给定一个字符串 s ，请你找出其中不含有重复字符的 最长子串 的长度。
//
// Constraints：
// 0 <= s.length <= 5 * 104
// s 由英文字母、数字、符号和空格组成
//
func lengthOfLongestSubstring(s string) int {
	if s == "" {
		return 0
	}
	totalCount := 1
	for i := 0; i < len(s); i++ {
		j := i + 1
		t := make(map[string]struct{})
		tCount := 1
		t[string(s[i])] = struct{}{}
		for j < len(s) {
			if _, exist := t[string(s[j])]; exist {
				break
			}
			t[string(s[j])] = struct{}{}
			tCount += 1
			j++
		}
		if tCount > totalCount {
			totalCount = tCount
		}
	}
	return totalCount
}

// In Chinese: 寻找两个正序数组的中位数
// Difficulty: difficult
//
// Description: 给定两个大小分别为 m 和 n 的正序（从小到大）数组nums1 和nums2。
// 请你找出并返回这两个正序数组的 中位数 。 算法的时间复杂度应该为 O(log (m+n)) 。
//
// Constraints：
// nums1.length == m
// nums2.length == n
// 0 <= m <= 1000
// 0 <= n <= 1000
// 1 <= m + n <= 2000
// -106 <= nums1[i], nums2[i] <= 106
//
func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
	var mergeArray []int
	i, j := 0, 0
	for i < len(nums1) && j < len(nums2) {
		if nums1[i] > nums2[j] {
			mergeArray = append(mergeArray, nums2[j])
			j++
		} else if nums1[i] == nums2[j] {
			mergeArray = append(mergeArray, nums1[i], nums2[j])
			i++
			j++
		} else {
			mergeArray = append(mergeArray, nums1[i])
			i++
		}
	}
	for i < len(nums1) {
		mergeArray = append(mergeArray, nums1[i])
		i++
	}
	for j < len(nums2) {
		mergeArray = append(mergeArray, nums2[j])
		j++
	}
	if len(mergeArray)%2 == 0 {
		l := len(mergeArray) / 2
		sum := mergeArray[l] + mergeArray[l-1]
		fmt.Println(sum)
		return float64(sum) / 2
	} else {
		return float64(mergeArray[len(mergeArray)/2])
	}
}

// In Chinese: 三数之和
// Difficulty: medium
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
	sort.Ints(nums)

	for i := 0; i < len(nums); i++ {
		if nums[i] > 0 {
			return
		}
		if i > 0 && nums[i-1] == nums[i] {
			continue
		}

		j, k := i+1, len(nums)-1
		for j < k {
			sum := nums[i] + nums[j] + nums[k]
			switch {
			case sum == 0:
				result = append(result, []int{nums[i], nums[j], nums[k]})
				for j < k && nums[j] == nums[k] {
					j++
				}
				for j < k && nums[k] == nums[k-1] {
					k--
				}
				j++
				k--
			case sum > 0:
				k--
			case sum < 0:
				j++
			}
		}
	}
	return
}

var phoneMap map[string]string = map[string]string{
	"2": "abc",
	"3": "def",
	"4": "ghi",
	"5": "jkl",
	"6": "mno",
	"7": "pqrs",
	"8": "tuv",
	"9": "wxyz",
}

// In Chinese: 电话号码的字母组合
// Difficulty: medium
//
// Tags:
// #backtracking
// #all the possible solutions
// #Brute-force search
// func letterCombinations(digits string) []string {
// 	if len(digits) == 0 {
// 		return []string{}
// 	}
// 	combinations = []string{}
// 	backtrack(digits, 0, "")
// 	return combinations
// }

// combinations is used by letterCombinations and backtrack,
// it is the final answer of letterCombinations.
var combinations []string

// backtrack finds all the possible combinations using backtracking algorithm.
func backtrack(digits string, index int, combination string) {
	// I just can't figure out what "index == len(digits)" stands for.
	if index == len(digits) {
		combinations = append(combinations, combination)
	} else {
		// convert and get digit
		digit := string(digits[index])
		// get the corresponding letters
		letters := phoneMap[digit]
		// count the letters
		lettersCount := len(letters)
		// iterate the letters of the index
		for i := 0; i < lettersCount; i++ {
			// enter the recursion with the next digit and do the combination
			backtrack(digits, index+1, combination+string(letters[i]))
		}
	}
}

// In Chinese: 最长回文子串
// Difficulty: medium
//
// A palindrome is a word, number, phrase, or other sequence of characters
// which reads the same backward as forward, such as "madam" or "aka" etc.
//
// Tags: #dynamic programming
func longestPalindrome(s string) string {
	n := len(s)
	if n < 2 {
		return s
	}

	maxLen := 1
	begin := 0
	var dp [1005][1005]bool
	for i := 0; i < n; i++ {
		dp[i][i] = true
	}
	for L := 2; L <= n; L++ {
		for i := 0; i < n; i++ {
			j := L + i - 1
			if j >= n {
				break
			}

			if s[i] != s[j] {
				dp[i][j] = false
			} else {
				if j-i < 3 {
					dp[i][j] = true
				} else {
					dp[i][j] = dp[i+1][j-1]
				}
			}
			if dp[i][j] && j-i+1 > maxLen {
				maxLen = j - i + 1
				begin = i
			}

		}
	}
	return s[begin : begin+maxLen]
}

// In Chinese: 正则表达式匹配
// Difficulty: Difficult
//
// Tags: #dynamic programming
func isMatch(s string, p string) bool {
	m, n := len(s), len(p)
	matches := func(i, j int) bool {
		if i == 0 {
			return false
		}
		if p[j-1] == '.' {
			return true
		}
		return s[i-1] == p[j-1]
	}

	f := make([][]bool, m+1)
	for i := 0; i < len(f); i++ {
		f[i] = make([]bool, n+1)
	}
	f[0][0] = true
	for i := 0; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if p[j-1] == '*' {
				f[i][j] = f[i][j] || f[i][j-2]
				if matches(i, j-1) {
					f[i][j] = f[i][j] || f[i-1][j]
				}
			} else if matches(i, j) {
				f[i][j] = f[i][j] || f[i-1][j-1]
			}
		}
	}
	return f[m][n]
}

// In Chinese: 盛最多水的容器
// Difficulty: medium
//
// Tags: #double-pointers
func maxArea(height []int) int {
	i, j := 0, len(height)-1
	max := float64(0)
	for i < j {
		s := math.Abs(float64(i)-float64(j)) * math.Min(float64(height[i]), float64(height[j]))
		if s > max {
			max = s
		}
		if height[i] <= height[j] {
			i++
		} else {
			j--
		}
	}
	return int(max)
}

// In Chinese: 删除链表的倒数第 N 个结点
// Difficulty: medium
//
// Description: 给你一个链表，删除链表的倒数第 n 个结点，并且返回链表的头结点。
//
// Constraints:
// 链表中结点的数目为 sz
// 1 <= sz <= 30
// 0 <= Node.val <= 100
// 1 <= n <= sz
//
func removeNthFromEnd(head *ListNode, n int) *ListNode {
	count := 0
	temp := head
	for temp != nil {
		count++
		temp = temp.Next
	}
	switch {
	case count == n: // delete head node
		head = head.Next
	case n == 1: // delete the last node (n==1)
		m := head
		i := 1
		for i < count && m != nil {
			if i == count-1 {
				m.Next = nil
				break
			}
			i++
			m = m.Next
		}
	default: // delete node that is in the middle of the link list
		m := head
		t := count - n - 1
		for i := 0; i < count; i++ {
			if i == t {
				m.Next = m.Next.Next
				break
			}
			m = m.Next
		}
	}
	return head
}

// another solution for removeNthFromEnd.
//
// Tags: #double pointers #dummy ListNode
//
//
func removeNthFromEnd2(head *ListNode, n int) *ListNode {
	dummy := &ListNode{Next: head}
	fast, slow := head, dummy
	for i := 0; i < n; i++ {
		fast = fast.Next
	}
	for ; fast != nil; fast = fast.Next {
		slow = slow.Next
	}
	slow.Next = slow.Next.Next
	return dummy.Next
}

// In Chinese: 有效的括号
// Difficulty: easy
//
// Description: 给定一个只包括 '('，')'，'{'，'}'，'['，']'的字符串 s ，判断字符串是否有效。
//
// 有效字符串需满足：
// 左括号必须用相同类型的右括号闭合。
// 左括号必须以正确的顺序闭合。
//
// Constraints:
// 1 <= s.length <= 104
// s 仅由括号 '()[]{}' 组成
//
// Tags: #stack
//
// "匹配"关键字可以考虑栈结构，正确的括号最终会使得栈肯定是空的，否则说明不正确
func isValid(s string) bool {
	mapping := map[string]string{
		"(": ")",
		"[": "]",
		"{": "}",
	}
	stack := make([]string, 0, len(s))
	for i := 0; i < len(s); i++ {
		_, exist := mapping[string(s[i])]
		if exist {
			stack = append(stack, string(s[i]))
			continue
		}
		if len(stack) == 0 {
			return false
		}
		popChar := stack[len(stack)-1]
		if mapping[popChar] != string(s[i]) {
			return false
		}
		stack = stack[:len(stack)-1]
	}
	return len(stack) == 0
}

// Solution 1 using stack structure to solve efficiently.
// Solution 2 woks in an inefficient way.
func isValid2(s string) bool {
	for strings.Contains(s, "{}") || strings.Contains(s, "[]") || strings.Contains(s, "()") {
		s = strings.Replace(s, "{}", "", 1)
		s = strings.Replace(s, "()", "", 1)
		s = strings.Replace(s, "[]", "", 1)
	}
	return s == ""
}

// In Chinese: 合并两个有序链表
// Difficulty: easy
//
// Description: 将两个升序链表合并为一个新的 升序 链表并返回。
// 新链表是通过拼接给定的两个链表的所有节点组成的。
//
// Constraints:
// 两个链表的节点数目范围是 [0, 50]
// -100 <= Node.val <= 100
// l1 和 l2 均按 非递减顺序 排列
//
// Tags: #double pointers #dummy ListNode
//
// Time complexity: O(n)
// Space complexity: O(1)
//
func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	// put an additional dummy node to the head, in case nil pointer deference
	var dummy = &ListNode{}
	tail := dummy
	for list2 != nil && list1 != nil {
		if list1.Val <= list2.Val {
			tail.Next = list1
			list1 = list1.Next
		} else {
			tail.Next = list2
			list2 = list2.Next
		}
		tail = tail.Next
	}
	if list1 != nil {
		tail.Next = list1
	}
	if list2 != nil {
		tail.Next = list2
	}
	return dummy.Next
}

// In Chinese: 括号生成
// Difficulty: medium
//
// Description: 数字 n 代表生成括号的对数，请你设计一个函数，用于能够生成所有可能的并且 有效的 括号组合。
//
// Constraints
// 1 <= n <= 8
//
// Tags: #backtracking #dynamic programming
//
func generateParenthesis(n int) []string {
	// I can't do it!
	return nil
}

type Location struct {
	Value int
	Loc   int
}

// In Chinese: 合并K个升序链表
// Difficulty: medium
//
// Description: 给你一个链表数组，每个链表都已经按升序排列。
// 请你将所有链表合并到一个升序链表中，返回合并后的链表。
//
// Constraints:
// k == lists.length
// 0 <= k <= 10^4
// 0 <= lists[i].length <= 500
// -10^4 <= lists[i][j] <= 10^4
// lists[i] 按 升序 排列
// lists[i].length 的总和不超过 10^4
//
// I can't believe I resolved it without referring
// to any official or comment explanations!
//
// TODO: To try some other resolutions.
func mergeKLists(lists []*ListNode) *ListNode {
	var dummy = &ListNode{}
	var tail = dummy
	for {
		var minValue *Location
		for i := 0; i < len(lists); i++ {
			if lists[i] == nil {
				continue
			}
			if minValue == nil {
				minValue = &Location{Value: lists[i].Val, Loc: i}
			} else {
				if lists[i].Val < minValue.Value {
					minValue.Value = lists[i].Val
					minValue.Loc = i
				}
			}
		}
		if minValue == nil {
			break
		}

		node := &ListNode{Val: lists[minValue.Loc].Val}
		lists[minValue.Loc] = lists[minValue.Loc].Next
		tail.Next = node
		tail = tail.Next
	}
	return dummy.Next
}

// In Chinese: 下一个排列
// Difficulty: medium
//
// Description: 整数数组的排列就是将其所有成员以序列或线性顺序排列。
// 整数数组的下一个排列是指其整数的下一个字典序更大的排列。更正式地，
// 如果数组的所有排列根据其字典顺序从小到大排列在一个容器中，那么数组的
// 下一个排列就是在这个有序容器中排在它后面的那个排列。如果不存在下一个更大的排列，
// 那么这个数组必须重排为字典序最小的排列（即，其元素按升序排列）。
//
// 例如，arr = [1,2,3] 的下一个排列是 [1,3,2] 。
// 类似地，arr = [2,3,1] 的下一个排列是 [3,1,2] 。
// 而 arr = [3,2,1] 的下一个排列是 [1,2,3] ，因为 [3,2,1] 不存在一个字典序更大的排列。
//
// NOTE:
//  - 给你一个整数数组 nums ，找出 nums 的下一个排列。
//  - 必须 原地 修改，只允许使用额外常数空间。
//
// Constraints:
// 1 <= nums.length <= 100
// 0 <= nums[i] <= 100
//
// What a stupid question!
func nextPermutation(nums []int) {
	// I can't do it! shit!
}

// In Chinese: 最长有效括号
// Difficulty: difficult
//
// Description: Given a string containing just the characters '(' and ')',
// find the length of the longest valid (well-formed) parentheses substring.
//
// Example:
// Input: s = ")()())"
// Output: 4
// Explanation: The longest valid parentheses substring is "()()".
//
// Constraints:
// 0 <= s.length <= 3 * 104
// s[i] is '(', or ')'.
//
// What a stupid question! shit
func longestValidParentheses(s string) int {
	return 0
}

// In Chinese: 在排序数组中查找元素的第一个和最后一个位置
// Difficulty: medium
//
// Description: Given an array of integers nums sorted in non-decreasing order,
// find the starting and ending position of a given target value.
// If target is not found in the array, return [-1, -1].
// You must write an algorithm with O(log n) runtime complexity.
//
// Example:
// Input: nums = [5,7,7,8,8,10], target = 8
// Output: [3,4]
//
// Constraints:
// 0 <= nums.length <= 105
// -109 <= nums[i] <= 109
// nums is a non-decreasing array.
// -109 <= target <= 109
//
// What a stupid question! shit!
func searchRange(nums []int, target int) []int {
	return nil
}

// In Chinese: 组合总和
// Difficulty: medium
//
// Description: Given an array of distinct integers candidates and
// a target integer target, return a list of all unique combinations
// of candidates where the chosen numbers sum to target.
// You may return the combinations in any order.
//
// The same number may be chosen from candidates an unlimited number of times.
// Two combinations are unique if the frequency of at least one of the
// chosen numbers is different.
//
// It is guaranteed that the number of unique combinations that sum up
// to target is less than 150 combinations for the given input.
//
// Example:
// Input: candidates = [2,3,6,7], target = 7
// Output: [[2,2,3],[7]]
//
// Explanation:
// 2 and 3 are candidates, and 2 + 2 + 3 = 7. Note that 2 can be used multiple times.
// 7 is a candidate, and 7 = 7.
// These are the only two combinations.
//
// Constraints:
// 1 <= candidates.length <= 30
// 1 <= candidates[i] <= 200
// All elements of candidates are distinct.
// 1 <= target <= 500
//
// What a stupid question! shit!
// func combinationSum(candidates []int, target int) [][]int {
// 	return nil
// }

// [Cautions]
// TODO: Some of top100 algorithms are omit temporarily here. Please accomplish them in time.
// [Cautions]

// permuteResult stores the outcome of permute.
// It should be putting into a field of a structure instead of a global variable,
// why doing like this is because of leetcode' restriction.
var permuteResult [][]int

// backTrackingPermute generates all the possible permutations for nums.
//
// It's similar to generateParenthesis, letterCombinations and so on.
func backTrackingPermute(target int, first int, output []int) {
	if first == target {
		permuteResult = append(permuteResult, output)
	}
	for i := first; i < target; i++ {
		output[first], output[i] = output[i], output[first]
		backTrackingPermute(target, first+1, output)
		output[first], output[i] = output[i], output[first]
	}
}

// In Chinese: 全排列
// Difficulty: medium
//
// Description: Given an array nums of distinct integers, return all the
// possible permutations. You can return the answer in any order.
//
// Example:
// Input: nums = [1,2,3]
// Output: [[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]]
//
// Constraints:
// 1 <= nums.length <= 6
// -10 <= nums[i] <= 10
// All the integers of nums are unique.
//
// Tags: #backtracking
func permute(nums []int) [][]int {
	// Oh god, I totally get confused  -_-|.
	//
	// The reason I don't empty it is that it's so difficult
	// that I even can't understand it, so I keep the code here
	// for future review.
	permuteResult = make([][]int, 0)
	backTrackingPermute(len(nums), 0, nums)
	return permuteResult
}

// #48
// In Chinese:旋转图像
// Difficulty: medium
//
// Description: You are given an "n x n" 2D matrix representing
// an image, rotate the image by 90 degrees (clockwise).
//
// You have to rotate the image in-place, which means you have
// to modify the input 2D matrix directly.
//
// DO NOT allocate another 2D matrix and do the rotation.
//
// Example:
// Input: matrix = [[1,2,3],[4,5,6],[7,8,9]]
// Output: [[7,4,1],[8,5,2],[9,6,3]]
//
// Constraints:
// n == matrix.length == matrix[i].length
// 1 <= n <= 20
// -1000 <= matrix[i][j] <= 1000
//
// Tags: #understood #another resolutions
//
// Hint: Try to resolve it with an assistant slice first,
// and sum up the formula you had found out, then kill this
// stupid question.
func rotate(matrix [][]int) {
	// I have understood this algorithm, in case I forget,
	// just keep it empty for the next review.
	return
}

// #49
// In Chinese:字母异位词分组
// Difficulty: medium
//
// Description: Given an array of strings strs, group the anagrams together.
// You can return the answer in any order.
//
// An Anagram is a word or phrase formed by rearranging the letters of a
// different word or phrase, typically using all the original letters exactly once.
//
//
// Example:
// Input: strs = ["eat","tea","tan","ate","nat","bat"]
// Output: [["bat"],["nat","tan"],["ate","eat","tea"]]
//
// Constraints:
// 1 <= strs.length <= 104
// 0 <= strs[i].length <= 100
// strs[i] consists of lowercase English letters.
//
// Hint: byte slice, order
func groupAnagrams(strs []string) [][]string {
	// Oh god, I can't resolve this.
	return nil
}

// #53
// In Chinese:最大子数组和
// Difficulty: medium
//
// Description: Given an integer array nums, find the contiguous
// subarray (containing at least one number) which has the largest
// sum and return its sum.
//
// A subarray is a contiguous part of an array.
//
// Example:
// Input: nums = [-2,1,-3,4,-1,2,1,-5,4]
// Output: 6
// Explanation: [4,-1,2,1] has the largest sum = 6.
//
// Constraints:
// 1 <= nums.length <= 105
// -104 <= nums[i] <= 104
//
// Tags: #dynamic programming #dynamic array
func maxSubArray(nums []int) int {
	// Oh god, I can't resolve this.
	// I still get confused after reading the explanation.
	return 0
}

// #55
// In Chinese: 跳跃游戏
// Difficulty: medium
//
// Description: You are given an integer array nums.
// You are initially positioned at the array's first index,
// and each element in the array represents your maximum jump
// length at that position.
//
// Return true if you can reach the last index, or false otherwise.
//
// Example:
// Input: nums = [2,3,1,1,4]
// Output: true
// Explanation: Jump 1 step from index 0 to 1, then 3 steps to the last index.
//
// Constraints:
// 1 <= nums.length <= 104
// 0 <= nums[i] <= 105
//
// Solving Strategy:
//  - Dont complicate it!
//  - It's all about the index computation stuff
//
func canJump(nums []int) bool {
	var longMost = nums[0]
	for i := 1; i < len(nums); i++ {
		if i > longMost {
			return false
		}
		longMost = max(longMost, i+nums[i])
		if longMost >= len(nums) {
			return true
		}
	}
	return true
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// #56
// In Chinese: 合并区间
// Difficulty: medium
//
//
// Description: Given an array of intervals where intervals[i] = [starti, endi],
// merge all overlapping intervals, and return an array of the non-overlapping
// intervals that cover all the intervals in the input.
//
//
// Example:
// Input: intervals = [[1,3],[2,6],[8,10],[15,18]]
// Output: [[1,6],[8,10],[15,18]]
// Explanation: Since intervals [1,3] and [2,6] overlap, merge them into [1,6].
//
//
// Tags: mathematics
func merge(intervals [][]int) (result [][]int) {
	if len(intervals) == 0 {
		return
	}

	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})
	result = make([][]int, 0, len(intervals))
	result = append(result, intervals[0])
	iOR := 0

	for i := 1; i < len(intervals); i++ {
		if intersects(result[iOR], intervals[i]) {
			result[iOR][1] = max(result[iOR][1], intervals[i][1])
		} else {
			result = append(result, intervals[i])
			iOR++
		}
	}
	return
}

func intersects(a, b []int) bool {
	if a[0] <= b[0] && b[0] <= a[1] {
		return true
	}
	return false
}

// #57
// In Chinese: 插入区间
// Difficulty: medium
//
// Description: You are given an array of non-overlapping intervals
// where intervals[i] = [starti, endi] represent the start and the end of
// the ith interval and intervals is sorted in ascending order by starti.
// You are also given an interval newInterval = [start, end] that represents
// the start and end of another interval.
//
// Insert newInterval into intervals such that intervals is still sorted in
// ascending order by starti and intervals still does not have any overlapping
// intervals (merge overlapping intervals if necessary).
//
// Return intervals after the insertion.
//
// Constraints:
// intervals is sorted by starti in ascending order.
//
// Example:
// Input: intervals = [[1,3],[6,9]], newInterval = [2,5]
// Output: [[1,5],[6,9]]
//
// Tags: mathematics
func insert(intervals [][]int, newInterval []int) [][]int {
	if len(intervals) == 0 {
		return [][]int{newInterval}
	}

	intervals = append(intervals, newInterval)
	return intervals
}

// #62
// In Chinese: 不同路径
// Difficulty: medium
//
// There is a robot on an [m x n] grid. The robot is initially located at
// the top-left corner (i.e., grid[0][0]). The robot tries to move to the
// bottom-right corner (i.e., grid[m - 1][n - 1]). The robot can only move
// either down or right at any point in time.
//
// Given the two integers m and n, return the number of possible unique paths
// that the robot can take to reach the bottom-right corner.
//
// The test cases are generated so that the answer will be less than or equal to 2 * 109.
//
// Examples:
// Input: m = 3, n = 7
// Output: 28
//
/*
从左上角到右下角的过程中，我们需要移动m+n−2次，其中有 m−1 次向下移动,n-1次向右移动。
总路径数等从m+n-2次移动中选择m-1次向下移动的方案数:
Cm-1          m+n-1         (m+n-2)(m+n-3)...n           (m+n-2)!
        =  (--------)  =   --------------------- =     -------------
Cm+n-2        m-1               (m-1)!                  (m-1)!(n-1)!
*/
func uniquePaths(m, n int) int {
	// what the heck is this question!?
	// R U serious??
	return int(new(big.Int).Binomial(int64(m+n-2), int64(n-1)).Int64())
}

// #64
// In Chinese: 最小路径和
// Difficulty: medium
//
// Description:
// Given a "m x n" grid filled with non-negative numbers,
// find a path from top left to bottom right, which minimizes
// the sum of all numbers along its path.
// Note: You can only move either down or right at any point in time.
//
// Examples:
// Input: grid = [[1,3,1],[1,5,1],[4,2,1]]
// Output: 7
// Explanation: Because the path 1 → 3 → 1 → 1 → 1 minimizes the sum.
//
// Tags: #dynamic programming
func minPathSum(grid [][]int) int {
	return 0
}

// #70
// In Chinese: 爬楼梯
// Difficulty: easy
//
// Descriptions:
// You are climbing a staircase. It takes n steps to reach the top.
// Each time you can either climb 1 or 2 steps. In how many distinct
// ways can you climb to the top?
//
// Constraints:
// 1 <= n <= 45
//
// Examples:
// Input: n = 3
// Output: 3
// Explanation: There are three ways to climb to the top.
// 1. 1 step + 1 step + 1 step
// 2. 1 step + 2 steps
// 3. 2 steps + 1 step
//
// Tags: #dynamic programming
// Formula: dp[i] = dp[i-2] + dp[i-1]
func climbStairs(n int) int {
	if n <= 2 {
		return n
	}
	q := 1
	p := 2
	r := 0
	for i := 3; i <= n; i++ {
		r = q + p
		q = p
		p = r
	}
	return r
}

// #75
// In Chinese: 颜色分类
// Difficulty: medium
//
// Descriptions:
// Given an array nums with n objects colored red, white, or blue,
// sort them in-place so that objects of the same color are adjacent,
// with the colors in the order red, white, and blue.
//
// We will use the integers 0, 1, and 2 to represent the color red,
// white, and blue, respectively.
//
// You must solve this problem without using the library's sort function.
//
// Constraints:
// n == nums.length
// 1 <= n <= 300
// nums[i] is either 0, 1, or 2.
//
// Examples:
// Input: nums = [2,0,2,1,1,0]
// Output: [0,0,1,1,2,2]
//
func sortColors(nums []int) {
	countZero, countOne, countTwo := 0, 0, 0
	for i := 0; i < len(nums); i++ {
		switch nums[i] {
		case 0:
			countZero++
		case 1:
			countOne++
		default:
			countTwo++
		}
	}
	for i := 0; i < countZero; i++ {
		nums[i] = 0
	}
	for i := countZero; i < countZero+countOne; i++ {
		nums[i] = 1
	}
	for i := countOne + countZero; i < len(nums); i++ {
		nums[i] = 2
	}
}

func subsets(nums []int) [][]int {
	result := make([][]int, 0)
	result = append(result, []int{}, nums)
	if len(nums) == 0 {
		return result
	}

	i := 1
	for i < len(nums) {
		left := 0
		right := left + i
		for left < len(nums) && right <= len(nums) {
			values := nums[left:right]
			result = append(result, values)
			left++
			right = left + i
		}
		i++
	}
	return result
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// #94
// In Chinese: 二叉树最大深度
// Difficulty: easy
// Given the root of a binary tree, return the inorder traversal of its nodes' values.
//
//
// Examples:
// Input: root = [1,null,2,3]
// Output: [1,3,2]
//
//
// Constraints:
// The number of nodes in the tree is in the range [0, 100].
// -100 <= Node.val <= 100
//
func inorderTraversal(root *TreeNode) []int {
	inorderResult = make([]int, 0)
	traverse(root)
	return inorderResult
}

var inorderResult []int

func traverse(root *TreeNode) {
	if root == nil {
		return
	}
	traverse(root.Left)
	inorderResult = append(inorderResult, root.Val)
	traverse(root.Right)
}

// #104
// In Chinese: 二叉树最大深度
// Difficulty: easy
//
// Given the root of a binary tree, return its maximum depth.
// A binary tree's maximum depth is the number of nodes along
// the longest path from the root node down to the farthest leaf node.
//
// Example:
// Input: root = [3,9,20,null,null,15,7]
// Output: 3
//
// Constraints:
// The number of nodes in the tree is in the range [0, 104].
// -100 <= Node.val <= 100
func maxDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}
	return max(maxDepth(root.Left), maxDepth(root.Right)) + 1
}

// #121
// In Chinese: 买卖股票的最佳时机
// Difficulty: easy
//
//
// Description:
// You are given an array prices where prices[i] is the price of a given
// stock on the ith day.
//
// You want to maximize your profit by choosing a single day to buy one
// stock and choosing a different day in the future to sell that stock.
//
// Return the maximum profit you can achieve from this transaction.
// If you cannot achieve any profit, return 0.
//
//
// Examples:
// Input: prices = [7,1,5,3,6,4]
// Output: 5
// Explanation: Buy on day 2 (price = 1) and sell on day 5 (price = 6),
// profit = 6-1 = 5. Note that buying on day 2 and selling on day 1
// is not allowed because you must buy before you sell.
//
// Constraints:
// 1 <= prices.length <= 105
// 0 <= prices[i] <= 104
func maxProfit(prices []int) int {
	if len(prices) == 0 {
		return 0
	}

	maxProfit := 0
	min := prices[0]
	for i := 1; i < len(prices); i++ {
		min = minInt(min, prices[i])
		maxProfit = max(maxProfit, prices[i]-min)
	}
	return maxProfit
}
