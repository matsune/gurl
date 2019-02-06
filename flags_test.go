package gurl

import (
	"reflect"
	"testing"
)

func Test_parseFlags(t *testing.T) {
	tests := []struct {
		name       string
		osArgs     []string
		wantFlags  *cmdFlags
		wantFields []string
		wantErr    bool
	}{
		{
			osArgs: []string{"gurl", "get", "localhost", "-H", "a:b", "-u", "user:pass", "-v"},
			wantFlags: &cmdFlags{
				Version: true,
				Headers: []string{"a:b"},
				Basic:   "user:pass",
			},
			wantFields: []string{"get", "localhost"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parseFlags(tt.osArgs)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantFlags) {
				t.Errorf("parseFlags() got = %v, want %v", got, tt.wantFlags)
			}
			if !reflect.DeepEqual(got1, tt.wantFields) {
				t.Errorf("parseFlags() got1 = %v, want %v", got1, tt.wantFields)
			}
		})
	}
}
