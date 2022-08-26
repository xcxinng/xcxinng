package algorithm

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
	combineResult = make([][]int, 0)
	combinePath = make([]int, 0)
	return nil
}

// 递归就有点像需要动态套不同for循环，如果不用递归，都不知道代码要咋写
// 这里的path
func combineBacktracking(n, k, startIndex int) {
	if len(combinePath) == k {
		combineResult = append(combineResult, combinePath)
		return
	}

	// add something to be deleted
	for i := startIndex; i < n-(k-len(combinePath))+1; i++ {
		combinePath = append(combinePath, i)
		combineBacktracking(n, k, i+1)
		combinePath = combinePath[:len(combinePath)-1]
	}
}

var (
	combineResult [][]int
	combinePath   []int
)
