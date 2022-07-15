package leetcode

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"unsafe"
)

func maxProfit(prices []int) int {
	var profit int

	for i := 0; i < len(prices)-1; i++ {
		if prices[i] < prices[i+1] {
			profit = profit + prices[i+1] - prices[i]
		}
	}

	return profit
}

/*
输入: nums = [1,2,3,4,5,6,7], k = 3
输出: [5,6,7,1,2,3,4]
解释:
向右轮转 1 步: [7,1,2,3,4,5,6]
向右轮转 2 步: [6,7,1,2,3,4,5]
向右轮转 3 步: [5,6,7,1,2,3,4]
*/

func rotate(nums []int, k int) {
	tmp := make([]int, len(nums))
	for i := 0; i < len(nums); i++ {
		index := (i + k) % len(nums)
		tmp[index] = nums[i]
	}
	copy(nums, tmp)
}

func reverse(nums []int) {
	tmp := 0
	l := len(nums)
	for i := 0; i < l/2; i++ {
		tmp = nums[i]
		nums[i] = nums[l-1-i]
		nums[l-i-1] = tmp
	}
}

func rotate2(nums []int, k int) {
	if k > len(nums) {
		k = k % len(nums)
	}
	// 左边数组长度len=k+1
	// 右边数组长度明显 len = len(nums) - k -1
	// 把切割得到的右边数组，copy到原先数组左边， copy(nums[:k+2],左边数组)
	// 把切割得到的左边数组，copy到原先数组右边， copy(nums[:k+2],右边数组)

	var (
		left  = make([]int, k+1)
		right = make([]int, len(nums)-k-1)
	)
	fmt.Println(k, len(left))
	copy(left, nums[:len(left)])
	copy(right, nums[len(left):])
	fmt.Println(left, right)

	copy(nums[:len(right)], right)
	copy(nums[len(right):], left)
}

//func containsDuplicate(nums []int) bool {
//	record := make(map[int]struct{},len(nums))
//
//	for _, num := range nums {
//		if _, exist := record[num]; exist {
//			return true
//		}
//		record[num] = struct{}{}
//	}
//	return false
//}

//func containsDuplicate(nums []int) bool {
//	if singleNumber(nums)
//}

func ByteArrayToInt(arr []byte) int64 {
	val := int64(0)
	size := len(arr)
	for i := 0; i < size; i++ {
		*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&val)) + uintptr(i))) = arr[i]
	}
	return val
}

func IntToByteArray(num int64) []byte {
	size := int(unsafe.Sizeof(num))
	arr := make([]byte, size)
	for i := 0; i < size; i++ {
		byt := *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&num)) + uintptr(i)))
		arr[i] = byt
	}
	return arr
}

//func singleNumber(nums []int) int {
//	record := make(map[int]int)
//	for i := 0; i < len(nums); i++ {
//		record[nums[i]] += 1
//	}
//	for value, count := range record {
//		if count == 1 {
//			return value
//		}
//	}
//	return 0
//}

func singleNumber(nums []int) int {
	for i := 1; i < len(nums); i++ {
		nums[0] = nums[0] ^ nums[i]
	}
	return nums[0]
}

func pow(m, n int) int {
	result := m
	for i := 1; i < n; i++ {
		result *= m
	}
	return result
}

//func intersect(nums1 []int, nums2 []int) []int {
//	sort.Ints(nums1)
//	sort.Ints(nums2)
//	var sets []int
//	var i, j int
//	for i < len(nums1) && j < len(nums2) {
//		switch {
//		case nums1[i] == nums2[j]:
//			sets = append(sets, nums1[i])
//			i++
//			j++
//			continue
//		case nums1[i] > nums2[j]:
//			j++
//		default:
//			i++
//		}
//	}
//	return sets
//}

func intersect(nums1 []int, nums2 []int) []int {
	tmp := make(map[int]int)
	for _, num := range nums1 {
		tmp[num] += 1
	}

	var res []int
	for _, num := range nums2 {
		if count, exist := tmp[num]; exist && count > 0 {
			res = append(res, num)
			tmp[num] -= 1
		}
	}
	return res
}

