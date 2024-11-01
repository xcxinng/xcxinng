package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"strconv"
	"strings"
)

/*
[Why do i write this demo]

The aim of this demonstration is to confirm the integrity of the HTTP message
format.

Given that HTTP is predicated on TCP, a protocol that provides a stream-oriented
transport layer, it is essential to have a standardized header format.
This standardization is crucial for both HTTP requests and responses, as it
includes a well-defined header that informs the recipient—whether a service or
client—how to properly parse the accompanying message body, if it exists.

As per RFC 7230, an HTTP request begins with a status line, which is terminated
with a carriage return and line feed sequence, '\r\n'. Subsequent to the status
line, if included, are the request headers presented as multiple lines.

Each header line is delineated by the same delimiter, '\r\n'. The sequence signifies
the end of the headers when a blank line, indicated by two consecutive '\r\n',
is encountered, implying the presence of a message body following this delimiter.

The presence of a Content-Length header in either an HTTP request or response
typically indicates that a message body follows. This header specifies the size
of the body in bytes, allowing the recipient to know how many bytes to expect and
where the body ends.
*/

// 从 TCP 连接中读取并解析 HTTP 请求或响应
func printHTTPMessage(conn net.Conn) {
	defer conn.Close()

	reader := textproto.NewReader(bufio.NewReader(conn))
	statusLine, err := reader.ReadLine()
	if err != nil {
		panic(err)
	}
	fmt.Println(statusLine)

	header, err := reader.ReadMIMEHeader()
	if err != nil {
		panic(err)
	}
	fmt.Println("mime header: ")
	fmt.Println(header)

	length, exist := header["Content-Length"]
	if exist {
		count, err := strconv.ParseInt(length[0], 10, 64)
		if err != nil {
			panic(err)
		}

		bodyBuf := make([]byte, count)
		readCount, err := reader.R.Read(bodyBuf)
		if err != nil {
			panic(err)
		}
		fmt.Println("read count: ", readCount)
		fmt.Println("request body: ")
		fmt.Println(string(bodyBuf))
	}
}

// 解析 HTTP 头部字段
func parseHeaders(data []byte) map[string]string {
	headers := make(map[string]string)
	i := 0

	for {
		// 查找头部结束位置（两个连续的 \r\n）
		end := i
		for end < len(data)-2 {
			if data[end] == '\r' && data[end+1] == '\n' && data[end+2] == '\r' && data[end+3] == '\n' {
				break
			}
			end++
		}

		if end >= len(data)-2 {
			break
		}

		line := data[i:end]
		i = end + 4

		colonIndex := bytes.IndexByte(line, ':')
		if colonIndex != -1 {
			key := strings.TrimSpace(string(line[:colonIndex]))
			value := strings.TrimSpace(string(line[colonIndex+1:]))
			headers[key] = value
		}
	}

	return headers
}

// 检查数据是否是 HTTP 头部
func isHTTPHeader(data []byte) bool {
	// HTTP 请求或响应的起始行必须以 "HTTP/" 开头
	return strings.HasPrefix(string(data), "HTTP/") || bytes.HasPrefix(data, []byte("GET "))
}

// 读取一个http请求 然后将协议内容print出来到标准输出
func runHttpPrinter() {
	// 监听本地 8080 端口
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("监听端口失败:", err)
	}
	defer listener.Close()

	log.Println("服务器正在监听 8080 端口...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("接受连接失败:", err)
			continue
		}

		go printHTTPMessage(conn)
	}
}

func main() {
	runHttpPrinter()
}
