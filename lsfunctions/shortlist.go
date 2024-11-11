package lsfunctions

import (
	"fmt"
	"io"
	"strings"
)

func DisplayShortList(w io.Writer, e []FileInfo) {
	// Process entries to create type []Entry
	entries, _ := processEntries(e)
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
			fmt.Fprint(w, name + strings.Repeat(" ", padding))
		}
		fmt.Fprintln(w)
	}
}

func GetShortFormat(entry Entry) string {
	// Color output
	entry = colorName(entry)

	return entry.Name
}

