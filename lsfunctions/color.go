package lsfunctions

import "strings"

// Define color codes
const (
	reset            = "\033[0m"
	turqoise         = "\033[1;38;2;42;161;179m"
	orangeBackground = "\033[48;2;192;28;20m"
	boldBlue         = "\033[1;38;5;01;34m"
	boldYellow       = "\033[1;38;2;162;115;76m"
	yellowBackground = "\033[48;2;162;115;76m"
	blackBack        = "\033[40m"
	blackText        = "\033[30m"
	red              = "\033[1;31m"
	cyan             = "\033[1;36m"
	green            = "\033[1;38;2;39;169;105m"
	greenBackground = "\033[42m"
)

func colorName(entry Entry, isTarget bool) Entry {
	// Color map for different file types
	colors := map[string]string{
		"setuid":  orangeBackground,
		"setgid":  yellowBackground + blackText,
		"dir":     boldBlue,
		"dev":     boldYellow + blackBack,
		"archive": red,
		"audio":   "\x1b[1;96m", // Light Cyan
		"image":   "\x1b[1;35m", // Magenta
		"crd":     "\x1b[1;38;5;8m",
		"css":     cyan,
		"exec":    green,
	}

	// Handle symbolic links
	if entry.Mode[0] == 'l' || entry.Mode[0] == 'l' && !isTarget {
		originalsize := len(entry.Name)
		entry.Name = strings.TrimSpace(entry.Name)
		padding := originalsize - len(entry.Name)
		entry.Name = turqoise + entry.Name + reset
		entry.Name = entry.Name + strings.Repeat(" ", padding)
		return entry
	}

	// Color code based on file type
	entry, fileType := getFileType(entry)
	if color, exists := colors[fileType]; exists {
		originalsize := len(entry.Name)
		entry.Name = strings.TrimSpace(entry.Name)
		padding := originalsize - len(entry.Name)

		entry.Name = color + entry.Name + reset
		entry.Name = entry.Name + strings.Repeat(" ", padding)
	}

	return entry
}
