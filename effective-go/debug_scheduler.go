package main

import (
	"fmt"
	"time"
)

// func main() {
// 	msgChan := make(chan Message, 3)
// 	// Create goroutines.
// 	for i := 0; i < 12; i++ {
// 		go process(msgChan)
// 	}

// 	msg := getMessages()
// 	for _, v := range msg {
// 		msgChan <- v
// 	}

// 	tm := time.NewTimer(time.Second)
// 	<-tm.C
// }

func getMessages() []Message {
	var ret []Message
	for i := 0; i < 100; i++ {
		ret = append(ret, Message{
			Id:      i,
			Message: fmt.Sprintf("message-%d", i),
		})
	}
	return ret
}

type Message struct {
	Id      int
	Message string
}

func process(m <-chan Message) {
	time.Sleep(time.Millisecond * 300)
	for {
		msg, ok := <-m
		if !ok { // channel closed
			return
		}
		// processing message
		fmt.Println(msg.Message)
	}
}
