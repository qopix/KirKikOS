package process

const (
	STATE_READY        = 0
	STATE_RUNNING      = 1
	STATE_BLOCKED      = 2
	STATE_ZOMBIE       = 3
	STATE_TERMINATED   = 4
	FLAG_USER          = 1 << 0
	MAX_PROCESSES      = 256
	MAX_PROCESS_NAME   = 64
	MAX_CWD            = 256
	MAX_USERNAME       = 32
	MAX_USERS          = 32
)

type Process struct {
	ID     uint32
	State  uint32
	Flags  uint32
	PC     uint64
	RSP    uint64
	RBP    uint64
	PML4   uint64
	Entry  uint64
	Argv   [256]byte
	Envv   [256]byte
	Stdin  uint32
	Stdout uint32
	Stderr uint32
	CWD    [MAX_CWD]byte
	Name   [MAX_PROCESS_NAME]byte
}

type User struct {
	UID       uint32
	GID       uint32
	Username  [MAX_USERNAME]byte
	HomeDir   [MAX_CWD]byte
	Shell     [MAX_CWD]byte
	Password  [MAX_CWD]byte
	Active    bool
}

var (
	procs      [MAX_PROCESSES]*Process
	currentPID uint32 = 1
	nextPID    uint32 = 2
	initPID    uint32 = 1
	users      [MAX_USERS]*User
	currentUID uint32 = 1000
	nextUID    uint32 = 1001
	userCount  int
)

func Init() {
	for i := range procs {
		procs[i] = nil
	}
	currentPID = initPID
	nextPID = 2
	newProcess(1, "init", "/")
	InitUsers()
	// Create shell with guest user's home directory
	user := GetCurrentUser()
	homeDir := "/"
	if user != nil {
		homeDir = fixedString(user.HomeDir[:])
	}
	newProcess(2, "shell", homeDir)
	currentPID = 2
}

func StartInit() {
	proc := procs[initPID]
	if proc != nil {
		proc.State = STATE_READY
	}
}

func RunInit() {
	proc := procs[initPID]
	if proc != nil {
		proc.State = STATE_RUNNING
	}
}

func CreateProcess(path string, argv []string, envp []string) *Process {
	name := path
	if len(argv) > 0 && argv[0] != "" {
		name = argv[0]
	}
	proc := newProcess(nextPID, name, CurrentCWD())
	if proc == nil {
		return nil
	}
	for i := 0; i < len(argv) && i < len(proc.Argv); i++ {
		copyFixed(proc.Argv[:], argv[i])
	}
	_ = envp
	return proc
}

func Exit(pid uint32) {
	if pid == 0 {
		pid = currentPID
	}
	proc := procs[pid]
	if proc != nil {
		proc.State = STATE_TERMINATED
	}
	if pid == currentPID && procs[initPID] != nil {
		currentPID = initPID
	}
}

func GetProcess(pid uint32) *Process {
	if pid < MAX_PROCESSES {
		return procs[pid]
	}
	return nil
}

func GetCurrentProcess() *Process {
	return procs[currentPID]
}

func ListProcesses() []*Process {
	var result []*Process
	for i := 1; i < MAX_PROCESSES; i++ {
		if procs[i] != nil {
			result = append(result, procs[i])
		}
	}
	return result
}

func SetCurrentProcess(pid uint32) bool {
	proc := procs[pid]
	if proc == nil {
		return false
	}
	currentPID = pid
	return true
}

func CurrentCWD() string {
	proc := procs[currentPID]
	if proc == nil {
		return "/"
	}
	return fixedString(proc.CWD[:])
}

func SetCurrentCWD(path string) bool {
	proc := procs[currentPID]
	if proc == nil {
		return false
	}
	copyFixed(proc.CWD[:], path)
	return true
}

func StateName(state uint32) string {
	switch state {
	case STATE_READY:
		return "ready"
	case STATE_RUNNING:
		return "running"
	case STATE_BLOCKED:
		return "blocked"
	case STATE_ZOMBIE:
		return "zombie"
	case STATE_TERMINATED:
		return "terminated"
	default:
		return "unknown"
	}
}

