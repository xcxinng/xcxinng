package algorithm

import (
	"fmt"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

func GenerateLinkList(n int) *ListNode {
	var head = &ListNode{}
	var t = head
	for i := n; i > 0; i-- {
		node := &ListNode{}
		t.Next = node
		t = t.Next
	}
	return head.Next
}

func GenerateLinkListWithArray(nums []int) *ListNode {
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

func PrintlnLinkListAsArray(l *ListNode, reverse bool) {
	var array []int
	for l != nil {
		array = append(array, l.Val)
		l = l.Next
	}
	if !reverse {
		fmt.Println(array)
	} else {
		var t = make([]int, len(array))
		y := 0
		for i := len(array) - 1; i >= 0; i-- {
			t[y] = array[i]
			y++
		}
		fmt.Println(t)
	}
}

// ReverseListNode construct a new ListNode whose order is the reverse of the order of head.
func ReverseListNode(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}

	newHead := ReverseListNode(head.Next)
	head.Next.Next = head
	head.Next = nil // 这一步如果不置空，反转之后的最后一个与倒数第二个节点会有环
	return newHead
}

func AddTwoNumbers(l1, l2 *ListNode) *ListNode {
	var (
		dummy = &ListNode{}
		carry = false
	)
	cur := dummy
	for l1 != nil || l2 != nil {
		var sum int
		if l1 == nil {
			sum = l2.Val
		} else if l2 == nil {
			sum = l1.Val
		} else {
			sum = l1.Val + l2.Val
		}
		if carry {
			sum += 1
		}
		if sum >= 10 {
			sum = sum % 10
			carry = true
		} else {
			carry = false
		}
		n := &ListNode{Val: sum}
		cur.Next = n
		cur = cur.Next
		if l1 != nil {
			l1 = l1.Next
		}
		if l2 != nil {
			l2 = l2.Next
		}
	}
	if carry {
		n := &ListNode{Val: 1}
		cur.Next = n
	}
	return dummy.Next

}

// 合并生序链表
func mergeTwoLists2(list1 *ListNode, list2 *ListNode) *ListNode {
	var dummy = &ListNode{}
	cur := dummy

	for list1 != nil && list2 != nil {
		newNode := &ListNode{}
		if list1.Val >= list2.Val {
			newNode.Val = list2.Val
			list2 = list2.Next
		} else {
			newNode.Val = list1.Val
			list1 = list1.Next
		}
		cur.Next = newNode
		cur = cur.Next
	}
	for list1 != nil {
		newNode := &ListNode{Val: list1.Val}
		cur.Next = newNode
		list1 = list1.Next
		cur = cur.Next
	}
	for list2 != nil {
		newNode := &ListNode{Val: list2.Val}
		cur.Next = newNode
		list2 = list2.Next
		cur = cur.Next
	}
	return dummy.Next
}

// 给链表排序
func sortList(head *ListNode) *ListNode {
	if head == nil {
		return nil
	}
	if head.Next == nil {
		return head
	}
	dummy := &ListNode{}
	cur := dummy
	cur.Next = &ListNode{Val: head.Val}
	cur = cur.Next
	head = head.Next

	for head != nil {
		newNode := &ListNode{Val: head.Val}
		if head.Val >= cur.Val {
			newCur := dummy.Next
			for newCur != nil {
				if newCur.Val >= head.Val {
					newNode.Next = newCur.Next
					newCur.Next = newNode
					cur = newNode
					break
				}
			}
		} else {
			newNode.Next = cur
			dummy.Next = newNode
		}

		cur = newNode
		head = head.Next
	}
	return dummy.Next
}

// 对链表进行两两节点交换，不能改节点值，只能原地交换
// [1,2,3,4] ==> [2,1,4,3]
// [1,2,3] ==> [2,1,3]
func swapPairs(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	dummy := &ListNode{Next: head}
	cur := dummy
	for cur.Next != nil && cur.Next.Next != nil {
		node1 := cur.Next
		node2 := cur.Next.Next

		// 注意顺序
		cur.Next = node2
		node1.Next = node2.Next
		node2.Next = node1
		// 游标每次指向当前已完成节点交换的链表的末端
		cur = node1
	}
	return dummy.Next
}

// 对链表进行K分组翻转
// [1,2,3,4,5],k=2 ==> [2,1,4,3,5]
// [1,2,3,4,5],k=3 ==> [3,2,1,4,5]
//
// 时间复杂度O(n)
// 空间复杂度O(n)
func reverseKGroup(head *ListNode, k int) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	var numbs []int
	for head != nil {
		numbs = append(numbs, head.Val)
		if len(numbs) == k {
			for j := len(numbs) - 1; j >= 0; j-- {
				node := &ListNode{Val: numbs[j]}
				cur.Next = node
				cur = cur.Next
			}
			numbs = nil

		} else if head.Next == nil {
			for j := 0; j < len(numbs); j++ {
				node := &ListNode{Val: numbs[j]}
				cur.Next = node
				cur = cur.Next
			}
			numbs = nil
		}
		head = head.Next
	}
	return dummy.Next
}

func rotateRight(head *ListNode, k int) *ListNode {
	if k == 0 || head == nil || head.Next == nil {
		return head
	}
	count := 0
	cur := head
	for cur.Next != nil {
		count++
		cur = cur.Next
	}
	count += 1
	if k == count {
		return head
	}
	cur.Next = head
	moveCount := (count - 1) - k%count
	cur = head
	for i := moveCount; i > 0; i-- {
		cur = cur.Next
	}
	newHead := cur.Next
	cur.Next = nil
	return newHead
}

func deleteDuplicates(head *ListNode) *ListNode {
	if head == nil {
		return nil
	}

	dummy := &ListNode{0, head}

	cur := dummy
	for cur.Next != nil && cur.Next.Next != nil {
		if cur.Next.Val == cur.Next.Next.Val {
			x := cur.Next.Val
			for cur.Next != nil && cur.Next.Val == x {
				cur.Next = cur.Next.Next
			}
		} else {
			cur = cur.Next
		}
	}

	return dummy.Next
}
