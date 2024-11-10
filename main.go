package main

import (
	"fmt"
	"os"

	"my-ls/lsfunctions"
)

func main() {
	// Parse flags from command line
	flags, paths := lsfunctions.ParseFlags(os.Args[1:])
	if len(paths) == 0 {
		paths = []string{"."}
	}
	// Sort paths alphabetically and case-insensitively
	paths, idx := lsfunctions.SortPaths(paths)

	for i, path := range paths {
		if (flags.Recursive && len(paths) > 1) || (i >= idx) {
			if len(paths) != 1 {
				fmt.Println()
				fmt.Printf("%s:\n", path)
			}
		}
		err := lsfunctions.ListPath(path, flags)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ls: %s: %v\n", path, err)
		}

	}
}
