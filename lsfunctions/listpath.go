package lsfunctions

import (
	"fmt"
	"os"
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
		DisplayLongFormat(os.Stdout, entries)
	} else {
		DisplayShortList(os.Stdout, entries)
	}
	if flags.Recursive {
		for _, entry := range entries {
			if entry.Info.IsDir() {
				if entry.Name == ".." || entry.Name == "." {
					continue
				}
				fmt.Println()
				newPath := joinPath(path, entry.Name)
				fmt.Printf("%s:\n", newPath)
				if err := ListPath(newPath, flags); err != nil {
					fmt.Fprintf(os.Stdout, "total 0\n")
					fmt.Fprintf(os.Stderr, "ls: cannot open directory '%s': Permission denied\n", newPath)
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
func readDir(path string, flags Flags) ([]FileDetails, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, fmt.Errorf("error accessing path %s: %w", path, err)
	}
	if !info.IsDir() && flags.Long {
		return handleNonDirectory(path, info, flags)
	}

	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	files, err := dir.ReadDir(-1)
	if err != nil {
		return nil, err
	}
	entries := make([]FileDetails, 0, len(files)+2)

	// Add entries for parents directory and current directory
	if flags.All {
		entries = append(entries, createDotEntry(path)...)
	}

	for _, file := range files {
		if !flags.All && strings.HasPrefix(file.Name(), ".") {
			continue
		}
		fileInfo, err := file.Info()
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not get info for %s: %v\n", file.Name(), err)
			continue
		}
		entry := createFileDetails(path, file.Name(), fileInfo)
		entries = append(entries, entry)
	}

	return sortEntries(entries, flags), nil
}

func handleNonDirectory(path string, info os.FileInfo, flags Flags) ([]FileDetails, error) {
	entry := FileDetails{Name: path, Info: info}
	if flags.Long {
		if info.Mode()&os.ModeSymlink != 0 {
			if linkTarget, err := os.Readlink(path); err == nil {
				// _, err := getLinkTargetType(path, linkTarget)
				// if err != nil {
				// 	if strings.Contains(err.Error(), "target not found") {
				// 		entry.IsBrokenLink = true
				// 		// entry.TargetInfo = TargetInfo{Name: linkTarget, Mode: newEntry.Mode, IsBrokenLink: true}
				// 	} 
				// }
				// entry.TargetInfo = TargetInfo{Name: linkTarget, Mode: newEntry.Mode}
				entry.LinkTarget = linkTarget
			}
		}
	}
	return []FileDetails{entry}, nil
}

func createFileDetails(path, name string, info os.FileInfo) FileDetails {
	entry := FileDetails{Name: name, Info: info}
	setEntryPath(path, &entry)
	if info.Mode()&os.ModeSymlink != 0 {
		newPath := joinPath(path, name)
		if linkTarget, err := os.Readlink(newPath); err == nil {
			// _, err := getLinkTargetType(path, linkTarget)
			// if err != nil {
			// 	if strings.Contains(err.Error(), "target not found") {
			// 		entry.IsBrokenLink = true
			// 		// entry.TargetInfo = TargetInfo{Name: linkTarget, Mode: newEntry.Mode, IsBrokenLink: true}
			// 	} 
			// }
			// entry.TargetInfo = TargetInfo{Name: linkTarget, Mode: newEntry.Mode}
			entry.LinkTarget = linkTarget
		}
	}
	return entry
}

func createDotEntry(path string) []FileDetails {
	var entries []FileDetails
	if currentInfo, err := os.Stat(path); err == nil {
		currentEntry := FileDetails{Name: ".", Info: currentInfo}
		entries = append(entries, currentEntry)
	}
	parentDir := getParentDir(path)
	if parentInfo, err := os.Stat(parentDir); err == nil {
		parentEntry := FileDetails{Name: "..", Info: parentInfo}
		entries = append(entries, parentEntry)
	}
	return entries
}
