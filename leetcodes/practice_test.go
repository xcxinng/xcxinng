package algorithm

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintReversedString(t *testing.T) {
	PrintReverseString("abcef")
	PrintReverseString("a")
	PrintReverseString("ab")
	PrintReverseString("")
	PrintReverseString("1234")
	// output:
	// fecba
	// a
	// ba
	//
	// 4321
}

func Test_isValidBST(t *testing.T) {
	type args struct {
		root *TreeNode
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			args: args{
				root: &TreeNode{
					Val:  5,
					Left: &TreeNode{Val: 1},
					Right: &TreeNode{
						Val:   4,
						Left:  &TreeNode{Val: 3},
						Right: &TreeNode{Val: 6},
					},
				}},
			want: false,
		},
		{
			args: args{root: &TreeNode{Val: 0}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidBST(tt.args.root); got != tt.want {
				t.Errorf("isValidBST() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bfsPrint(t *testing.T) {
	// node := &TreeNode{
	// 	Val:  1,
	// 	Left: &TreeNode{Val: 2},
	// 	Right: &TreeNode{
	// 		Val:   3,
	// 		Left:  &TreeNode{Val: 4},
	// 		Right: &TreeNode{Val: 5},
	// 	},
	// }
	// bfsPrint(node)
	// output: 12345

	fmt.Println()

	node := &TreeNode{
		Val: 5,
		Left: &TreeNode{
			Val:   1,
			Left:  &TreeNode{Val: 100},
			Right: &TreeNode{Val: 200}},
		Right: &TreeNode{
			Val:   4,
			Left:  &TreeNode{Val: 3},
			Right: &TreeNode{Val: 6},
		},
	}
	bfsPrint(node)
	// output:
	// 51436
}

func Test_practiceIsSymmetric(t *testing.T) {
	node := &TreeNode{
		Val: 1,
		Left: &TreeNode{
			Val:   2,
			Left:  &TreeNode{Val: 2},
			Right: nil},
		Right: &TreeNode{
			Val:   2,
			Left:  &TreeNode{Val: 2},
			Right: nil,
		}}
	fmt.Println(practiceIsSymmetric(node))
}

func Test_findOcurrences(t *testing.T) {
	type args struct {
		text   string
		first  string
		second string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			args: args{text: "She is a good girl she is a student", first: "is", second: "a"},
			want: []string{"good", "student"},
		},
		{
			args: args{text: "we will we will rock you", first: "we", second: "will"},
			want: []string{"we", "rock"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findOcurrences(tt.args.text, tt.args.first, tt.args.second); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findOcurrences() = %v, want %v", got, tt.want)
			}
		})
	}
}

type PortValues struct {
	Device string
	Port   string
	Values []V
}

type V struct {
	Value int
	Ts    int
}

type ValueIndex struct {
	Value      int
	FirstIndex int
}

func TestDeduplicate(t *testing.T) {
	_, ok := isOk([]V{
		{Value: 1, Ts: 1},
		{Value: 1, Ts: 2},
		{Value: 0, Ts: 3},
		{Value: 0, Ts: 4},
		{Value: 1, Ts: 5},
		{Value: 1, Ts: 6},
	})
	assert.EqualValues(t, false, ok)

	_, ok = isOk([]V{
		{Value: 1, Ts: 1},
		{Value: 1, Ts: 2},
		{Value: 1, Ts: 3},
		{Value: 1, Ts: 4},
		{Value: 0, Ts: 5},
		{Value: 0, Ts: 6},
		{Value: 0, Ts: 7},
		{Value: 0, Ts: 8},
		{Value: 1, Ts: 9},
		{Value: 1, Ts: 10},
		{Value: 1, Ts: 11},
		{Value: 1, Ts: 12},
	})

	assert.EqualValues(t, true, ok)

	_, ok = isOk([]V{
		{Value: 1, Ts: 1},
		{Value: 1, Ts: 2},
		{Value: 0, Ts: 3},
		{Value: 0, Ts: 4},
		{Value: 1, Ts: 5},
		{Value: 1, Ts: 6},
		{Value: 1, Ts: 7},
		{Value: 1, Ts: 8},
		{Value: 1, Ts: 9},
		{Value: 1, Ts: 10},
		{Value: 1, Ts: 11},
		{Value: 1, Ts: 12},
	})

	assert.EqualValues(t, false, ok)
}

type Stack struct {
	v []ValueIndex
}

func (s *Stack) Pop() ValueIndex {
	e := s.v[len(s.v)-1]
	s.v = s.v[:len(s.v)-1]
	return e
}
func (s *Stack) Push(e ValueIndex) {
	s.v = append(s.v, e)
}
func (s *Stack) Len() int {
	return len(s.v)
}

func newStack() *Stack {
	return &Stack{v: make([]ValueIndex, 0)}
}

func isOk(elems []V) ([]int, bool) {
	// elems sorted by timestamp
	var proceedOne = -10
	stack := newStack()
	counting := make(map[int]int)
	for i, v := range elems {

		if proceedOne == -10 { // the first elem
			proceedOne = v.Value
			stack.Push(ValueIndex{Value: v.Value, FirstIndex: i})
			counting[i] = 1
			continue
		}

		e := stack.Pop()
		if v.Value == e.Value {
			counting[e.FirstIndex]++
			stack.Push(e)
		} else {
			stack.Push(e)
			stack.Push(ValueIndex{Value: v.Value, FirstIndex: i})
			counting[i] = 1

		}
	}
	if len(counting) != 3 {
		return nil, false
	}
	var resp []int
	for _, v := range stack.v {
		if counting[v.FirstIndex] <= 3 { //假设间隔要大于1小时
			return nil, false
		}
		resp = append(resp, counting[v.FirstIndex])

	}
	return resp, true
}
