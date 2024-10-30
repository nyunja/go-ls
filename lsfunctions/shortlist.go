package lsfunctions

import "fmt"

func DisplayShortList(entries []FileInfo) {
	for _, entry := range entries {
		fmt.Println(entry.Name)
	}
}
