package algorithm

import (
	"fmt"
	"math"
	"sort"
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

type ListNode struct {
	Val  int
	Next *ListNode
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
func letterCombinations(digits string) []string {
	if len(digits) == 0 {
		return []string{}
	}
	combinations = []string{}
	backtrack(digits, 0, "")
	return combinations
}

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

// In Chinese: 合并两个有序链表
// Difficulty: easy
//
// Description: 将两个升序链表合并为一个新的 升序 链表并返回。新链表是通过拼接给定的两个链表的所有节点组成的。
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
// Description: 给你一个链表数组，每个链表都已经按升序排列。请你将所有链表合并到一个升序链表中，返回合并后的链表。
//
// Constraints:
// k == lists.length
// 0 <= k <= 10^4
// 0 <= lists[i].length <= 500
// -10^4 <= lists[i][j] <= 10^4
// lists[i] 按 升序 排列
// lists[i].length 的总和不超过 10^4
//
// I can't believe I resolved it without referring to any official or comment explanations!
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
func combinationSum(candidates []int, target int) [][]int {
	return nil
}

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
	// I can't figure it out, so sad -_-|.
	permuteResult = make([][]int, 0) // ensure result is clean
	backTrackingPermute(len(nums), 0, nums)
	return permuteResult
}
