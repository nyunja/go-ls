package lsfunctions

import "testing"

func Test_getFileType(t *testing.T) {
	type args struct {
		entry Entry
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "test 1", args: args{entry: Entry{Name: "test.pdf", Mode: "-rw-r--r--"}}, want: "pdf"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFileType(tt.args.entry); got != tt.want {
				t.Errorf("getFileType() = %v, want %v", got, tt.want)
			}
		})
	}
}
