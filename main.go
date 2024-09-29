package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// type Options struct {
// 	l        bool
// 	a        bool
// 	progName string
// 	path     string
// }

// var options Options

func main() {
	longFormat := flag.Bool("l", false, "Use long listing format.")
	allFiles := flag.Bool("a", false, "Show hidden files.")
	recursiveDir := flag.Bool("R", false, "List subdirectories recursively.")
	timeFlag := flag.Bool("t", false, "List files in descending order of time (i.e. newest first)")
	reverser := flag.Bool("r", false, "List in reverse order.")

	flag.Parse()

	fmt.Printf("Long format: %v\n", *longFormat)
	fmt.Printf("Show all files: %v\n", *allFiles)
	fmt.Printf("List subdirectories recursively: %v\n", *recursiveDir)
	fmt.Printf("Order time: %v\n", *timeFlag)
	fmt.Printf("Order in reverse: %v\n", *reverser)

	args := flag.Args()
	if len(args) == 0 {
		files, err := os.ReadDir(".")
		if err != nil {
			log.Fatal(err)
		}
		if *longFormat == true {
			for _, file := range files {
                fmt.Println(file.Name())
            }
		} else {
			for _, file := range files {
				fmt.Printf("%s ", file.Name())
			}
			fmt.Println()
		}
	
		
    }
	// fmt.Printf("Other arguments: %v\n", args)
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

// func ls(args []string) error {

// }
