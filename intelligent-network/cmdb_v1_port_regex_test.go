package intelligentnetwork

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"
)

func TestReplace(t *testing.T) {
	inputString := "[1-1],[2-2],[3-3],[4-4],[5-5],[6-6],[7-7],[8-8],[9-9]"

	// 调用函数替换字符串
	result := ReplaceStringsWithSameNumber(inputString)

	fmt.Println("替换前:", inputString)
	fmt.Println("替换后:", result)
}

func TestFinalGenerateRegexp(t *testing.T) {
	type args struct {
		start int
		end   int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// 从1至各个跨度
		{args: args{start: 1, end: 9}, want: "[1-9]"},
		{args: args{start: 1, end: 19}, want: "[1-9]|1[0-9]"},
		{args: args{start: 1, end: 29}, want: "[1-9]|1[0-9]|2[0-9]"},
		{args: args{start: 1, end: 39}, want: "[1-9]|[1-2][0-9]|3[0-9]"},
		{args: args{start: 1, end: 49}, want: "[1-9]|[1-3][0-9]|4[0-9]"},

		// 跨度10
		{args: args{start: 1, end: 15}, want: "[1-9]|1[0-5]"},
		{args: args{start: 16, end: 24}, want: "1[6-9]|2[0-4]"},
		{args: args{start: 25, end: 36}, want: "2[5-9]|3[0-6]"},
		{args: args{start: 37, end: 48}, want: "3[7-9]|4[0-8]"},
		{args: args{start: 49, end: 58}, want: "49|5[0-8]"},

		// 跨度20
		{args: args{start: 1, end: 21}, want: "[1-9]|1[0-9]|2[0-1]"},
		{args: args{start: 1, end: 20}, want: "[1-9]|1[0-9]|20"},
		{args: args{start: 11, end: 31}, want: "1[1-9]|2[0-9]|3[0-1]"},
		{args: args{start: 10, end: 30}, want: "1[0-9]|2[0-9]|30"},
		{args: args{start: 20, end: 40}, want: "2[0-9]|3[0-9]|40"},
		{args: args{start: 21, end: 41}, want: "2[1-9]|3[0-9]|4[0-1]"},
		{args: args{start: 30, end: 50}, want: "3[0-9]|4[0-9]|50"},
		{args: args{start: 31, end: 51}, want: "3[1-9]|4[0-9]|5[0-1]"},

		// 跨度30
		{args: args{start: 1, end: 31}, want: "[1-9]|[1-2][0-9]|3[0-1]"},
		{args: args{start: 11, end: 41}, want: "1[1-9]|[2-3][0-9]|4[0-1]"},
		{args: args{start: 10, end: 40}, want: "1[0-9]|[2-3][0-9]|40"},
		{args: args{start: 20, end: 50}, want: "2[0-9]|[3-4][0-9]|50"},
		{args: args{start: 21, end: 51}, want: "2[1-9]|[3-4][0-9]|5[0-1]"},

		// 跨度40
		{args: args{start: 1, end: 41}, want: "[1-9]|[1-3][0-9]|4[0-1]"},
		{args: args{start: 11, end: 51}, want: "1[1-9]|[2-4][0-9]|5[0-1]"},
		{args: args{start: 11, end: 50}, want: "1[1-9]|[2-4][0-9]|50"},

		// 需求cases
		{args: args{start: 1, end: 2}, want: "[1-2]"},
		{args: args{start: 1, end: 4}, want: "[1-4]"},
		{args: args{start: 1, end: 15}, want: "[1-9]|1[0-5]"},
		{args: args{start: 1, end: 24}, want: "[1-9]|1[0-9]|2[0-4]"},
		{args: args{start: 1, end: 42}, want: "[1-9]|[1-3][0-9]|4[0-2]"},
		{args: args{start: 1, end: 44}, want: "[1-9]|[1-3][0-9]|4[0-4]"},
		{args: args{start: 1, end: 46}, want: "[1-9]|[1-3][0-9]|4[0-6]"},
		{args: args{start: 16, end: 19}, want: "1[6-9]"},
		{args: args{start: 20, end: 24}, want: "2[0-4]"},
		{args: args{start: 21, end: 24}, want: "2[1-4]"},
		{args: args{start: 33, end: 52}, want: "3[3-9]|4[0-9]|5[0-2]"},
		{args: args{start: 43, end: 46}, want: "4[3-6]"},
		{args: args{start: 44, end: 48}, want: "4[4-8]"},
		{args: args{start: 45, end: 48}, want: "4[5-8]"},
		{args: args{start: 49, end: 52}, want: "49|5[0-2]"},
		{args: args{start: 49, end: 56}, want: "49|5[0-6]"},

		// 异常cases
		{args: args{start: 0, end: 56}, want: "", wantErr: true},
		{args: args{start: 56, end: 0}, want: "", wantErr: true},
		{args: args{start: 0, end: 0}, want: "", wantErr: true},
		{args: args{start: 56, end: 56}, want: "", wantErr: true},
		{args: args{start: 56, end: 16}, want: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FinalGenerateRegexp(tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("FinalGenerateRegexp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FinalGenerateRegexp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGeneratePortRegexp(t *testing.T) {
	type args struct {
		ratePrefix    string
		stacking      bool
		portSplitting bool
		totalSlot     int
		portRanges    string
	}
	tests := []struct {
		name                string
		args                args
		want                string
		wantErr             bool
		portShouldBeMatched []string
		portShouldNotMatch  []string
	}{
		{
			name: "[无堆叠]-[无拆分]-[三槽位]-[连续范围]",
			args: args{
				ratePrefix:    "10GE",
				stacking:      false,
				portSplitting: false,
				totalSlot:     3,
				portRanges:    "1-42",
			},
			want:                "^10GE1/0/([1-9]|[1-3][0-9]|4[0-2])$",
			portShouldBeMatched: []string{"10GE1/0/1", "10GE1/0/41"},
			portShouldNotMatch:  []string{"10GE1/0/0", "10GE1/0/43"},
		},
		{
			name: "[无堆叠]-[无拆分]-[三槽位]-[非连续范围]",
			args: args{
				ratePrefix:    "TwentyFiveGigE",
				stacking:      true,
				portSplitting: false,
				totalSlot:     3,
				portRanges:    "1-24,33-52",
			},
			want:                "^TwentyFiveGigE([1-2])/0/([1-9]|1[0-9]|2[0-4]|3[3-9]|4[0-9]|5[0-2])$",
			portShouldBeMatched: []string{"TwentyFiveGigE1/0/1", "TwentyFiveGigE2/0/52"},
			portShouldNotMatch:  []string{"TwentyFiveGigE1/0/25", "TwentyFiveGigE2/0/53"},
		},
		{
			name: "[有堆叠]-[无拆分]-[三槽位]-[连续范围]",
			args: args{
				ratePrefix:    "10GE",
				stacking:      true,
				portSplitting: false,
				totalSlot:     3,
				portRanges:    "1-42",
			},
			want:                "^10GE([1-2])/0/([1-9]|[1-3][0-9]|4[0-2])$",
			portShouldBeMatched: []string{"10GE1/0/1", "10GE2/0/1", "10GE1/0/42", "10GE2/0/42"},
			portShouldNotMatch:  []string{"10GE3/0/1", "10GE2/0/43", "10GE1/0/43"},
		},
		{
			name: "[有堆叠]-[有拆分]-[三槽位]-[连续范围]",
			args: args{
				ratePrefix:    "10GE",
				stacking:      true,
				portSplitting: true,
				totalSlot:     3,
				portRanges:    "1-42",
			},
			want:                "^10GE([1-2])/0/([1-9]|[1-3][0-9]|4[0-2]):[1-4]$",
			portShouldBeMatched: []string{"10GE1/0/1:1", "10GE2/0/1:2", "10GE1/0/42:3", "10GE2/0/42:4"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := GeneratePortRegexp(tt.args.ratePrefix, tt.args.stacking, tt.args.portSplitting, tt.args.totalSlot, tt.args.portRanges)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePortRegexp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GeneratePortRegexp() = %v, want %v", got, tt.want)
			}

			expr, _ := regexp.Compile(got)
			for _, v := range tt.portShouldBeMatched {
				if !expr.MatchString(v) {
					t.Errorf("regexp:%s should match port:%s", got, v)
				}
			}

			for _, v := range tt.portShouldNotMatch {
				if expr.MatchString(v) {
					t.Errorf("regexp:%s should not match port:%s", got, v)
				}
			}

		})
	}
}

func TestPortRanges_GetCount(t *testing.T) {
	tests := []struct {
		name    string
		p       PortRanges
		want    int
		wantErr bool
	}{
		{p: "1-10", want: 10},
		{p: "1-10,15-20", want: 16},
		{p: "1-24,33-52", want: (52 - 33 + 1) + 24},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.GetCount()
			if (err != nil) != tt.wantErr {
				t.Errorf("PortRanges.GetCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PortRanges.GetCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmulatePostHandler(t *testing.T) {
	var data = []byte(`{"downlink_port":{"rate_prefix":"Twenty-FiveGigE","slot_number":3,"ranges":"1-24,33-52","stacking":true,"port_splitting":false},"uplink_port":{"rate_prefix":"HundredGige","slot_number":3,"ranges":"53-56"},"reserved_port":{"rate_prefix":"Twenty-FiveGigE","slot_number":3,"ranges":"25-32"},"total_port":56}`)
	var parameter = ModelPortQuotaParameter{}
	err := json.Unmarshal(data, &parameter)
	if err != nil {
		t.Error(err)
	}

	resp, err := EmulatePostHandler(parameter)
	if err != nil {
		t.Error(err)
	}
	t.Log(PrettyStruct(resp))
}

// func PrettyStruct(data interface{}) string {
// 	val, err := json.MarshalIndent(data, "", "    ")
// 	if err != nil {
// 		return ""
// 	}
// 	return string(val)
// }
