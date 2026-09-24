; KirkikOS amd64 Bootstrap Assembly
; Multiboot2 header + 32-bit entry that switches the CPU into long mode
; before handing off to the 64-bit Go kernel.
;
; IMPORTANT: GRUB always hands control to a Multiboot2 kernel with the CPU
; in 32-bit protected mode, paging disabled - regardless of whether the
; ELF itself is 64-bit. A kernel that starts with `bits 64` code at _start
; will be mis-decoded by the CPU and crash/triple-fault immediately.
; This stub does the standard PM32 -> long mode transition first.

section .multiboot_header
align 8
header_start:
    dd 0xe85250d6                ; magic
    dd 0                         ; architecture: 0 = i386 (required by the
                                  ; Multiboot2 spec even for a 64-bit kernel -
                                  ; there is no "amd64" value here)
    dd header_end - header_start ; header length
    dd 0x100000000 - (0xe85250d6 + 0 + (header_end - header_start)) ; checksum

    ; end tag
    dw 0
    dw 0
    dd 8
header_end:

section .bss
align 16
stack_bottom:
    resb 4096 * 16      ; 64 KiB stack
stack_top:

align 4096
p4_table:
    resb 4096
p3_table:
    resb 4096
p2_table:
    resb 4096

section .text
bits 32

global _start
extern main

_start:
    cli
    mov esp, stack_top

    ; EAX = multiboot2 magic, EBX = pointer to boot info (from GRUB)
    mov edi, ebx        ; keep boot info pointer around in a callee-saved reg

    call check_multiboot
    call check_cpuid
    call check_long_mode

    call set_up_page_tables
    call enable_paging

    lgdt [gdt64.pointer]
    jmp gdt64.code:long_mode_start

; --- sanity checks -----------------------------------------------------

check_multiboot:
    cmp eax, 0x36d76289
    jne .no_multiboot
    ret
.no_multiboot:
    mov al, "0"
    jmp error

check_cpuid:
    pushfd
    pop eax
    mov ecx, eax
    xor eax, 1 << 21
    push eax
    popfd
    pushfd
    pop eax
    push ecx
    popfd
    cmp eax, ecx
    je .no_cpuid
    ret
.no_cpuid:
    mov al, "1"
    jmp error

check_long_mode:
    mov eax, 0x80000000
    cpuid
    cmp eax, 0x80000001
    jb .no_long_mode

    mov eax, 0x80000001
    cpuid
    test edx, 1 << 29
    jz .no_long_mode
    ret
.no_long_mode:
    mov al, "2"
    jmp error

; --- long mode setup -----------------------------------------------------

set_up_page_tables:
    ; P4[0] -> P3 table
    mov eax, p3_table
    or eax, 0b11            ; present + writable
    mov [p4_table], eax

    ; P3[0] -> P2 table
    mov eax, p2_table
    or eax, 0b11
    mov [p3_table], eax

    ; identity-map the first 1 GiB with 2 MiB pages
    xor ecx, ecx
.map_p2_table:
    mov eax, 0x200000       ; 2 MiB
    mul ecx
    or eax, 0b10000011      ; present + writable + huge page
    mov [p2_table + ecx * 8], eax

    inc ecx
    cmp ecx, 512
    jne .map_p2_table
    ret

enable_paging:
    ; load P4 into CR3
    mov eax, p4_table
    mov cr3, eax

    ; enable PAE
    mov eax, cr4
    or eax, 1 << 5
    mov cr4, eax

    ; enable SSE (CR4.OSFXSR = bit 9, CR4.OSXMMEXCPT = bit 10)
    mov eax, cr4
    or eax, (1 << 9) | (1 << 10)
    mov cr4, eax

    ; set the long mode bit in EFER
    mov ecx, 0xC0000080
    rdmsr
    or eax, 1 << 8
    wrmsr

    ; enable paging
    mov eax, cr0
    or eax, 1 << 31
    mov cr0, eax
    ret

error:
    ; print "ERR: X" in white-on-red to the VGA text buffer, then hang
    mov dword [0xb8000], 0x4f524f45
    mov dword [0xb8004], 0x4f3a4f52
    mov dword [0xb8008], 0x4f204f20
    mov byte  [0xb800a], al
    hlt
    jmp $

section .rodata
align 8
gdt64:
    dq 0                                          ; null descriptor
.code: equ $ - gdt64
    dq (1<<43) | (1<<44) | (1<<47) | (1<<53)       ; executable, code, present, long-mode
.pointer:
    dw $ - gdt64 - 1
    dq gdt64

section .text
bits 64
long_mode_start:
    mov ax, 0
    mov ss, ax
    mov ds, ax
    mov es, ax
    mov fs, ax
    mov gs, ax

    ; Enable FPU/SSE in 64-bit mode
    ; Clear CR0.EM (bit 2), set CR0.MP (bit 1), clear CR0.TS (bit 3)
    mov rax, cr0
    and rax, ~(1 << 2)    ; clear EM
    or rax, (1 << 1)      ; set MP
    and rax, ~(1 << 3)    ; clear TS
    mov cr0, rax

    ; Enable SSE in CR4 (OSFXSR = bit 9, OSXMMEXCPT = bit 10)
    mov rax, cr4
    or rax, (1 << 9) | (1 << 10)
    mov cr4, rax

    ; Initialize FPU
    fninit

    ; rdi already holds the multiboot info pointer (was edi above);
    ; pass (magic-ish placeholder, boot info ptr) into Go's main()
    xor rsi, rsi
    mov rsi, rdi
    xor rdi, rdi
    call main
.hang:
    cli
    hlt
    jmp .hang
