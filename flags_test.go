package gurl

import (
	"net/url"
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
			name:       "empty",
			osArgs:     []string{},
			wantFlags:  &cmdFlags{},
			wantFields: []string{},
		},
		{
			name:    "unknown flag",
			osArgs:  []string{"gurl", "-z", "aaa"},
			wantErr: true,
		},
		{
			name:       "no args",
			osArgs:     []string{"gurl"},
			wantFlags:  &cmdFlags{},
			wantFields: []string{},
		},
		{
			name:       "1 arg",
			osArgs:     []string{"gurl", "a"},
			wantFlags:  &cmdFlags{},
			wantFields: []string{"a"},
		},
		{
			name:       "2 args",
			osArgs:     []string{"gurl", "a", "b"},
			wantFlags:  &cmdFlags{},
			wantFields: []string{"a", "b"},
		},
		{
			name:       "3 args",
			osArgs:     []string{"gurl", "a", "b", "c"},
			wantFlags:  &cmdFlags{},
			wantFields: []string{"a", "b", "c"},
		},
		{
			name:       "version short flag",
			osArgs:     []string{"gurl", "-v"},
			wantFlags:  &cmdFlags{Version: true},
			wantFields: []string{},
		},
		{
			name:       "version long flag",
			osArgs:     []string{"gurl", "--version"},
			wantFlags:  &cmdFlags{Version: true},
			wantFields: []string{},
		},
		{
			name:       "interactive short flag",
			osArgs:     []string{"gurl", "-i"},
			wantFlags:  &cmdFlags{Interactive: true},
			wantFields: []string{},
		},
		{
			name:       "interactive long flag",
			osArgs:     []string{"gurl", "--interactive"},
			wantFlags:  &cmdFlags{Interactive: true},
			wantFields: []string{},
		},
		{
			name:       "oneline short flag",
			osArgs:     []string{"gurl", "-o"},
			wantFlags:  &cmdFlags{Oneline: true},
			wantFields: []string{},
		},
		{
			name:       "oneline long flag",
			osArgs:     []string{"gurl", "--oneline"},
			wantFlags:  &cmdFlags{Oneline: true},
			wantFields: []string{},
		},
		{
			name:       "basic key value flag",
			osArgs:     []string{"gurl", "-u", "user:pass"},
			wantFlags:  &cmdFlags{Basic: "user:pass"},
			wantFields: []string{},
		},
		{
			name:       "1 header key-value flag",
			osArgs:     []string{"gurl", "-H", "key:val"},
			wantFlags:  &cmdFlags{Headers: []string{"key:val"}},
			wantFields: []string{},
		},
		{
			name:       "2 header key-values flag",
			osArgs:     []string{"gurl", "--header", "key1:val1", "-H", "key2:val2"},
			wantFlags:  &cmdFlags{Headers: []string{"key1:val1", "key2:val2"}},
			wantFields: []string{},
		},
		{
			name:       "json short flag",
			osArgs:     []string{"gurl", "-j", "{'a':'b'}"},
			wantFlags:  &cmdFlags{JSON: "{'a':'b'}"},
			wantFields: []string{},
		},
		{
			name:       "json long flag",
			osArgs:     []string{"gurl", "--json", "{'a':'b'}"},
			wantFlags:  &cmdFlags{JSON: "{'a':'b'}"},
			wantFields: []string{},
		},
		{
			name:       "xml short flag",
			osArgs:     []string{"gurl", "-x", "<a><b></b></a>"},
			wantFlags:  &cmdFlags{XML: "<a><b></b></a>"},
			wantFields: []string{},
		},
		{
			name:       "xml long flag",
			osArgs:     []string{"gurl", "--xml", "<a><b></b></a>"},
			wantFlags:  &cmdFlags{XML: "<a><b></b></a>"},
			wantFields: []string{},
		},
		{
			name:       "1 form key-value flag",
			osArgs:     []string{"gurl", "-f", "key:val"},
			wantFlags:  &cmdFlags{Form: []string{"key:val"}},
			wantFields: []string{},
		},
		{
			name:       "2 form key-values flag",
			osArgs:     []string{"gurl", "--form", "key1:val1", "-f", "key2:val2"},
			wantFlags:  &cmdFlags{Form: []string{"key1:val1", "key2:val2"}},
			wantFields: []string{},
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

func Test_cmdFlags_headers(t *testing.T) {
	tests := []struct {
		name    string
		Headers []string
		want    map[string][]string
		wantErr bool
	}{
		{
			Headers: []string{"a:b", "a:c", "b:d"},
			want: map[string][]string{
				"a": []string{"b", "c"},
				"b": []string{"d"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdFlags{
				Headers: tt.Headers,
			}
			got, err := f.headers()
			if (err != nil) != tt.wantErr {
				t.Errorf("cmdFlags.headers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cmdFlags.headers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cmdFlags_bodyData(t *testing.T) {
	type fields struct {
		JSON string
		XML  string
		Form []string
	}
	tests := []struct {
		name    string
		fields  fields
		want    BodyData
		wantErr bool
	}{
		{
			name: "json",
			fields: fields{
				JSON: "{'a':'b'}",
			},
			want: JSONData("{'a':'b'}"),
		},
		{
			name: "xml",
			fields: fields{
				XML: "<a><b></b></a>",
			},
			want: XMLData("<a><b></b></a>"),
		},
		{
			name: "form",
			fields: fields{
				Form: []string{"a:b", "c:d"},
			},
			want: EncodedData(url.Values{
				"a": []string{"b"},
				"c": []string{"d"},
			}),
		},

		{
			name: "invalid form",
			fields: fields{
				Form: []string{"a:b:ddd"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdFlags{
				JSON: tt.fields.JSON,
				XML:  tt.fields.XML,
				Form: tt.fields.Form,
			}
			got, err := f.bodyData()
			if (err != nil) != tt.wantErr {
				t.Errorf("cmdFlags.bodyData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cmdFlags.bodyData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cmdFlags_basic(t *testing.T) {
	tests := []struct {
		name  string
		Basic string
		want  *Basic
	}{
		{
			Basic: "user:pass",
			want: &Basic{
				User:     "user",
				Password: "pass",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdFlags{
				Basic: tt.Basic,
			}
			if got := f.basic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cmdFlags.basic() = %v, want %v", got, tt.want)
			}
		})
	}
}
