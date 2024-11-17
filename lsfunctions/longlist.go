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
	// colors := map[string]string{
	// 	"setuid":  orangeBackground,
	// 	"setgid":  yellowBackground + blackText,
	// 	"dir":     boldBlue,
	// 	"dev":     boldYellow + blackBack,
	// 	"archive": red,
	// 	"audio":   "\x1b[1;96m", // Light Cyan
	// 	"image":   "\x1b[1;35m", // Magenta
	// 	"crd":     "\x1b[1;38;5;8m",
	// 	"css":     cyan,
	// 	"exec":    green,
	// 	"t": "\033[38;2;38;162;105m",
	// }
	// if file, ok := specialFiles[s]; ok {
	// 	if color, exists := colors[file]; exists {
	// 		return color + s + reset
	// 	}
	// }
	absPath := resolveRelativePath(path, s)
	info, err := os.Lstat(absPath)
	var newEntry Entry
	if os.IsNotExist(err) {
		// debug broken link
		fmt.Printf("Broken link: %s -> %s\n", path, s)
		return "link is broken: " + err.Error()
	} else if err != nil {
		return "cannot access link target: " + err.Error()
	} else {
		permissions, err := formatPermissionsWithACL(absPath, info.Mode())
		if err!= nil {
            return "cannot format permissions: " + err.Error()
        }
		newEntry = Entry{Name: s, Mode: permissions, Path: absPath}
		// return "target link exists and is of type: " + info.Mode().String()
	}
	// newEntry := Entry{Name: s, Mode: info.Mode().String()}
	colorEntry := colorName(newEntry, true)
	// entry, fileType := getFileType(newEntry)

	// if color, exists := colors[fileType]; exists {
	// 	return color + entry.Name + reset
	// }
	return colorEntry.Name
}

func resolveRelativePath(linkPath, target string) string {
	// linkDir, err := os.Getwd()
	// if err!= nil {
    //     fmt.Println("unable to get current working directory: ", err.Error())
    // }
	// fmt.Println("working directory: ", linkDir)
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

	// linkPath = strings.Join(stack, "/")
	// if strings.HasPrefix(linkPath, "..") {
	// 	return strings.Join(stack, "/")
	// }
	// Join normalized segments
	return "/" + strings.Join(stack, "/")
}
