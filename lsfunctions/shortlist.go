package lsfunctions

import (
	"fmt"
)

func DisplayShortList(entries []FileInfo) {
	nameCol := getMaxWidth(entries)
	for _, entry := range entries {
		fmt.Printf("%-*s", nameCol, GetShortFormat(entry))
	}
	fmt.Println()
}

func getMaxWidth(entries []FileInfo) int {
	nameCol := 0
	for _, entry := range entries {
		if len(entry.Name) > nameCol {
			nameCol = len(entry.Name)
		}
	}
	return nameCol
}

func GetShortFormat(entry FileInfo) string {
	name := entry.Name
	info := entry.Info
	// Add color blue for directories
	if info.IsDir() {
		return "\x1b[34m" + name + "\x1b[0m"
	}
	// Color coding based on file type
	switch getFileType(entry) {
	case "text":
		name = "\x1b[97m" + name + "\x1b[0m" // White
	case "pdf":
		name = "\x1b[91m" + name + "\x1b[0m" // Light Red
	case "word":
		name = "\x1b[94m" + name + "\x1b[0m" // Light Blue
	case "excel":
		name = "\x1b[92m" + name + "\x1b[0m" // Light Green
	case "powerpoint":
		name = "\x1b[93m" + name + "\x1b[0m" // Light Yellow
	case "archive":
		name = "\x1b[31m" + name + "\x1b[0m" // Red
	case "audio":
		name = "\x1b[96m" + name + "\x1b[0m" // Light Cyan
	case "video":
		name = "\x1b[95m" + name + "\x1b[0m" // Light Magenta
	case "image":
		name = "\x1b[35m" + name + "\x1b[0m" // Magenta
	case "go":
		name = "\x1b[36m" + name + "\x1b[0m" // Cyan
	case "python":
		name = "\x1b[33m" + name + "\x1b[0m" // Yellow
	case "javascript":
		name = "\x1b[33m" + name + "\x1b[0m" // Yellow
	case "html":
		name = "\x1b[91m" + name + "\x1b[0m" // Light Red
	case "css":
		name = "\x1b[36m" + name + "\x1b[0m" // Cyan
	case "link":
		name = "\x1b[38;5;51m" + name + "\x1b[0m"
	case "exec":
		name = "\x1b[38;5;46m" + name + "\x1b[0m" // Add color green for executables
	}
	return name
}
