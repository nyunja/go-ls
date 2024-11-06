package lsfunctions

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func SortPaths(paths []string) ([]string, int) {
	nonDirIdx := 0
	sort.Slice(paths, func(i, j int) bool {
		pathI, pathJ := paths[i],paths[j]
		m,n := strings.LastIndex(paths[i], "/"), strings.LastIndex(paths[j], "/")
		if m == -1 || n == -1 {
			return false
		}
		if m == len(paths[i]) -1 {
			pathI = pathI[:m]
			m = strings.LastIndex(pathI, "/")
		}
		if n == len(paths[i]) -1 {
			pathJ = pathJ[:n]
			n = strings.LastIndex(pathJ, "/")
		}
		pathI, pathJ = paths[i][m+1:], paths[j][n+1:]
        return strings.ToLower(pathI) < strings.ToLower(pathJ)
    })
	sort.SliceStable(paths, func(i, j int) bool {
		infoI, err := os.Lstat(paths[i])
		// infoJ, err := os.Lstat(paths[j])
		if err!= nil {
            return false
        }
		if !infoI.IsDir() {
			nonDirIdx++
		}
		return !infoI.IsDir()

	})
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
		DisplayLongFormat(entries)
	} else {
		DisplayShortList(entries)
	}
	if flags.Recursive {
		for _, entry := range entries {
			if entry.Info.IsDir() {
				if entry.Name == ".." || entry.Name == "." {
					continue
				}
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
	if flags.Long {
		if info, err := os.Lstat(path); err == nil {
			if target, err := os.Readlink(path); err == nil {
				entry := FileInfo{Name: path, Info: info, LinkTarget: target}
				return []FileInfo{entry}, nil
			}
			if !info.IsDir() {
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
		mode := file.Mode().String()
		switch mode[0] {
		case 'l', 'L':
			linkTarget, err := os.Readlink(path + file.Name())
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
	if lastIndexSep == len(path) -1 {
		path = path[:lastIndexSep]
		lastIndexSep = strings.LastIndex(path, "/")
	}
	if lastIndexSep == 0 {
		return "/"
	}
	return path[:lastIndexSep]
}

// Sort entries using quicksort
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

// // Sort entries
// func sortEntries(entries []FileInfo, flags Flags) []FileInfo {
// 	sort.SliceStable(entries, func(i, j int) bool {
// 		if flags.Time {
// 			return entries[i].Info.ModTime().After(entries[j].Info.ModTime())
// 		}
// 		s1 := strings.ToLower(entries[i].Name)
// 		s2 := strings.ToLower(entries[j].Name)
// 		if cleanName(s1) == cleanName(s2) {
// 			return entries[i].Name < entries[j].Name
// 		}
// 		return cleanName(s1) < cleanName(s2)
// 	})
// 	return entries
// }

// Clean string to remove -, _, and. from the name.
func cleanName(name string) string {
    return strings.Map(func(r rune) rune {
        if r == '-' || r == '_' || r == '.' {
            return -1
        }
        return r
    }, name)
}
