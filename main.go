package main

import (
	"fmt"
	"os"

	ls "my-ls/lsfunctions"
)

func main() {
	// Parse flags from command line
	flags, paths, err := ls.ParseFlags(os.Args[1:])
	if err != nil {
        fmt.Fprintf(os.Stderr, "ls: %v\n", err)
        return
    }
	if len(paths) == 0 {
		paths = []string{"."}
	}
	// Sort paths alphabetically and case-insensitively
	paths, idx := ls.SortPaths(paths)

	for i, path := range paths {
		if (flags.Recursive && len(paths) > 1) || (i >= idx) {
			if len(paths) != 1 {
				fmt.Println()
				fmt.Printf("%s:\n", path)
			}
			if i >= idx && !ls.ShowTotals {
				ls.ShowTotals = true
			}
		}
		err := ls.ListPath(path, flags)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ls: %s: %v\n", path, err)
		}
	}
}
