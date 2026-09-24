; Memory management assembly stubs
bits 64

section .text

global asmCr3
asmCr3:
    mov cr3, rdi
    ret

global asmEax
asmEax:
    mov eax, 0
    ret

global asmSetCr0
asmSetCr0:
    mov cr0, rdi
    ret