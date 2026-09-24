; KirkikOS amd64 Assembly Stubs
bits 64

section .text

; HLT instruction
global asmHlt
asmHlt:
    hlt
    ret

; Port I/O
global outb
outb:
    mov dx, di
    mov al, sil
    out dx, al
    ret

global inb
inb:
    mov dx, di
    in al, dx
    ret

global outl
outl:
    mov dx, di
    mov eax, esi
    out dx, eax
    ret

global inl
inl:
    mov dx, di
    in eax, dx
    ret