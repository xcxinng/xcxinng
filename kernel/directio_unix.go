//go:build linux
// +build linux

package kernel

import (
	"os"
	"syscall"
)

// OpenFile is a modified version of os.OpenFile which sets O_DIRECT
func OpenFile(name string, flag int, perm os.FileMode) (file *os.File, err error) {
	// the key to the so-called direct IO is the O_DIRECT flag.
	return os.OpenFile(name, syscall.O_DIRECT|flag, perm)
}
