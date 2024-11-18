package lsfunctions

import (
	"bytes"
	"os"
	"testing"
)

func TestDisplayShortList(t *testing.T) {
	type args struct {
		e []FileDetails
	}
	tests := []struct {
		name  string
		args  args
		wantW string
	}{
		{name: "test with no entries", args: args{e: []FileDetails{}}, wantW: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			DisplayShortList(w, tt.args.e)
			os.Stdout.WriteString("here: " + w.String())
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("DisplayShortList() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
