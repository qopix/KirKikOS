; KirkikOS amd64 I/O Assembly Stubs
bits 64

section .text

; Port I/O functions
global outb
outb:
    mov dx, di          ; port parameter
    mov al, sil         ; value parameter
    out dx, al
    ret

global inb
inb:
    mov dx, di          ; port parameter
    in al, dx
    ret

global outl
outl:
    mov dx, di          ; port parameter
    mov eax, esi        ; value parameter
    out dx, eax
    ret

global inl
inl:
    mov dx, di          ; port parameter
    in eax, dx
    ret

; CPU control
global asmHlt
asmHlt:
    hlt
    ret

global asmCli
asmCli:
    cli
    ret

global asmSti
asmSti:
    sti
    ret

; CPUID
global asmCpuid
asmCpuid:
    push rbx
    cpuid
    mov [rsi], eax
    mov [rsi+4], ebx
    mov [rsi+8], ecx
    mov [rsi+12], edx
    pop rbx
    ret