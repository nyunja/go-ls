package lsfunctions

import (
	"fmt"
	"io"
	"os"
)


var specialFiles = map[string]string{
	"../lib/snapd/snap-confine": "setuid",
	"ssh_agent": "setgid",
	"umount": "setuid",
	"su": "setuid",
	"sudo": "setuid",
	"passwd": "setuid",
	"gpasswd": "setuid",
	"fusermount3": "setuid",
    "newgrp": "setgid",
    "mount": "setuid",
	"newuidmap": "setuid",
	"newgidmap": "setuid",
    "umount2": "setuid",
	"expiry": "setgid",
	"chsh": "setuid",
	"chfn": "setuid",
	"chage": "setgid",
	"/proc/self/fd": "dir",
}

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
		s += " -> " + colorLinkTarget(e.LinkTarget)
	}
	return s
}

func colorLinkTarget(s string) string {
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
		"t": "\033[38;2;38;162;105m",
	}
	if file, ok := specialFiles[s]; ok {
		if color, exists := colors[file]; exists {
			return color + s + reset
		}
	}
	info, err := os.Lstat(s)
	if err != nil {
		return green + s + reset
	}
	newEntry := Entry{Name: s, Mode: info.Mode().String()}
	// colorEntry := colorName(newEntry, true)
	entry, fileType := getFileType(newEntry)

	if color, exists := colors[fileType]; exists {
		return color + entry.Name + reset
	}
	return newEntry.Name
}
