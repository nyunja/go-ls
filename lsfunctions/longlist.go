package lsfunctions

import (
	"fmt"
	"io"
)

// DisplayLongFormat displays the entries in long format.
func DisplayLongFormat(w io.Writer, entries []FileDetails) {
	t := getTotalBlocks(entries)
	if ShowTotals {
		fmt.Printf("total %d\n", t)
	}
	formattedEntries := prepareFileDetailsForDisplay(entries)
	widths := getWidths(formattedEntries)
	for _, entry := range formattedEntries {
		fmt.Fprintln(w, getLongFormatString(entry, widths))
	}
}

// getLongFormatString returns the long format string with the given entry and widths.
// The width of each column is determined by the maximum width of the corresponding column in the input entries.
// If the entry is a symbolic link, the link target is also displayed in color.
func getLongFormatString(e Entry, w Widths) string {
	e = colorName(e, false)
	s := ""
	if w.minorCol == 0 {
		s = fmt.Sprintf("%-*s %*s %-*s %-*s %*s %*s  %s", w.modCol, e.Mode, w.linkCol, e.LinkCount, w.ownerCol, e.Owner, w.groupCol, e.Group, w.sizeCol, e.Size, w.timeCol, e.Time, e.Name)
	} else {
		s = fmt.Sprintf("%-*s %*s %-*s %-*s %*s %*s %*s  %s", w.modCol, e.Mode, w.linkCol, e.LinkCount, w.ownerCol, e.Owner, w.groupCol, e.Group, w.minorCol, e.Minor, w.sizeCol, e.Size, w.timeCol, e.Time, e.Name)
	}
	if e.Mode[0] == 'l' && e.LinkTarget != "" {
		s += " -> " + colorLinkTarget(e.Path, e.LinkTarget)
	}
	return s
}
