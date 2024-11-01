package algorithm

import (
	"fmt"
	"reflect"
	"testing"
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
