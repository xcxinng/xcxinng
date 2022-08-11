package algorithm

import (
	"fmt"
	"testing"
)

func TestThreeSum(t *testing.T) {
	fmt.Println(threeSum([]int{-1, -1, 0, 0, 1, 1, 2}))
}

func TestLetterCombinations(t *testing.T) {
	fmt.Println(letterCombinations("23"))
}

func TestLongestPalindrome(t *testing.T) {
	fmt.Println(longestPalindrome("kkabadd"))
}

func TestIsMatch(t *testing.T) {
	fmt.Println(isMatch("aa", "a*"))
}

func TestMaxArea(t *testing.T) {
	fmt.Println(maxArea([]int{1, 8, 6, 2, 5, 4, 8, 3, 7}))
}

func TestRemoveNthFromEnd(t *testing.T) {
	head := &ListNode{Val: 1, Next: &ListNode{Val: 2, Next: &ListNode{Val: 3, Next: &ListNode{Val: 4}}}}
	// head := &ListNode{Val: 1, Next: &ListNode{Val: 2}}
	head = removeNthFromEnd(head, 2)
	for head != nil {
		fmt.Println(head.Val)
		head = head.Next
	}
}

func TestIsValid(t *testing.T) {
	fmt.Println(isValid("(("))
	fmt.Println(isValid("))"))
	fmt.Println(isValid("(([))"))
	fmt.Println(isValid("(([]))"))
	// output:
	// false
	// false
	// false
	// true
}

func TestMergeTwoLists(t *testing.T) {
	list := mergeTwoLists(&ListNode{Val: 1, Next: &ListNode{Val: 2}}, &ListNode{Val: 3, Next: &ListNode{Val: 4}})
	for list != nil {
		fmt.Println(list.Val)
		list = list.Next
	}
	fmt.Println()
	list = mergeTwoLists(&ListNode{Val: 1, Next: &ListNode{Val: 2, Next: &ListNode{Val: 4}}},
		&ListNode{Val: 1, Next: &ListNode{Val: 3, Next: &ListNode{Val: 4}}})
	for list != nil {
		fmt.Println(list.Val)
		list = list.Next
	}
	list = mergeTwoLists(nil, &ListNode{})
	for list != nil {
		fmt.Println(list.Val)
		list = list.Next
	}
	// output:
	// 1
	// 2
	// 3
	// 4
	//
	// 1
	// 1
	// 2
	// 3
	// 4
	// 4
}

func TestMergeKLists(t *testing.T) {
	// [[],[-10,-9,-8,-7,-2,-1,0,1],[-4],[]]
	res := mergeKLists([]*ListNode{
		// case 1
		// {Val: 1, Next: &ListNode{Val: 4, Next: &ListNode{Val: 5}}},
		// {Val: 1, Next: &ListNode{Val: 3, Next: &ListNode{Val: 4}}},
		// {Val: 2, Next: &ListNode{Val: 6}},
		// output:
		// 1
		// 1
		// 2
		// 3
		// 4
		// 4
		// 5
		// 6

		// case84
		nil,
		{Val: -10, Next: &ListNode{Val: -9, Next: &ListNode{Val: -8, Next: &ListNode{Val: -7, Next: &ListNode{Val: 0, Next: &ListNode{Val: 1}}}}}},
		{Val: -4},
		nil,
		// output:
		// -10
		// -9
		// -8
		// -7
		// -4
		// 0
		// 1
	})
	for res != nil {
		fmt.Println(res.Val)
		res = res.Next
	}
}

func TestLongestValidParentheses(t *testing.T) {
	fmt.Println(longestValidParentheses("(()"))
	fmt.Println(longestValidParentheses(")()())"))
	fmt.Println(longestValidParentheses("()(()"))
	// output:
	// 2
	// 4
	// 2
}

func TestSearchRange(t *testing.T) {
	fmt.Println(searchRange([]int{1, 4}, 4))
	fmt.Println(searchRange([]int{1, 3}, 1))
	fmt.Println(searchRange([]int{2, 2}, 2))
	fmt.Println(searchRange([]int{1, 2, 3, 4, 3}, 3))
	fmt.Println(searchRange([]int{5, 7, 7, 8, 8, 10}, 8))
	fmt.Println(searchRange([]int{5, 7, 7, 8, 8, 10}, 6))
	// output:
	// [1 1]
	// [0 0]
	// [0 1]
	// [2 4]
	// [3 4]
	// [-1 -1]
}

func TestPermute(t *testing.T) {
	fmt.Println(permute([]int{1, 2, 3}))
	fmt.Println(permute([]int{0, 1}))
	fmt.Println(permute([]int{1}))
	// output:
	// [[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]]
	// [[0,1],[1,0]]
	// [[1]]
}
