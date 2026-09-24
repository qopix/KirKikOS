package shell

import (
	"kirkikos/fs"
)

func Zip(args []string) {
	if len(args) == 0 {
		printString("Usage: zip [-z|-uz] <archive> [files...]\n")
		printString("  -z   create zip archive\n")
		printString("  -uz  extract zip archive\n")
		return
	}
	
	if args[0] == "-z" {
		if len(args) < 3 {
			printString("Usage: zip -z <archive> <file1> [file2...]\n")
			return
		}
		archive := args[1]
		files := args[2:]
		createZip(archive, files)
	} else if args[0] == "-uz" {
		if len(args) < 2 {
			printString("Usage: zip -uz <archive> [destination]\n")
			return
		}
		archive := args[1]
		dest := "."
		if len(args) > 2 {
			dest = args[2]
		}
		extractZip(archive, dest)
	} else {
		printString("Unknown option: ")
		printString(args[0])
		printString("\nUsage: zip [-z|-uz] <archive> [files...]\n")
	}
}

func createZip(archive string, files []string) {
	printString("Creating archive: ")
	printString(archive)
	printString("\n")
	
	for _, file := range files {
		resolved := resolve(file)
		item := fs.GetRootFS().Find(resolved)
		if item == nil {
			printString("  zip: ")
			printString(file)
			printString(": file not found\n")
			continue
		}
		
		if item.IsDir {
			printString("  adding: ")
			printString(file)
			printString("/ (directory)\n")
			addDirToZip(archive, resolved)
		} else {
			printString("  adding: ")
			printString(file)
			printString(" (")
			printUint(uint64(item.Size), 0)
			printString(" bytes)\n")
			addFileToZip(archive, resolved, item)
		}
	}
	
	printString("Archive created successfully.\n")
}

func addDirToZip(archive, dir string) {
	entries := fs.GetRootFS().ListDirectory(dir)
	if entries == nil {
		return
	}
	for _, entry := range entries {
		name := getEntryName(entry)
		fullPath := dir
		if fullPath != "/" {
			fullPath += "/"
		}
		fullPath += name
		
		item := fs.GetRootFS().Find(fullPath)
		if item == nil {
			continue
		}
		
		if entry.Attributes&fs.ATTR_DIRECTORY != 0 {
			addDirToZip(archive, fullPath)
		} else {
			addFileToZip(archive, fullPath, item)
		}
	}
}

func addFileToZip(archive, path string, item *fs.Entry) {
	// In a real implementation, this would write to a zip file
	// For now, we store in the filesystem with .zip extension
	zipPath := archive + "/" + path
	fs.GetRootFS().WriteFile(zipPath, item.Data[:item.Size])
}

func extractZip(archive, dest string) {
	printString("Extracting archive: ")
	printString(archive)
	printString(" to ")
	printString(dest)
	printString("\n")
	
	// List contents of zip archive
	entries := fs.GetRootFS().ListDirectory(archive)
	if entries == nil {
		printString("zip: archive not found\n")
		return
	}
	
	for _, entry := range entries {
		name := getEntryName(entry)
		fullPath := archive + "/" + name
		
		item := fs.GetRootFS().Find(fullPath)
		if item == nil {
			continue
		}
		
		destPath := dest
		if dest != "/" && dest != "." {
			destPath += "/"
		}
		destPath += name
		
		if entry.Attributes&fs.ATTR_DIRECTORY != 0 {
			printString("  creating: ")
			printString(destPath)
			printString("/\n")
			fs.GetRootFS().Mkdir(destPath)
		} else {
			printString("  extracting: ")
			printString(destPath)
			printString(" (")
			printUint(uint64(item.Size), 0)
			printString(" bytes)\n")
			fs.GetRootFS().WriteFile(destPath, item.Data[:item.Size])
		}
	}
	
	printString("Archive extracted successfully.\n")
}

func getEntryName(e fs.FAT32DirEntry) string {
	name := ""
	for i := 0; i < 8 && e.Name[i] != 0 && e.Name[i] != ' '; i++ {
		name += string(e.Name[i])
	}
	if e.Extension[0] != 0 && e.Extension[0] != ' ' {
		ext := ""
		for i := 0; i < 3 && e.Extension[i] != 0 && e.Extension[i] != ' '; i++ {
			ext += string(e.Extension[i])
		}
		if ext != "" {
			name += "." + ext
		}
	}
	if e.Attributes&fs.ATTR_DIRECTORY != 0 {
		name += "/"
	}
	return name
}