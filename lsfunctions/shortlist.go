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
	maxwith := 0

	for i := 0; i < len(entries); i++ {
		if len(entries[i].Name) > maxwith {
			maxwith = len(entries[i].Name)
		}
	}

	for i := 0; i < len(entries); i++ {
		entries[i].Name = entries[i].Name + strings.Repeat(" ", (maxwith-len(entries[i].Name)+5))
	}

	// Prepare the list of file names and check if any entry is too long
	for _, entry := range entries {
		formattedName := GetShortFormat(entry)
		if len(formattedName) > 100 {
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

	// Print each row with aligned columns
	for _, row := range columns {
		for _, name := range row {
			fmt.Fprint(w, name)
		}
		fmt.Fprintln(w)
	}
}

func GetShortFormat(entry Entry) string {
	// Color output
	entry = colorName(entry, false)

	return entry.Name
}
