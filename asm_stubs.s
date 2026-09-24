; Stub functions for baremetal kernel
; Provides minimal implementations for runtime symbols that are not provided

bits 64

section .text

; Memory operations
global memset
memset:
    mov rax, rdi
    mov eax, esi
    mov rcx, rdx
    rep stosb
    mov rax, rdi
    ret

global memcpy
memcpy:
    mov rax, rdi
    mov rcx, rdx
    rep movsb
    ret

global memmove
memmove:
    mov rax, rdi
    cmp rdi, rsi
    jb .forward
    lea rcx, [rsi+rdx]
    cmp rdi, rcx
    jae .forward
    lea rdi, [rdi+rdx-1]
    lea rsi, [rsi+rdx-1]
    std
    rep movsb
    cld
    ret
.forward:
    mov rcx, rdx
    rep movsb
    ret

global memzero
memzero:
    mov rax, rdi
    xor eax, eax
    mov rcx, rdx
    rep stosb
    ret

global runtime.alloc
runtime.alloc:
    mov rax, [rel runtime_heap_current]
    test rax, rax
    jnz .has_base
    lea rax, [rel runtime_heap_start]
    mov [rel runtime_heap_current], rax
.has_base:
    mov rcx, rdi
    add rcx, 15
    and rcx, -16
    add rax, rcx
    cmp rdx, [rel runtime_heap_end]
    ja .alloc_fail
    mov [rel runtime_heap_current], rdx
    ret
.alloc_fail:
    xor eax, eax
    ret

section .bss
align 16
global runtime_heap_start
runtime_heap_start:
    resb 8*1024*1024
global runtime_heap_end
runtime_heap_end:
global runtime_heap_current
runtime_heap_current:
    resq 1

section .text

global runtime.realloc
runtime.realloc:
    test rdi, rdi
    jz runtime.alloc
    mov rax, rdi
    ret

global runtime.free
runtime.free:
    ret

; RDTSC stub
global asmReadRdtsc
asmReadRdtsc:
    rdtsc
    ret

; Pause stub
global asmPause
asmPause:
    pause
    ret

; UEFI call stubs (not used in baremetal)
global uefiCall1
uefiCall1:
    xor eax, eax
    ret

global uefiCall2
uefiCall2:
    xor eax, eax
    ret

global uefiCall3
uefiCall3:
    xor eax, eax
    ret

global uefiCall4
uefiCall4:
    xor eax, eax
    ret

global uefiCall5
uefiCall5:
    xor eax, eax
    ret

; Random stub
global getrandom
getrandom:
    xor eax, eax
    mov eax, 0          ; return 0 (no random data)
    ret

; Putchar stub (not used, VGA handles output)
global write
write:
    mov rax, rdx        ; return number of bytes "written"
    ret

; Signal handling stubs
global tinygo_register_fatal_signals
tinygo_register_fatal_signals:
    ret

global tinygo_caught_signal
tinygo_caught_signal:
    ret

global raise
raise:
    xor eax, eax
    ret

; Futex stubs
global tinygo_futex_wake_all
tinygo_futex_wake_all:
    ret

; Abort stub
global abort
abort:
    cli
    hlt
    jmp $