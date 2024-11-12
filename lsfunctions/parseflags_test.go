package lsfunctions

import (
	"reflect"
	"testing"
)

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		wantFlags      Flags
		wantParsedArgs []string
	}{
		{name: "test with no flags", args: []string{"file1", "file2"}, wantFlags: Flags{}, wantParsedArgs: []string{"file1", "file2"}},
		{name: "test with single flag", args: []string{"-l", "file1", "file2"}, wantFlags: Flags{Long: true}, wantParsedArgs: []string{"file1", "file2"}},
		{name: "test with multiple flags", args: []string{"-laR", "file1", "file2"}, wantFlags: Flags{Long: true, All: true, Recursive: true}, wantParsedArgs: []string{"file1", "file2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFlags, gotParsedArgs, _ := ParseFlags(tt.args)
			if !reflect.DeepEqual(gotFlags, tt.wantFlags) {
				t.Errorf("ParseFlags() gotFlags = %v, want %v", gotFlags, tt.wantFlags)
			}
			if !reflect.DeepEqual(gotParsedArgs, tt.wantParsedArgs) {
				t.Errorf("ParseFlags() gotParsedArgs = %v, want %v", gotParsedArgs, tt.wantParsedArgs)
			}
		})
	}
}
