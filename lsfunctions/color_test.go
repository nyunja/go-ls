package lsfunctions

import (
	"reflect"
	"testing"
)

func Test_colorName(t *testing.T) {
	type args struct {
		entry Entry
	}
	tests := []struct {
		name string
		args args
		want Entry
	}{
		{name: "test 1", args: args{entry: Entry{Name: "test.txt", Mode: "-rw-r--r--"}}, want: Entry{Name: "test.txt", Mode: "-rw-r--r--"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := colorName(tt.args.entry, false); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("colorName() = %q, want %q", got.Name, tt.want.Name)
			}
		})
	}
}
