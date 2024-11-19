package lsfunctions

import (
	"strings"
)

// Define color codes
const (
	reset            = "\033[0m"
	turqoise         = "\033[1;38;2;42;161;179m"
	orangeBackground = "\033[48;2;192;28;20m"
	boldBlue         = "\033[1;38;5;01;34m"
	yellow           = "\033[38;2;162;115;76m"
	boldYellow       = "\033[1;38;2;162;115;76m"
	yellowBackground = "\033[48;2;162;115;76m"
	blackBack        = "\033[40m"
	blackText        = "\033[30m"
	red              = "\033[1;31m"
	cyan             = "\033[1;36m"
	green            = "\033[1;38;2;39;169;105m"
	greenBackground  = "\033[42m"
	magentaBold      = "\033[1;35m"
	socket           = "\033[1;38;2;163;71;181m"
)

// colorName applies color codes to different file types and symbolic links.
// Returns the colorized name of the file or symbolic link.
// If the entry is a symbolic link and the isTarget flag is true, the link target is also colorized.
func colorName(entry Entry, isTarget bool) Entry {
	// Color map for different file types
	colors := map[string]string{
		"world-writable": magentaBold,
		"pipe":           yellow + blackBack,
		"socket":         socket,
		"sticky":         greenBackground + blackText,
		"setuid":         orangeBackground,
		"setgid":         yellowBackground + blackText,
		"dir":            boldBlue,
		"dev":            boldYellow + blackBack,
		"archive":        red,
		"audio":          "\x1b[1;96m", // Light Cyan
		"image":          "\x1b[1;35m", // Magenta
		"crd":            "\x1b[1;38;5;8m",
		"css":            cyan,
		"exec":           green,
	}
	// Handle target links that are symbolic links
	if isTarget {
		if strings.Split(entry.Mode, "-")[0] == "lrwx" {
			entry.Name = addColorAndPadding(boldYellow + blackBack, entry.Name, reset)
			return entry
		}
	}
	// Handle symbolic links
	if (entry.Mode[0] == 'l' || entry.Mode[0] == 'L') && !isTarget {
		entry.Name = addColorAndPadding(turqoise, entry.Name, reset)
		return entry
	}

	// Color code based on file type
	entry, fileType := getFileType(entry)
	if color, exists := colors[fileType]; exists {
		entry.Name = addColorAndPadding(color, entry.Name, reset)
	}

	return entry
}

// addColorAndPadding adds color codes and padding to the file or symbolic link name.
// The color and reset codes are applied to the name using ANSI escape sequences.
// The padding is added to the end of the name to ensure consistent formatting.
// It returns the colored and padded name.
func addColorAndPadding(color, name, reset string) string {
	originalsize := len(name)
	name = strings.TrimSpace(name)
	padding := originalsize - len(name)
	name = color + name + reset
	name = name + strings.Repeat(" ", padding)
	return name
}

// This function resolves the target of symbolic links and applies color formatting to the link target.
// It takes the path to the parent directory and the symbolic link path as input,
// and returns the colored link target.
func colorLinkTarget(path, s string) string {
	newEntry, err := getLinkTargetType(path, s)
	if err != nil {
		return s
	}
	colorEntry := colorName(newEntry, true)

	return colorEntry.Name
}
