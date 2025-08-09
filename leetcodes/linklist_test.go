package algorithm

import (
	"testing"
)

func Test_add2(t *testing.T) {
	type args struct {
		l1 *ListNode
		l2 *ListNode
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// {args: args{
		// 	l1: GenerateLinkListWithArray([]int{2, 1, 7}),
		// 	l2: GenerateLinkListWithArray([]int{5, 5}),
		// }},
		{args: args{
			l1: GenerateLinkListWithArray([]int{1, 7}),
			l2: GenerateLinkListWithArray([]int{5, 5}),
		}},
		// {args: args{
		// 	l1: GenerateLinkListWithArray([]int{2, 4, 3}),
		// 	l2: GenerateLinkListWithArray([]int{5, 6, 4}),
		// }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AddTwoNumbers(tt.args.l1, tt.args.l2)
			PrintlnLinkListAsArray(got, false)
		})
	}
}

func Test_reverseListNode(t *testing.T) {
	type args struct {
		head *ListNode
	}
	tests := []struct {
		name string
		args args
		want *ListNode
	}{
		{
			args: args{
				head: GenerateLinkListWithArray([]int{1, 2, 3}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintlnLinkListAsArray(ReverseListNode(tt.args.head), false)
		})
	}
}

func Test_mergeTwoLists2(t *testing.T) {
	type args struct {
		list1 *ListNode
		list2 *ListNode
	}
	tests := []struct {
		name string
		args args
		want *ListNode
	}{
		{args: args{
			list1: GenerateLinkListWithArray([]int{1, 2, 3}),
			list2: GenerateLinkListWithArray([]int{4, 5, 6}),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergeTwoLists2(tt.args.list1, tt.args.list2)
			PrintlnLinkListAsArray(got, false)
		})
	}
}

func Test_sortList(t *testing.T) {
	type args struct {
		head *ListNode
	}
	tests := []struct {
		name string
		args args
		want *ListNode
	}{
		{args: args{
			head: GenerateLinkListWithArray([]int{4, 2, 1, 3}),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sortList(tt.args.head)
			PrintlnLinkListAsArray(got, false)
		})
	}
}

func Test_swapPairs(t *testing.T) {
	type args struct {
		head *ListNode
	}
	tests := []struct {
		name string
		args args
		want *ListNode
	}{
		{
			args: args{
				head: GenerateLinkListWithArray([]int{1, 2, 3, 4}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := swapPairs(tt.args.head)
			PrintlnLinkListAsArray(got, false)
		})
	}
}

func Test_reverseKGroup(t *testing.T) {
	type args struct {
		head *ListNode
		k    int
	}
	tests := []struct {
		name string
		args args
		want *ListNode
	}{
		{args: args{head: GenerateLinkListWithArray([]int{1, 2, 3, 4, 5}), k: 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := reverseKGroup(tt.args.head, tt.args.k); !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("reverseKGroup() = %v, want %v", got, tt.want)
			// }
			got := reverseKGroup(tt.args.head, tt.args.k)
			PrintlnLinkListAsArray(got, false)
		})
	}
}

func Test_rotateRight(t *testing.T) {
	type args struct {
		head *ListNode
		k    int
	}
	tests := []struct {
		name string
		args args
		want *ListNode
	}{
		{
			args: args{head: GenerateLinkListWithArray([]int{1, 2, 3}), k: 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rotateRight(tt.args.head, tt.args.k)
			PrintlnLinkListAsArray(got, false)
		})
	}
}

func Test_deleteDuplicates(t *testing.T) {
	type args struct {
		head *ListNode
	}
	tests := []struct {
		name string
		args args
		want *ListNode
	}{
		{args: args{head: GenerateLinkListWithArray([]int{1, 1, 1, 2, 3})}},
		{args: args{head: GenerateLinkListWithArray([]int{1, 2, 2, 3, 3, 4, 5})}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := deleteDuplicates(tt.args.head)
			PrintlnLinkListAsArray(got, false)
		})
	}
}
