package main

import (
	"fmt"
	"math/rand"
	"time"
)

/*
first version
*/
//func main() {
//	c := make(chan string)
//	go boring("util", c)
//	for i := 0; i < 5; i++ {
//		fmt.Printf("you say:%s\n", <-c)
//	}
//	fmt.Println("you are boring,i'm leaving")
//}

//func boring(msg string, ch chan string) {
//	for i := 0; ; i++ {
//		ch <- fmt.Sprintf("%s %d", msg, i)
//		rand.Seed(time.Now().Unix())
//		time.Sleep(time.Duration(rand.Intn(1e3)))
//	}
//}

/*
second version
*/
//func main() {
//	c := boring("Joe")
//	d := boring("Ann")
//	for i := 0; i < 5; i++ {
//		fmt.Println(<-c)
//		fmt.Println(<-d)
//	}
//	fmt.Println("you guys are boring,i'm leaving")
//}
//
//// channels are first-class values, just like strings or integers.
//func boring(msg string) chan string {
//	c := make(chan string)
//	go func() {
//		for i := 0; ; i++ {
//			time.Sleep(time.Millisecond * time.Duration(rand.Intn(2e3)))
//			c <- fmt.Sprintf("%s %d", msg, i)
//		}
//	}()
//	return c
//}

/*
multiplexing
*/
//func fanIn(input1, input2 <-chan string) <-chan string {
//	c := make(chan string)
//	go func() {
//		for {
//			c <- <-input1
//		}
//	}()
//	go func() {
//		for {
//			c <- <-input2
//		}
//	}()
//	return c
//}

//func main() {
//	c := fanIn(boring("Joe"), boring("Ann")) //boring func is the same as the previous pattern
//	for i := 0; i < 6; i++ {
//		fmt.Println(<-c)
//	}
//	fmt.Println("you're boring, I'm leaving.")
//}

/*
restore sequencing
*/

type Message struct {
	str  string
	wait chan bool // act like a signaler
}

// produce boring message, expose service by a channel
//func boring(msg string) chan Message {
//	c := make(chan Message)
//	waitForIt := make(chan bool)
//	go func() {
//		for i := 0; ; i++ {
//			message := Message{str: fmt.Sprintf("%s %d", msg, i), wait: waitForIt}
//			c <- message
//			time.Sleep(time.Duration(rand.Intn(2e3)) * time.Millisecond)
//			<-waitForIt
//		}
//	}()
//	return c
//}

// fanIn use to get multiple message in a sync way.
// 1.only get one message every time.
// 2.pass channel on a channel.
//func fanIn(input1, input2 <-chan Message) <-chan Message {
//	c := make(chan Message)
//	go func() {
//		for {
//			c <- <-input1
//		}
//	}()
//	go func() {
//		for {
//			c <- <-input2
//		}
//	}()
//	return c
//}
//
//func main() {
//	// Though blocking time is random in boring func, but actually Joe and Ann
//	// are both blocked until tell it is ok to go ahead. In which case, they
//	// work in a sequencing way.
//	c := fanIn(boring("Joe"), boring("Ann")) // 牛逼的地方在于boring函数对channel的使用
//	for i := 0; i < 5; i++ {
//		msg1 := <-c
//		fmt.Println(msg1.str)
//		msg2 := <-c
//		fmt.Println(msg2.str)
//		msg1.wait <- true
//		msg2.wait <- true
//	}
//}

/*
use select statement to change fanIn
*/
//func fanIn(input1, input2 <-chan Message) <-chan Message {
//	c := make(chan Message)
//	// no more 2 goroutines, only one is enough
//	go func() {
//		for {
//			select {
//			case m := <-input2:
//				c <- m
//			case m := <-input1:
//				c <- m
//			}
//		}
//	}()
//	return c
//}

/*
timeout using select
*/
//func main() {
//	c := boring("Joe")
//	timeout := time.After(time.Second * 5)
//	for i := 0; ; i++ {
//		select {
//		case m := <-c:
//			fmt.Println(m)
//		case <-timeout:
//			fmt.Println("timeout")
//			return
//		}
//	}
//}

