package algorithm

import (
	"sort"
	"strconv"
	"strings"
)

// #77 Combinations
//
// Descriptions:
// Given two integers n and k, return all possible
// combinations of k numbers chosen from the range [1, n].
// You may return the answer in any order.
//
// Examples:
// Input: n = 4, k = 2
// Output: [[1,2],[1,3],[1,4],[2,3],[2,4],[3,4]]
// Explanation: There are 4 choose 2 = 6 total combinations.
// Note that combinations are unordered, i.e., [1,2] and [2,1] are considered to be the same combination.
//
// Constraints:
// 1 <= n <= 20
// 1 <= k <= n
//
// [explanation]: https://programmercarl.com/0077.组合.html#回溯法三部曲
func combine(n int, k int) [][]int {
	result = make([][]int, 0)
	path = make([]int, 0)
	backtracking(n, k, 1)
	return result
}

var (
	result [][]int
	path   []int
)

// TODO: enhance this algorithm
// refer [combination enhancement] for more detail.
//
// [combination enhancement]: https://programmercarl.com/0077.组合优化.html
func backtracking(n, k, startIndex int) {
	if len(path) == k {
		// why not:
		// result = append(result, path)
		// instead of:
		t := make([]int, k)
		copy(t, path)
		// Because the growth of the path slice will never happen again after
		// its first length reach to 2, and every time result.append(path),
		// it just add the same slice head into result, thus, all the references
		// are led to the same underlying array.
		//
		// That's why the elements in the result are always the same finally.
		//
		// And be cautious with it! I just can't believe I fell into the trap again.
		result = append(result, t)
		return
	}

	// The path slice here is like a stack structure,
	// only modify the element on the tail(top) of the slice.
	for i := startIndex; i <= n; i++ {
		path = append(path, i) // push
		backtracking(n, k, i+1)
		path = path[:len(path)-1] // pop
	}
}

// #216
// In Chinese: 组合总和 III
//
//
// Description:
// Find all valid combinations of k numbers that sum up to n such that
// the following conditions are true:
//
//      Only numbers 1 through 9 are used.
//      Each number is used at most once.
//
// Return a list of all possible valid combinations. The list must not
// contain the same combination twice, and the combinations may be
// returned in any order.
//
//
// Example1:
// Input: k = 3, n = 7
// Output: [[1,2,4]]
// Explanation:
// 1 + 2 + 4 = 7
// There are no other valid combinations.
//
// Example2:
// Input: k = 3, n = 9
// Output: [[1,2,6],[1,3,5],[2,3,4]]
// Explanation:
// 1 + 2 + 6 = 9
// 1 + 3 + 5 = 9
// 2 + 3 + 4 = 9
// There are no other valid combinations.
//
//
// Constraints:
// 2 <= k <= 9
// 1 <= n <= 60
func combinationSum3(k, n int) [][]int {
	resultSum3 = make([][]int, 0)
	pathSum3 = make([]int, 0)
	sum3Backtracking(k, n, 1, 0)
	return resultSum3
}

var (
	pathSum3   []int
	resultSum3 [][]int
)

func sum3Backtracking(k, n, startIndex, currentSum int) {
	if len(pathSum3) == k {
		if currentSum == n {
			t := make([]int, k)
			copy(t, pathSum3)
			resultSum3 = append(resultSum3, t)
		}
		return
	}

	for i := startIndex; i <= 9; i++ {
		pathSum3 = append(pathSum3, i)
		sum3Backtracking(k, n, i+1, currentSum+i)
		pathSum3 = pathSum3[:len(pathSum3)-1]
	}
}

func letterCombinations(digits string) []string {
	letterResult = make([]string, 0)
	if len(digits) == 0 {
		return letterResult
	}
	letterBacktracking(digits, 0, "")
	return letterResult
}

var letterResult []string

func letterBacktracking(digits string, index int, combination string) {
	if index == len(digits) {
		letterResult = append(letterResult, combination)
		return
	}
	letters := phoneMap[string(digits[index])]
	for i := 0; i < len(letters); i++ {
		// Q: Why doesn't it do the backtracking of combination ?
		// Somebody can answer me?
		letterBacktracking(digits, index+1, combination+string(letters[i]))
	}
}

var (
	comResult [][]int
	comPath   []int
)

