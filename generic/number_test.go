package generic

import (
	"encoding/json"
	"testing"
)

func TestValidate(t *testing.T) {
	type args struct {
		av *AttributeValue
		t  ValueType
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[int] 正常输入",
			args: args{av: &AttributeValue{
				DefaultValue: int64(1),
				Min:          int64(1),
				Max:          int64(100),
				ValueType:    TypeInt,
			}},
		},
		{
			name: "[float] 正常输入",
			args: args{av: &AttributeValue{
				DefaultValue: float64(0.5),
				Min:          float64(0.0),
				Max:          float64(1.0),
				ValueType:    TypeFloat,
			}},
		},
		{
			name: "[float] 混合使用1",
			args: args{av: &AttributeValue{
				DefaultValue: int64(10),
				Min:          float64(1.1),
				Max:          float64(100.1),
				ValueType:    TypeFloat,
			}},
		},
		{
			name: "[float] 混合使用2",
			args: args{av: &AttributeValue{
				DefaultValue: int64(10),
				Min:          float64(1.1),
				Max:          float64(100.1),
				ValueType:    TypeInt,
			}},
		},
		{
			name: "[int] 混合使用1",
			args: args{av: &AttributeValue{
				DefaultValue: float64(1.99),
				Min:          int64(1),
				Max:          int64(100),
				ValueType:    TypeInt,
			}},
		},
		{
			name: "[int] 混合使用1",
			args: args{av: &AttributeValue{
				DefaultValue: float64(1.99),
				Min:          int64(1),
				Max:          int64(100),
				ValueType:    TypeFloat,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.args.av.ValueType {
			case TypeFloat:
				if err := Validate[float64](tt.args.av, tt.args.t); (err != nil) != tt.wantErr {
					t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				} else {
					data, _ := json.Marshal(tt.args.av)
					t.Log(string(data))
				}
			case TypeInt:
				if err := Validate[int64](tt.args.av, tt.args.t); (err != nil) != tt.wantErr {
					t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				} else {
					data, _ := json.Marshal(tt.args.av)
					t.Log(string(data))
				}
			}
		})
	}
}

func TestAttributeValue_NormalizeNumberType(t *testing.T) {
	type fields struct {
		DefaultValue interface{}
		EnumValues   []string
		Max          interface{}
		Min          interface{}
		ValueType    ValueType
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			fields: fields{
				DefaultValue: int64(100),
				Max:          float64(200.9),
				Min:          float64(10.9),
				ValueType:    TypeFloat,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := &AttributeValue{
				DefaultValue: tt.fields.DefaultValue,
				EnumValues:   tt.fields.EnumValues,
				Max:          tt.fields.Max,
				Min:          tt.fields.Min,
				ValueType:    tt.fields.ValueType,
			}
			if err := av.NormalizeNumberType(); (err != nil) != tt.wantErr {
				t.Errorf("AttributeValue.NormalizeNumberType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
