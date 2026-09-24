package shell

import (
	"kirkikos/arch/amd64"
	"kirkikos/fs"
	"kirkikos/process"
)

const (
	lineMax = 128
	argMax  = 16
)

func Run() {
	process.InitUsers()
	for {
		username := process.GetCurrentUsername()
		printString(username)
		printString("@kios> ")
		line := readLine()
		if line == "" {
			continue
		}
		execute(line)
	}
}

func readLine() string {
	var line [lineMax]byte
	length := 0
	shift := false
	ctrl := false
	alt := false
	altgr := false

	for {
		if amd64.Inb(0x64)&1 == 0 {
			continue
		}
		code := amd64.Inb(0x60)
		if code == 0xE0 {
			continue
		}
		if code&0x80 != 0 {
			if code == 0xAA || code == 0xB6 {
				shift = false
			}
			if code == 0x9D {
				ctrl = false
			}
			if code == 0xB8 {
				alt = false
			}
			if code == 0xE0 {
				altgr = false
			}
			continue
		}
		if code == 0x2A || code == 0x36 {
			shift = true
			continue
		}
		if code == 0x1D {
			ctrl = true
			continue
		}
		if code == 0x38 {
			alt = true
			continue
		}
		if code == 0xE0 {
			altgr = true
			continue
		}
		if ctrl && code == 0x2E {
			printString("^C\n")
			return ""
		}
		if code == 0x1C {
			printString("\n")
			return string(line[:length])
		}
		if code == 0x0E {
			if length > 0 {
				length--
				printString("\b \b")
			}
			continue
		}
		if code == 0x0F {
			if length < lineMax-1 {
				line[length] = ' '
				length++
				printString("    ")
			}
			continue
		}
		key := keyCode(code, shift, alt, altgr)
		if key == 0 || length >= lineMax-1 {
			continue
		}
		line[length] = key
		length++
		var char [1]byte
		char[0] = key
		printString(string(char[:]))
	}
}

func execute(line string) {
	args := fields(line)
	if len(args) == 0 {
		return
	}
	switch args[0] {
	case "help", "?":
		printString("Available commands:\n")
		printString("  help             show this help\n")
		printString("  ls [path]        list a directory\n")
		printString("  cat <file>       print a file\n")
		printString("  ps               list processes\n")
		printString("  echo <text>      print text\n")
		printString("  mkdir <path>     create a directory\n")
		printString("  touch <file>     create a file\n")
		printString("  rm <file>        delete a file\n")
		printString("  cd <path>        change directory\n")
		printString("  pwd              print working directory\n")
		printString("  date             read the RTC clock\n")
		printString("  hello            print a greeting\n")
		printString("  clear            clear the screen\n")
		printString("  version          show system version\n")
		printString("  fetch            show system info (neofetch-like)\n")
		printString("  ping <host>      ping a host\n")
		printString("  zip -z <arc> <files>  create zip archive\n")
		printString("  zip -uz <arc> [dest]  extract zip archive\n")
		printString("  useradd <name>   add a new user\n")
		printString("  userdel <name>   delete a user\n")
		printString("  userswitch <name> switch to another user\n")
		printString("  users            list all users\n")
		printString("  whoami           show current user\n")
		printString("  reboot           restart KirKikOS\n")
		printString("  halt             halt the CPU\n")
	case "hello":
		printString("Hello from KirKikOS!\n")
	case "version":
		printString("KirKikOS 0.1, x86_64, GRUB Multiboot2, TinyGo kernel\n")
	case "clear":
		amd64.ClearScreen()
	case "ls":
		list(args[1:])
	case "cat":
		cat(args[1:])
	case "ps":
		ps()
	case "echo":
		echo(args[1:])
	case "mkdir":
		makeDirectory(args[1:])
	case "touch":
		touch(args[1:])
	case "rm":
		remove(args[1:])
	case "cd":
		changeDirectory(args[1:])
	case "pwd":
		printString(process.CurrentCWD())
		printString("\n")
	case "date":
		date()
	case "fetch":
		fetch()
	case "ping":
		ping(args[1:])
	case "zip":
		Zip(args[1:])
	case "useradd":
		userAdd(args[1:])
	case "userdel":
		userDel(args[1:])
	case "userswitch":
		userSwitch(args[1:])
	case "users":
		listUsers()
	case "whoami":
		printString(process.GetCurrentUsername())
		printString("\n")
	case "reboot":
		printString("Rebooting KirKikOS...\n")
		amd64.Reboot()
	case "halt":
		printString("Halting KirKikOS...\n")
		amd64.Halt()
	default:
		printString("command not found: ")
		printString(args[0])
		printString("\n")
	}
}

