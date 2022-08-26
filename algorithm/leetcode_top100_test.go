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

func TestRotate(t *testing.T) {
	m := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	rotate(m)
	fmt.Println(m)
	// Output: [[7,4,1 ],[8,5,2],[9,6,3]]
}

func TestGroupAnagrams(t *testing.T) {
	// fmt.Println(groupAnagrams([]string{"eat", "tea", "tan", "ate", "nat", "bat"}))
	// fmt.Println(groupAnagrams([]string{""}))
	// fmt.Println(groupAnagrams([]string{"a"}))
	// fmt.Println(groupAnagrams([]string{"ill", "duh"}))
	fmt.Println(groupAnagrams([]string{
		"tho", "tin", "erg", "end", "pug", "ton", "alb", "mes", "job", "ads", "soy", "toe",
		"tap", "sen", "ape", "led", "rig", "rig", "con", "wac", "gog", "zen", "hay", "lie", "pay", "kid",
		"oaf", "arc", "hay", "vet", "sat", "gap", "hop", "ben", "gem", "dem", "pie", "eco", "cub", "coy",
		"pep", "wot", "wee"}))
	// output:
	// [["bat"],["nat","tan"],["ate","eat","tea"]]  (order can be different)
	// [[""]]
	// [["a"]]
	// [["duh"] ["ill"]]
	// [["wee"],["pep"],["cub"],["eco"],["dem"],["gap"],["vet"],["job"],["ben"],["toe"],["hay","hay"],["mes"],["ads"],
	// ["alb"],["wot"],["gem"],["oaf"],["hop"],["ton"],["pug"],["end"],["con"],["coy"],["sat"],["soy"],["pay"],["tin"],
	// ["pie"],["ape"],["tho"],["erg"],["sen"],["rig","rig"],["tap"],["wac"],["gog"],["led"],["zen"],["arc"],["lie"],["kid"]]
}

func TestVal(t *testing.T) {
	fmt.Println(maxSubArray([]int{-2, 1, -3, 4, -1, 2, 1, -5, 4}))
	// fmt.Println(maxSubArray([]int{1}))
	// fmt.Println(maxSubArray([]int{5, 4, -1, 7, 8}))
	// output:
	// 6
	// 1
	// 23
}

func TestCanJump(t *testing.T) {
	fmt.Println(canJump([]int{2, 3, 1, 1, 4}))
	fmt.Println(canJump([]int{3, 2, 1, 0, 4}))
	// output:
	// true
	// false
}

func TestMerge(t *testing.T) {
	fmt.Println(merge([][]int{{1, 4}, {2, 6}, {10, 13}, {15, 18}}))
	// output:
	// [[1 6] [10 13] [15 18]]
}

func TestInsert(t *testing.T) {
	fmt.Println(insert([][]int{{1, 2}, {3, 5}, {6, 7}, {8, 10}, {12, 16}}, []int{4, 8}))
	fmt.Println(insert([][]int{{1, 5}}, []int{6, 8}))
}

func TestMaxProfit(t *testing.T) {
	fmt.Println(maxProfit([]int{1, 2}))
	fmt.Println(maxProfit([]int{7, 1, 5, 3, 6, 4}))
	fmt.Println(maxProfit([]int{7, 6, 4, 3, 1}))
}

func TestClimbStairs(t *testing.T) {
	fmt.Println(climbStairs(2))
	fmt.Println(climbStairs(3))
	fmt.Println(climbStairs(4))
	// output:
	// 2
	// 3
	// 5
}

func TestSortColors(t *testing.T) {
	num := []int{1, 0, 1, 0, 2, 1, 2, 2, 0}
	sortColors(num)
	fmt.Println(num)
	// output:
	// [0 0 0 1 1 1 2 2 2]
}

func TestSubset(t *testing.T) {
	fmt.Println(subsets([]int{1}))
	fmt.Println(subsets([]int{1, 2}))
	fmt.Println(subsets([]int{1, 2, 3}))
	// output:
	// [[] [1]]
	// [[] [1 2] [1] [2]]
	// [[],[1],[2],[1,2],[3],[1,3],[2,3],[1,2,3]]
}
