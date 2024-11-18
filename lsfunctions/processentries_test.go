package lsfunctions

import (
	"reflect"
	"testing"
)

func Test_processEntries(t *testing.T) {
	mockEntries := []FileDetails{
		{Name: "main.go", Info: mockFileInfo{name: "main.go", mode: 0o100644, size: 4}},
		{Name: "ted", Info: mockFileInfo{name: "ted", mode: 0o40755, size: 10}},
		{Name: "go.mod", Info: mockFileInfo{name: "go.mod", mode: 0o644, size: 7}},
	}
	newEntries := []Entry{
		{Name: "main.go", Mode: "-rw-r--r--", Size: "4", Time: "Jan  1  0001"},
		{Name: "ted", Mode: "-rwxr-xr-x", Size: "10", Time: "Jan  1  0001"},
		{Name: "go.mod", Mode: "-rw-r--r--", Size: "7", Time: "Jan  1  0001"},
	}
	tests := []struct {
		name    string
		entries []FileDetails
		want    []Entry
		want1   Widths
	}{
		{name: "test1", entries: mockEntries, want: newEntries, want1: Widths{sizeCol: 2, timeCol: 12, modCol: 10}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := processEntries(tt.entries)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processEntries() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("processEntries() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
