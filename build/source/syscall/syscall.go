package syscall

import (
	"unsafe"
)

// System call numbers for KirkikOS
const (
	SYS_EXIT = 0
	SYS_READ = 1
	SYS_WRITE = 2
	SYS_OPEN = 3
	SYS_CLOSE = 4
	SYS_GETPID = 5
	SYS_GETCWD = 6
	SYS_CHDIR = 7
)

// Write to file descriptor
func Write(fd int, data []byte) int {
	if len(data) == 0 {
		return 0
	}
	// Simple syscall wrapper
	result := rawSyscall(SYS_WRITE, uint64(fd), uint64(unsafe.Pointer(&data[0])), uint64(len(data)))
	return int(result)
}

// Exit process
func Exit(code int) {
	rawSyscall(SYS_EXIT, uint64(code), 0, 0)
	// Should not return
	for {
		asmHlt()
	}
}

// Assembly stubs
func rawSyscall(num, a1, a2, a3 uint64) uint64
func asmHlt()