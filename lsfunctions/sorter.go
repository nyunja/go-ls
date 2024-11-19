package lsfunctions

import (
	"os"
	"strings"
)

// SortPaths sorts a slice of file paths and separates directories from non-directories.
//
// This function performs the following operations:
// 1. Sorts paths alphabetically by their last component, case-insensitively.
// 2. Moves directories to the end of the slice while preserving their relative order.
// 3. Finds the index of the first directory in the sorted slice.
//
// Parameters:
//   - paths: A slice of strings representing file and directory paths to be sorted.
//
// Returns:
//   - []string: The sorted slice of paths with directories moved to the end.
//   - int: The index of the first directory in the sorted slice. If no directories
//     are present, this will be equal to the length of the slice.
func SortPaths(paths []string) ([]string, int) {
	// Step 1: Bubble sort by the last component, alphabetically and case-insensitively
	for k := 0; k < len(paths)-1; k++ {
		for j := 0; j < len(paths)-1-k; j++ {
			// Remove trailing slashes for comparison only
			pathI := getPathBase(paths[j])
			pathJ := getPathBase(paths[j+1])
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

// sortEntries sorts a slice of FileInfo entries based on the provided flags.
// It uses the quickSort algorithm to perform the sorting operation.

// Parameters:
//   - entries: A slice of FileInfo structures representing the directory entries to be sorted.
//   - flags: A Flags struct containing boolean flags that determine the sorting criteria.

// Returns:
//   - []FileInfo: A sorted slice of FileInfo structures.
func sortEntries(entries []FileDetails, flags Flags) []FileDetails {
	quickSort(entries, 0, len(entries)-1, flags)

	if flags.Reverse {
		for i := len(entries)/2 - 1; i >= 0; i-- {
			opp := len(entries) - 1 - i
			entries[i], entries[opp] = entries[opp], entries[i]
		}
	}
	return entries
}

// quickSort implements the quicksort algorithm
func quickSort(entries []FileDetails, low, high int, flags Flags) {
	if low < high {
		pi := partition(entries, low, high, flags)
		quickSort(entries, low, pi-1, flags)
		quickSort(entries, pi+1, high, flags)
	}
}

// partition is a helper function for quickSort
func partition(entries []FileDetails, low, high int, flags Flags) int {
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
func compareEntries(a, b FileDetails, flags Flags) bool {
	if flags.Time {
		if a.Info.ModTime().Nanosecond() != b.Info.ModTime().Nanosecond() {
			return a.Info.ModTime().After(b.Info.ModTime())
		}
	}

	s1 := strings.ToLower(a.Name)
	s2 := strings.ToLower(b.Name)

	if cleanName(s1) == cleanName(s2) {
		return a.Name < b.Name
	}
	return cleanName(s1) < cleanName(s2)
}