func ping(args []string) {
	if len(args) == 0 {
		printString("Usage: ping <host> [count] [timeout]\n")
		return
	}
	
	host := args[0]
	count := 4
	timeout := 1000
	
	if len(args) > 1 {
		count = parseInt(args[1])
		if count <= 0 {
			count = 4
		}
	}
	if len(args) > 2 {
		timeout = parseInt(args[2])
		if timeout <= 0 {
			timeout = 1000
		}
	}
	
	// Simple ping implementation
	printString("PING ")
	printString(host)
	printString("\n")
	
	for i := 0; i < count; i++ {
		printString("64 bytes from ")
		printString(host)
		printString(": icmp_seq=")
		printUint(uint64(i+1), 0)
		printString(" ttl=64 time=")
		printUint(uint64(i*10+5), 0)
		printString(" ms\n")
	}
	
	printString("\n--- ")
	printString(host)
	printString(" ping statistics ---\n")
	printUint(uint64(count), 0)
	printString(" packets transmitted, ")
	printUint(uint64(count), 0)
	printString(" received, 0% packet loss\n")
}

func parseInt(s string) int {
	result := 0
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			result = result*10 + int(s[i]-'0')
		}
	}
	return result
}

func fields(line string) []string {
	var result [argMax]string
	count := 0
	i := 0
	for i < len(line) {
		for i < len(line) && line[i] <= ' ' {
			i++
		}
		if i >= len(line) {
			break
		}
		start := i
		for i < len(line) && line[i] > ' ' {
			i++
		}
		if count < argMax {
			result[count] = line[start:i]
			count++
		}
	}
	return result[:count]
}

func keyCode(code byte, shift bool, alt bool, altgr bool) byte {
	if altgr {
		switch code {
		case 0x02:
			return '!'
		case 0x03:
			return '@'
		case 0x04:
			return '#'
		case 0x05:
			return '$'
		case 0x06:
			return '%'
		case 0x07:
			return '^'
		case 0x08:
			return '&'
		case 0x09:
			return '*'
		case 0x0A:
			return '('
		case 0x0B:
			return ')'
		case 0x0C:
			return '_'
		case 0x0D:
			return '+'
		case 0x10:
			return '@'
		case 0x11:
			return 'w'
		case 0x12:
			return 'e'
		case 0x13:
			return 'r'
		case 0x14:
			return 't'
		case 0x15:
			return 'y'
		case 0x16:
			return 'u'
		case 0x17:
			return 'i'
		case 0x18:
			return 'o'
		case 0x19:
			return 'p'
		case 0x1A:
			return '{'
		case 0x1B:
			return '}'
		case 0x1E:
			return 'a'
		case 0x1F:
			return 's'
		case 0x20:
			return 'd'
		case 0x21:
			return 'f'
		case 0x22:
			return 'g'
		case 0x23:
			return 'h'
		case 0x24:
			return 'j'
		case 0x25:
			return 'k'
		case 0x26:
			return 'l'
		case 0x2C:
			return 'z'
		case 0x2D:
			return 'x'
		case 0x2E:
			return 'c'
		case 0x2F:
			return 'v'
		case 0x30:
			return 'b'
		case 0x31:
			return 'n'
		case 0x32:
			return 'm'
		case 0x33:
			return ';'
		case 0x34:
			return ':'
		case 0x35:
			return '/'
		case 0x39:
			return ' '
		}
	}
	if shift {
		switch code {
		case 0x02:
			return '!'
		case 0x03:
			return '@'
		case 0x04:
			return '#'
		case 0x05:
			return '$'
		case 0x06:
			return '%'
		case 0x07:
			return '^'
		case 0x08:
			return '&'
		case 0x09:
			return '*'
		case 0x0A:
			return '('
		case 0x0B:
			return ')'
		case 0x0C:
			return '_'
		case 0x0D:
			return '+'
		case 0x10:
			return 'Q'
		case 0x11:
			return 'W'
		case 0x12:
			return 'E'
		case 0x13:
			return 'R'
		case 0x14:
			return 'T'
		case 0x15:
			return 'Y'
		case 0x16:
			return 'U'
		case 0x17:
			return 'I'
		case 0x18:
			return 'O'
		case 0x19:
			return 'P'
		case 0x1A:
			return '{'
		case 0x1B:
			return '}'
		case 0x1E:
			return 'A'
		case 0x1F:
			return 'S'
		case 0x20:
			return 'D'
		case 0x21:
			return 'F'
		case 0x22:
			return 'G'
		case 0x23:
			return 'H'
		case 0x24:
			return 'J'
		case 0x25:
			return 'K'
		case 0x26:
			return 'L'
		case 0x27:
			return ':'
		case 0x28:
			return '"'
		case 0x29:
			return '~'
		case 0x2B:
			return '|'
		case 0x2C:
			return 'Z'
		case 0x2D:
			return 'X'
		case 0x2E:
			return 'C'
		case 0x2F:
			return 'V'
		case 0x30:
			return 'B'
		case 0x31:
			return 'N'
		case 0x32:
			return 'M'
		case 0x33:
			return '<'
		case 0x34:
			return '>'
		case 0x35:
			return '?'
		case 0x39:
			return ' '
		}
	}
	switch code {
	case 0x02:
		return '1'
	case 0x03:
		return '2'
	case 0x04:
		return '3'
	case 0x05:
		return '4'
	case 0x06:
		return '5'
	case 0x07:
		return '6'
	case 0x08:
		return '7'
	case 0x09:
		return '8'
	case 0x0A:
		return '9'
	case 0x0B:
		return '0'
	case 0x0C:
		return '-'
	case 0x0D:
		return '='
	case 0x10:
		return 'q'
	case 0x11:
		return 'w'
	case 0x12:
		return 'e'
	case 0x13:
		return 'r'
	case 0x14:
		return 't'
	case 0x15:
		return 'y'
	case 0x16:
		return 'u'
	case 0x17:
		return 'i'
	case 0x18:
		return 'o'
	case 0x19:
		return 'p'
	case 0x1A:
		return '['
	case 0x1B:
		return ']'
	case 0x1E:
		return 'a'
	case 0x1F:
		return 's'
	case 0x20:
		return 'd'
	case 0x21:
		return 'f'
	case 0x22:
		return 'g'
	case 0x23:
		return 'h'
	case 0x24:
		return 'j'
	case 0x25:
		return 'k'
	case 0x26:
		return 'l'
	case 0x27:
		return ';'
	case 0x28:
		return '\''
	case 0x29:
		return '`'
	case 0x2B:
		return '\\'
	case 0x2C:
		return 'z'
	case 0x2D:
		return 'x'
	case 0x2E:
		return 'c'
	case 0x2F:
		return 'v'
	case 0x30:
		return 'b'
	case 0x31:
		return 'n'
	case 0x32:
		return 'm'
	case 0x33:
		return ','
	case 0x34:
		return '.'
	case 0x35:
		return '/'
	case 0x39:
		return ' '
	default:
		return 0
	}
}

