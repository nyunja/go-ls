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

// FileInfo struct to store file information from readDir function
type FileInfo struct {
	Name string
	Info os.FileInfo
}

func DisplayLongFormat(entries []FileInfo) {
	var totalBlocks int64
	for _, entry := range entries {
		if stat, ok := entry.Info.Sys().(*syscall.Stat_t); ok {
			totalBlocks += stat.Blocks
		}
	}
	fmt.Printf("total %d\n", totalBlocks/2)
	sizeCol, ownerCol, groupCol, linkCol, timeCol, modCol := getColumnWidth(entries)
	for _, entry := range entries {
		fmt.Println(GetLongFormatString(entry, sizeCol, ownerCol, groupCol, linkCol, timeCol, modCol))
	}
}

func GetLongFormatString(entry FileInfo, sizeCol, ownerCol, groupCol, linkCol, timeCol, modCol int) string {
	info := entry.Info
	mode := info.Mode()
	modeStr := mode.String()
	size := info.Size()
	modTime := info.ModTime()
	name := entry.Name
	if strings.Contains(name, " ") {
		name = "'" + name + "'"
	}
	// Add color blue for directories
	if info.IsDir() {
		name = "\x1b[34m" + name + "\x1b[0m"
	}
	// Color coding based on file type
	switch getFileType(entry) {
	case "text":
		name = "\x1b[97m" + name + "\x1b[0m" // White
	case "pdf":
		name = "\x1b[91m" + name + "\x1b[0m" // Light Red
	case "word":
		name = "\x1b[94m" + name + "\x1b[0m" // Light Blue
	case "excel":
		name = "\x1b[92m" + name + "\x1b[0m" // Light Green
	case "powerpoint":
		name = "\x1b[93m" + name + "\x1b[0m" // Light Yellow
	case "archive":
		name = "\x1b[31m" + name + "\x1b[0m" // Red
	case "audio":
		name = "\x1b[96m" + name + "\x1b[0m" // Light Cyan
	case "video":
		name = "\x1b[95m" + name + "\x1b[0m" // Light Magenta
	case "image":
		name = "\x1b[35m" + name + "\x1b[0m" // Magenta
	case "go":
		name = "\x1b[36m" + name + "\x1b[0m" // Cyan
	case "python":
		name = "\x1b[33m" + name + "\x1b[0m" // Yellow
	case "javascript":
		name = "\x1b[33m" + name + "\x1b[0m" // Yellow
	case "html":
		name = "\x1b[91m" + name + "\x1b[0m" // Light Red
	case "css":
		name = "\x1b[36m" + name + "\x1b[0m" // Cyan
	case "link":
		name = "\x1b[38;5;100m" + name + "\x1b[0m"
		if strings.HasPrefix(mode.String(), "L") {
			modeStr = "l" + modeStr[1:]
		}
		target, err := os.Readlink(info.Name()) // Get target link
		if err == nil {
			_, err := os.Stat(target)
			if err == nil {
				name += " -> " + target
			} else {
				name += " -> " + target + " (Broken link)"
			}
		}
	case "exec":
		name = "\x1b[32m" + name + "\x1b[0m" // Add color green for executables
	}
	
	// Format symbolic links
	// if mode&fs.ModeSymlink != 0 {
	// 	target, err := os.Readlink(info.Name()) // Get target link
	// 	if err == nil {
	// 		_, err := os.Stat(target)
	// 		if err == nil {
	// 			name += " -> " + target
	// 		} else {
	// 			name += " -> " + target + " (Broken link)"
	// 		}
	// 	}
	// }
	

	var owner, group string
	var linkCount uint64

	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		uid := stat.Uid
		gid := stat.Gid
		linkCount = stat.Nlink
		owner = strconv.FormatUint(uint64(uid), 10)
		group = strconv.FormatUint(uint64(gid), 10)
	} else {
		fmt.Printf("error getting syscall info")
		return ""
	}
	if u, err := user.LookupId(owner); err == nil {
		owner = u.Username
	}
	if g, err := user.LookupGroupId(group); err == nil {
		group = g.Name
	}

	timeString := formatTime(modTime)

	sizeStr := toString(size)

	s := fmt.Sprintf("%-*s %*d %-*s %-*s %*s %*s  %s", modCol, modeStr, linkCol, linkCount, ownerCol, owner, groupCol, group, sizeCol, sizeStr, timeCol, timeString, name)
	return s
}

func getColumnWidth(entries []FileInfo) (int, int, int, int, int, int) {
	var owner, group string
	var linkCount uint64
	sizeCol, groupCol, ownerCol, linkCol, timeCol, modCol := 0, 0, 0, 0, 0, 0

	for _, entry := range entries {
		info := entry.Info
		if stat, ok := info.Sys().(*syscall.Stat_t); ok {
			uid := stat.Uid
			gid := stat.Gid
			linkCount = stat.Nlink
			owner = strconv.FormatUint(uint64(uid), 10)
			group = strconv.FormatUint(uint64(gid), 10)
		} else {
			fmt.Printf("error getting syscall info")
			return sizeCol, ownerCol, groupCol, linkCol, timeCol, modCol
		}
		if u, err := user.LookupId(owner); err == nil {
			owner = u.Username
		}
		if g, err := user.LookupGroupId(group); err == nil {
			group = g.Name
		}

		modStr := info.Mode().String()
		if len(modStr) > modCol {
			modCol = len(modStr)
		}

		if len(owner) > ownerCol {
			ownerCol = len(owner)
		}
		if len(group) > groupCol {
			groupCol = len(group)
		}
		linkStr := toString(linkCount)
		if len(linkStr) > linkCol {
			linkCol = len(linkStr)
		}
		sizeStr := toString(entry.Info.Size())
		if len(sizeStr) > sizeCol {
			sizeCol = len(sizeStr)
		}
		timeString := formatTime(info.ModTime())
		if len(timeString) > timeCol {
			timeCol = len(timeString)
		}

	}
	return sizeCol, ownerCol, groupCol, linkCol, timeCol, modCol
}

func toString(size interface{}) string {
	return fmt.Sprintf("%v", size)
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

func getFileType(entry FileInfo) string {
	mod := entry.Info.Mode().String()
	switch mod[0] {
	case 'l', 'L':
		return "link"
	}
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
