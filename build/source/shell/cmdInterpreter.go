package shell

import (
	"kirkikos/fs"
	"kirkikos/process"
)

func ParseCommand(line string) (cmd string, args []string) {
	line = trimSpace(line)
	if line == "" {
		return "", nil
	}
	
	parts := splitArgs(line)
	if len(parts) == 0 {
		return "", nil
	}
	return parts[0], parts[1:]
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && s[start] <= ' ' {
		start++
	}
	for end > start && s[end-1] <= ' ' {
		end--
	}
	return s[start:end]
}

func splitArgs(s string) []string {
	var result []string
	var current []byte
	inQuote := false
	quoteChar := byte(0)
	escaped := false
	
	for i := 0; i < len(s); i++ {
		c := s[i]
		
		if escaped {
			current = append(current, c)
			escaped = false
			continue
		}
		
		if c == '\\' {
			escaped = true
			continue
		}
		
		if (c == '"' || c == '\'') && !inQuote {
			inQuote = true
			quoteChar = c
			continue
		}
		
		if c == quoteChar && inQuote {
			inQuote = false
			quoteChar = 0
			continue
		}
		
		if c <= ' ' && !inQuote {
			if len(current) > 0 {
				result = append(result, string(current))
				current = nil
			}
			continue
		}
		
		current = append(current, c)
	}
	
	if len(current) > 0 {
		result = append(result, string(current))
	}
	
	return result
}

func ResolvePath(path string) string {
	if path == "" {
		return process.CurrentCWD()
	}
	
	path = expandTilde(path)
	
	if path[0] == '/' {
		return fs.CleanPath(path)
	}
	
	cwd := process.CurrentCWD()
	if cwd == "/" {
		return fs.CleanPath("/" + path)
	}
	return fs.CleanPath(cwd + "/" + path)
}

func expandTilde(path string) string {
	if len(path) > 0 && path[0] == '~' {
		home := process.GetCurrentHomeDir()
		if len(path) == 1 {
			return home
		}
		if path[1] == '/' {
			return home + path[1:]
		}
	}
	return path
}

func CompletePath(prefix string) []string {
	if prefix == "" {
		return nil
	}
	
	dir := "/"
	filePrefix := prefix
	lastSlash := -1
	for i := 0; i < len(prefix); i++ {
		if prefix[i] == '/' {
			lastSlash = i
		}
	}
	
	if lastSlash >= 0 {
		dir = prefix[:lastSlash+1]
		filePrefix = prefix[lastSlash+1:]
	}
	
	dir = ResolvePath(dir)
	entries := fs.GetRootFS().ListDirectory(dir)
	if entries == nil {
		return nil
	}
	
	var matches []string
	for _, entry := range entries {
		name := getEntryName(entry)
		if len(name) >= len(filePrefix) && name[:len(filePrefix)] == filePrefix {
			if entry.Attributes&fs.ATTR_DIRECTORY != 0 {
				matches = append(matches, dir+name+"/")
			} else {
				matches = append(matches, dir+name)
			}
		}
	}
	
	return matches
}
	
func ValidatePath(path string) error {
	if path == "" {
		return nil
	}
	
	resolved := ResolvePath(path)
	item := fs.GetRootFS().Find(resolved)
	if item == nil {
		return nil
	}
	return nil
}

func IsDir(path string) bool {
	resolved := ResolvePath(path)
	return fs.GetRootFS().IsDirectory(resolved)
}

func Exists(path string) bool {
	resolved := ResolvePath(path)
	return fs.GetRootFS().Exists(resolved)
}