package lsfunctions

import "strings"

func colorName(entry Entry) Entry {
	// Define color codes
	const (
		reset      = "\033[0m"
		turqoise = "\033[1;38;2;42;161;179m"
		orangeBackground = "\033[48;2;192;28;20m"
		boldBlue   = "\033[1;38;5;01;34m"
		boldYellow = "\033[1;38;2;162;115;76m"
		blackBack = "\033[40m"
		red        = "\033[1;31m"
		cyan       = "\033[1;36m"
		green      = "\033[1;38;2;39;169;105m"
	)
	// Color map for different file types
	colors := map[string]string{
		"archive":    red,
		"audio":      "\x1b[1;96m", // Light Cyan
		"image":      "\x1b[1;35m", // Magenta
		"crd":       "\x1b[1;38;5;8m",
		"css":        cyan,
		"exec":       green,
	}

	// Handle symbolic links
	if entry.Mode[0] == 'l' {
		entry.Name = turqoise+ entry.Name + reset
		return entry
	}
	// Handle setuid files
	if entry.Mode[0] == 'u' {
		entry.Mode = swapU(entry.Mode)
		entry.Name = orangeBackground + entry.Name + reset
		return entry
	}
	// Handle directories
	if entry.Mode[0] == 'd' {
		entry.Name = boldBlue + entry.Name + reset
		if strings.Contains(entry.Mode, "t") {
			entry.Mode = swapT(entry.Mode)
		}
		return entry
	}
	// Handle device files
	if entry.Mode[0] == 'D' {
		entry.Mode = "b" + strings.TrimPrefix(entry.Mode, "D")
		entry.Name = boldYellow + blackBack + entry.Name + reset
		return entry
	}

	// Handle device files
	if entry.Mode[0] == 'D' {
		entry.Mode = "b" + strings.TrimPrefix(entry.Mode, "D")
		entry.Name = boldYellow + entry.Name + reset
		return entry
	}

	// Ensure mode is properly formatted
	if len(entry.Mode) != 10 {
		entry.Mode = "b" + entry.Mode
		if len(entry.Mode) != 10 {
			entry.Mode += "-"
		}
	}

	// Color code based on file type
	fileType := getFileType(entry)
	if color, exists := colors[fileType]; exists {
		entry.Name = color + entry.Name + reset
	}

	return entry
}
