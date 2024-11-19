package lsfunctions

import (
	"fmt"
	"strings"
	"syscall"
	"time"
)

// This function resolves relative paths against a link path.
// It follows the rules of resolving relative paths in a file system.
// It returns the resolved path.
func resolveRelativePath(linkPath, target string) string {
	linkDir := ""
	lastSlash := strings.LastIndex(linkPath, "/")
	if lastSlash != -1 {
		linkDir = linkPath[:lastSlash]
	}
	// Return path as is if it's already absolute
	if strings.HasPrefix(target, "/") {
		return target
	}
	joinedPath := joinPath(linkDir, target)

	// Normalize combined path
	segments := strings.Split(joinedPath, "/")
	var stack []string

	for _, segment := range segments {
		if segment == ".." {
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		} else if segment != "." && segment != "" {
			stack = append(stack, segment)
		}
	}
	return "/" + strings.Join(stack, "/")
}

// Clean string to remove -, _, and. from the name.
func cleanName(name string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case '-', '_', '.', '[', '#', ']', '{', '}', '|', '\\', ':', ';', '<', '>', ',', '?', '!', '@', '$', '%', '^', '&', '*', '(', ')', '~', '`', '"', '\'', '=', '+', '/':
			return -1
		}
		return r
	}, name)
}

// getParentDir returns the parent directory path of the given path.
// It handles various edge cases such as root directory, paths without separators,
// and paths ending with a separator.
func getParentDir(path string) string {
	if path == "/" {
		return "/"
	}
	path = strings.TrimSuffix(path, "/")
	lastIndexSep := strings.LastIndex(path, "/")

	if lastIndexSep == 0 {
		return "/"
	}
	if lastIndexSep == -1 {
		return ".."
	}
	return path[:lastIndexSep]
}

// setEntryPath sets the full path for a FileInfo entry
func setEntryPath(baseDir string, entry *FileDetails) {
	if entry.Name == "." {
		entry.Path = baseDir
	} else if entry.Name == ".." {
		entry.Path = getParentDir(baseDir)
	} else {
		entry.Path = joinPath(baseDir, entry.Name)
	}
}

// joinPath joins two paths into a single path.
func joinPath(dir, file string) string {
	dir = strings.TrimSuffix(dir, "/")
	file = strings.TrimPrefix(file, "/")
	if dir == "" {
		return "/" + file
	}
	return dir + "/" + file
}

// getPathBase extracts the last component from a given path string.
// If the path does not contain any slashes, it returns the original string.
func getPathBase(s string) string {
	l := strings.LastIndex(s, "/")
	if l == -1 {
		return s
	}
	return s[l+1:]
}

// getMax returns the maximum of two integers.
func getMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// formatTime formats a given time based on whether it's in the current year or not.
// For times in the current year, it returns the format "Jan _2 15:04".
// For times in previous years, it returns the format "Jan _2 2006".
func formatTime(modTime time.Time) string {
	now := time.Now()
	if modTime.Year() == now.Year() {
		return modTime.Format("Jan _2 15:04")
	}
	return modTime.Format("Jan _2  2006")
}

// getTotalBlocks calculates the total number of 512-byte blocks in the filesystem.
func getTotalBlocks(entries []FileDetails) TotalBlocks {
	var t TotalBlocks
	for _, entry := range entries {
		if stat, ok := entry.Info.Sys().(*syscall.Stat_t); ok {
			t += TotalBlocks(stat.Blocks)
		}
	}
	return t / 2
}

// addQuotes adds quotes to the input string if it contains spaces or special characters.
func addQuotes(s string) string {
	if strings.Contains(s, " ") || hasSpecialChar(s) {
		s = fmt.Sprintf(`'%s'`, s)
	}
	return s
}

// hasSpecialChar checks if the input string contains any special characters.
func hasSpecialChar(s string) bool {
	if len(s) == 0 {
		return false
	}
	specialChars := []rune{'[', '#', ']', '{', '}', '|', '\\', ':', ';', '<', '>', ',', '?', '!', '@', '$', '%', '^', '&', '*', '(', ')', '~', '`', '"', '\'', '=', '+'}
	for _, ch := range specialChars {
		if strings.ContainsRune(s, ch) {
			return true
		}
	}
	return false
}

// major returns the major device number of the given device.
func major(dev uint64) uint64 {
	return (dev >> 8) & 0xff
}

// minor returns the minor device number of the given device.
func minor(dev uint64) uint64 {
	return dev & 0xff
}
