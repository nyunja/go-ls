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
		if flags.Recursive {
			if flags.Long {
				ls.ShowTotals = true
			}
			if i > 0 {
				fmt.Println()
			}
			fmt.Printf("%s:\n", path)
		} else if i >= idx {
			if flags.Long {
				ls.ShowTotals = true
			}
			if i > 0 {
				fmt.Println()
			}
			if len(paths) > 1 {
				fmt.Printf("%s:\n", path)
			}
			// if flags.Recursive || len(paths) >= 1 {
			// 	if i > 0 {
			// 		fmt.Println()
			// 	}
			// 	if len(paths) > 1 {
			// 		fmt.Printf("%s:\n", path)
			// 	}
			// }
			// ls.ShowTotals = true
		}
		err := ls.ListPath(path, flags)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ls: cannot access '%s': No such file or directory\n", path)
		}
	}
}
