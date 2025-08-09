package main

import (
	"fmt"
	"io"
	"os"
)

var Dollar = 0

func main() {
	// create myfile
	myfile, err := os.Create("myfile.txt")
	if err != nil {
		panic(err)
	}

	// write to myfile
	escritor := io.Writer(myfile)
	nb, err := escritor.Write([]byte("Ola "))
	if err != nil {
		panic(err)
	}
	fmt.Println("bytes written: ", nb)
	defer myfile.Close()

	offset, err := myfile.Seek(0, io.SeekCurrent)
	if err != nil {
		panic(err)
	}
	fmt.Println("offset after the write: ", offset)

	newOffset, err := myfile.Seek(0, 0)
	if err != nil {
		panic(err)
	}
	fmt.Println("new offset: ", newOffset)

	// read file
	buffer := make([]byte, 1024)
	rb, err := myfile.Read(buffer)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buffer[:rb]))

	// io.Writer(mf).Write()

	// go func() {
	// 	fmt.Println("dollar = ", Dollar)
	// }()
	// go func() {
	// 	Dollar = 1
	// }()

	// time.Sleep(time.Second * 2)
}
