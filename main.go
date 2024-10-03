package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
)

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
	if len(parsedArgs) == 0 {
		parsedArgs = []string{"."}
	}
	if !*longFormat {
		displayShortList(parsedArgs)
	} else {
		displayLongList(parsedArgs)
		return
	}
	fmt.Printf("Other arguments: %v\n", parsedArgs)
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

func displayShortList(paths []string) {
	var noFileList []string
	var filesList []string
	var dirList []string
	for _, path := range paths {
		fi, err := os.Stat(path)
		// fmt.Println(fi.Mode())
		if err != nil {
			s := fmt.Sprintf("ls: %v: no file or directory\n", path)
			noFileList = append(noFileList, s)
			continue
		}
		if !fi.IsDir() {
			filesList = append(filesList, fi.Name())
			continue
		} else {
			dirList = addDirList(dirList, path)
		}
		// Get list of files in the directory
	}
	for _, f := range noFileList {
		fmt.Println(f)
	}
	for _, f := range filesList {
		fmt.Println(f)
	}
	for _, f := range dirList {
		fmt.Println(f)
	}
}

func addDirList(dirList []string, path string) []string {
	file, err := os.Open(path)
	if err != nil {
		return dirList
	}
	fileNames, err := file.Readdirnames(0)
	if err != nil {
		return dirList
	}
	dirList = append(dirList, "\n"+path+":")
	sort.Strings(fileNames)
	dirList = append(dirList, fileNames...)
	return dirList
}

func addLongDirList(dirList []string, path string) []string {
	entries, err := os.ReadDir(path)
	if err != nil {
		return dirList
	}
	dirList = append(dirList, "\n"+path+":")
	for _, entry := range entries {
		info,err := entry.Info()
		if err!= nil {
			fmt.Printf("error reading entry: %v", err)
            continue
        }
		size := calcSize(info.Size())
		s := fmt.Sprintf("%v 1 johnotieno0 bocal %s %v %v %v:%v %s", info.Mode(), size, info.ModTime().Month().String()[0:3], info.ModTime().Day(), info.ModTime().Hour(), info.ModTime().Minute(), info.Name())
		dirList = append(dirList, s)
	}
	return dirList
}

func displayLongList(paths []string) {
	fmt.Println("here is a short list display ", paths)
	var noFileList []string
	var filesList []string
	var dirList []string
	for _, path := range paths {
		fi, err := os.Stat(path)
		fmt.Println("here")
		// fmt.Println(fi.Mode())
		if err != nil {
			s := fmt.Sprintf("ls: %v: no file or directory\n", path)
			noFileList = append(noFileList, s)
			continue
		}
		if !fi.IsDir() {
			size := calcSize(fi.Size())
			s := fmt.Sprintf("%v 1 johnotieno0 bocal %s %v %v %v:%v %s", fi.Mode(), size, fi.ModTime().Month().String()[0:3], fi.ModTime().Day(), fi.ModTime().Hour(), fi.ModTime().Minute(), fi.Name())
			filesList = append(filesList, s)
			continue
		} else {
			dirList = addLongDirList(dirList, path)
		}
		// Get list of files in the directory
	}
	for _, f := range noFileList {
		fmt.Println(f)
	}
	for _, f := range filesList {
		fmt.Println(f)
	}
	for _, f := range dirList {
		fmt.Println(f)
	}
}

func calcSize(s int64) string {
	// unit := "B"
	return fmt.Sprintf("%v", s)
}
