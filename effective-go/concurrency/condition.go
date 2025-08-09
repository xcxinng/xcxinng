package main

import (
	"fmt"
	"sync"
	"time"
)

var queue chan []Metric

// 消费者（错误示例：忙等待）
func consumer() {
	metricBuffer := make([]Metric, 0, 5000)
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case ms := <-queue:
			metricBuffer = append(metricBuffer, ms...)
			if len(metricBuffer) > 5000 {
				Push(metricBuffer[:])           // 仅复制slice header，引用同个底层数组
				metricBuffer = metricBuffer[:0] // 保留底层数组，减少内存分配
			}

		case <-ticker.C:
			Push(metricBuffer[:])
			metricBuffer = metricBuffer[:0] //清空
		}
	}
}

func Push(ms []Metric) error {
	_ = ms
	// 模拟网络IO
	time.Sleep(time.Second)
	return nil
}

var cond = sync.NewCond(&mu) // 条件变量

// 生产者
func producer() {
	for i := 0; ; i++ {
		mu.Lock()
		queue = append(queue, i)
		fmt.Println("生产任务:", i)
		cond.Signal() // 通知一个消费者
		mu.Unlock()
		time.Sleep(time.Second)
	}
}

// // 消费者
// func consumer(id int) {
//     for {
//         mu.Lock()
//         // 循环检查条件是否满足，避免虚假唤醒
//         for len(queue) == 0 {
//             cond.Wait() // 释放锁并进入阻塞，直到被 Signal/Broadcast 唤醒
//         }
//         // 走到这里时，锁已经被重新获取，且队列非空
//         task := queue[0]
//         queue = queue[1:]
//         fmt.Printf("消费者 %d 处理任务: %d\n", id, task)
//         mu.Unlock()
//     }
// }
