package lsfunctions

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// listPath lists the contents of a specified directory path based on the given flags.
//
// Parameters:
//   - path: A string representing the directory path to list.
//   - flags: A Flags struct containing boolean flags to control the behavior of the listing.
//
// Returns:
//   - error: An error if there was a problem reading the directory or displaying its contents.
//     Returns nil if the operation was successful.
func ListPath(path string, flags Flags) error {
	entries, err := readDir(path, flags)
	if err != nil {
		return err
	}
	if flags.Long {
		DisplayLongFormat(entries)
	} else {
		DisplayShortList(entries)
	}
	if flags.Recursive {
		for _, entry := range entries {
			if entry.Info.IsDir() {
				fmt.Println()
				newPath := filepath.Join(path, entry.Name)
				fmt.Printf("%s:\n", newPath)
				if err := ListPath(newPath, flags); err != nil {
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

	// Add etries for parents directory and current directory
	if flags.All {
		if currentInfo, err := os.Stat(path); err == nil {
			entries = append(entries, FileInfo{Name: ".", Info: currentInfo})
		}
		parentDir := getParentDir(path)
		if parentInfo, err := os.Stat(parentDir); err == nil {
			entries = append(entries, FileInfo{Name: "..", Info: parentInfo})
		}
	}
	for _, file := range files {
		if !flags.All && strings.HasPrefix(file.Name(), ".") {
			continue
		}
		entries = append(entries, FileInfo{Name: file.Name(), Info: file})
	}
	sort.SliceStable(entries, func(i, j int) bool {
		if flags.Time {
			return entries[i].Info.ModTime().After(entries[j].Info.ModTime())
		}
		s1 := strings.ToLower(entries[i].Name)
		s2 := strings.ToLower(entries[j].Name)
		if cleanName(s1) == cleanName(s2) {
			return entries[i].Name < entries[j].Name
		}
		return cleanName(s1) < cleanName(s2)
	})
	if flags.Reverse {
		for i := len(entries)/2 - 1; i >= 0; i-- {
			opp := len(entries) - 1 - i
			entries[i], entries[opp] = entries[opp], entries[i]
		}
	}
	return entries, nil
}

// Get parent directory
func getParentDir(path string) string {
	if path == "/" {
		return "/"
	}
	lastIndexSep := strings.LastIndex(path, "/") 
	if lastIndexSep == -1 {
		return "/"
	}
	if lastIndexSep == len(path) -1 {
		path = path[:lastIndexSep]
		lastIndexSep = strings.LastIndex(path, "/")
	}
	if lastIndexSep == 0 {
		return "/"
	}
	return path[:lastIndexSep]
}

// Clean string to remove -, _, and. from the name.
func cleanName(name string) string {
    return strings.Map(func(r rune) rune {
        if r == '-' || r == '_' || r == '.' {
            return -1
        }
        return r
    }, name)
}
