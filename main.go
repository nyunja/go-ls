package main

import (
	"fmt"
	"os"

	"my-ls/lsfunctions"
)

func main() {
	// Parse flags from command line
	flags, args := lsfunctions.ParseFlags(os.Args[1:])
	if len(args) == 0 {
		args = []string{"."}
	}
	for i, path := range args {
		err := lsfunctions.ListPath(path, flags)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ls: %s: %v\n", path, err)
		}
		if i < len(args)-1 {
			fmt.Println()
		}
	}
}


// func calcSize(s int64) string {
// 	// unit := "B"
// 	return fmt.Sprintf("%v", s)
// }
