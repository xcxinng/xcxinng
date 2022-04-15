package main

import (
	"awesomeGolang/util"
	"fmt"
)

func main() {
	// study golang linkname
	fmt.Println(util.CallByteToString([]byte("hello world!")))
}
