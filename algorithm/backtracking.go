package algorithm

import "sort"

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
// 给你一个 无重复元素 的整数数组 candidates 和一个目标整数 target，
// 找出 candidates 中可以使数字和为目标数 target 的 所有 不同组合，
// 并以列表形式返回。你可以按 任意顺序 返回这些组合。
//
// candidates 中的 同一个 数字可以 无限制重复被选取 。
// 如果至少一个数字的被选数量不同，则两种组合是不同的。
//
// 对于给定的输入，保证和为 target 的不同组合数少于 150 个。
//
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
// 仔细分析与 [组合总和] 的区别
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
//
// 回文串是正着读和反着读都一样的字符串。
//
// Example:
// 输入：s = "aab"
// 输出：[["a","a","b"],["aa","b"]]
//
// "所有可能" ==> 回溯 ===> 穷举
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