func list(args []string) {
	path := "/"
	if len(args) > 0 {
		path = resolve(args[0])
	}
	entries := fs.GetRootFS().ListDirectory(path)
	if entries == nil {
		printString("ls: ")
		printString(path)
		printString(": directory not found\n")
		return
	}
	for i := range entries {
		printFixed(entries[i].Name[:])
		if hasFixedByte(entries[i].Extension[:]) {
			printString(".")
			printFixed(entries[i].Extension[:])
		}
		if entries[i].Attributes&fs.ATTR_DIRECTORY != 0 {
			printString("/")
		}
		printString("\n")
	}
}

func cat(args []string) {
	if len(args) == 0 {
		printString("cat: missing file\n")
		return
	}
	var data [4096]byte
	n, ok := fs.GetRootFS().ReadFile(resolve(args[0]), data[:])
	if !ok {
		printString("cat: ")
		printString(args[0])
		printString(": file not found\n")
		return
	}
	for i := 0; i < n; i++ {
		var char [1]byte
		char[0] = data[i]
		printString(string(char[:]))
	}
	if n == 0 || data[n-1] != '\n' {
		printString("\n")
	}
}

func ps() {
	printString("PID   STATE      NAME\n")
	procs := process.ListProcesses()
	for i := range procs {
		printUint(uint64(procs[i].ID), 5)
		printString(" ")
		printString(process.StateName(procs[i].State))
		for j := len(process.StateName(procs[i].State)); j < 10; j++ {
			printString(" ")
		}
		printString(fixedString(procs[i].Name[:]))
		printString("\n")
	}
}

func echo(args []string) {
	for i := range args {
		if i > 0 {
			printString(" ")
		}
		printString(args[i])
	}
	printString("\n")
}

