package generic

import (
	"errors"
	"reflect"
)

// SumIntsOrFloats sums the values of map m. It supports both int64 and float64
// as types for map values.
func SumIntsOrFloats[K comparable, V int64 | float64](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

type AttributeValue struct {
	DefaultValue interface{}
	EnumValues   []string
	Max          interface{}
	Min          interface{}
	ValueType    ValueType
}

type Item struct {
	Value reflect.Value
	K     reflect.Kind
}

func (av *AttributeValue) NormalizeNumberType() error {
	store := map[string]ValueType{}
	// 逐个断言，哪个成功了就计数加一
	// 最后，计数最大必须等于 ValueType，否则认为错误；计数小的，需要进行类型转换，转换成计数大的
	if av.DefaultValue != nil {
		_, ok := av.DefaultValue.(float64)
		if !ok {
			_, ok1 := av.DefaultValue.(int64)
			if ok1 {
				store["default"] = TypeInt
			} else {
				return errors.New("error")
			}
		} else {
			store["default"] = TypeFloat
		}
	}
	if av.Max != nil {
		_, ok := av.Max.(float64)
		if !ok {
			_, ok1 := av.Max.(int64)
			if ok1 {
				store["max"] = TypeInt
			} else {
				return errors.New("error")
			}

		} else {
			store["max"] = TypeFloat
		}
	}
	if av.Min != nil {
		_, ok := av.Min.(float64)
		if !ok {
			_, ok2 := av.Min.(int64)
			if ok2 {
				store["min"] = TypeInt
			} else {
				return errors.New("error")
			}
		} else {
			store["min"] = TypeFloat
		}
	}

	expectedValueType := av.ValueType
	for k, actualType := range store {
		if actualType != expectedValueType {
			switch k {
			case "default":
				if actualType == TypeFloat { // 期望是int,但实际是float
					av.DefaultValue = int64(av.DefaultValue.(float64))
				} else {
					// 期望是float,但实际是int
					av.DefaultValue = float64(av.DefaultValue.(int64))
				}

			case "max":
				if actualType == TypeFloat {
					av.Max = int64(av.Max.(float64))
				} else {
					av.Max = float64(av.Max.(int64))
				}
			case "min":
				if actualType == TypeFloat { // 期望是int,但实际是float
					av.Min = int64(av.Min.(float64))
				} else {
					av.Min = float64(av.Min.(int64))
				}
			}
		}
	}

	return nil
}

type ValueType string

const (
	TypeInt   ValueType = "int"
	TypeFloat ValueType = "float"
)

func Validate[V int64 | float64](av *AttributeValue, t ValueType) error {
	if err := av.NormalizeNumberType(); err != nil {
		return err
	}
	if av.Max != nil && av.Min != nil {
		gMax, gMin := av.Max.(V), av.Min.(V)
		if gMax < gMin {
			return errors.New("minimum should less than or equal to the maximum")
		}
	}
	if av.DefaultValue != nil {
		gdv := av.DefaultValue.(V)
		av.DefaultValue = gdv

		if av.Max != nil {
			gMax := av.Max.(V)
			if gdv > gMax {
				return errors.New("default value must less than or equal to the maximum")
			}
			av.Max = gMax
		}

		if av.Min != nil {
			gMin := av.Min.(V)
			if gdv < gMin {
				return errors.New("default value must greater than or equal to the minimum")
			}
			av.Min = gMin
		}
	}

	if av.Max != nil {
		av.Max = av.Max.(V)
	}

	if av.Min != nil {
		av.Min = av.Min.(V)
	}
	return nil
}
