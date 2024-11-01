package kernel

import (
	"io"
	"log"
	"os"
	"syscall"
)

func DirectIO() {
	// open file to get the file descriptor
	filename := "yourfile.txt"
	// ensure all io operations are performed from the end of the file,
	// and append only
	file, err := os.OpenFile(filename, os.O_CREATE|io.SeekEnd, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
		return
	}

	// [fcntl syscall]
	//
	// #include <fcntl.h>
	// int fcntl(int fd, int op, ... /* arg */ );
	//
	_, _, e1 := syscall.Syscall(syscall.SYS_FCNTL, file.Fd(), syscall.F_NOCACHE, 1)
	if e1 != 0 {
		file.Close()
		file = nil
		log.Fatalf("Failed to set F_NOCACHE: %s", e1)
	}
}
