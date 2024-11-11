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