// 组合总和
//
// 给你一个无重复元素的整数数组 candidates 和一个目标整数 target，
// 找出 candidates 中可以使数字和为目标数 target 的 所有 不同组合，
// 并以列表形式返回。你可以按 任意顺序 返回这些组合。
//
// candidates 中的同一个数字可以无限制重复被选取 。
// 如果至少一个数字的被选数量不同，则两种组合是不同的。
//
// 对于给定的输入，保证和为 target 的不同组合数少于 150 个。
//
// 分析:
// 同个数字可以被选取，递归 i
// 同个集合内的组合，需要startIndex
func combinationSum(candidates []int, target int) [][]int {
	comResult = make([][]int, 0)
	comBacktracking(candidates, target, 0, 0)
	return comResult
}

func comBacktracking(candidates []int, target, sum, startIndex int) {
	if sum >= target {
		if sum == target {
			t := make([]int, len(comPath))
			copy(t, comPath)
			comResult = append(comResult, t)
		}
		return
	}

	for i := startIndex; i < len(candidates); i++ {
		sum += candidates[i]
		comPath = append(comPath, candidates[i])
		comBacktracking(candidates, target, sum, i)
		sum -= candidates[i]
		comPath = comPath[:len(comPath)-1]
	}
}

var (
	com2Result [][]int
	com2Path   []int
)

// 组合总和2
//
// 给定一个候选人编号的集合 candidates 和一个目标数 target，
// 找出 candidates 中所有可以使数字和为 target 的组合。
// candidates 中的每个数字在每个组合中只能使用 一次。
//
// 注意：解集不能包含重复的组合。
//
// 分析:
// 数字只能被使用一次，递归 i+1
// 同个集合内的组合， 需要startIndex
// 有重复数字但不能含重复组合，需要排序然后通过 n[i] == n[i-1] 去重
func combinationSum2(candidates []int, target int) [][]int {
	com2Result = make([][]int, 0)
	sort.Ints(candidates)
	com2Backtracking(candidates, target, 0, 0)
	return com2Result
}
func com2Backtracking(candidates []int, target, sum, startIndex int) {
	if sum >= target {
		if sum == target {
			t := make([]int, len(com2Path))
			copy(t, com2Path)
			com2Result = append(com2Result, t)
		}
		return
	}

	for i := startIndex; i < len(candidates); i++ {
		if i > startIndex && candidates[i] == candidates[i-1] {
			continue
		}
		sum += candidates[i]
		com2Path = append(com2Path, candidates[i])
		com2Backtracking(candidates, target, sum, i+1)
		sum -= candidates[i]
		com2Path = com2Path[:len(com2Path)-1]
	}
}

// 分割回文串
//
// 给你一个字符串s，请你将s分割成一些子串，使每个子串都是回文串。返回s所有可能的分割方案。
// 回文串是正着读和反着读都一样的字符串。
//
// Example:
// 输入：s = "aab"
// 输出：[["a","a","b"],["aa","b"]]
//
// 分析
// 分割 ===  s[startIndex:i]
// 由N叉树容易得，递归 i+1
func partition(s string) [][]string {
	palindromeRes = make([][]string, 0)
	palindromeBacktracking(s, 0)
	return palindromeRes
}

var (
	palindromeRes  [][]string
	palindromePath []string
)

func palindromeBacktracking(s string, startIndex int) {
	if startIndex >= len(s) {
		t := make([]string, len(palindromePath))
		copy(t, palindromePath)
		palindromeRes = append(palindromeRes, t)
		return
	}

	for i := startIndex; i < len(s); i++ {
		if !isPalindrome(s, startIndex, i) {
			continue
		}
		// Note: everytime trim a substring, it will be s[startIndex:i+1]
		// here are some examples:
		// 0:1, a
		// 0:2, aa
		// 0:3, aaa
		//
		// 1:2, b
		// 1:3, bb
		// 1:4, bbb
		palindromePath = append(palindromePath, s[startIndex:i+1])
		palindromeBacktracking(s, i+1)
		palindromePath = palindromePath[:len(palindromePath)-1]
	}
}

func isPalindrome(s string, start, end int) bool {
	for start < end {
		if s[start] != s[end] {
			return false
		}
		start++
		end--
	}
	return true
}

