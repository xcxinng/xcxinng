package algorithm

// Complexity:
// Time : O(n^2)
// Space: O(1)
//
// Note: nums[0] is used as a guard, so the caller should keep it empty.
func straightInsertionSort(nums []int) {
	var j = 0
	for i := 2; i < len(nums); i++ {
		if nums[i] < nums[i-1] {
			nums[0] = nums[i]                        // write it to the nums[0] avoiding being overwritten
			nums[i] = nums[i-1]                      // move the large one backward
			for j := i - 2; nums[0] < nums[j]; j-- { // move them all backward
				nums[j+1] = nums[j]
			}
			nums[j+1] = nums[0] // write nums[0] to the right place
		}
	}
}

// BubbleSort is a classic bubble sort algorithm.
//
// Time Complexity: O(n^2)
// Space Complexity: O(0)
//
// Description:
// In bubble sort, the smaller numbers will be moved leftward,
// and the larger numbers rightward.
func BubbleSort(nums []int) {
	for i := 0; i < len(nums); i++ {
		for j := 0; j < len(nums)-1; j++ {
			if nums[j] > nums[j+1] {
				nums[j+1], nums[j] = nums[j], nums[j+1]
			}
		}
	}
}

// partition returns the location of the pivot.
//
// numbers[i] <= nums[pivot], low <= i <= pivotIndex;
// numbers[i] >= nums[pivot], pivotIndex <= i <= high.
func partition(nums []int, low int, high int) (pivotIndex int) {
	pivotValue := nums[high]
	pivotIndex = low
	for low < high {
		if nums[low] <= pivotValue {
			nums[pivotIndex], nums[low] = nums[low], nums[pivotIndex]
			pivotIndex++
		}
		low++
	}
	nums[pivotIndex], nums[high] = pivotValue, nums[pivotIndex]
	return
}

func _quickSort(nums []int, low, high int) {
	if low >= high {
		return
	}
	// "divide and conquer"
	pivotLoc := partition(nums, low, high)
	_quickSort(nums, low, pivotLoc-1)
	_quickSort(nums, pivotLoc+1, high)
}

// QuickSort
//
// Time complexity: O(nlogn)
// Space complexity: O(1) but it has stack overhead, at the worst situation: O(n)
//
// It's well-known that QuickSort has the best average performance among all
// the sort algorithms with time complexity O(nlogn).
func QuickSort(nums []int) {
	_quickSort(nums, 0, len(nums)-1)
}

// SimpleSelectSort
// Time complexity: O(n^2)
// Space complexity: O(1)
func SimpleSelectSort(nums []int) {
	for i := 0; i < len(nums); i++ {
		j := min(nums, i, len(nums))
		if i != j {
			nums[i], nums[j] = nums[j], nums[i]
		}
	}
}

func min(nums []int, left, right int) (index int) {
	v := nums[left]
	index = left
	for i := left; i < right; i++ {
		if nums[i] < v {
			index = i
			v = nums[i]
		}
	}
	return
}

// heapAdjust adjusts nums[s,m] to meet the definition of the heap:
// Ki >= K2i
// Ki >= K2i+1
// 1<=  i  <=  len(nums)/2
//
// Before adjust start, nums[s+1]....nums[m] should be a heap already.
//
// s stands for the unsorted node(index) which will be adjusted by heapAdjust.
// m stands for the number of the heap's node.
//
// Reference from <data structure>:
// LChild(i) = 2*i (2i<=n)
// RChild(i) = 2*i + 1  (2i+1<=n)
func heapAdjust(nums []int, s, m int) {
	rc := nums[s] // stores the value of which should be adjusted

	// Q: Is the purpose of this loop is to choose the right index for nums[s] ??
	// A: Yes, it is.
	//
	// Q: Do you really know the point of *2 ?
	// A: j times 2, it means heading into the next child tree,
	//    the direction from which is the nums[s] node to the
	//    termination node.
	for j := 2 * s; j <= m; j *= 2 {
		if j < m && nums[j] < nums[j+1] {
			// compare nums[j] with nums[j+1] (both are child nodes for node s)
			// then choose the large one.
			j++ // j = j + 1
		}

		if !(rc < nums[j]) {
			// If rc is smaller than nums[j], meaning the j is not the valid index,
			// and continue to head into j's child tree.
			//
			// Otherwise, j is the valid index, and should break this loop
			break
		}

		// rc smaller than nums[j] then do:
		nums[s] = nums[j]
		s = j
	}
	nums[s] = rc // do a swap between the last nums[j] and the original nums[s]
}

// HeapSort sort nums with HeapSort algorithm.
// Time complexity: O(nlogn)
// Space complexity: O(1)
//
// 摘抄自严蔚敏的《数据结构-C语言版》P280:
//
// 解决堆排序需要解决2个问题：
// (1)如何由一个无序序列建成一个堆
// (2)如何在输出堆顶元素之后，调整剩余元素使成为一个新的堆
//
// " 从一个无序序列建堆的过程就是一个反复"筛选"的过程。"
func HeapSort(nums []int) {
	if len(nums) == 0 {
		return
	}

	// adjust nums into a heap (Ki>=K2i,Ki>=K2i+1)
	// As heap is like a complete binary tree, we just need to
	// adjust interval [0,len/2]
	for i := (len(nums) - 1) / 2; i >= 0; i-- {
		heapAdjust(nums, i, len(nums)-1)
	}

	for i := len(nums) - 1; i > 0; i-- {
		// swap the maximum nums[i] with nums[0];
		nums[0], nums[i] = nums[i], nums[0]
		// secondly, heapify nums[0] to nums[i-1] into a heap.
		heapAdjust(nums, 0, i-1)
	}
}
