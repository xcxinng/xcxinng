package main

import (
	"time"
)

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			time.Sleep(time.Millisecond * 500)
			out <- n
		}
		close(out)
	}()
	return out
}

func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			time.Sleep(time.Second)
			out <- n * n
		}
		close(out)
	}()
	return out
}

//func main() {
//	now := time.Now()
//	// Set up the pipeline.
//	c := gen(2, 3)
//	out := sq(c)
//
//	// Consume the output.
//	fmt.Println(<-out) // 4
//	fmt.Println(<-out) // 9
//	fmt.Println(time.Now().Sub(now).Milliseconds())
//}
