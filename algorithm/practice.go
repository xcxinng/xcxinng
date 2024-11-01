package algorithm

import (
	"fmt"
	"math"
	"strings"
)

// Happen to hear this interviewing question from my current leader,
// who had interviewed a more than 5-year-experience gaming developer.
//
// Surprisingly, the interviewee didn't give my leader a satisfying
// solution, by which I was shocked.
//
// From my point of view, a developer who has a little fundamental
// coding knowledge would kill this problem in a quick way.
//
// Description: Given a string s, print each character of s in a
// reverse order and MUST solve it by a recursion manner.
func PrintReverseString(s string) {
	if len(s) < 2 {
		fmt.Println(s)
		return
	}
	recurseString(s, 0)
}

// the wheel that prints each character of s in a reverse order.
func recurseString(s string, i int) {
	// the condition to break the recursion
	if i == len(s) {
		return
	}
	// enter the recursion
	recurseString(s, i+1)

	// operation for string s
	fmt.Print(string(s[i]))
}

var (
	value int64 = math.MinInt64
	count int
)

func isValidBST(root *TreeNode) bool {
	count = 0
	value = math.MinInt64
	recurseTree(root)
	return count == 0
}

func recurseTree(root *TreeNode) {
	if root == nil {
		return
	}
	recurseTree(root.Left)

	if root.Val <= int(value) {
		count++
	}
	value = int64(root.Val)

	recurseTree(root.Right)
}

func bfsPrint(root *TreeNode) {
	var queue []*TreeNode
	if root == nil {
		return
	}

	// If root was a graph, an extra map was needed to
	// recognize that a node whether has been accessed.
	queue = append(queue, root)
	for len(queue) != 0 {
		// pop an element from queue
		node := queue[0]
		queue = queue[1:]
		fmt.Println(node.Val)

		// push into queue if not nil
		if node.Left != nil {
			queue = append(queue, node.Left)
		}
		if node.Right != nil {
			queue = append(queue, node.Right)
		}
	}
}

// The core of BFS algorithm is that it uses a queue to store
// all the nodes that will be visited and leverages the "FIFO"
// charactertistic of a queue.
func levelOrder(root *TreeNode) [][]int {
	result := [][]int{}
	if root == nil {
		return result
	}
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		nextLevelQueue := []*TreeNode{}
		values := []int{}
		// clear the queue of current level
		for i := 0; i < len(queue); i++ {
			node := queue[i]
			values = append(values, node.Val)

			// put all children that connected with the nodes in the current level
			// into the nextLevelQueue, which represents the next level queue
			// containing all nodes that will be traversed.
			if node.Left != nil {
				nextLevelQueue = append(nextLevelQueue, node.Left)
			}
			if node.Right != nil {
				nextLevelQueue = append(nextLevelQueue, node.Right)
			}
		}
		// update queue
		queue = nextLevelQueue
		result = append(result, values)
	}
	return result
}

func practiceIsSymmetric(root *TreeNode) bool {
	return check2(root, root)
}

func check2(l *TreeNode, r *TreeNode) bool {
	if l == nil && r == nil {
		return true
	}
	if l == nil || r == nil {
		return false
	}
	return l.Val == r.Val && check2(l.Left, r.Right) && check2(l.Right, r.Left)
}

func findOcurrences(text string, first string, second string) []string {
	var ret []string
	words := strings.Split(text, " ")
	if len(words) < 3 {
		return nil
	}

	for i := 0; i < len(words)-2; i++ {
		if words[i] == first && words[i+1] == second {
			ret = append(ret, words[i+2])
		}
	}
	return ret
}
