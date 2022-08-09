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
