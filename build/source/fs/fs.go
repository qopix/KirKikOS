package fs

const (
	FAT32BootSignature = 0x55AA
	FAT32EOC           = 0x0FFFFFF8
	FAT32EntrySize     = 4

	ATTR_READONLY  = 0x01
	ATTR_HIDDEN    = 0x02
	ATTR_SYSTEM    = 0x04
	ATTR_VOLUME_ID = 0x08
	ATTR_DIRECTORY = 0x10
	ATTR_ARCHIVE   = 0x20
	ATTR_LONG_NAME = 0x0F

	maxEntries   = 64
	maxPathLen   = 128
	maxFileLen   = 4096
)

type FAT32BootSector struct {
	BootJump          [3]byte
	OEMName           [8]byte
	BytesPerSector    uint16
	SectorsPerCluster uint8
	ReservedSectors   uint16
	NumFATs           uint8
	RootEntries       uint16
	TotalSectors16    uint16
	Media             uint8
	FATSize16         uint16
	SectorsPerTrack   uint16
	NumHeads          uint16
	HiddenSectors     uint32
	TotalSectors32    uint32
	FATSize32         uint32
	ExtendedFlags     uint16
	FileSystemVersion uint16
	RootCluster       uint32
	FSInfoSector      uint16
	BackupBootSector  uint16
	Reserved           [12]byte
	DriveNumber        uint8
	Reserved2          uint8
	BootSignature      uint8
	VolumeID           uint32
	VolumeLabel        [11]byte
	FileSystemType     [8]byte
}

type FAT32DirEntry struct {
	Name               [8]byte
	Extension          [3]byte
	Attributes         uint8
	Reserved           uint8
	CreationTimeFine   uint8
	CreationTime       uint16
	CreationDate       uint16
	LastAccessDate     uint16
	FirstClusterHigh   uint16
	LastWriteTime      uint16
	LastWriteDate      uint16
	FirstClusterLow    uint16
	FileSize           uint32
}

type Entry struct {
	Path    [maxPathLen]byte
	PathLen int
	Data    [maxFileLen]byte
	Size    int
	IsDir   bool
}

type FAT32FS struct {
	bootSector        *FAT32BootSector
	drive             uint8
	sectorsPerCluster uint8
	bytesPerSector    uint16
	clusterSize       uint32
	rootCluster       uint32
	fatStart          uint32
	dataStart         uint32
	entries           [maxEntries]Entry
	count             int
}

var rootFS *FAT32FS

func Init() {
	rootFS = &FAT32FS{
		sectorsPerCluster: 1,
		bytesPerSector:    512,
		clusterSize:       512,
		rootCluster:       2,
	}
	rootFS.addDir("/")
	rootFS.addDir("/bin")
	rootFS.addDir("/cmd")
	rootFS.addDir("/etc")
	rootFS.addDir("/home")
	rootFS.addDir("/home/usr")
	rootFS.addDir("/home/usr/dta")
	rootFS.addDir("/source")
	rootFS.addFile("/bin/hello", nil)
	rootFS.addFile("/bin/ls", nil)
	rootFS.addFile("/bin/cat", nil)
	rootFS.addFile("/bin/ps", nil)
	rootFS.addFile("/bin/init", nil)
	rootFS.addFile("/bin/mkdir", nil)
	rootFS.addFile("/bin/touch", nil)
	rootFS.addFile("/bin/rm", nil)
	rootFS.addFile("/bin/cd", nil)
	rootFS.addFile("/bin/pwd", nil)
	rootFS.addFile("/bin/date", nil)
	rootFS.addFile("/bin/fetch", nil)
	rootFS.addFile("/bin/useradd", nil)
	rootFS.addFile("/bin/userdel", nil)
	rootFS.addFile("/bin/userswitch", nil)
	rootFS.addFile("/bin/users", nil)
	rootFS.addFile("/bin/whoami", nil)
	rootFS.addFile("/bin/reboot", nil)
	rootFS.addFile("/bin/halt", nil)
	rootFS.addFile("/bin/clear", nil)
	rootFS.addFile("/bin/version", nil)
	rootFS.addFile("/bin/help", nil)
	rootFS.addFile("/bin/echo", nil)
	rootFS.addFile("/cmd/hello", nil)
	rootFS.addFile("/cmd/ls", nil)
	rootFS.addFile("/cmd/cat", nil)
	rootFS.addFile("/cmd/ps", nil)
	rootFS.addFile("/cmd/init", nil)
	rootFS.addFile("/cmd/mkdir", nil)
	rootFS.addFile("/cmd/touch", nil)
	rootFS.addFile("/cmd/rm", nil)
	rootFS.addFile("/cmd/cd", nil)
	rootFS.addFile("/cmd/pwd", nil)
	rootFS.addFile("/cmd/date", nil)
	rootFS.addFile("/cmd/fetch", nil)
	rootFS.addFile("/cmd/useradd", nil)
	rootFS.addFile("/cmd/userdel", nil)
	rootFS.addFile("/cmd/userswitch", nil)
	rootFS.addFile("/cmd/users", nil)
	rootFS.addFile("/cmd/whoami", nil)
	rootFS.addFile("/cmd/reboot", nil)
	rootFS.addFile("/cmd/halt", nil)
	rootFS.addFile("/cmd/clear", nil)
	rootFS.addFile("/cmd/version", nil)
	rootFS.addFile("/cmd/help", nil)
	rootFS.addFile("/cmd/echo", nil)
	rootFS.addFile("/source/kernel.go", []byte("// KirKikOS Kernel Source\npackage main\n"))
	rootFS.addFile("/source/shell.go", []byte("// KirKikOS Shell Source\npackage shell\n"))
	rootFS.addFile("/source/fs.go", []byte("// KirKikOS Filesystem Source\npackage fs\n"))
	rootFS.addFile("/source/process.go", []byte("// KirKikOS Process Source\npackage process\n"))
	rootFS.addFile("/etc/motd", []byte("Welcome to KirKikOS.\n"))
	rootFS.addFile("/etc/os-release", []byte("NAME=\"KirKikOS\"\nVERSION=\"0.1\"\n"))
	rootFS.addFile("/home/usr/dta/readme.txt", []byte("KirKikOS data directory.\nUse 'help' to list commands.\n"))
	rootFS.addFile("/home/usr/dta/passwd", []byte("guest:*:1000:1000:Guest:/home/usr/dta:/bin/sh\n"))
	rootFS.addFile("/home/usr/dta/shadow", []byte("guest:*:18000:0:99999:7:::\n"))
}

