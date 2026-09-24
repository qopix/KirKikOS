# KirKikOS Makefile
# Build system for KirKikOS operating system

GO      := /home/master/go/go/bin/go
TINYGO  := GOROOT=/home/master/go/go PATH=/home/master/go/go/bin:$$PATH tinygo
NASM    := nasm
OBJCOPY := objcopy
RM      := rm -rf
MKDIR   := mkdir -p

ASMFLAGS := -f elf64

SRC_DIR    := src
ROOTFS_DIR := rootfs
BUILD_DIR  := build
ISO_DIR    := iso/boot

KERNEL_BIN := kernel.bin
KERNEL_ELF := kernel.elf

BOOT_ASM  := $(SRC_DIR)/boot/boot.s
IO_ASM    := $(SRC_DIR)/arch/amd64/io.s
STUBS_ASM := asm_stubs.s

BOOT_OBJ  := $(BUILD_DIR)/boot.o
IO_OBJ    := $(BUILD_DIR)/io.o
STUBS_OBJ := $(BUILD_DIR)/asm_stubs.o

.PHONY: all clean iso kernel userspace help run source

all: kernel userspace iso

kernel: $(KERNEL_BIN)

$(KERNEL_BIN): $(BOOT_OBJ) $(IO_OBJ) $(STUBS_OBJ) linker.ld linux-amd64.json
	@echo "Building Go kernel with TinyGo..."
	cd $(SRC_DIR) && CGO_ENABLED=0 $(TINYGO) build \
		-target=../linux-amd64.json \
		-o ../$(BUILD_DIR)/kernel.o .
	@echo "Linking kernel ELF..."
	ld -T linker.ld -o $(KERNEL_ELF) $(BUILD_DIR)/boot.o $(BUILD_DIR)/io.o $(BUILD_DIR)/asm_stubs.o $(BUILD_DIR)/kernel.o --undefined=_start --no-gc-sections --no-undefined
	@echo "✓ Kernel built: $(KERNEL_ELF)"

$(BUILD_DIR)/boot.o: $(BOOT_ASM)
	@$(MKDIR) $(BUILD_DIR)
	@echo "Assembling boot.s..."
	$(NASM) $(ASMFLAGS) $< -o $@

$(BUILD_DIR)/io.o: $(IO_ASM)
	@$(MKDIR) $(BUILD_DIR)
	@echo "Assembling io.s..."
	$(NASM) $(ASMFLAGS) $< -o $@
	$(OBJCOPY) --redefine-sym outb=kirkikos/arch/amd64.outb \
		--redefine-sym inb=kirkikos/arch/amd64.inb \
		--redefine-sym outl=kirkikos/arch/amd64.outl \
		--redefine-sym inl=kirkikos/arch/amd64.inl \
		--redefine-sym asmHlt=kirkikos/arch/amd64.asmHlt \
		--redefine-sym asmCli=kirkikos/arch/amd64.asmCli \
		--redefine-sym asmSti=kirkikos/arch/amd64.asmSti $@
	$(OBJCOPY) --redefine-sym asmCpuid=kirkikos/arch/amd64.asmCpuid $@

$(BUILD_DIR)/asm_stubs.o: $(STUBS_ASM)
	@$(MKDIR) $(BUILD_DIR)
	@echo "Assembling asm_stubs.s..."
	$(NASM) $(ASMFLAGS) $< -o $@

