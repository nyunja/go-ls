package lsfunctions

import (
	"fmt"
	"os"
	"strings"
)

// getPath extracts the last component from a given path string.//+
// If the path does not contain any slashes, it returns the original string.//+
// //+
// Parameters://+
//   - s: A string representing the input path.//+
//
// //+
// Returns://+
//   - A string representing the last component of the input path.//+
func getPath(s string) string {
	l := strings.LastIndex(s, "/")
	if l == -1 {
		return s
	}
	return s[l+1:]
}

// SortPaths sorts a slice of file paths and separates directories from non-directories.//+
// //+
// This function performs the following operations://+
// 1. Sorts paths alphabetically by their last component, case-insensitively.//+
// 2. Moves directories to the end of the slice while preserving their relative order.//+
// 3. Finds the index of the first directory in the sorted slice.//+
// //+
// Parameters://+
//   - paths: A slice of strings representing file and directory paths to be sorted.//+
//
// //+
// Returns://+
//   - []string: The sorted slice of paths with directories moved to the end.//+
//   - int: The index of the first directory in the sorted slice. If no directories//+
//     are present, this will be equal to the length of the slice.//+
func SortPaths(paths []string) ([]string, int) {
	// Step 1: Bubble sort by the last component, alphabetically and case-insensitively
	for k := 0; k < len(paths)-1; k++ {
		for j := 0; j < len(paths)-1-k; j++ {
			// Remove trailing slashes for comparison only
			pathI := getPath(paths[j])
			pathJ := getPath(paths[j+1])
			// fmt.Printf("before: %s %s\n", pathI, pathJ)
			// Compare alphabetically
			if strings.ToLower(pathI) > strings.ToLower(pathJ) {
				paths[j], paths[j+1] = paths[j+1], paths[j] // Swap
			}
			// fmt.Printf("after: %v %s, %s\n", paths, pathI, pathJ)
		}
	}

	// Step 2: Bubble sort to move directories to the back while preserving order
	for k := 0; k < len(paths)-1; k++ {
		for j := 0; j < len(paths)-1-k; j++ {
			infoI, errI := os.Lstat(paths[j])
			infoJ, errJ := os.Lstat(paths[j+1])

			if errI != nil || errJ != nil {
				continue // Ignore errors
			}

			// Move directories to the back
			if infoI.IsDir() && !infoJ.IsDir() {
				paths[j], paths[j+1] = paths[j+1], paths[j] // Swap
			}
		}
	}

	// Step 3: Find the index of the first non-directory
	nonDirIdx := len(paths) // Default to the end if no non-directories are found
	for i, path := range paths {
		info, err := os.Lstat(path)
		if err != nil {
			continue
		}
		if info.IsDir() {
			nonDirIdx = i
			break
		}
	}

	return paths, nonDirIdx
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
				newPath :=joinPath(path, entry.Name)
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
	var target string
	if info, err := os.Lstat(path); err == nil {
		if flags.Long {
			if target, err = os.Readlink(path); err == nil {
				entry := FileInfo{Name: path, Info: info, LinkTarget: target}
				return []FileInfo{entry}, nil
			}
		}
		if !info.IsDir() {
			if _, err = os.Readlink(path); err != nil {
				entry := FileInfo{Name: path, Info: info}
				return []FileInfo{entry}, nil
			}
		}
	}

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
		entry := FileInfo{Name: file.Name(), Info: file}
		// Get the Rdev for device files
		mode := file.Mode().String()
		switch mode[0] {
		case 'l', 'L':
			newPath := joinPath(path, file.Name())
			linkTarget, err := os.Readlink(newPath)
			if err == nil {
				entry.LinkTarget = linkTarget
			}
		}
		entries = append(entries, entry)
	}

	// Sort entries
	entries = sortEntries(entries, flags)

	if flags.Reverse {
		for i := len(entries)/2 - 1; i >= 0; i-- {
			opp := len(entries) - 1 - i
			entries[i], entries[opp] = entries[opp], entries[i]
		}
	}
	return entries, nil
}

// getParentDir returns the parent directory path of the given path.
//
// It handles various edge cases such as root directory, paths without separators,
// and paths ending with a separator.
//
// Parameters:
//   - path: A string representing the input path for which to find the parent directory.
//
// Returns:
//   - A string representing the parent directory path.
//     Returns "/" for the root directory, ".." for paths without separators,
//     and the appropriate parent path for other cases.
func getParentDir(path string) string {
	if path == "/" {
		return "/"
	}
	lastIndexSep := strings.LastIndex(path, "/")
	if lastIndexSep == -1 {
		return ".."
	}
	if lastIndexSep == len(path)-1 {
		path = path[:lastIndexSep]
		lastIndexSep = strings.LastIndex(path, "/")
	}

	if lastIndexSep == 0 {
		return "/"
	}
	if lastIndexSep == -1 {
		return ".."
	}
	return path[:lastIndexSep]
}

// sortEntries sorts a slice of FileInfo entries based on the provided flags.
// It uses the quickSort algorithm to perform the sorting operation.

// Parameters:
//   - entries: A slice of FileInfo structures representing the directory entries to be sorted.//+
//   - flags: A Flags struct containing boolean flags that determine the sorting criteria.//+

// Returns:
//   - []FileInfo: A sorted slice of FileInfo structures.
func sortEntries(entries []FileInfo, flags Flags) []FileInfo {
	quickSort(entries, 0, len(entries)-1, flags)
	return entries
}

// quickSort implements the quicksort algorithm
func quickSort(entries []FileInfo, low, high int, flags Flags) {
	if low < high {
		pi := partition(entries, low, high, flags)
		quickSort(entries, low, pi-1, flags)
		quickSort(entries, pi+1, high, flags)
	}
}

// partition is a helper function for quickSort
func partition(entries []FileInfo, low, high int, flags Flags) int {
	pivot := entries[high]
	i := low - 1

	for j := low; j < high; j++ {
		if compareEntries(entries[j], pivot, flags) {
			i++
			entries[i], entries[j] = entries[j], entries[i]
		}
	}

	entries[i+1], entries[high] = entries[high], entries[i+1]
	return i + 1
}

// compareEntries compares two FileInfo entries based on the sorting criteria
func compareEntries(a, b FileInfo, flags Flags) bool {
	if flags.Time {
		return a.Info.ModTime().After(b.Info.ModTime())
	}

	s1 := strings.ToLower(a.Name)
	s2 := strings.ToLower(b.Name)

	if cleanName(s1) == cleanName(s2) {
		return a.Name < b.Name
	}
	return cleanName(s1) < cleanName(s2)
}

// Clean string to remove -, _, and. from the name.
func cleanName(name string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case '-', '_', '.', '[', '#', ']', '{', '}', '|', '\\', ':', ';', '<', '>', ',', '?', '!', '@', '$', '%', '^', '&', '*', '(', ')', '~', '`', '"', '\'', '=', '+', '/':
			return -1
		}
		return r
	}, name)
}

func joinPath(parts...string) string {
	res := ""
	for i, part := range parts {
		if i == len(parts)-1 && !strings.HasPrefix(part, "/") {
			res += "/"
		}
		res += part
	}
	return res
}