func plusOne(digits []int) []int {
	if digits[len(digits)-1] != 9 {
		digits[len(digits)-1] += 1
		return digits
	}
	preNumHasCarried := true
	for i := len(digits) - 1; i >= 0; i-- {
		if preNumHasCarried {
			digits[i] += 1
			if digits[i] == 10 {
				digits[i] = 0
			} else {
				preNumHasCarried = false
			}
		}
	}
	if digits[0] == 0 {
		t := make([]int, len(digits)+1)
		copy(t[1:], digits)
		t[0] = 1
		return t
	}
	return digits
}

/*
给定一个数组 nums，编写一个函数将所有 0 移动到数组的末尾，同时保持非零元素的相对顺序。
请注意 ，必须在不复制数组的情况下原地对数组进行操作。
nums = [0,1,0,3,12]
输出: [1,3,12,0,0]
*/
func moveZeroes(nums []int) {
	if len(nums) == 1 {
		return
	}

	var count int
	var jReadyToSetValue bool
	for i, j := 0, 0; i < len(nums) && j < len(nums); i++ {
		if nums[i] == 0 && !jReadyToSetValue {
			count++
			jReadyToSetValue = true
			j = i
		}

		if nums[i] != 0 && jReadyToSetValue {
			nums[j] = nums[i]
			nums[i] = 0
			if nums[j+1] == 0 {
				j++
			} else {
				j = i
			}
		}
	}
	if count > 0 {
		for i := len(nums) - 1; i > len(nums)-count; i-- {
			nums[i] = 0
		}
	}
}

//func twoSum(nums []int, target int) []int {
//	for i := 0; i < len(nums); i++ {
//		for j := i + 1; j < len(nums); j++ {
//			if nums[i]+nums[j] == target {
//				return []int{i, j}
//			}
//		}
//	}
//	return nil
//}
func twoSum(nums []int, target int) []int {
	var j = len(nums) - 1
	for i := 0; i < len(nums) && j >= 0; i++ {
		fmt.Println(i, j)
		if nums[i]+nums[j] == target {
			return []int{i, j}
		}
		j -= 1
	}
	return nil
}

type loc struct {
	x, y int
}

// 输入x,y坐标，返回宫数
func getCell(i, j int) int {
	if 0 <= i && i < 3 {
		if 0 <= j && j < 3 {
			return 1
		} else if 3 <= j && j < 6 {
			return 2
		} else {
			return 3
		}
	}
	if 3 <= i && i < 6 {
		if 0 <= j && j < 3 {
			return 4
		} else if 3 <= j && j < 6 {
			return 5
		} else {
			return 6
		}
	}
	if 6 <= i && i < 9 {
		if 0 <= j && j < 3 {
			return 7
		} else if 3 <= j && j < 6 {
			return 8
		} else {
			return 9
		}
	}
	return -1
}

func isValidSudoku(board [][]byte) bool {
	const dot = byte('.')
	//fmt.Println("dot=", dot)
	columnMap := make(map[int]map[byte]struct{}, 9)
	rowMap := make(map[int]map[byte]struct{}, 9)
	cellMap := make(map[int]map[byte]struct{}, 9)
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			row := board[i][j]

			if row != dot {
				//fmt.Printf("row>>>board[%d][%d]=%v\n", i, j, row)
				if rowMap[i] == nil {
					rowMap[i] = map[byte]struct{}{row: {}}
				} else {
					_, exist := rowMap[i][row]
					if exist {
						//fmt.Printf("%d already in rowMap:%+v\n", row, rowMap)
						return false
					} else {
						rowMap[i][row] = struct{}{}
					}
				}

				cell := getCell(i, j) - 1
				if cellMap[cell] == nil {
					cellMap[cell] = map[byte]struct{}{}
				}
				_, exist := cellMap[cell][row]
				if exist {
					return false
				} else {
					cellMap[cell][row] = struct{}{}
				}
			}

			column := board[j][i]
			if column != dot {
				//fmt.Printf("column>>board[%d][%d]=%v\n", j, i, column)
				if columnMap[i] == nil {
					columnMap[i] = map[byte]struct{}{column: {}}
				} else {
					_, exist := columnMap[i][column]
					if exist {
						//fmt.Printf("%d already in columnMap:%+v\n", column, columnMap)
						return false
					} else {
						columnMap[i][column] = struct{}{}
					}
				}
			}

		}
	}
	//fmt.Printf("%+v\n", columnMap)
	return true
}

func reverseString(s []byte) {
	//for i := 0; i < len(s)/2; i++ {
	//	s[i] = s[i] ^ s[len(s)-1-i]
	//	s[len(s)-1-i] = s[i] ^ s[len(s)-1-i]
	//	s[i] = s[i] ^ s[len(s)-1-i]
	//}
}

