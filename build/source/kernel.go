package main

import (
	"kirkikos/arch/amd64"
	"kirkikos/fs"
	"kirkikos/process"
	"kirkikos/shell"
)

func kmain(magic uint64, bootInfo uint64) {
	amd64.ClearScreen()
	printLine("KirKikOS 0.1")
	printLine("Bootloader: GRUB (Multiboot2)")
	printLine("Architecture: x86_64 (amd64)")
	printLine("Kernel: TinyGo")
	printLine("")
	printLine("Type 'help' for the command list.")
	printLine("")
	fs.Init()
	process.Init()
	shell.Run()
}

func main() {
	kmain(0, 0)
}

func printLine(value string) {
	for i := 0; i < len(value); i++ {
		amd64.PutChar(value[i])
	}
	amd64.PutChar('\n')
}
