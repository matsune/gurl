package gurl

import (
	"net/http"
	"testing"
)

func TestOptions_oneliner(t *testing.T) {
	basename := "gurl"

	type fields struct {
		Method       string
		URL          string
		Basic        *Basic
		CustomHeader http.Header
		Body         BodyData
	}

	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "xml data",
			fields: fields{
				Method: "post",
				URL:    "github.com",
				Basic: &Basic{
					User:     "user",
					Password: "pass",
				},
				CustomHeader: http.Header{
					"a": []string{"b", "c"},
				},
				Body: XMLData("<a><b></b></a>"),
			},
			want: `gurl post "github.com" -u user: -H a:b -H a:c -x '<a><b></b></a>'`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &Options{
				Method:       tt.fields.Method,
				URL:          tt.fields.URL,
				Basic:        tt.fields.Basic,
				CustomHeader: tt.fields.CustomHeader,
				Body:         tt.fields.Body,
			}
			got, err := opts.oneliner(basename)
			if (err != nil) != tt.wantErr {
				t.Errorf("Options.oneliner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Options.oneliner() = %v, want %v", got, tt.want)
			}
		})
	}
}
