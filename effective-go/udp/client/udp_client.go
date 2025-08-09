package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	// 定义服务器的地址和端口
	serverAddr := net.UDPAddr{
		Port: 8080,
		IP:   net.ParseIP("127.0.0.1"),
	}

	conn, err := net.DialUDP("udp", nil, &serverAddr)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	// 要发送的两个5字节的消息
	messages := []string{"Hello", "World"}

	// 发送消息
	for _, message := range messages {
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
		fmt.Printf("Sent message: %s\n", message)

		// 等待一段时间，确保消息被单独发送
		time.Sleep(1 * time.Second)
	}
}
