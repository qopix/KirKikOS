package interrupts

import (
	"kirkikos/arch/amd64"
)

const (
	// Interrupt vectors
	DivideError = 0
	Timer = 32
	Keyboard = 33
	Com1 = 33
	Com2 = 35
	Spurious = 255

	// Syscall
	Syscall = 128
)

var (
	// Interrupt handlers
	handlers [256]uintptr
)

type InterruptHandler func(vector uint64)

func Init() {
	// Set up IDT
	setupIDT()
	
	// Remap PIC
	remapPIC()
	
	// Install default handlers
	installDefaultHandlers()
}

func setupIDT() {
	// Initialize IDT entries
	// This is a simplified version
}

func remapPIC() {
	// Remap master PIC to vectors 32-39
	// Remap slave PIC to vectors 40-47
	// Standard IO ports
	const (
		PIC1_COMMAND = 0x20
		PIC1_DATA    = 0x21
		PIC2_COMMAND = 0xA0
		PIC2_DATA    = 0xA1
	)

	// Save masks
	m1 := amd64.Inb(PIC1_DATA)
	m2 := amd64.Inb(PIC2_DATA)

	// Start initialization
	amd64.Outb(PIC1_COMMAND, 0x11)
	amd64.Outb(PIC2_COMMAND, 0x11)
	amd64.Outb(PIC1_DATA, 0x20) // Master offset 32
	amd64.Outb(PIC2_DATA, 0x28) // Slave offset 40
	amd64.Outb(PIC1_DATA, 0x04) // Cascade
	amd64.Outb(PIC2_DATA, 0x02)
	amd64.Outb(PIC1_DATA, 0x05)
	amd64.Outb(PIC2_DATA, 0x05)
	amd64.Outb(PIC1_DATA, m1) // Restore masks
	amd64.Outb(PIC2_DATA, m2)
}

func installDefaultHandlers() {
	// Install timer handler
	RegisterHandler(Timer, timerHandler)
	RegisterHandler(Keyboard, keyboardHandler)
}

func RegisterHandler(vector uint64, handler InterruptHandler) {
	if vector < 256 {
		handlers[vector] = uintptr(handler)
	}
}

func timerHandler(vector uint64) {
	// Timer tick - would trigger scheduler
	// Send EOI
	amd64.Outb(0x20, 0x20)
}

func keyboardHandler(vector uint64) {
	// Handle keyboard input
	key := amd64.Inb(0x60)
	// Process key
	amd64.Outb(0x20, 0x20) // EOI
}

func Enable() {
	amd64.EnableInterrupts()
}

func Disable() {
	amd64.DisableInterrupts()
}

func HandleInterrupt(vector uint64) {
	if vector < 256 && handlers[vector] != 0 {
		handler := handlers[vector]
		// Call handler
		asmCallHandler(handler, vector)
	}
}

func asmCallHandler(handler uintptr, vector uint64)