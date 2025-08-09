package main

/*
#include <fcntl.h>
#include <sys/mman.h>
#include <unistd.h>
#include <stdlib.h>
#include <errno.h>

// 封装 shm_open 并返回 errno
int c_shm_open(const char *name, int flags, mode_t mode, int *err) {
    int fd = shm_open(name, flags, mode);
    if (fd == -1) {
        *err = errno;
    }
    return fd;
}

// 封装 shm_unlink 并返回 errno
int c_shm_unlink(const char *name, int *err) {
    int ret = shm_unlink(name);
    if (ret == -1) {
        *err = errno;
    }
    return ret;
}

// 辅助函数：获取当前线程的 errno
int get_errno() {
    return errno;
}
*/
import "C"

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	SIZE     = 4096
	SHM_NAME = "/my_mac_shm"
)

func main() {
	// 转换共享内存名称到 C 字符串
	cName := C.CString(SHM_NAME)
	defer C.free(unsafe.Pointer(cName))

	// 1. 打开共享内存对象（只读模式）
	var cErr C.int
	shmFd := C.c_shm_open(cName, C.O_RDONLY, 0666, &cErr)
	if shmFd == -1 {
		fmt.Printf("打开共享内存失败: %v\n", syscall.Errno(cErr))
		return
	}
	defer C.close(shmFd)

	// 2. 映射共享内存
	fd := int(shmFd)
	ptr, err := syscall.Mmap(
		fd,
		0,
		SIZE,
		syscall.PROT_READ,
		syscall.MAP_SHARED,
	)
	if err != nil {
		fmt.Printf("内存映射失败: %v\n", err)
		return
	}
	defer syscall.Munmap(ptr)

	// 3. 读取数据（处理可能的空终止符）
	var message []byte
	for _, b := range ptr {
		if b == 0 {
			break
		}
		message = append(message, b)
	}
	fmt.Printf("Go 读取到的数据: %s\n", string(message))

	// 4. 清理共享内存（通常由写入方操作）
	var unlinkErr C.int
	if ret := C.c_shm_unlink(cName, &unlinkErr); ret == -1 {
		fmt.Printf("删除共享内存失败: %v\n", syscall.Errno(unlinkErr))
	}
}