userspace:
	@echo "Building user-space commands..."
	cd $(ROOTFS_DIR) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		$(GO) build -o ../$(BUILD_DIR)/hello bin/hello.go
	cd $(ROOTFS_DIR) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		$(GO) build -o ../$(BUILD_DIR)/ls bin/ls.go
	cd $(ROOTFS_DIR) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		$(GO) build -o ../$(BUILD_DIR)/cat bin/cat.go
	cd $(ROOTFS_DIR) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		$(GO) build -o ../$(BUILD_DIR)/ps bin/ps.go
	cd $(ROOTFS_DIR) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		$(GO) build -o ../$(BUILD_DIR)/init bin/init.go
	cd $(ROOTFS_DIR) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		$(GO) build -o ../$(BUILD_DIR)/mkdir bin/mkdir.go
	cd $(ROOTFS_DIR) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		$(GO) build -o ../$(BUILD_DIR)/touch bin/touch.go
	cd $(ROOTFS_DIR) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		$(GO) build -o ../$(BUILD_DIR)/rm bin/rm.go
	cd $(ROOTFS_DIR) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		$(GO) build -o ../$(BUILD_DIR)/cd bin/cd.go
	cd $(ROOTFS_DIR) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		$(GO) build -o ../$(BUILD_DIR)/pwd bin/pwd.go
	cd $(ROOTFS_DIR) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		$(GO) build -o ../$(BUILD_DIR)/date bin/date.go
	cd $(ROOTFS_DIR) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		$(GO) build -o ../$(BUILD_DIR)/fetch bin/fetch.go
	cd $(ROOTFS_DIR) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		$(GO) build -o ../$(BUILD_DIR)/useradd bin/useradd.go
	cd $(ROOTFS_DIR) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		$(GO) build -o ../$(BUILD_DIR)/userdel bin/userdel.go
	cd $(ROOTFS_DIR) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		$(GO) build -o ../$(BUILD_DIR)/userswitch bin/userswitch.go
	cd $(ROOTFS_DIR) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		$(GO) build -o ../$(BUILD_DIR)/users bin/users.go
	cd $(ROOTFS_DIR) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		$(GO) build -o ../$(BUILD_DIR)/whoami bin/whoami.go
	cd $(ROOTFS_DIR) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		$(GO) build -o ../$(BUILD_DIR)/reboot bin/reboot.go
	cd $(ROOTFS_DIR) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		$(GO) build -o ../$(BUILD_DIR)/halt bin/halt.go
	cd $(ROOTFS_DIR) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		$(GO) build -o ../$(BUILD_DIR)/clear bin/clear.go
	cd $(ROOTFS_DIR) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		$(GO) build -o ../$(BUILD_DIR)/version bin/version.go
	cd $(ROOTFS_DIR) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		$(GO) build -o ../$(BUILD_DIR)/help bin/help.go
	cd $(ROOTFS_DIR) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		$(GO) build -o ../$(BUILD_DIR)/echo bin/echo.go
	@echo "Installing commands..."
	cp $(BUILD_DIR)/hello $(ROOTFS_DIR)/bin/
	cp $(BUILD_DIR)/ls $(ROOTFS_DIR)/bin/
	cp $(BUILD_DIR)/cat $(ROOTFS_DIR)/bin/
	cp $(BUILD_DIR)/ps $(ROOTFS_DIR)/bin/
	cp $(BUILD_DIR)/init $(ROOTFS_DIR)/bin/
	cp $(BUILD_DIR)/mkdir $(ROOTFS_DIR)/bin/
	cp $(BUILD_DIR)/touch $(ROOTFS_DIR)/bin/
	cp $(BUILD_DIR)/rm $(ROOTFS_DIR)/bin/
	cp $(BUILD_DIR)/cd $(ROOTFS_DIR)/bin/
	cp $(BUILD_DIR)/pwd $(ROOTFS_DIR)/bin/
	cp $(BUILD_DIR)/date $(ROOTFS_DIR)/bin/
	cp $(BUILD_DIR)/fetch $(ROOTFS_DIR)/bin/
	cp $(BUILD_DIR)/useradd $(ROOTFS_DIR)/bin/
	cp $(BUILD_DIR)/userdel $(ROOTFS_DIR)/bin/
	cp $(BUILD_DIR)/userswitch $(ROOTFS_DIR)/bin/
	cp $(BUILD_DIR)/users $(ROOTFS_DIR)/bin/
	cp $(BUILD_DIR)/whoami $(ROOTFS_DIR)/bin/
	cp $(BUILD_DIR)/reboot $(ROOTFS_DIR)/bin/
	cp $(BUILD_DIR)/halt $(ROOTFS_DIR)/bin/
	cp $(BUILD_DIR)/clear $(ROOTFS_DIR)/bin/
	cp $(BUILD_DIR)/version $(ROOTFS_DIR)/bin/
	cp $(BUILD_DIR)/help $(ROOTFS_DIR)/bin/
	cp $(BUILD_DIR)/echo $(ROOTFS_DIR)/bin/
	@echo "✓ User-space commands built"

