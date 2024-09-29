package main

import (
	"flag"
	"fmt"
	"os"
)

// type Options struct {
// 	l        bool
// 	a        bool
// 	progName string
// 	path     string
// }

var (
	// Declare flag formats
	longFormat   = flag.Bool("l", false, "Use long listing format.")
	allFiles     = flag.Bool("a", false, "Show hidden files.")
	recursiveDir = flag.Bool("R", false, "List subdirectories recursively.")
	timeFlag     = flag.Bool("t", false, "List files in descending order of time (i.e. newest first)")
	reverser     = flag.Bool("r", false, "List in reverse order.")
)

func main() {

	// Parse flags from command line
	args := os.Args[1:]
	parsedArgs := parseFlags(args)

	// fmt.Printf("Long format: %v\n", *longFormat)
	// fmt.Printf("Show all files: %v\n", *allFiles)
	// fmt.Printf("List subdirectories recursively: %v\n", *recursiveDir)
	// fmt.Printf("Order time: %v\n", *timeFlag)
	// fmt.Printf("Order in reverse: %v\n", *reverser)

	var paths []string
	// if len(parsedArgs) > 1 {
	// 	fmt.Println("Usage: go run . [options] [path]\n[options] are flags\n[path] is the path to the directory whose contents you want to list. This is optional.")
	// 	return
	// }
	if len(parsedArgs) == 0 {
		paths = []string{"."}
		// files, err := os.ReadDir(".")
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// if *longFormat {
		// 	for _, file := range files {
		//         fmt.Println(file.Name())
		//     }
		// } else {
		// 	for _, file := range files {
		// 		fmt.Printf("%s ", file.Name())
		// 	}
		// 	fmt.Println()
		// }

	} else {
		paths = parsedArgs
	}
	fmt.Printf("Other arguments: %v\n", parsedArgs)
	fmt.Println(paths)
	// options.progName = os.Args[0]
	// args := os.Args[1:]
	// for _, arg := range args {
	// 	switch arg {
	// 	case "-l":
	// 		options.l = true
	// 	case "-a":
	// 		options.a = true
	// 	}
	// }
	// fmt.Println(options)
}

func parseFlags(args []string) (parsedArgs []string) {
	for _, arg := range args {
		if len(arg) > 1 && arg[0] == '-' {
			switch arg {
			case "--reverse":
				*reverser = true
			case "--long":
				*longFormat = true
			case "--all":
				*allFiles = true
			case "--recursive":
				*recursiveDir = true
			case "--time":
				*timeFlag = true
			default:
				for _, flag := range arg[1:] {
					switch flag {
					case 'l':
						*longFormat = true
					case 'a':
						*allFiles = true
					case 'R':
						*recursiveDir = true
					case 't':
						*timeFlag = true
					case 'r':
						*reverser = true
					}
				}
			}
		} else {
			parsedArgs = append(parsedArgs, arg)
		}
	}
	return parsedArgs
}

// func ls(args []string) error {

// }
