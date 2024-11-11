package lsfunctions

import (
	"bytes"
	"os"
	"testing"
	"time"
)

func TestDisplayShortList(t *testing.T) {
	var mockEntries = []FileInfo{
		{Name: "main.go", Info: mockFileInfo{name: "main.go", mode: 0100644, modTime: time.Now()}},
		{Name: "ted", Info: mockFileInfo{name: "ted", mode: 040755, modTime: time.Now()}},
		{Name: "go.mod", Info: mockFileInfo{name: "go.mod", mode: 0644, modTime: time.Now()}},
	}
	type args struct {
		e []FileInfo
	}
	tests := []struct {
		name  string
		args  args
		wantW string
	}{
		{name: "test with no entries", args: args{e: []FileInfo{}}, wantW: ""},
		// {name: "test with one entry", args: args{e: []FileInfo{{Name: "file1", Info: mockFileInfo{name: "file1",mode: 0644},}}}, wantW: "file1"},
		{name: "test with multiple entries", args: args{e: mockEntries}, wantW: "file1\nfile2\nfile3"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			DisplayShortList(w, tt.args.e)
			os.Stdout.WriteString("here: "+ w.String())
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("DisplayShortList() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
