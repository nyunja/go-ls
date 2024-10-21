package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Flag struct to store parsed flag and its value
type Flags struct {
	Long      bool
	All       bool
	Recursive bool
	Reverse   bool
	Time      bool
}

// FileInfo struct to store file information from readDir function
type FileInfo struct {
	name string
	info os.FileInfo
}

func main() {
	// Parse flags from command line
	flags, args := parseFlags(os.Args[1:])
	if len(args) == 0 {
		args = []string{"."}
	}
	for i, path := range args {
		err := listPath(path, flags)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ls: %s: %v\n", path, err)
		}
		if i < len(args)-1 {
			fmt.Println()
		}
	}
}

// listPath lists the contents of a specified directory path based on the given flags.
//
// Parameters:
//   - path: A string representing the directory path to list.
//   - flags: A Flags struct containing boolean flags to control the behavior of the listing.
//
// Returns:
//   - error: An error if there was a problem reading the directory or displaying its contents.
//     Returns nil if the operation was successful.
func listPath(path string, flags Flags) error {
	entries, err := readDir(path, flags)
	if err != nil {
		return err
	}
	if !flags.Long {
		displayShortList(entries)
	} else {
		displayLongFormat(entries)
	}
	if flags.Recursive {
		for _, entry := range entries {
			if entry.info.IsDir() {
				fmt.Println()
				newPath := filepath.Join(path, entry.name)
				if err := listPath(newPath, flags); err != nil {
					fmt.Fprintf(os.Stderr, "ls: %s: %v\n", newPath, err)
				}
			}
		}
	}
	return nil
}

// readDir reads the contents of a directory and returns a slice of FileInfo structures.
// It applies filtering and sorting based on the provided flags.
//
// Parameters:
//   - path: A string representing the directory path to read.
//   - flags: A Flags struct containing boolean flags to control the behavior of the function.
//
// Returns:
//   - []FileInfo: A slice of FileInfo structures containing information about the directory entries.
//   - error: An error if there was a problem reading the directory or its contents.
func readDir(path string, flags Flags) ([]FileInfo, error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	files, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}
	var entries []FileInfo
	for _, file := range files {
		if !flags.All && strings.HasPrefix(file.Name(), ".") {
			continue
		}
		entries = append(entries, FileInfo{name: file.Name(), info: file})
	}
	sort.SliceStable(entries, func(i, j int) bool {
		if flags.Time {
			return entries[i].info.ModTime().After(entries[j].info.ModTime())
		}
		return entries[i].name < entries[j].name
	})
	if flags.Reverse {
		for i := len(entries)/2 - 1; i >= 0; i-- {
			opp := len(entries) - 1 - i
			entries[i], entries[opp] = entries[opp], entries[i]
		}
	}
	return entries, nil
}
// parseFlags parses command-line arguments to extract flags and non-flag arguments.
// It supports both long format (e.g., "--long") and short format (e.g., "-l") flags.
//
// Parameters:
//   - args: A slice of strings representing the command-line arguments to be parsed.
//
// Returns:
//   - flags: A Flags struct containing boolean values for each recognized flag.
//   - parsedArgs: A slice of strings containing the non-flag arguments.
func parseFlags(args []string) (flags Flags, parsedArgs []string) {
	for _, arg := range args {
		if len(arg) > 1 && arg[0] == '-' {
			switch arg {
			case "--reverse":
				flags.Reverse = true
			case "--long":
				flags.Long = true
			case "--all":
				flags.All = true
			case "--recursive":
				flags.Recursive = true
			case "--time":
				flags.Time = true
			default:
				for _, flag := range arg[1:] {
					switch flag {
					case 'l':
						flags.Long = true
					case 'a':
						flags.All = true
					case 'R':
						flags.Recursive = true
					case 't':
						flags.Time = true
					case 'r':
						flags.Reverse = true
					}
				}
			}
		} else {
			parsedArgs = append(parsedArgs, arg)
		}
	}
	return flags, parsedArgs
}

func displayShortList(entries []FileInfo) {
	for _, entry := range entries {
		fmt.Println(entry.name)
	}
}

func getLongFormatString(info fs.FileInfo) string {
	mode := info.Mode()
	size := info.Size()
	modTime := info.ModTime()
	name := info.Name()
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

	s := fmt.Sprintf("%s %2d %s %s %8d %s %s", mode, linkCount, owner, group, size, timeString, name)
	if name[0] == '.' {
		s = "h" + s
	}
	return s
}

func formatTime(modTime time.Time) string {
	now := time.Now()
	if modTime.Year() == now.Year() {
		return modTime.Format("Jan _2 15:04")
	}
	return modTime.Format("Jan _2 2006")
}

func displayLongFormat(entries []FileInfo) {
	var totalBlocks int64
	for _, entry := range entries {
		if stat, ok := entry.info.Sys().(*syscall.Stat_t); ok {
			totalBlocks += stat.Blocks
		}
	}
	fmt.Printf("total %d\n", totalBlocks/2)
	for _, entry := range entries {
		fmt.Println(getLongFormatString(entry.info))
	}
}

// func calcSize(s int64) string {
// 	// unit := "B"
// 	return fmt.Sprintf("%v", s)
// }