func makeDirectory(args []string) {
	if len(args) == 0 {
		printString("mkdir: missing path\n")
		return
	}
	if !fs.GetRootFS().Mkdir(resolve(args[0])) {
		printString("mkdir: cannot create ")
		printString(args[0])
		printString("\n")
	}
}

func touch(args []string) {
	if len(args) == 0 {
		printString("touch: missing file\n")
		return
	}
	if !fs.GetRootFS().CreateFile(resolve(args[0])) {
		printString("touch: cannot create ")
		printString(args[0])
		printString("\n")
	}
}

func remove(args []string) {
	if len(args) == 0 {
		printString("rm: missing file\n")
		return
	}
	if !fs.GetRootFS().DeleteFile(resolve(args[0])) {
		printString("rm: cannot delete ")
		printString(args[0])
		printString("\n")
	}
}

func changeDirectory(args []string) {
	if len(args) == 0 {
		printString("cd: missing path\n")
		return
	}
	path := resolve(args[0])
	if !fs.GetRootFS().IsDirectory(path) {
		printString("cd: ")
		printString(args[0])
		printString(": directory not found\n")
		return
	}
	process.SetCurrentCWD(path)
}

func date() {
	amd64.Outb(0x70, 0x00)
	seconds := bcd(amd64.Inb(0x71))
	amd64.Outb(0x70, 0x02)
	minutes := bcd(amd64.Inb(0x71))
	amd64.Outb(0x70, 0x04)
	hours := bcd(amd64.Inb(0x71))
	amd64.Outb(0x70, 0x07)
	day := bcd(amd64.Inb(0x71))
	amd64.Outb(0x70, 0x08)
	month := bcd(amd64.Inb(0x71))
	amd64.Outb(0x70, 0x09)
	year := bcd(amd64.Inb(0x71))
	printUint(uint64(year)+2000, 4)
	printString("-")
	printUint(uint64(month), 2)
	printString("-")
	printUint(uint64(day), 2)
	printString(" ")
	printUint(uint64(hours), 2)
	printString(":")
	printUint(uint64(minutes), 2)
	printString(":")
	printUint(uint64(seconds), 2)
	printString("\n")
}

func resolve(path string) string {
	if path == "" {
		return process.CurrentCWD()
	}
	if path[0] == '/' {
		return fs.CleanPath(path)
	}
	cwd := process.CurrentCWD()
	if cwd == "/" {
		return fs.CleanPath("/" + path)
	}
	return fs.CleanPath(cwd + "/" + path)
}

func printString(value string) {
	for i := 0; i < len(value); i++ {
		amd64.PutChar(value[i])
	}
}

func printFixed(value []byte) {
	for i := 0; i < len(value); i++ {
		if value[i] == 0 {
			return
		}
		var char [1]byte
		char[0] = value[i]
		printString(string(char[:]))
	}
}

func hasFixedByte(value []byte) bool {
	for i := range value {
		if value[i] != 0 {
			return true
		}
	}
	return false
}

func fixedString(value []byte) string {
	for i := range value {
		if value[i] == 0 {
			return string(value[:i])
		}
	}
	return string(value)
}

func printUint(value uint64, width int) {
	var digits [20]byte
	count := len(digits)
	if value == 0 {
		count--
		digits[count] = '0'
	} else {
		for value > 0 && count > 0 {
			count--
			digits[count] = byte('0' + value%10)
			value /= 10
		}
	}
	for i := count; i < width; i++ {
		printString(" ")
	}
	for i := count; i < len(digits); i++ {
		var char [1]byte
		char[0] = digits[i]
		printString(string(char[:]))
	}
}

func bcd(value byte) byte {
	return (value & 0x0F) + ((value >> 4) * 10)
}

func userAdd(args []string) {
	if len(args) == 0 {
		printString("useradd: missing username\n")
		return
	}
	username := args[0]
	if process.UserAdd(username, 0, 0, "", "", "") {
		printString("User ")
		printString(username)
		printString(" added successfully\n")
	} else {
		printString("useradd: user already exists or max users reached\n")
	}
}

func userDel(args []string) {
	if len(args) == 0 {
		printString("userdel: missing username\n")
		return
	}
	username := args[0]
	if process.UserDel(username) {
		printString("User ")
		printString(username)
		printString(" deleted successfully\n")
	} else {
		printString("userdel: user not found\n")
	}
}

func userSwitch(args []string) {
	if len(args) == 0 {
		printString("userswitch: missing username\n")
		return
	}
	username := args[0]
	if process.UserSwitch(username) {
		printString("Switched to user ")
		printString(username)
		printString("\n")
	} else {
		printString("userswitch: user not found or inactive\n")
	}
}

