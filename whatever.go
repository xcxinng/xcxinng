package main

import (
	"fmt"
	"time"
	// "os/exec"
)

func main() {
	// cmd := exec.Command("set_ipfix.py")
	// result, err := cmd.CombinedOutput()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(string(result))

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("recover with error: ", err)
			}
		}()
		apiHandler()
	}()
	time.Sleep(time.Second * 3)
}

func apiHandler() {
	go doTask()
	fmt.Println(100)
}

func doTask() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover with error: ", err)
		}
	}()
	time.Sleep(time.Second * 2)
	panic("task failed")
}

// func useSyscall()  {
// 	file,err := os.Open("test.txt")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer file.Close()

// 	offset, err = file.Seek(0,2)
// 	if err!=nil {
// 		panic(err)
// 	}
// 	buf := bufio.NewScanner(file)
// 	buf.Scan()
// }
