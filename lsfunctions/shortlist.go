package lsfunctions

import (
	"fmt"
	"strings"
)

func DisplayShortList(entries []FileInfo) {
	var fileNameEntries []string
	displayColumn := false

	// Prepare the list of file names and check if any entry is too long
	for _, entry := range entries {
		formattedName := GetShortFormat(entry)
		if len(formattedName) > 80 {
			displayColumn = true
		}
		fileNameEntries = append(fileNameEntries, formattedName)
	}

	// Display in single-column format if any entry is too long
	if displayColumn {
		for _, name := range fileNameEntries {
			fmt.Println(name)
		}
		return
	}

	// For fewer than or equal to 8 items, display in a single line
	if len(fileNameEntries) <= 8 {
		fmt.Println(strings.Join(fileNameEntries, "  "))
		return
	}

	// Arrange items into columns
	const itemsPerRow = 4
	var columns [][]string
	for i := 0; i < len(fileNameEntries); i += itemsPerRow {
		end := i + itemsPerRow
		if end > len(fileNameEntries) {
			end = len(fileNameEntries)
		}
		columns = append(columns, fileNameEntries[i:end])
	}

	// Calculate the maximum width of each column
	columnWidths := make([]int, itemsPerRow)
	for _, row := range columns {
		for i, name := range row {
			if len(name) > columnWidths[i] {
				columnWidths[i] = len(name)
			}
		}
	}

	// Print each row with aligned columns
	for _, row := range columns {
		for i, name := range row {
			padding := columnWidths[i] - len(name) + 2 // Add extra space between columns
			fmt.Print(name + strings.Repeat(" ", padding))
		}
		fmt.Println()
	}
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
