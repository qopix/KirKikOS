package memory

import (
	"unsafe"
)

const (
	// Page sizes
	PageSize = 4096

	// Memory regions
	KernelBase    = 0xFFFFFFFF80000000
	KernelHeap    = 0xFFFFFFFFC0000000
	UserBase      = 0x0000000040000000
	UserHeap      = 0x0000000080000000
	StackBase     = 0x000000007FFF0000

	// Page table flags
	Present = 1 << 0
	Writable = 1 << 1
	User = 1 << 2
	WT = 1 << 3
	CacheDisable = 1 << 4
	Accessed = 1 << 5
	Dirty = 1 << 6
	LargePage = 1 << 7
	Global = 1 << 8

	// Kernel page directory
	KernelPD = 0xFFFFFFFFC0000000
)

var (
	totalMemory uint64
	availableMemory uint64
	heapStart uint64
	heapEnd uint64
	heapCurrent uint64
)

func Init(bootInfo uint64) {
	heapStart = KernelHeap
	heapEnd = KernelHeap + 64*1024*1024
	heapCurrent = heapStart

	setupPageTables()
	enablePaging()
}

func setupPageTables() {
	pd := (*[1024]uint64)(unsafe.Pointer(uintptr(KernelPD)))
	for i := range pd {
		pd[i] = 0
	}

	pd[0] = 0x00000000 | LargePage | Present | Writable

	for i := 256; i < 512; i++ {
		page := uint64(i-256) << 22
		pd[i] = page | LargePage | Present | Writable
	}
}

func enablePaging() {
	asmCr3(uintptr(KernelPD))
	cr0 := asmEax()
	cr0 |= 0x80000000
	asmSetCr0(cr0)
}

func AllocPage() uint64 {
	if heapCurrent >= heapEnd {
		return 0
	}
	page := heapCurrent
	heapCurrent += PageSize
	return page
}

func FreePage(addr uint64) {
}

func GetMemoryInfo() (total, available uint64) {
	return totalMemory, availableMemory
}

func asmCr3(pd uintptr)
func asmEax() uint32
func asmSetCr0(value uint32)