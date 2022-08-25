package algorithm

import (
	"fmt"
	"testing"
)

var nums = []int{49, 38, 65, 70, 13, 1}

func TestStraightInsertionSort(t *testing.T) {
	nums = []int{0, 49, 38, 65, 70, 13, 1}
	straightInsertionSort(nums)
	fmt.Println(nums)
	// output:
	// [whatever,1 13 38 49 65 70]
}

func TestBubbleSort(t *testing.T) {
	BubbleSort(nums)
	fmt.Println(nums)
}

func TestQuickSort(t *testing.T) {
	QuickSort(nums)
	fmt.Println(nums)
	// output:
	// [1 13 38 49 65 70]
}

func TestSimpleSelectSort(t *testing.T) {
	SimpleSelectSort(nums)
	fmt.Println(nums)
	nums = []int{5, 4, 3, 2, 1}
	SimpleSelectSort(nums)
	fmt.Println(nums)
	nums = []int{2, 0, 1}
	SimpleSelectSort(nums)
	fmt.Println(nums)
	// output:
	// [1 13 38 49 65 70]
	// [1 2 3 4 5]
	// [1]
}

func TestHeapSort(t *testing.T) {
	// n := []int{4, 9, 10, 0, -4, 7}
	n := []int{10, 4, 2, 5}
	HeapSort(n)
	fmt.Println(n)
	fmt.Println(3 / 2)
	// n = []int{1000, 100, 10, 0, -1000}
	// HeapSort(n)
	// fmt.Println(n)
	// n = []int{1, 1, 1, 1}
	// HeapSort(n)
	// fmt.Println(n)
	// n = []int{}
	// HeapSort(n)
	// fmt.Println(n)
	// n = []int{111}
	// HeapSort(n)
	// fmt.Println(n)
	// output:
	// [-4 0 4 7 9 10]
	// [-1000 0 10 100 1000]
	// [1 1 1 1]
	// []
	// [111]
}
