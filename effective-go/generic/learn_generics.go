package main

import (
	"fmt"
	"log"
	"strconv"

	"golang.org/x/exp/constraints"
)

func SumInts(data map[string]int64) int64 {
	var sum int64
	for _, i := range data {
		sum += i
	}
	return sum
}

func SumFloats(data map[string]float64) float64 {
	var sum float64
	for _, i := range data {
		sum += i
	}
	return sum
}

// func SumIntsOrFloats[V int64 | float64](m map[string]V) V {
// 	var s V
// 	for _, v := range m {
// 		s += v
// 	}
// 	return s
// }

func convertToIntPrt[T uint8 | uint16 | uint32](t T) *int {
	var a int
	a = int(t)
	return &a
}

func learnGenerics() {
	//var (
	//	ints   = map[string]int64{"1st": 1, "2nd": 2}
	//	floats = map[string]float64{"1st": 1, "2nd": 2}
	//)
	//
	//fmt.Printf("non-generics: SumInts:%v,SumFloats:%v\n",
	//	SumInts(ints),
	//	SumFloats(floats))
	//
	//fmt.Printf("generics: SumInts:%v,SumFloats:%v\n",
	//	SumIntsOrFloats[int64](ints),
	//	SumIntsOrFloats[float64](floats))

	// you can omit the type argument when calling the generic function
	//fmt.Printf("generic that omit type argument: SumInts:%v\n", SumIntsOrFloats(ints))
	//fmt.Printf("generic that imit type argument: SumFloats:%v\n", SumIntsOrFloats(floats))

	//callSumNumbers()
	// fmt.Println(convertAny2String(myFloat(1.0)))
}

// Number 这种用法应该是1.18才支持的,可以认为这种一种范型interface,经过测试，仅支持内建类型
type Number interface {
	int64 | float64
}

// ...T 与 ...interface 不同之处在于前者会保证该次输入的所有元素类型一致，后者则不会检查类型
// 这也是 1.18 新增
// 类型约束通常是一个类型集合，但在编译时，编译器把它当作一个具体的类型对待
// 具体的类型取决于调用者指定的类型或根据实参推断出来
func printLn[T any](elems ...T) {
	for _, elem := range elems {
		fmt.Println(elem)
	}
}

type StringerFloat interface {
	float32 | ~float64
	String() string
}

func convertAny2String[F StringerFloat](value F) string {
	return value.String()

}

func convertAndDoStuff[T constraints.Unsigned](s string, setter func(c T)) {
	c, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		log.Println(err)
	}
	setter(T(c))
}
