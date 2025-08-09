package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	requests := 50    // 总请求数
	concurrency := 15 // 并发数
	baseURL := "http://localhost:8080/api"

	sem := make(chan bool, concurrency) // 并发控制通道
	results := make(chan string, requests)

	start := time.Now()

	// 启动请求
	for i := 1; i <= requests; i++ {
		wg.Add(1)
		sem <- true

		go func(id int) {
			defer wg.Done()
			defer func() { <-sem }()

			startReq := time.Now()
			resp, err := http.Get(baseURL)
			duration := time.Since(startReq).Round(time.Millisecond)

			var result string
			if err != nil {
				result = fmt.Sprintf("[%02d] ❌ Error: %v", id, err)
			} else {
				status := "✅ OK"
				if resp.StatusCode == http.StatusTooManyRequests {
					status = "⛔ RATE LIMITED"
				}
				result = fmt.Sprintf("[%02d] %s - %s (Status: %d, Duration: %v)",
					id, startReq.Format("15:04:05.000"), status, resp.StatusCode, duration)
			}

			results <- result
		}(i)
	}

	// 结果收集协程
	go func() {
		wg.Wait()
		close(results)
	}()

	// 打印结果
	for res := range results {
		fmt.Println(res)
	}

	// 性能统计
	elapsed := time.Since(start)
	fmt.Printf("\n== 测试结果 ==\n")
	fmt.Printf("总请求数: %d\n", requests)
	fmt.Printf("并发数: %d\n", concurrency)
	fmt.Printf("总耗时: %s\n", elapsed.Round(time.Millisecond))
	fmt.Printf("平均RPS: %.1f\n", float64(requests)/elapsed.Seconds())
}