// 复原IP地址
//
// 有效IP地址正好由四个整数（每个整数位于 0 到 255 之间组成，且不能含有前导 0），整数之间用 '.' 分隔。
// 例如："0.1.2.201" 和 "192.168.1.1" 是 有效 IP 地址，但是 "0.011.255.245"、"192.168.1.312"
// 和 "192.168@1.1" 是无效IP地址。
//
// 给定一个只包含数字的字符串 s ，用以表示一个 IP 地址，返回所有可能的有效 IP 地址，这些地址可以通过
// 在 s 中插入'.'来形成。你不能重新排序或删除s中的任何数字。你可以按任何顺序返回答案。
//
// Constraints：
// 1 <= s.length <= 20
// s 仅由数字组成
//
// 整体思路也是走的切割路线，与切割回文串不同的是，这里多了对IP的合法性校验
func restoreIpAddresses(s string) []string {
	ipResult = make([]string, 0)
	ipBacktracking(s, 0)
	return ipResult
}

var (
	ipResult []string
	ipPath   []string
)

func ipBacktracking(s string, startIndex int) {
	if startIndex == len(s) {
		if len(ipPath) == 4 {
			ipResult = append(ipResult, strings.Join(ipPath, "."))
		}
		return
	}
	for i := startIndex; i < len(s); i++ {
		t := s[startIndex : i+1]
		if len(t) > 1 && strings.HasPrefix(t, "0") {
			continue
		}
		if len(t) > 3 {
			continue
		}
		num, _ := strconv.ParseInt(t, 10, 32)
		if num > 255 {
			continue
		}
		ipPath = append(ipPath, t)
		ipBacktracking(s, i+1)
		ipPath = ipPath[:len(ipPath)-1]
	}
}

/*

==============贯穿整个回溯篇幅===============

"在树形结构中子集问题是要收集所有节点的结果，而组合问题是收集叶子节点的结果"

==============确实很重要的一点===============

*/

// 子集
//
// 给你一个整数数组 nums ，数组中的元素 互不相同 。返回该数组所有可能的子集（幂集）。
// 解集 不能 包含重复的子集。你可以按 任意顺序 返回解集。
//
// Constraint：
// 1 <= nums.length <= 10
// -10 <= nums[i] <= 10
// nums 中的所有元素 互不相同
//
// 记住： 组合问题和分割问题都是收集树的叶子节点，而子集问题是找树的所有节点！
func subsets(nums []int) [][]int {
	subsetResult = make([][]int, 0)
	subsetPath = make([]int, 0)
	subsetBacktracking(nums, 0)
	return subsetResult
}

var (
	subsetResult [][]int
	subsetPath   []int
)

// 子集：求所有结果
func subsetBacktracking(num []int, startIndex int) {
	t := make([]int, len(subsetPath))
	copy(t, subsetPath)
	subsetResult = append(subsetResult, t)
	for i := startIndex; i < len(num); i++ {
		subsetPath = append(subsetPath, num[i])
		subsetBacktracking(num, i+1)
		subsetPath = subsetPath[:len(subsetPath)-1]
	}
}

// 子集II
//
// 给你一个整数数组 nums ，其中可能包含重复元素，请你返回该数组所有可能的子集（幂集）。
// 解集 不能 包含重复的子集。返回的解集中，子集可以按 任意顺序 排列。
//
// 提示：
// 1 <= nums.length <= 10
// -10 <= nums[i] <= 10
//
// 就多了个去重逻辑： 排序 + n[i] == n[i-1]
func subsetsWithDup(nums []int) [][]int {
	subsetDupPath = make([]int, 0)
	subsetDupResult = make([][]int, 0)
	sort.Ints(nums)
	subsetDupBacktracking(nums, 0)
	return subsetDupResult
}

var (
	subsetDupResult [][]int
	subsetDupPath   []int
)

func subsetDupBacktracking(nums []int, startIndex int) {
	t := make([]int, len(subsetDupPath))
	copy(t, subsetDupPath)
	subsetDupResult = append(subsetDupResult, t)
	for i := startIndex; i < len(nums); i++ {
		if i > startIndex && nums[i] == nums[i-1] {
			continue
		}
		subsetDupPath = append(subsetDupPath, nums[i])
		subsetDupBacktracking(nums, i+1)
		subsetDupPath = subsetDupPath[:len(subsetDupPath)-1]
	}
}

