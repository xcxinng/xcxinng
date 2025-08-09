package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

func main() {
	// 创建Redis客户端
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	defer rdb.Close()

	// 添加OpenTelemetry监控
	if err := redisotel.InstrumentTracing(rdb); err != nil {
		panic(err)
	}
	if err := redisotel.InstrumentMetrics(rdb); err != nil {
		panic(err)
	}

	// 初始化限流器
	limiter := redis_rate.NewLimiter(rdb)

	// 创建Gin引擎
	r := gin.Default()

	// 添加限流中间件
	r.Use(func(c *gin.Context) {
		// 使用客户端IP作为限流键
		ip := c.ClientIP()
		key := "rate_limit:" + ip

		// 每秒最多处理10个请求
		res, err := limiter.Allow(context.Background(), key, redis_rate.PerSecond(10))
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "Rate limiter error"})
			return
		}

		// 设置限流信息到响应头
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", res.Limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", res.Remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", res.ResetAfter/time.Millisecond))

		// 如果达到限制则返回429
		if res.Allowed == 0 {
			c.AbortWithStatusJSON(429, gin.H{
				"status":  "error",
				"message": "Too many requests",
				"detail": fmt.Sprintf("Rate limit exceeded. %d requests allowed per second. Retry after %v",
					res.Limit, res.RetryAfter.Round(time.Millisecond)),
			})
			return
		}

		c.Next()
	})

	// 示例API端点
	r.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Request processed successfully",
			"data": gin.H{
				"timestamp": time.Now().UnixMilli(),
				"client_ip": c.ClientIP(),
			},
		})
	})

	// 启动服务器
	fmt.Println("Server running on :8080")
	r.Run(":8080")
}
