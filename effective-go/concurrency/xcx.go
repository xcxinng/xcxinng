package main

import (
	"fmt"
	"math/rand"
	"time"
)

type ExtendMessage struct {
	str       string    // message
	waitForIt chan bool // sync with other goroutine
	done      chan bool // whether to quit
	feedback  chan bool // tell other goroutine the current goroutine quit successfully
}

func eBoring(msg string) chan ExtendMessage {
	c := make(chan ExtendMessage)
	d := make(chan bool)
	f := make(chan bool)
	w := make(chan bool)
	go func() {
		for i := 0; ; i++ {
			select {
			case <-d:
				fmt.Println("eBoring is done,bye")
				f <- true
				return
			default:
				time.Sleep(time.Duration(rand.Intn(2e3)) * time.Millisecond)
				c <- ExtendMessage{done: d, waitForIt: w, str: fmt.Sprintf("%s %d", msg, i), feedback: f}
				<-w
			}
		}
	}()
	return c
}

// func main() {
// 	msg := eBoring("Joe")
// 	for i := 0; i < 5; i++ {
// 		m := <-msg
// 		fmt.Println(m.str)
// 		m.waitForIt <- true

// 		if i == 4 {
// 			m.done <- true // tell eBoring to quit
// 			<-m.feedback   // wait for its quit feedback
// 			fmt.Println("ok, i got your feedback,see ya")
// 		}
// 	}
// }