func GetRootFS() *FAT32FS {
	return rootFS
}

func CleanPath(path string) string {
	if path == "" {
		return "/"
	}
	if path[0] != '/' {
		path = "/" + path
	}
	var parts [32]string
	count := 0
	start := 1
	for start < len(path) {
		for start < len(path) && path[start] == '/' {
			start++
		}
		if start >= len(path) {
			break
		}
		end := start
		for end < len(path) && path[end] != '/' {
			end++
		}
		part := path[start:end]
		if part == "." {
		} else if part == ".." {
			if count > 0 {
				count--
			}
		} else {
			parts[count] = part
			count++
		}
		start = end
	}
	if count == 0 {
		return "/"
	}
	result := ""
	for i := 0; i < count; i++ {
		result += "/" + parts[i]
	}
	return result
}

func (fs *FAT32FS) ReadCluster(cluster uint32, buffer []byte) bool {
	if fs == nil || cluster < 2 || int(cluster-2) >= fs.count {
		return false
	}
	source := fs.entries[cluster-2].Data
	copy(buffer, source[:fs.entries[cluster-2].Size])
	return true
}

func (fs *FAT32FS) WriteCluster(cluster uint32, buffer []byte) bool {
	if fs == nil || cluster < 2 || int(cluster-2) >= fs.count {
		return false
	}
	target := &fs.entries[cluster-2]
	copy(target.Data[:], buffer)
	target.Size = len(buffer)
	if target.Size > maxFileLen {
		target.Size = maxFileLen
	}
	return true
}

func (fs *FAT32FS) GetNextCluster(cluster uint32) uint32 {
	if fs == nil || int(cluster) >= fs.count {
		return 0
	}
	return cluster + 1
}

func (fs *FAT32FS) ListDirectory(path string) []FAT32DirEntry {
	if fs == nil {
		return nil
	}
	path = CleanPath(path)
	var result []FAT32DirEntry
	for i := 0; i < fs.count; i++ {
		item := &fs.entries[i]
		if item.pathString() == path || parentPath(item.pathString()) != path {
			continue
		}
		var out FAT32DirEntry
		name := baseName(item.pathString())
		copyFixed(out.Name[:], name)
		if dot := indexOf(name, '.'); dot >= 0 && len(name)-dot-1 <= 3 {
			copyFixed(out.Extension[:], name[dot+1:])
			for j := 0; j < dot && j < 8; j++ {
				out.Name[j] = out.Name[j+1]
			}
		}
		out.Attributes = ATTR_ARCHIVE
		if item.IsDir {
			out.Attributes = ATTR_DIRECTORY
		}
		out.FileSize = uint32(item.Size)
		result = append(result, out)
	}
	return result
}

func (fs *FAT32FS) ReadFile(path string, buffer []byte) (int, bool) {
	if fs == nil {
		return 0, false
	}
	item := fs.find(CleanPath(path))
	if item == nil || item.IsDir {
		return 0, false
	}
	if buffer != nil {
		if len(buffer) < item.Size {
			copy(buffer, item.Data[:len(buffer)])
			return len(buffer), true
		}
		copy(buffer, item.Data[:item.Size])
	}
	return item.Size, true
}

