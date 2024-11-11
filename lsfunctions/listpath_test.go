package lsfunctions

import (
	"reflect"
	"testing"
)

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
