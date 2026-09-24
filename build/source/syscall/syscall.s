; Syscall assembly stub
bits 64

section .text

; System call entry point
global syscall_entry
syscall_entry:
    ; Save registers
    push rdi
    push rsi
    push rdx
    push rcx
    push r8
    push r9
    push r10
    push r11
    
    ; Set up arguments for Go handler
    ; rdi = syscall number
    ; rsi = arg1
    ; rdx = arg2
    ; rcx = arg3
    ; r8 = arg4
    ; r9 = arg5
    
    mov rdi, rax  ; syscall number
    mov rsi, rdi  ; arg1
    mov rdx, rsi  ; arg2
    mov rcx, rdx  ; arg3
    mov r8, rcx   ; arg4
    mov r9, r8    ; arg5
    
    ; Call Go handler
    extern handleSyscall
    call handleSyscall
    
    ; Restore registers
    pop r11
    pop r10
    pop r9
    pop r8
    pop rcx
    pop rdx
    pop rsi
    pop rdi
    
    ; Return
    ret

; Assembly stub for calling syscall handlers
global asmSyscallHandler
asmSyscallHandler:
    ; Arguments: rdi=handler, rsi=num, rdx=a1, rcx=a2, r8=a3, r9=a4, stack=a5
    ; We need to call the handler with (a1, a2, a3, a4, a5)
    mov rax, rdi  ; handler
    mov rdi, rdx  ; a1
    mov rsi, rcx  ; a2
    mov rdx, r8   ; a3
    mov rcx, r9   ; a4
    mov r8, [rsp] ; a5 from stack
    
    call rax
    ret