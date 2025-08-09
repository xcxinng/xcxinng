package main

import (
	"time"
)

var a int64 = 10

// func main() {
// 	go func() {
// 		// atomic 并不能解决排序问题
// 		// atomic 类似于是一个内存操作事务，事务中多个内存操作完成了，其结果才对其他goroutine可见
// 		println(atomic.LoadInt64(&a))
// 	}()
// 	go func() {
// 		// atomic 并不能解决排序问题
// 		atomic.AddInt64(&a, 10)
// 	}()
// 	time.Sleep(time.Millisecond * 100)
// }

func main() {
	unbufferedCh := make(chan struct{})
	go func() {
		<-unbufferedCh
		println(a)
	}()
	go func() {
		a += 10
		unbufferedCh <- struct{}{}
	}()
	time.Sleep(time.Millisecond * 100)
}
