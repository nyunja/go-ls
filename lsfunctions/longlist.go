package lsfunctions

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func major(dev uint64) uint64 {
	return (dev >> 8) & 0xff
}

func minor(dev uint64) uint64 {
	return dev & 0xff
}

func DisplayLongFormat(entries []FileInfo) {
	t := getTotalBlocks(entries)
	if ShowTotals {
		fmt.Printf("total %d\n", t)
	}
	// widths, ugl := getColumnWidth(entries)
	newEntries, w := processEntries(entries)
	for _, entry := range newEntries {
		fmt.Println(GetLongFormatString2(entry, w))
	}
	// for _, entry := range entries {
	// fmt.Println(GetLongFormatString(entry, widths, ugl))
	// }
}

func formatName(s string) string {
	if strings.Contains(s, " ") {
		if strings.ContainsAny(s, "'") {
			s = fmt.Sprintf(`"%s"`, s)
		} else {
			s = fmt.Sprintf(`'%s'`, s)
		}
	}
	return s
}

func swapT(s string) string {
	if len(s) <= 1 {
		return s
	}
	res := []rune{}
	for i, ch := range s {
		if ch == 't' {
			continue
		} else if i == len(s)-1 {
			res = append(res, 't')
		} else {
			res = append(res, ch)
		}
	}
	return string(res)
}

func getTotalBlocks(entries []FileInfo) TotalBlocks {
	var t TotalBlocks
	for _, entry := range entries {
		if stat, ok := entry.Info.Sys().(*syscall.Stat_t); ok {
			t += TotalBlocks(stat.Blocks)
		}
	}
	return t / 2
}



func getFileType(entry Entry) string {
	mod := entry.Mode

	if strings.ContainsAny(mod, "x") {
		return "exec"
	}
	name := entry.Name
	tokens := strings.Split(strings.ToLower(name), ".")
	ext := tokens[len(tokens)-1]
	switch ext {
	case ".txt", ".md", ".log":
		return "text"
	case ".pdf":
		return "pdf"
	case ".doc", ".docx":
		return "word"
	case ".xls", ".xlsx":
		return "excel"
	case ".ppt", ".pptx":
		return "powerpoint"
	case ".zip", ".tar", ".gz", ".7z", "deb":
		return "archive"
	case ".mp3", ".wav", ".flac":
		return "audio"
	case ".mp4", ".avi", ".mkv":
		return "video"
	case ".jpg", ".jpeg", ".png", ".gif":
		return "image"
	case ".py":
		return "python"
	case ".js":
		return "javascript"
	case ".html", ".htm":
		return "html"
	case ".css":
		return "css"
	default:
		return "other"
	}
}

func processEntries(entries []FileInfo) ([]Entry, Widths) {
	var newEntries []Entry
	var w Widths
	for _, entry := range entries {
		var f Entry
		info := entry.Info
		mode := info.Mode()
		f.Name = formatName(entry.Name)
		f.Mode = mode.String()
		if strings.HasPrefix(f.Mode, "L") {
			f.Mode = "l" + f.Mode[1:]
		}
		f.IsDirectory = info.IsDir()
		f.LinkTarget = entry.LinkTarget
		f.Time = formatTime(info.ModTime())
		// Get size string
		f.Size = fmt.Sprintf("%d", info.Size())
		if mode&os.ModeDevice != 0 {
			if stat, ok := entry.Info.Sys().(*syscall.Stat_t); ok {
				entry.Rdev = stat.Rdev
			}
			major := major(entry.Rdev)
			minor := minor(entry.Rdev)
			f.Size = fmt.Sprintf("%d", major)
			f.Minor = fmt.Sprintf("%d,", minor)
		}
		var owner, group string

		if stat, ok := info.Sys().(*syscall.Stat_t); ok {
			uid := stat.Uid
			gid := stat.Gid
			f.LinkCount = fmt.Sprintf("%d", stat.Nlink)
			owner = strconv.FormatUint(uint64(uid), 10)
			group = strconv.FormatUint(uint64(gid), 10)
		} else {
			fmt.Printf("error getting syscall info")
			// return widths, ugl
		}
		if u, err := user.LookupId(owner); err == nil {
			f.Owner = u.Username
		}
		if g, err := user.LookupGroupId(group); err == nil {
			f.Group = g.Name
		}
		f = colorName(f)
		w.modCol = getMax(w.modCol, len(f.Mode))
		w.groupCol = getMax(w.groupCol, len(f.Group))
		w.ownerCol = getMax(w.ownerCol, len(f.Owner))
		w.sizeCol = getMax(w.sizeCol, len(f.Size))
		w.minorCol = getMax(w.minorCol, len(f.Minor))
		w.timeCol = getMax(w.timeCol, len(f.Time))
		w.linkCol = getMax(w.linkCol, len(f.LinkCount))
		newEntries = append(newEntries, f)
	}
	return newEntries, w
}

func getMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func GetLongFormatString2(e Entry, w Widths) string {
	s := ""
	if w.minorCol == 0 {
		s = fmt.Sprintf("%-*s %*s %-*s %-*s %*s %*s %s", w.modCol, e.Mode, w.linkCol, e.LinkCount, w.ownerCol, e.Owner, w.groupCol, e.Group, w.sizeCol, e.Size, w.timeCol, e.Time, e.Name)
	} else {
		s = fmt.Sprintf("%-*s %*s %-*s %-*s %*s %*s %*s %s", w.modCol, e.Mode, w.linkCol, e.LinkCount, w.ownerCol, e.Owner, w.groupCol, e.Group, w.minorCol, e.Minor, w.sizeCol, e.Size, w.timeCol, e.Time, e.Name)
	}
	if e.Mode[0] == 'l' && e.LinkTarget != "" {
		s += " -> " + e.LinkTarget
	}
	return s
}

// formatTime formats a given time based on whether it's in the current year or not.
// For times in the current year, it returns the format "Jan _2 15:04".
// For times in previous years, it returns the format "Jan _2 2006".
//
// Parameters:
//   - modTime: A time.Time value representing the modification time to be formatted.
//
// Returns:
//   - string: A formatted string representation of the input time.
func formatTime(modTime time.Time) string {
	now := time.Now()
	if modTime.Year() == now.Year() {
		return modTime.Format("Jan _2 15:04")
	}
	return modTime.Format("Jan _2  2006")
}