iso: $(KERNEL_ELF)
	@echo "Creating bootable ISO with embedded rootfs..."
	$(MKDIR) $(ISO_DIR)
	$(OBJCOPY) -O binary $(KERNEL_ELF) $(ISO_DIR)/$(KERNEL_BIN)
	@echo "Creating FAT32 filesystem from rootfs..."
	$(MKDIR) $(BUILD_DIR)/fat32
	dd if=/dev/zero of=$(BUILD_DIR)/fat32/rootfs.img bs=1M count=32 2>/dev/null
	mkfs.fat -F 32 $(BUILD_DIR)/fat32/rootfs.img 2>/dev/null
	mkdir -p $(BUILD_DIR)/fat32/mnt
	sudo mount -o loop $(BUILD_DIR)/fat32/rootfs.img $(BUILD_DIR)/fat32/mnt 2>/dev/null || mount -o loop $(BUILD_DIR)/fat32/rootfs.img $(BUILD_DIR)/fat32/mnt 2>/dev/null
	cp -r $(ROOTFS_DIR)/* $(BUILD_DIR)/fat32/mnt/
	sync
	sudo umount $(BUILD_DIR)/fat32/mnt 2>/dev/null || umount $(BUILD_DIR)/fat32/mnt 2>/dev/null
	cp $(BUILD_DIR)/fat32/rootfs.img $(ISO_DIR)/rootfs.img
	cp limine.conf $(ISO_DIR)/
	@echo "✓ ISO created in $(ISO_DIR)/ with embedded FAT32 rootfs"

clean:
	@echo "Cleaning build artifacts..."
	$(RM) $(BUILD_DIR)
	$(RM) $(ISO_DIR)
	$(RM) $(KERNEL_BIN) $(KERNEL_ELF)
	@echo "✓ Clean complete"

source:
	@echo "Creating source archive..."
	$(MKDIR) $(BUILD_DIR)/source
	cp -r $(SRC_DIR)/* $(BUILD_DIR)/source/
	cp -r $(ROOTFS_DIR)/* $(BUILD_DIR)/source/
	cp Makefile $(BUILD_DIR)/source/
	cp linker.ld $(BUILD_DIR)/source/
	cp linux-amd64.json $(BUILD_DIR)/source/
	cp baremetal.json $(BUILD_DIR)/source/ 2>/dev/null || true
	cp limine.conf $(BUILD_DIR)/source/ 2>/dev/null || true
	@echo "✓ Source copied to $(BUILD_DIR)/source/"

help:
	@echo "KirKikOS Build System"
	@echo "Available targets:"
	@echo "  all       - Build kernel, user-space, and ISO (default)"
	@echo "  kernel    - Build only the kernel"
	@echo "  userspace - Build only user-space commands"
	@echo "  iso       - Create bootable ISO"
	@echo "  clean     - Remove all build artifacts"
	@echo "  run       - Run kernel in QEMU"
	@echo "  source    - Create source directory with all sources"
	@echo "  help      - Show this help message"

run: kernel source
	@echo "Running KirKikOS in QEMU via GRUB..."
	@$(MKDIR) $(ISO_DIR)/boot/grub
	@cp $(KERNEL_ELF) $(ISO_DIR)/boot/kernel.elf
	@echo "set timeout=0" > $(ISO_DIR)/boot/grub/grub.cfg
	@echo "set default=0" >> $(ISO_DIR)/boot/grub/grub.cfg
	@echo "menuentry \"KirKikOS\" {" >> $(ISO_DIR)/boot/grub/grub.cfg
	@echo "multiboot2 /boot/kernel.elf" >> $(ISO_DIR)/boot/grub/grub.cfg
	@echo "boot" >> $(ISO_DIR)/boot/grub/grub.cfg
	@echo "}" >> $(ISO_DIR)/boot/grub/grub.cfg
	@grub-mkrescue -o kirKikOS.iso $(ISO_DIR) 2>&1 | tail -5
	@qemu-system-x86_64 -cdrom kirKikOS.iso -m 256M 2>&1

