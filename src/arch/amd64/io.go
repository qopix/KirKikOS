package amd64

import (
	"unsafe"
)

// Port I/O
func Outb(port uint16, value byte) {
	outb(port, value)
}

func Inb(port uint16) byte {
	return inb(port)
}

func Outl(port uint16, value uint32) {
	outl(port, value)
}

func Inl(port uint16) uint32 {
	return inl(port)
}

// Assembly stubs (defined in io.s)
func outb(port uint16, value byte)
func inb(port uint16) byte
func outl(port uint16, value uint32)
func inl(port uint16) uint32
func asmHlt()
func asmCli()
func asmSti()

// Serial port I/O (COM1)
const COM1 = 0x3F8

func SerialInit() {
	Outb(COM1+1, 0x00)
	Outb(COM1+3, 0x80)
	Outb(COM1+0, 0x01)
	Outb(COM1+1, 0x00)
	Outb(COM1+3, 0x03)
	Outb(COM1+2, 0xC7)
	Outb(COM1+4, 0x0B)
}

// VGA text mode
const (
	vgaMemory   = 0xB8000
	vgaWidth    = 80
	vgaHeight   = 25
)

var (
	vgaX, vgaY int
)

func PutChar(c byte) {
	if c == '\n' {
		PutChar('\r')
	}
	vgaPutChar(c)
}

func vgaPutChar(c byte) {
	if c == '\n' {
		vgaY++
		vgaX = 0
		return
	}
	if c == '\r' {
		vgaX = 0
		return
	}

	if vgaY >= vgaHeight {
		scrollVGA()
		vgaY = vgaHeight - 1
	}

	offset := (vgaY*vgaWidth + vgaX) * 2
	*(*uint16)(unsafe.Pointer(uintptr(vgaMemory + offset))) = uint16(c) | 0x0700
	vgaX++

	if vgaX >= vgaWidth {
		vgaX = 0
		vgaY++
	}

	updateCursor()
}

func scrollVGA() {
	for y := 1; y < vgaHeight; y++ {
		for x := 0; x < vgaWidth; x++ {
			srcOffset := (y*vgaWidth + x) * 2
			dstOffset := ((y-1)*vgaWidth + x) * 2
			*(*uint16)(unsafe.Pointer(uintptr(vgaMemory + dstOffset))) = 
				*(*uint16)(unsafe.Pointer(uintptr(vgaMemory + srcOffset)))
		}
	}

	for x := 0; x < vgaWidth; x++ {
		offset := ((vgaHeight-1)*vgaWidth + x) * 2
		*(*uint16)(unsafe.Pointer(uintptr(vgaMemory + offset))) = 0x0700
	}
}

func updateCursor() {
	cursor := uint16(vgaY*vgaWidth + vgaX)
	Outb(0x3D4, 14)
	Outb(0x3D5, byte(cursor >> 8))
	Outb(0x3D4, 15)
	Outb(0x3D5, byte(cursor & 0xFF))
}

func ClearScreen() {
	for i := 0; i < vgaWidth*vgaHeight; i++ {
		*(*uint16)(unsafe.Pointer(uintptr(vgaMemory+i*2))) = 0x0720
	}
	vgaX = 0
	vgaY = 0
	updateCursor()
}

func Reboot() {
	Outb(0x64, 0xFE)
	for {
		asmHlt()
	}
}

func Halt() {
	for {
		asmHlt()
	}
}

func DisableInterrupts() {
	asmCli()
}

func EnableInterrupts() {
	asmSti()
}

type CPUIDRegs struct {
	EAX uint32
	EBX uint32
	ECX uint32
	EDX uint32
}

func Cpuid(eax uint32) uint32 {
	var regs CPUIDRegs
	asmCpuid(eax, &regs)
	return regs.EAX
}

func CpuidRegs(eax uint32) CPUIDRegs {
	var regs CPUIDRegs
	asmCpuid(eax, &regs)
	return regs
}

func asmCpuid(eax uint32, regs *CPUIDRegs)