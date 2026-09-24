package syscall

import (
	"kirkikos/arch/amd64"
	"kirkikos/fs"
	"kirkikos/process"
)

const (
	// System call numbers
	SYS_EXIT = 0
	SYS_FORK = 1
	SYS_READ = 2
	SYS_WRITE = 3
	SYS_OPEN = 4
	SYS_CLOSE = 5
	SYS_EXECVE = 6
	SYS_GETPID = 7
	SYS_GETPPID = 8
	SYS_BRK = 9
	SYS_MMAP = 10
	SYS_MUNMAP = 11
	SYS_CHDIR = 12
	SYS_GETCWD = 13
	SYS_STAT = 14
	SYS_FSTAT = 15
	SYS_LSEEK = 16
	SYS_MKDIR = 17
	SYS_RMDIR = 18
	SYS_UNLINK = 19
	SYS_LINK = 20
	SYS_SYMLINK = 21
	SYS_READLINK = 22
	SYS_CHMOD = 23
	SYS_CHOWN = 24
	SYS_KILL = 25
	SYS_SIGNAL = 26
	SYS_WAIT = 27
	SYS_TIME = 28
	SYS_GETTIMEOFDAY = 29
	SYS_SOCKET = 30
	SYS_BIND = 31
	SYS_LISTEN = 32
	SYS_ACCEPT = 33
	SYS_CONNECT = 34
	SYS_SENDTO = 35
	SYS_RECVFROM = 36
	SYS_SHUTDOWN = 37
)

var (
	syscallHandlers [256]uintptr
)

func Init() {
	// Register syscall handlers
	RegisterHandler(SYS_EXIT, exitHandler)
	RegisterHandler(SYS_READ, readHandler)
	RegisterHandler(SYS_WRITE, writeHandler)
	RegisterHandler(SYS_OPEN, openHandler)
	RegisterHandler(SYS_CLOSE, closeHandler)
	RegisterHandler(SYS_GETPID, getpidHandler)
	RegisterHandler(SYS_CHDIR, chdirHandler)
	RegisterHandler(SYS_GETCWD, getCWDHandler)
	RegisterHandler(SYS_MKDIR, mkdirHandler)
	RegisterHandler(SYS_RMDIR, rmdirHandler)
	RegisterHandler(SYS_UNLINK, unlinkHandler)
	RegisterHandler(SYS_TIME, timeHandler)
	RegisterHandler(SYS_SOCKET, socketHandler)
	RegisterHandler(SYS_CONNECT, connectHandler)
	RegisterHandler(SYS_SENDTO, sendtoHandler)
	RegisterHandler(SYS_RECVFROM, recvfromHandler)
}

func RegisterHandler(num uint64, handler uintptr) {
	if num < 256 {
		syscallHandlers[num] = handler
	}
}

func HandleSyscall(num uint64, arg1 uint64, arg2 uint64, arg3 uint64, arg4 uint64, arg5 uint64) uint64 {
	if num < 256 && syscallHandlers[num] != 0 {
		handler := syscallHandlers[num]
		return asmSyscallHandler(handler, num, arg1, arg2, arg3, arg4, arg5)
	}
	return ^uint64(0) // -EINVAL
}

// Syscall handlers
func exitHandler(arg1, arg2, arg3, arg4, arg5 uint64) uint64 {
	process.Exit(uint32(arg1))
	return 0
}

func readHandler(arg1, arg2, arg3, arg4, arg5 uint64) uint64 {
	// fd, buf, count
	return 0
}

func writeHandler(arg1, arg2, arg3, arg4, arg5 uint64) uint64 {
	// fd, buf, count
	fd := arg1
	buf := arg2
	count := arg3

	if fd == 1 || fd == 2 {
		// Write to console
		data := (*[1 << 28]byte)(unsafe.Pointer(uintptr(buf)))[:count]
		for _, b := range data {
			amd64.PutChar(b)
		}
		return count
	}
	return ^uint64(0) // -EBADF
}

func openHandler(arg1, arg2, arg3, arg4, arg5 uint64) uint64 {
	// path, flags, mode
	return 0
}

func closeHandler(arg1, arg2, arg3, arg4, arg5 uint64) uint64 {
	// fd
	return 0
}

func getpidHandler(arg1, arg2, arg3, arg4, arg5 uint64) uint64 {
	proc := process.GetCurrentProcess()
	if proc != nil {
		return uint64(proc.ID)
	}
	return 0
}

func chdirHandler(arg1, arg2, arg3, arg4, arg5 uint64) uint64 {
	// path
	proc := process.GetCurrentProcess()
	if proc != nil {
		path := (*string)(unsafe.Pointer(uintptr(arg1)))
		copy(proc.CWD[:], *path)
		return 0
	}
	return ^uint64(0)
}

func getCWDHandler(arg1, arg2, arg3, arg4, arg5 uint64) uint64 {
	// buf, size
	proc := process.GetCurrentProcess()
	if proc != nil {
		buf := (*[1 << 28]byte)(unsafe.Pointer(uintptr(arg1)))
		size := arg2
		cwd := string(proc.CWD[:])
		copy(buf[:size], cwd)
		return uint64(len(cwd))
	}
	return ^uint64(0)
}

func mkdirHandler(arg1, arg2, arg3, arg4, arg5 uint64) uint64 {
	// path, mode
	path := (*string)(unsafe.Pointer(uintptr(arg1)))
	return boolToInt(fs.GetRootFS().Mkdir(*path))
}

func rmdirHandler(arg1, arg2, arg3, arg4, arg5 uint64) uint64 {
	// path
	path := (*string)(unsafe.Pointer(uintptr(arg1)))
	return boolToInt(fs.GetRootFS().Rmdir(*path))
}

func unlinkHandler(arg1, arg2, arg3, arg4, arg5 uint64) uint64 {
	// path
	path := (*string)(unsafe.Pointer(uintptr(arg1)))
	return boolToInt(fs.GetRootFS().DeleteFile(*path))
}

func timeHandler(arg1, arg2, arg3, arg4, arg5 uint64) uint64 {
	// tloc
	return 0
}

func socketHandler(arg1, arg2, arg3, arg4, arg5 uint64) uint64 {
	// domain, type, protocol
	return 0
}

func connectHandler(arg1, arg2, arg3, arg4, arg5 uint64) uint64 {
	// sockfd, addr, addrlen
	return 0
}

func sendtoHandler(arg1, arg2, arg3, arg4, arg5 uint64) uint64 {
	// sockfd, buf, len, flags, addr
	return 0
}

func recvfromHandler(arg1, arg2, arg3, arg4, arg5 uint64) uint64 {
	// sockfd, buf, len, flags, addr
	return 0
}

func boolToInt(b bool) uint64 {
	if b {
		return 0
	}
	return ^uint64(0)
}

func asmSyscallHandler(handler uintptr, num, a1, a2, a3, a4, a5 uint64) uint64

import "unsafe"