func Schedule() {
	for offset := uint32(1); offset < MAX_PROCESSES; offset++ {
		pid := (currentPID + offset) % MAX_PROCESSES
		if pid == 0 {
			continue
		}
		proc := procs[pid]
		if proc != nil && proc.State == STATE_READY {
			if current := procs[currentPID]; current != nil && current.State == STATE_RUNNING {
				current.State = STATE_READY
			}
			proc.State = STATE_RUNNING
			currentPID = pid
			return
		}
	}
}

func newProcess(id uint32, name string, cwd string) *Process {
	if id == 0 || id >= MAX_PROCESSES || procs[id] != nil {
		return nil
	}
	proc := &Process{
		ID:     id,
		State:  STATE_READY,
		Flags:  FLAG_USER,
		Stdin:  0,
		Stdout: 1,
		Stderr: 2,
	}
	copyFixed(proc.Name[:], name)
	copyFixed(proc.CWD[:], cwd)
	procs[id] = proc
	if id >= nextPID {
		nextPID = id + 1
	}
	return proc
}

func copyFixed(dst []byte, src string) {
	for i := 0; i < len(dst) && i < len(src); i++ {
		dst[i] = src[i]
	}
}

func fixedString(src []byte) string {
	for i := 0; i < len(src); i++ {
		if src[i] == 0 {
			return string(src[:i])
		}
	}
	return string(src)
}

func InitUsers() {
	for i := range users {
		users[i] = nil
	}
	userCount = 0
	nextUID = 1001
	currentUID = 1000
	addUser("guest", 1000, 1000, "/home/usr/dta", "/bin/sh", "")
	addUser("root", 0, 0, "/root", "/bin/sh", "")
}

func addUser(username string, uid, gid uint32, homeDir, shell, password string) *User {
	if userCount >= MAX_USERS {
		return nil
	}
	for i := 0; i < MAX_USERS; i++ {
		if users[i] != nil && fixedString(users[i].Username[:]) == username {
			return nil
		}
	}
	user := &User{
		UID:  uid,
		GID:  gid,
		Active: true,
	}
	copyFixed(user.Username[:], username)
	copyFixed(user.HomeDir[:], homeDir)
	copyFixed(user.Shell[:], shell)
	copyFixed(user.Password[:], password)
	for i := 0; i < MAX_USERS; i++ {
		if users[i] == nil {
			users[i] = user
			userCount++
			if uid >= nextUID {
				nextUID = uid + 1
			}
			return user
		}
	}
	return nil
}

func UserAdd(username string, uid, gid uint32, homeDir, shell, password string) bool {
	if uid == 0 {
		uid = nextUID
	}
	if gid == 0 {
		gid = uid
	}
	if homeDir == "" {
		homeDir = "/home/" + username
	}
	if shell == "" {
		shell = "/bin/sh"
	}
	return addUser(username, uid, gid, homeDir, shell, password) != nil
}

func UserDel(username string) bool {
	for i := 0; i < MAX_USERS; i++ {
		if users[i] != nil && fixedString(users[i].Username[:]) == username {
			users[i] = nil
			userCount--
			return true
		}
	}
	return false
}

func UserSwitch(username string) bool {
	for i := 0; i < MAX_USERS; i++ {
		if users[i] != nil && fixedString(users[i].Username[:]) == username && users[i].Active {
			currentUID = users[i].UID
			proc := procs[currentPID]
			if proc != nil {
				copyFixed(proc.CWD[:], fixedString(users[i].HomeDir[:]))
			}
			return true
		}
	}
	return false
}

func GetCurrentUser() *User {
	for i := 0; i < MAX_USERS; i++ {
		if users[i] != nil && users[i].UID == currentUID {
			return users[i]
		}
	}
	return nil
}

func GetCurrentUID() uint32 {
	return currentUID
}

func GetCurrentUsername() string {
	user := GetCurrentUser()
	if user != nil {
		return fixedString(user.Username[:])
	}
	return "unknown"
}

func GetCurrentHomeDir() string {
	user := GetCurrentUser()
	if user != nil {
		return fixedString(user.HomeDir[:])
	}
	return "/"
}

func ListUsers() []*User {
	var result []*User
	for i := 0; i < MAX_USERS; i++ {
		if users[i] != nil {
			result = append(result, users[i])
		}
	}
	return result
}