func (fs *FAT32FS) WriteFile(path string, data []byte) bool {
	if fs == nil {
		return false
	}
	path = CleanPath(path)
	if path == "/" {
		return false
	}
	item := fs.find(path)
	if item != nil {
		if item.IsDir {
			return false
		}
		copy(item.Data[:], data)
		item.Size = len(data)
		if item.Size > maxFileLen {
			item.Size = maxFileLen
		}
		return true
	}
	return fs.addFile(path, data)
}

func (fs *FAT32FS) CreateFile(path string) bool {
	return fs.WriteFile(path, nil)
}

func (fs *FAT32FS) DeleteFile(path string) bool {
	if fs == nil {
		return false
	}
	path = CleanPath(path)
	index := fs.findIndex(path)
	if index < 0 || fs.entries[index].IsDir {
		return false
	}
	fs.remove(index)
	return true
}

func (fs *FAT32FS) Mkdir(path string) bool {
	if fs == nil {
		return false
	}
	path = CleanPath(path)
	if path == "/" || fs.find(path) != nil {
		return false
	}
	parent := parentPath(path)
	parentItem := fs.find(parent)
	if parentItem == nil || !parentItem.IsDir {
		return false
	}
	return fs.addDir(path)
}

func (fs *FAT32FS) Rmdir(path string) bool {
	if fs == nil {
		return false
	}
	path = CleanPath(path)
	if path == "/" {
		return false
	}
	index := fs.findIndex(path)
	if index < 0 || !fs.entries[index].IsDir {
		return false
	}
	for i := 0; i < fs.count; i++ {
		if i != index && parentPath(fs.entries[i].pathString()) == path {
			return false
		}
	}
	fs.remove(index)
	return true
}

func (fs *FAT32FS) Mount() bool {
	return fs != nil
}

func (fs *FAT32FS) Unmount() bool {
	return true
}

func (fs *FAT32FS) Exists(path string) bool {
	return fs != nil && fs.find(CleanPath(path)) != nil
}

func (fs *FAT32FS) IsDirectory(path string) bool {
	item := fs.find(CleanPath(path))
	return item != nil && item.IsDir
}

func (fs *FAT32FS) Find(path string) *Entry {
	return fs.find(CleanPath(path))
}

func (fs *FAT32FS) addDir(path string) bool {
	return fs.add(path, true, nil)
}

func (fs *FAT32FS) addFile(path string, data []byte) bool {
	return fs.add(path, false, data)
}

func (fs *FAT32FS) add(path string, isDir bool, data []byte) bool {
	if fs.count >= maxEntries {
		return false
	}
	path = CleanPath(path)
	if path == "/" {
		if fs.count != 0 {
			return false
		}
		item := &fs.entries[fs.count]
		copyFixed(item.Path[:], path)
		item.PathLen = len(path)
		item.IsDir = true
		fs.count++
		return true
	}
	if fs.find(path) != nil {
		return false
	}
	parent := parentPath(path)
	if parent != "/" && (fs.find(parent) == nil || !fs.find(parent).IsDir) {
		return false
	}
	item := &fs.entries[fs.count]
	copyFixed(item.Path[:], path)
	item.PathLen = len(path)
	item.IsDir = isDir
	if data != nil {
		copy(item.Data[:], data)
		item.Size = len(data)
		if item.Size > maxFileLen {
			item.Size = maxFileLen
		}
	}
	fs.count++
	return true
}

func (fs *FAT32FS) find(path string) *Entry {
	index := fs.findIndex(path)
	if index < 0 {
		return nil
	}
	return &fs.entries[index]
}

func (fs *FAT32FS) findIndex(path string) int {
	path = CleanPath(path)
	for i := 0; i < fs.count; i++ {
		if fs.entries[i].pathString() == path {
			return i
		}
	}
	return -1
}

func (fs *FAT32FS) remove(index int) {
	for i := index; i < fs.count-1; i++ {
		fs.entries[i] = fs.entries[i+1]
	}
	fs.count--
}

func (e *Entry) pathString() string {
	return string(e.Path[:e.PathLen])
}

func parentPath(path string) string {
	if path == "/" {
		return "/"
	}
	for i := len(path) - 1; i > 0; i-- {
		if path[i] == '/' {
			return path[:i]
		}
	}
	return "/"
}

func baseName(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[i+1:]
		}
	}
	return path
}

func indexOf(value string, char byte) int {
	for i := 0; i < len(value); i++ {
		if value[i] == char {
			return i
		}
	}
	return -1
}

func copyFixed(dst []byte, src string) {
	for i := 0; i < len(dst) && i < len(src); i++ {
		dst[i] = src[i]
	}
}
