package lsfunctions

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// var specialFiles = map[string]string{
// 	"../lib/snapd/snap-confine": "setuid",
// 	"ssh_agent":                 "setgid",
// 	"umount":                    "setuid",
// 	"su":                        "setuid",
// 	"sudo":                      "setuid",
// 	"passwd":                    "setuid",
// 	"gpasswd":                   "setuid",
// 	"fusermount3":               "setuid",
// 	"newgrp":                    "setgid",
// 	"mount":                     "setuid",
// 	"newuidmap":                 "setuid",
// 	"newgidmap":                 "setuid",
// 	"umount2":                   "setuid",
// 	"expiry":                    "setgid",
// 	"chsh":                      "setuid",
// 	"chfn":                      "setuid",
// 	"chage":                     "setgid",
// 	"/proc/self/fd":             "dir",
// }

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
		s += " -> " + colorLinkTarget(e.Path, e.LinkTarget)
	}
	return s
}

func colorLinkTarget(path, s string) string {
	absPath := resolveRelativePath(path, s)
	// Handle shared folder link targets
	if strings.HasPrefix(s, "../share") {
		absPath = "/usr" + absPath
	}
	info, err := os.Lstat(absPath)
	var newEntry Entry
	if os.IsNotExist(err) {
		return s
	} else if err != nil {
		return s
	} else {
		permissions, err := formatPermissionsWithACL(absPath, info.Mode())
		if err!= nil {
            return "cannot format permissions: " + err.Error()
        }
		newEntry = Entry{Name: s, Mode: permissions, Path: absPath}
	}
	colorEntry := colorName(newEntry, true)

	return colorEntry.Name
}

func resolveRelativePath(linkPath, target string) string {
	linkDir := ""
	lastSlash := strings.LastIndex(linkPath, "/")
	if lastSlash != -1 {
		linkDir = linkPath[:lastSlash]
	}
	// Return path as is if it's already absolute
	if strings.HasPrefix(target, "/") {
		return target
	}
	joinedPath := joinPath(linkDir, target)

	// Normalize combined path
	segments := strings.Split(joinedPath, "/")
	var stack []string

	for _, segment := range segments {
		if segment == ".." {
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		} else if segment != "." && segment != "" {
			stack = append(stack, segment)
		}
	}
	return "/" + strings.Join(stack, "/")
}
