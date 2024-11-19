package lsfunctions

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

// prepareFileDetailsForDisplay converts a list of FileDetails into a list of Entry.
// If the file is a device file, is also gets the major and minor device numbers from the Rdev field.
// It also adds quotes around the file names if they contain spaces.
// It also gets the user and group names from the owner and group fields.
func prepareFileDetailsForDisplay(entries []FileDetails) []Entry {
	var formattedEntries []Entry
	for _, entry := range entries {
		var f Entry
		info := entry.Info
		mode := info.Mode()
		f.Name = addQuotes(entry.Name)
		f.Mode, _ = formatPermissionsWithACL(entry.Path, mode)
		f.Path = entry.Path
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
		}
		if u, err := user.LookupId(owner); err == nil {
			f.Owner = u.Username
		}
		if g, err := user.LookupGroupId(group); err == nil {
			f.Group = g.Name
		}

		formattedEntries = append(formattedEntries, f)
	}
	return formattedEntries
}

// getWidths calculates the maximum width for each column in the long format output.
// It considers the mode, link count, owner, group, size, minor, and time columns.
func getWidths(entries []Entry) Widths {
	var w Widths
	for _, f := range entries {
		w.modCol = getMax(w.modCol, len(f.Mode))
		w.groupCol = getMax(w.groupCol, len(f.Group))
		w.ownerCol = getMax(w.ownerCol, len(f.Owner))
		w.sizeCol = getMax(w.sizeCol, len(f.Size))
		w.minorCol = getMax(w.minorCol, len(f.Minor))
		w.timeCol = getMax(w.timeCol, len(f.Time))
		w.linkCol = getMax(w.linkCol, len(f.LinkCount))
	}
	return w
}
