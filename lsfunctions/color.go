package lsfunctions

import "strings"

func colorName(entry Entry) Entry {
	// Define color codes
	const (
		reset      = "\033[0m"
		boldBlue   = "\033[38;5;01;34m"
		boldYellow = "\033[38;5;01;33m"
		red        = "\033[31m"
		cyan       = "\033[36m"
		green      = "\033[38;5;46m"
	)
	// Color map for different file types
	colors := map[string]string{
		"text":       "\x1b[97m",
		"pdf":        "\x1b[91m", // Light Red
		"word":       "\x1b[94m", // Light Blue
		"excel":      "\x1b[92m", // Light Green
		"powerpoint": "\x1b[93m", // Light Yellow
		"archive":    red,
		"audio":      "\x1b[96m", // Light Cyan
		"video":      "\x1b[95m", // Light Magenta
		"image":      "\x1b[35m", // Magenta
		"go":         cyan,
		"python":     "\x1b[33m", // Yellow
		"javascript": "\x1b[33m", // Yellow
		"html":       "\x1b[91m", // Light Red
		"css":        cyan,
		"exec":       green,
	}

	// Handle symbolic links
	if entry.Mode[0] == 'l' {
		entry.Name = "\033[38;5;01;34m" + entry.Name + "\033[0m"
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
		entry.Name = boldYellow + entry.Name + reset
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
