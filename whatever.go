package main

import (
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command("set_ipfix.py")
	result, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(result))
}
