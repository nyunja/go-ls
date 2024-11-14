package lsfunctions

import (
	"fmt"
	"io"
	"os"
)

func DisplayLongFormat(w io.Writer, entries []FileInfo) {
	t := getTotalBlocks(entries)
	if ShowTotals {
		fmt.Printf("total %d\n", t)
	}
	newEntries, widths := processEntries(entries)
	for _, entry := range newEntries {
		fmt.Fprintln(w, GetLongFormatString2(entry, widths))
	}
}

func GetLongFormatString2(e Entry, w Widths) string {
	e = colorName(e, false)
	s := ""
	if w.minorCol == 0 {
		s = fmt.Sprintf("%-*s %*s %-*s %-*s %*s %*s  %s", w.modCol, e.Mode, w.linkCol, e.LinkCount, w.ownerCol, e.Owner, w.groupCol, e.Group, w.sizeCol, e.Size, w.timeCol, e.Time, e.Name)
	} else {
		s = fmt.Sprintf("%-*s %*s %-*s %-*s %*s %*s %*s  %s", w.modCol, e.Mode, w.linkCol, e.LinkCount, w.ownerCol, e.Owner, w.groupCol, e.Group, w.minorCol, e.Minor, w.sizeCol, e.Size, w.timeCol, e.Time, e.Name)
	}
	if e.Mode[0] == 'l' && e.LinkTarget != "" {
		s += " -> " + formatLinkTarget(e.LinkTarget)
	}
	return s
}

func formatLinkTarget(s string) string {
	info, err := os.Lstat(s)
	if err != nil {
		return green + s + reset
	}
	mod := info.Mode()
	fmt.Printf("%q\n", mod)
	return s
	// switch mod[0] {
	// case 'g':
	// 	return yellowBackground + "\033[38;3;40m" + s + reset
	// case 'u':
	// 	return orangeBackground + s + reset
	// default:
	// 	return green + s + reset
	// }
}