// heellhhoo
func firstUniqChar(s string) int {
	var left, right int

	for i := 0; i < len(s); i++ {
		left = i - 1
		right = i + 1
		hasLeft, hasRight := false, false
		for left >= 0 {
			if s[left] == s[i] {
				hasLeft = true
				break
			}
			left -= 1
		}
		for right < len(s) {
			if s[right] == s[i] {
				hasRight = true
				break
			}
			right += 1
		}

		//fmt.Printf("[%d] s[%d]=%s,hasLeft=%v,hasRight=%v\n", i, i, string(s[i]), hasLeft, hasRight)
		if !hasRight && !hasLeft {
			return i
		}
	}
	return -1
}

func isAnagram(s string, t string) bool {
	a := make(map[string]int)
	b := make(map[string]int)

	for _, word := range s {
		a[string(word)] += 1
	}
	for _, word := range t {
		b[string(word)] += 1
	}

	return reflect.DeepEqual(a, b)
}

var reg = regexp.MustCompile(`[A-Za-z\d]`)

func isPalindrome(s string) bool {
	var i, j = 0, len(s) - 1
	for i <= j {
		left, right := strings.ToLower(string(s[i])), strings.ToLower(string(s[j]))
		mL, mR := reg.MatchString(left), reg.MatchString(right)
		//fmt.Println(left, mL, right, mR)
		//i++
		//j--
		switch {
		case !mR:
			j--
		case !mL:
			i++
		case mL && mR && left != right:
			return false
		case mL && mR && left == right:
			i++
			j--
		}
		//if !mL {
		//	i++
		//}
		//if !mR {
		//	j--
		//}
		//if mL && mR && left != right {
		//	return false
		//}
		//if mL && mR && left == right {
		//	i++
		//	j--
		//}

	}
	return true
}

type ListNode struct {
	Val  int
	Next *ListNode
}

func removeNthFromEnd(head *ListNode, n int) *ListNode {
	count := 0
	var node = head
	for node != nil {
		count++
		node = node.Next
	}
	if count == 1 && n == 1 {
		return nil
	}

	node = head
	i := 1
	var target = count - n
	if target == 0 {
		node = head.Next
		head = node
		return node
	}
	for node != nil {
		if i == target {
			if node.Next.Next == nil {
				node.Next = nil
			} else {
				node.Next = node.Next.Next
			}
			break
		}

		node = node.Next
		i++
	}
	return head
}

func reverseList(head *ListNode) *ListNode {
	var tmp []*ListNode
	listNode := head
	for listNode != nil {
		tmp = append(tmp, listNode)
		listNode = listNode.Next
	}
	var newHead = tmp[len(tmp)-1]
	t := newHead
	for i := len(tmp) - 1; i >= 0; i-- {
		t.Val = tmp[i].Val
		if i == 0 {
			t.Next = nil
		} else {
			t.Next = tmp[i-1]
			t = tmp[i-1]
		}
	}
	return newHead
}

func NewListNode(l int) *ListNode {
	var head = &ListNode{Val: 1}
	tmp := head
	for i := 1; i < l; i++ {
		a := &ListNode{Val: i + 1, Next: nil}
		tmp.Next = a
		tmp = a
	}

	return head
}
func NewListNode2(begin, end int) *ListNode {
	var head = &ListNode{Val: begin}
	tmp := head
	for i := begin; i < end; i++ {
		a := &ListNode{Val: i + 1, Next: nil}
		tmp.Next = a
		tmp = a
	}

	return head
}

func printListNode(h *ListNode) {
	head := h
	for head != nil {
		fmt.Print(head.Val)
		head = head.Next
	}
	fmt.Println()
}

func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	if list1 == nil {
		return list2
	}
	if list2 == nil {
		return list1
	}

	var newHead = new(ListNode)
	if list1.Val > list2.Val {
		newHead = list2
	} else {
		newHead = list1
	}
	tmp := newHead
	var left, right = list1, list2
	for left != nil && right != nil {
		var node = new(ListNode)
		if left.Val > right.Val {
			node.Val = right.Val
			right = right.Next
		} else {
			node.Val = left.Val
			left = left.Next
		}
		tmp.Next = node
		tmp = node
	}
	for left != nil {
		node := new(ListNode)
		node.Val = left.Val
		tmp.Next = node
		tmp = node
		left = left.Next
	}
	for right != nil {
		node := new(ListNode)
		node.Val = right.Val
		tmp.Next = node
		tmp = node
		right = right.Next
	}
	newHead = newHead.Next
	return newHead
}

