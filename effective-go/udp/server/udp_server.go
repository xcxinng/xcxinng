package main

import (
	"fmt"
	"net"
)

func main() {
	addr := net.UDPAddr{
		Port: 8080,
		IP:   net.ParseIP("127.0.0.1"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer conn.Close()

	// [background]： 消息大小为5字节， client主动发送2个消息
	//
	// [server]
	// case1: 如果用小于消息大小的缓冲区(e.g. 4B)读取UDP消息时，消息剩余的数据会被kernel清除
	// case2: 如果用大于消息大小的缓冲区(e.g. 8B)读取UDP消息时，一次最多只能读取一个消息大小的数据
	//
	// 所以UDP是自带消息边界的，kernel知道消息的起始和结束位置;
	// 而TCP的边界通常靠应用层决定，kernel不过多干预;

	// buffer := make([]byte, 4)
	buffer := make([]byte, 8)
	for i := 0; i < 2; i++ {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}
		fmt.Printf("Received %d bytes from %s: %v\n", n, addr.String(), string(buffer[:n]))
	}
}
