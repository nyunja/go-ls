package lsfunctions

import (
	"strings"
)

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