func ListNodeIsPalindrome(head *ListNode) bool {
	var tmp []*ListNode
	listNode := head
	for listNode != nil {
		tmp = append(tmp, listNode)
		listNode = listNode.Next
	}
	nHead := head
	for i := 0; i < len(tmp)/2; i++ {
		if nHead.Val != tmp[len(tmp)-1-i].Val {
			return false
		}
		nHead = nHead.Next
	}
	return true
}

func hasCycle(head *ListNode) bool {
	t := make(map[uintptr]struct{})
	node := head
	for node != nil {
		_, exist := t[uintptr(unsafe.Pointer(node.Next))]
		if exist && node.Next != nil {
			return true
		}
		t[uintptr(unsafe.Pointer(node.Next))] = struct{}{}
		node = node.Next
	}
	return false
}

func merge(nums1 []int, m int, nums2 []int, n int) {
	var i, j int
	tmp := make([]int, 0, m+n)
	for i < m && j < n {
		if nums1[i] < nums2[j] {
			tmp = append(tmp, nums1[i])
			i++
		} else {
			tmp = append(tmp, nums2[j])
			j++
		}
	}
	if i < m {
		tmp = append(tmp, nums1[i:m]...)
	}
	if j < n {
		tmp = append(tmp, nums2[j:n]...)
	}
	for i := 0; i < len(nums1); i++ {
		nums1[i] = tmp[i]
	}
	return
}

func isBadVersion(i int) bool {
	return i == 3
}

func strStr(haystack string, needle string) int {
	if needle == "" {
		return 0
	}
	if haystack == needle {
		return 0
	}
	i := 0
	for i < len(haystack) {
		if haystack[i] == needle[0] {
			// i = 9 ,and haystack[i] = p
			end := i + len(needle)
			if end >= len(haystack)-1 {
				end = len(haystack) - 1
			}

			if i == end && haystack[i:] == needle {
				return i
			}
			if end == len(haystack)-1 && haystack[i:] == needle {
				return i
			}
			if haystack[i:end] == needle {
				return i
			}
		}
		i++
	}
	return -1
}

func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return strs[0]
	}
	min := len(strs[0])
	minIndex := 0
	for i := 0; i < len(strs); i++ {
		if len(strs[i]) < min {
			min = len(strs[i])
			minIndex = i
		}
	}
	fmt.Println(min, minIndex)
	var prefix = make([]byte, 0, min)
	for i := 0; i < min; i++ {
		fmt.Println(i)
		var tmp byte
		same := true
		for j := 0; j < len(strs); j++ {
			if tmp == 0 {
				tmp = strs[j][i]
			}
			if strs[j][i] != tmp {
				same = false
				break
			}
			//fmt.Println(j, i, string(strs[j][i]), string(tmp))
		}
		if same {
			prefix = append(prefix, strs[minIndex][i])
		} else {
			break
		}
	}
	return string(prefix)
}

func myAtoi(s string) int {
	if s == "" {
		return -1
	}
	//negative := false
	if s[0] == '-' {
		s = s[1:]
		//negative = true
	} else if s[0] == '+' {
		s = s[1:]
	}

	nums := make([]int, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' && len(nums) == 0 {
			continue
		}
		if s[i] == ' ' && len(nums) > 0 {
			break
		}
	}

	return 0
}

type ProductAError struct {
	prd     string
	message string
	fn      string
}

func (p ProductAError) Error() string {
	if p.fn != "" {
		return fmt.Sprintf("product:%s,error:%q,fn:%s", p.prd, p.message, p.fn)
	}
	return fmt.Sprintf("product:%s,error:%q", p.prd, p.message)
}

func PrdAError(msg string, fn string) error {
	return &ProductAError{prd: "Apple", message: msg, fn: fn}
}

func mathPow(bit int) int {
	num := 1
	for i := 0; i < bit; i++ {
		num *= 10
	}
	return num
}

func PrdAMultiEmptyError(names ...string) error {
	return &ProductAError{prd: "Apple", message: fmt.Sprintf("one of them is empty: %v", names)}
}