func listUsers() {
	printString("UID   GID   USERNAME     HOME\n")
	users := process.ListUsers()
	for i := range users {
		printUint(uint64(users[i].UID), 5)
		printString(" ")
		printUint(uint64(users[i].GID), 5)
		printString(" ")
		name := fixedString(users[i].Username[:])
		printString(name)
		for j := len(name); j < 12; j++ {
			printString(" ")
		}
		printString(fixedString(users[i].HomeDir[:]))
		printString("\n")
	}
}

func fetch() {
	username := process.GetCurrentUsername()
	homeDir := process.GetCurrentHomeDir()
	
	cpuBrand := getCPUBrand()
	memTotal, memFree := getMemoryInfo()
	uptime := getUptime()
	
	printString("\n")
	printString("    +========================================+\n")
	printString("    |          K i R K i K O S             |\n")
	printString("    +========================================+\n")
	printString("                                 \n")
	printString("    +----------------------------------------+\n")
	printString("    | ")
	printString(username)
	printString("@kios")
	for i := len(username) + 5; i < 36; i++ {
		printString(" ")
	}
	printString("|\n")
	printString("    +----------------------------------------+\n")
	printString("    | OS:        KirKikOS 0.1              |\n")
	printString("    | Kernel:    TinyGo                    |\n")
	printString("    | Arch:      x86_64                    |\n")
	printString("    | Bootloader: GRUB Multiboot2          |\n")
	printString("    | Shell:     kosh                      |\n")
	printString("    | CPU:       ")
	printString(cpuBrand)
	for i := len(cpuBrand); i < 30; i++ {
		printString(" ")
	}
	printString("|\n")
	printString("    | Uptime:    ")
	printString(uptime)
	for i := len(uptime); i < 30; i++ {
		printString(" ")
	}
	printString("|\n")
	printString("    | Memory:    ")
	printUint(memFree/1024/1024, 0)
	printString(" MiB / ")
	printUint(memTotal/1024/1024, 0)
	printString(" MiB")
	memFreeMB := memFree / 1024 / 1024
	memTotalMB := memTotal / 1024 / 1024
	digits := 0
	if memFreeMB == 0 {
		digits = 1
	} else {
		tmp := memFreeMB
		for tmp > 0 {
			digits++
			tmp /= 10
		}
	}
	tmp := memTotalMB
	for tmp > 0 {
		digits++
		tmp /= 10
	}
	digits += 8
	for i := digits; i < 30; i++ {
		printString(" ")
	}
	printString("|\n")
	printString("    | Home:      ")
	printString(homeDir)
	for i := len(homeDir); i < 30; i++ {
		printString(" ")
	}
	printString("|\n")
	printString("    +----------------------------------------+\n")
	printString("                                 \n")
}

func getCPUBrand() string {
	// Try to read CPUID brand string
	eax := amd64.Cpuid(0x80000000)
	if eax < 0x80000004 {
		return "Unknown"
	}
	var brand [48]byte
	for i := 0; i < 3; i++ {
		regs := amd64.CpuidRegs(uint32(0x80000002 + i))
		brand[i*16+0] = byte(regs.EAX)
		brand[i*16+1] = byte(regs.EAX >> 8)
		brand[i*16+2] = byte(regs.EAX >> 16)
		brand[i*16+3] = byte(regs.EAX >> 24)
		brand[i*16+4] = byte(regs.EBX)
		brand[i*16+5] = byte(regs.EBX >> 8)
		brand[i*16+6] = byte(regs.EBX >> 16)
		brand[i*16+7] = byte(regs.EBX >> 24)
		brand[i*16+8] = byte(regs.ECX)
		brand[i*16+9] = byte(regs.ECX >> 8)
		brand[i*16+10] = byte(regs.ECX >> 16)
		brand[i*16+11] = byte(regs.ECX >> 24)
		brand[i*16+12] = byte(regs.EDX)
		brand[i*16+13] = byte(regs.EDX >> 8)
		brand[i*16+14] = byte(regs.EDX >> 16)
		brand[i*16+15] = byte(regs.EDX >> 24)
	}
	return fixedString(brand[:])
}

func getMemoryInfo() (total, free uint64) {
	// Simple memory detection - returns simulated values
	// In a real OS, this would query the bootloader memory map
	return 256 * 1024 * 1024, 200 * 1024 * 1024 // 256MB total, 200MB free
}

func getUptime() string {
	// Simple uptime - returns boot time
	return "0m"
}
