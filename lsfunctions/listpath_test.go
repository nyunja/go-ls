package lsfunctions

import (
	"os"
	"reflect"
	"testing"
	"time"
)

// mockFileInfo implements os.FileInfo for testing purposes
type mockFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	sys     interface{}
}

func (m mockFileInfo) Name() string       { return m.name }
func (m mockFileInfo) Size() int64        { return m.size }
func (m mockFileInfo) Mode() os.FileMode  { return m.mode }
func (m mockFileInfo) ModTime() time.Time { return m.modTime }
func (m mockFileInfo) IsDir() bool        { return m.mode.IsDir() }
func (m mockFileInfo) Sys() interface{}   { return m.sys }

func TestSortPaths(t *testing.T) {
	type args struct {
		paths []string
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		want1 int
	}{
		{name: "test 1", args: args{paths: []string{"./ted", "main.go", "go.mod"}}, want: []string{"go.mod", "main.go", "./ted"}, want1: 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := SortPaths(tt.args.paths)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortPaths() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SortPaths() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_getPath(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{name: "test 1", args: "./ted", want: "ted"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPath(tt.args); got != tt.want {
				t.Errorf("getPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readDir(t *testing.T) {
	type args struct {
		path  string
		flags Flags
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{name: "test 1", args: args{path: "../ted", flags: Flags{Long: true}}, want: 3, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readDir(tt.args.path, tt.args.flags)
			if (err != nil) != tt.wantErr {
				t.Errorf("readDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("readDir() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func Test_cleanName(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{name: "test 1", args: "_t-ed", want: "ted"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanName(tt.args); got != tt.want {
				t.Errorf("cleanName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sortEntries(t *testing.T) {
	// Mock a []FileInfo for testing
	var mockEntries = []FileInfo{
		{Name: "main.go", Info: mockFileInfo{name: "main.go", mode: 0100644, modTime: time.Now()}},
		{Name: "ted", Info: mockFileInfo{name: "ted", mode: 040755, modTime: time.Now()}},
		{Name: "go.mod", Info: mockFileInfo{name: "go.mod", mode: 0644, modTime: time.Now()}},
	}
	var sortedEntries = []FileInfo{
		{Name: "go.mod", Info: mockFileInfo{name: "go.mod", mode: 0644, modTime: time.Now()}},
		{Name: "main.go", Info: mockFileInfo{name: "main.go", mode: 0100644, modTime: time.Now()}},
		{Name: "ted", Info: mockFileInfo{name: "ted", mode: 040755, modTime: time.Now()}},
	}
	type args struct {
		entries []FileInfo
		flags   Flags
	}
	tests := []struct {
		name string
		args args
		want []FileInfo
	}{
		{name: "test 1", args: args{entries: mockEntries, flags: Flags{}}, want: sortedEntries},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sortEntries(tt.args.entries, tt.args.flags)
			for i, entry := range got {
				if !reflect.DeepEqual(entry.Name, tt.want[i].Name) {
					t.Errorf("sortEntries() entry.Name = %v, want %v", entry.Name, tt.want[i].Name)
				}
			}
		})
	}
}

func Test_quickSort(t *testing.T) {
	// Mock a []FileInfo for testing
	var mockEntries = []FileInfo{
		{Name: "main.go", Info: mockFileInfo{name: "main.go", mode: 0100644, modTime: time.Now()}},
		{Name: "ted", Info: mockFileInfo{name: "ted", mode: 040755, modTime: time.Now()}},
		{Name: "go.mod", Info: mockFileInfo{name: "go.mod", mode: 0644, modTime: time.Now()}},
	}
	var sortedEntries = []FileInfo{
		{Name: "go.mod", Info: mockFileInfo{name: "go.mod", mode: 0644, modTime: time.Now()}},
		{Name: "main.go", Info: mockFileInfo{name: "main.go", mode: 0100644, modTime: time.Now()}},
		{Name: "ted", Info: mockFileInfo{name: "ted", mode: 040755, modTime: time.Now()}},
	}
	type args struct {
		entries []FileInfo
		low     int
		high    int
		flags   Flags
		want    []FileInfo
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "test 1", args: args{entries: mockEntries, low: 0, high: len(mockEntries) - 1, flags: Flags{}, want: sortedEntries}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quickSort(tt.args.entries, tt.args.low, tt.args.high, tt.args.flags)
		})
		for i, entry := range tt.args.entries {
			if !reflect.DeepEqual(entry.Name, tt.args.want[i].Name) {
				t.Errorf("sortEntries() entry.Name = %v, want %v", entry.Name, tt.args.want[i].Name)
			}
		}
	}
}

func Test_compareEntries(t *testing.T) {
	// Mock a []FileInfo for testing
	var mockEntries = []FileInfo{
		{Name: "main.go", Info: mockFileInfo{name: "main.go", mode: 0100644, modTime: time.Now()}},
		{Name: "ted", Info: mockFileInfo{name: "ted", mode: 040755, modTime: time.Now()}},
		{Name: "go.mod", Info: mockFileInfo{name: "go.mod", mode: 0644, modTime: time.Now()}},
	}
	type args struct {
		a     FileInfo
		b     FileInfo
		flags Flags
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "test 1", args: args{a: mockEntries[0], b: mockEntries[1], flags: Flags{}}, want: true},
		{name: "test 2", args: args{a: mockEntries[0], b: mockEntries[2], flags: Flags{}}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compareEntries(tt.args.a, tt.args.b, tt.args.flags); got != tt.want {
				t.Errorf("compareEntries() = %v, want %v", got, tt.want)
			}
		})
	}
}
