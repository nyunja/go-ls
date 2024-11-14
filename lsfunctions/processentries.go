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

func processEntries(entries []FileInfo) ([]Entry, Widths) {
	var newEntries []Entry
	var w Widths
	for _, entry := range entries {
		var f Entry
		info := entry.Info
		mode := info.Mode()
		f.Name = addQuotes(entry.Name)
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

// swapT modifies a string by removing all 't' characters except for the last one,
// which is moved to the end of the string if it's not already there.
//
// Parameters:
//   - s: The input string to be modified.
//
// Returns:
//   - string: The modified string with 't' characters removed and potentially one 't' added at the end.
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

func swapU(s string) string {
	if len(s) <= 1 {
		return s
	}
	s = strings.Replace(s, "u", "-",1)
	return strings.Replace(s, "x", "s", 1)
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

func addQuotes(s string) string {
	if strings.Contains(s, " ") || hasSpecialChar(s){
		s = fmt.Sprintf(`'%s'`, s)
	}
	return s
}

func hasSpecialChar(s string) bool {
	if len(s) == 0 {
        return false
    }
	specialChars := []rune{ '[','#', ']', '{', '}', '|', '\\', ':', ';', '<', '>', ',', '?', '!', '@', '$', '%', '^', '&', '*', '(', ')', '~', '`', '"', '\'', '=', '+'}
	for _, ch := range specialChars {
        if strings.ContainsRune(s, ch) {
            return true
        }
    }
    return false
}

func major(dev uint64) uint64 {
	return (dev >> 8) & 0xff
}

func minor(dev uint64) uint64 {
	return dev & 0xff
}
