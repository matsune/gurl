package gurl

import (
	"reflect"
	"testing"
)

func Test_parseArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    *cmdArgs
		wantErr bool
	}{
		{
			name:    "nil",
			args:    nil,
			wantErr: true,
		},
		{
			name:    "empty",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "unknown flag",
			args:    []string{"gurl", "-z", "aaa"},
			wantErr: true,
		},
		{
			name: "no args",
			args: []string{"gurl"},
			want: &cmdArgs{flags: &cmdFlags{}, fields: []string{}},
		},
		{
			name: "1 arg",
			args: []string{"gurl", "a"},
			want: &cmdArgs{flags: &cmdFlags{}, fields: []string{"a"}},
		},
		{
			name: "2 args",
			args: []string{"gurl", "a", "b"},
			want: &cmdArgs{flags: &cmdFlags{}, fields: []string{"a", "b"}},
		},
		{
			name: "3 args",
			args: []string{"gurl", "a", "b", "c"},
			want: &cmdArgs{flags: &cmdFlags{}, fields: []string{"a", "b", "c"}},
		},
		{
			name: "version short flag",
			args: []string{"gurl", "-v"},
			want: &cmdArgs{flags: &cmdFlags{Version: true}, fields: []string{}},
		},
		{
			name: "version long flag",
			args: []string{"gurl", "--version"},
			want: &cmdArgs{flags: &cmdFlags{Version: true}, fields: []string{}},
		},
		{
			name: "interactive short flag",
			args: []string{"gurl", "-i"},
			want: &cmdArgs{flags: &cmdFlags{Interactive: true}, fields: []string{}},
		},
		{
			name: "interactive long flag",
			args: []string{"gurl", "--interactive"},
			want: &cmdArgs{flags: &cmdFlags{Interactive: true}, fields: []string{}},
		},
		{
			name: "oneline short flag",
			args: []string{"gurl", "-o"},
			want: &cmdArgs{flags: &cmdFlags{Oneline: true}, fields: []string{}},
		},
		{
			name: "oneline long flag",
			args: []string{"gurl", "--oneline"},
			want: &cmdArgs{flags: &cmdFlags{Oneline: true}, fields: []string{}},
		},
		{
			name: "basic key value flag",
			args: []string{"gurl", "-u", "user:pass"},
			want: &cmdArgs{flags: &cmdFlags{Basic: "user:pass"}, fields: []string{}},
		},
		{
			name: "1 header key-value flag",
			args: []string{"gurl", "-H", "key:val"},
			want: &cmdArgs{flags: &cmdFlags{Headers: []string{"key:val"}}, fields: []string{}},
		},
		{
			name: "2 header key-values flag",
			args: []string{"gurl", "--header", "key1:val1", "-H", "key2:val2"},
			want: &cmdArgs{flags: &cmdFlags{Headers: []string{"key1:val1", "key2:val2"}}, fields: []string{}},
		},
		{
			name: "json short flag",
			args: []string{"gurl", "-j", "{'a':'b'}"},
			want: &cmdArgs{flags: &cmdFlags{JSON: "{'a':'b'}"}, fields: []string{}},
		},
		{
			name: "json long flag",
			args: []string{"gurl", "--json", "{'a':'b'}"},
			want: &cmdArgs{flags: &cmdFlags{JSON: "{'a':'b'}"}, fields: []string{}},
		},
		{
			name: "xml short flag",
			args: []string{"gurl", "-x", "<a><b></b></a>"},
			want: &cmdArgs{flags: &cmdFlags{XML: "<a><b></b></a>"}, fields: []string{}},
		},
		{
			name: "xml long flag",
			args: []string{"gurl", "--xml", "<a><b></b></a>"},
			want: &cmdArgs{flags: &cmdFlags{XML: "<a><b></b></a>"}, fields: []string{}},
		},
		{
			name: "1 form key-value flag",
			args: []string{"gurl", "-f", "key:val"},
			want: &cmdArgs{flags: &cmdFlags{Form: []string{"key:val"}}, fields: []string{}},
		},
		{
			name: "2 form key-values flag",
			args: []string{"gurl", "--form", "key1:val1", "-f", "key2:val2"},
			want: &cmdArgs{flags: &cmdFlags{Form: []string{"key1:val1", "key2:val2"}}, fields: []string{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseArgs(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