/*
quit channel and with cleanup
*/
//func boring(msg string, quit chan string) chan string {
//	c := make(chan string)
//	go func() {
//		for i := 0; ; i++ {
//			select {
//			case m := <-quit:
//				fmt.Printf("got message:%q\n", m)
//				cleanup()
//				fmt.Println("ok, i quit")
//				quit <- "i'm done" // tell the caller it's done
//				return
//			case c <- fmt.Sprintf("%s %d", msg, i):
//				time.Sleep(time.Millisecond * time.Duration(rand.Intn(2e3)))
//			}
//		}
//	}()
//	return c
//}
//
//func cleanup() {
//	time.Sleep(time.Second)
//}
//
//func main() {
//	quit := make(chan string)
//	c := boring("Joe", quit)
//	for i := 3; i > 0; i-- {
//		fmt.Println(<-c)
//	}
//	quit <- "you can quit now"
//	fmt.Println("message from boring:", <-quit)
//}

/*
daisy-chain, a crazy extend version about "round-trip" we talk about above
*/

//func f(left, right chan int) {
//	left <- 1 + <-right
//}
//func main() {
//	const n = 10000
//	leftMost := make(chan int)
//	right := leftMost
//	left := leftMost
//	for i := 0; i < n; i++ {
//		right = make(chan int)
//		go f(left, right)
//		left = right
//	}
//	go func(c chan int) {
//		c <- 1
//	}(right)
//
//	fmt.Println(<-leftMost)
//}

/*
A fake framework
*/

type Result string

type Search func(kind string) Result

func fakeSearch(kind string) Search {
	return func(query string) Result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return Result(fmt.Sprintf("%s result for %q\n", kind, query))
	}
}

var (
	Web    = fakeSearch("web")
	Web1   = fakeSearch("web1")
	Web2   = fakeSearch("web2")
	Image  = fakeSearch("image")
	Image1 = fakeSearch("image1")
	Image2 = fakeSearch("image2")
	Video  = fakeSearch("video")
	Video1 = fakeSearch("video1")
	Video2 = fakeSearch("video2")
)

func Google(m string) []Result {
	/*
		search 2.0
	*/
	//var (
	//	results []Result
	//	chs     = make(chan Result)
	//)
	//go func() { chs <- Web(m) }()
	//go func() { chs <- Image(m) }()
	//go func() { chs <- Video(m) }()
	//for i := 0; i < 3; i++ {
	//	results = append(results, <-chs)
	//}
	//return results

	/*
		search 2.1: timeout restriction
	*/
	var (
		results []Result
		chs     = make(chan Result)
	)
	go func() { chs <- Web(m) }()
	go func() { chs <- Image(m) }()
	go func() { chs <- Video(m) }()

	timeout := time.After(time.Millisecond * 80)
	for i := 0; i < 3; i++ {
		select {
		case <-timeout:
			fmt.Println("search timeout")
		case res := <-chs:
			results = append(results, res)
		}
	}
	return results
}

//func main() {
//	rand.Seed(time.Now().Unix())
//	start := time.Now()
//	results := Google("golang")
//	elapsed := time.Since(start)
//	fmt.Println(results)
//	fmt.Println("time used::", elapsed)
//}

/*
replicate mode
*/

// First work in replicate mode, query from multiple replicates and return
// the fastest response.
func First(query string, replicas ...Search) Result {
	c := make(chan Result)
	searchReplica := func(i int) { c <- replicas[i](query) }
	for i := range replicas {
		go searchReplica(i)
	}
	return <-c
}

//
//func main() {
//	rand.Seed(time.Now().Unix())
//	start := time.Now()
//	results := First("golang", fakeSearch("replica 1"), fakeSearch("replica 2"))
//	elapsed := time.Since(start)
//	fmt.Println(results)
//	fmt.Println("time used::", elapsed)
//}

/*
final version which put all the magic tools together
*/

//func main() {
//	results := make(chan Result, 3)
//	responses := make([]Result, 0, 3)
//	rand.Seed(time.Now().UnixNano())
//	start := time.Now()
//
//	go func() { results <- First("golang", Web, Web1, Web2) }()
//	go func() { results <- First("golang", Image, Image1, Image2) }()
//	go func() { results <- First("golang", Video, Video1, Video2) }()
//	timeout := time.After(time.Millisecond * 80)
//	for i := 0; i < 3; i++ {
//		select {
//		case result := <-results:
//			responses = append(responses, result)
//		case <-timeout:
//			fmt.Println("Google search timeout")
//			return
//		}
//	}
//	fmt.Println(responses)
//	fmt.Println("time used on search:", time.Since(start))
//}