// 递增子序列
//
// 给你一个整数数组 nums ，找出并返回所有该数组中不同的递增子序列，
// 递增子序列中 至少有两个元素 。你可以按 任意顺序 返回答案。
// 数组中可能含有重复元素，如出现两个整数相等，也可以视作递增序列的一种特殊情况。
//
// 提示：
// 1 <= nums.length <= 15
// -100 <= nums[i] <= 100
//
func findSubsequences(nums []int) [][]int {
	subsequenceResult = make([][]int, 0)
	subsequencePath = make([]int, 0)
	subsequenceBacktracking(nums, 0, 0)
	return subsequenceResult
}

var (
	subsequenceResult [][]int
	subsequencePath   []int
)

func subsequenceBacktracking(nums []int, startIndex int, sum int) {
	if len(subsequencePath) > 1 {
		t := make([]int, len(subsequencePath))
		copy(t, subsequencePath)
		subsequenceResult = append(subsequenceResult, t)
	}

	// the key to deduplicate: A digit in the same layer can not be used repeatedly.
	// So, here, use a map to deduplicate.
	history := make(map[int]struct{})
	for i := startIndex; i < len(nums); i++ {
		if _, exist := history[nums[i]]; exist {
			continue
		}
		history[nums[i]] = struct{}{}

		// ensure subsequencePath is an ascendant slice
		if len(subsequencePath) > 0 && nums[i] < subsequencePath[len(subsequencePath)-1] {
			continue
		}
		subsequencePath = append(subsequencePath, nums[i])
		subsequenceBacktracking(nums, i+1, sum+nums[i])
		subsequencePath = subsequencePath[:len(subsequencePath)-1]
	}
}

// 全排列
//
// 给定一个不含重复数字的数组 nums ，返回其 所有可能的全排列 。你可以 按任意顺序 返回答案
//
// 把N叉树画出来，不难的
func permute(nums []int) [][]int {
	permuteResult = make([][]int, 0)
	permutePath = make([]int, 0)
	permuteBacktracking(nums)
	return permuteResult
}

var (
	permuteResult [][]int
	permutePath   []int
)

// 关键在于每次迭代时for循环的集合元素过滤
// 可以使用used数组记录（空间换时间），我这里直接for遍历查找，用时间换空间
func permuteBacktracking(nums []int) {
	if len(permutePath) == len(nums) {
		t := make([]int, len(nums))
		copy(t, permutePath)
		permuteResult = append(permuteResult, t)
		return
	}
	for i := 0; i < len(nums); i++ {
		has := false
		for j := 0; j < len(permutePath); j++ {
			if permutePath[j] == nums[i] {
				has = true
				break
			}
		}
		if has {
			continue
		}
		permutePath = append(permutePath, nums[i])
		permuteBacktracking(nums)
		permutePath = permutePath[:len(permutePath)-1]
	}
}

// 全排列II
//
// 给定一个可包含重复数字的序列 nums ，按任意顺序 返回所有不重复的全排列。
//
//
func permuteUnique(nums []int) [][]int {
	permuteUniqueResult = make([][]int, 0)
	permuteUniquePath = make([]int, 0)
	usedIndex := make(map[int]struct{})
	sort.Ints(nums)
	permuteUniqueBacktracking(nums, usedIndex)
	return permuteUniqueResult
}

var (
	permuteUniqueResult [][]int
	permuteUniquePath   []int
)

func permuteUniqueBacktracking(nums []int, usedIndex map[int]struct{}) {
	if len(permuteUniquePath) == len(nums) {
		t := make([]int, len(nums))
		copy(t, permuteUniquePath)
		permuteUniqueResult = append(permuteUniqueResult, t)
		return
	}
	for i := 0; i < len(nums); i++ {
		_, exist := usedIndex[i-1]
		if i > 0 && nums[i] == nums[i-1] && !exist {
			continue
		}
		if _, exist = usedIndex[i]; exist {
			continue
		}
		usedIndex[i] = struct{}{}
		permuteUniquePath = append(permuteUniquePath, nums[i])
		permuteUniqueBacktracking(nums, usedIndex)
		permuteUniquePath = permuteUniquePath[:len(permuteUniquePath)-1]
		delete(usedIndex, i)
	}
}
