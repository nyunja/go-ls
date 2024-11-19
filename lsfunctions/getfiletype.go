package lsfunctions

import (
	"fmt"
	"os"
	"strings"
)

// getFileType determines the file type based on the file mode and file extension.
// Returns the entry with the file type.
// If the file type cannot be determined, it returns other.
func getFileType(entry Entry) (Entry, string) {
	mod := entry.Mode
	// Handel pipe files
	if entry.Mode[0] == 'p' {
		return entry, "pipe"
	}
	if entry.Mode[0] == 's' {
		return entry, "socket"
	}
	// Handle setuid files
	if entry.Mode[3] == 's' || entry.Mode[3] == 'S' {
		return entry, "setuid"
	}
	// Handle setuid files
	if entry.Mode[6] == 's' || entry.Mode[6] == 'S' {
		return entry, "setgid"
	}
	// Handle directories
	if entry.Mode[0] == 'd' {
		if strings.HasSuffix(entry.Mode, "t") || strings.HasSuffix(entry.Mode, "T") {
			return entry, "sticky"
		}
		return entry, "dir"
	}
	// Handle device files
	if entry.Mode[0] == 'b' || entry.Mode[0] == 'c' {
		return entry, "dev"
	}

	if strings.ContainsAny(mod, "x") {
		return entry, "exec"
	}
	if strings.Count(mod, "w") == 3 {
		return entry, "worldw-ritable"
	}
	name := strings.Trim(entry.Name, "'")
	tokens := strings.Split(strings.ToLower(name), ".")
	ext := "." + tokens[len(tokens)-1]
	switch ext {
	case ".log":
		return entry, "text"
	case ".pdf":
		return entry, "pdf"
	case ".rar", ".zip", ".tar", ".gz", ".7z", "deb":
		return entry, "archive"
	case ".mp3", ".wav", ".flac":
		return entry, "audio"
	case ".webp", ".jpg", ".jpeg", ".png", ".gif":
		return entry, "image"
	case ".crdownload":
		return entry, "crd"
	default:
		return entry, "other"
	}
}

// This function resolves the target of symbolic links by concatenating the parent directory path and the symbolic link path.
// It takes the parent directory path and the symbolic link path as input,
// and returns the resolved target path.
func getLinkTargetType(path, s string) (Entry, error) {
	absPath := resolveRelativePath(path, s)
	if strings.HasPrefix(s, "../share") && !strings.HasPrefix(absPath, "/usr") {
		absPath = "/usr" + absPath
	}
	info, err := os.Lstat(absPath)
	var newEntry Entry
	if os.IsNotExist(err) {
		return newEntry, fmt.Errorf("target not found: %s", s)
	} else if err != nil {
		return newEntry, fmt.Errorf("error getting target info: %s", err.Error())
	} else {
		permissions, err := formatPermissionsWithACL(absPath, info.Mode())
		if err != nil {
			return newEntry, fmt.Errorf("cannot format permissions: %s", err.Error())
		}
		newEntry = Entry{Name: s, Mode: permissions, Path: absPath}
	}
	return newEntry, nil
}
