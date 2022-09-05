package algorithm

import (
	"fmt"
	"testing"
)

func TestCombine(t *testing.T) {
	fmt.Println(combine(1, 1))
	fmt.Println(combine(4, 2))
	// output:
	// [[1]]
	// [[1 2] [1 3] [1 4] [2 3] [2 4] [3 4]]
}

func TestCombinationSum3(t *testing.T) {
	fmt.Println(combinationSum3(3, 7))
	fmt.Println(combinationSum3(3, 9))
	// output:
	// [[1 2 4]]
	// [[1 2 6] [1 3 5] [2 3 4]]
}

func TestLetterCombinations(t *testing.T) {
	fmt.Println(letterCombinations("23"))
}

func TestCombinationSum(t *testing.T) {
	fmt.Println(combinationSum([]int{2, 3, 5}, 8))
}

func TestCombinationSum2(t *testing.T) {
	fmt.Println(combinationSum2([]int{2, 2}, 2))
	fmt.Println(combinationSum2([]int{10, 1, 2, 7, 6, 1, 5}, 8))
}

func TestPartition(t *testing.T) {
	fmt.Println(partition("aab"))
}
