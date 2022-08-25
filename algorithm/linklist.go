package algorithm

import "fmt"

type ListNode struct {
	Val  int
	Next *ListNode
}

func generateLinkList(n int) *ListNode {
	var head = &ListNode{}
	var t = head
	for i := n; i > 0; i-- {
		node := &ListNode{}
		t.Next = node
		t = t.Next
	}
	return head.Next
}

func generateLinkListWithArray(nums []int) *ListNode {
	var head = &ListNode{}
	var t = head
	index := 0
	for i := len(nums); i > 0 && index < len(nums); i-- {
		node := &ListNode{Val: nums[index]}
		t.Next = node
		t = t.Next
		index++
	}
	return head.Next
}

func PrintlnLinkList(l *ListNode) {
	for l != nil {
		fmt.Println(l.Val)
		l = l.Next
	}
}
