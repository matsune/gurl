package gurl

import (
	"reflect"
	"testing"
)

func Test_isMethod(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want bool
	}{
		{
			str:  "get",
			want: true,
		},
		{
			str:  "GET",
			want: true,
		},
		{
			str:  "Get",
			want: true,
		},
		{
			str:  "head",
			want: true,
		},
		{
			str:  "HEAD",
			want: true,
		},
		{
			str:  "post",
			want: true,
		},
		{
			str:  "POST",
			want: true,
		},
		{
			str:  "put",
			want: true,
		},
		{
			str:  "PUT",
			want: true,
		},
		{
			str:  "patch",
			want: true,
		},
		{
			str:  "PATCH",
			want: true,
		},
		{
			str:  "delete",
			want: true,
		},
		{
			str:  "Delete",
			want: true,
		},
		{
			str:  "connect",
			want: true,
		},
		{
			str:  "CONNECT",
			want: true,
		},
		{
			str:  "options",
			want: true,
		},
		{
			str:  "OPTIONS",
			want: true,
		},
		{
			str:  "trace",
			want: true,
		},
		{
			str:  "TRACe",
			want: true,
		},
		{
			str:  "ge",
			want: false,
		},
		{
			str:  "GETT",
			want: false,
		},
		{
			str:  "_get",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isMethod(tt.str); got != tt.want {
				t.Errorf("isMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}
func Test_splitKV(t *testing.T) {
	tests := []struct {
		name      string
		str       string
		wantKey   string
		wantValue string
		wantErr   bool
	}{
		{
			name:      "key:value",
			str:       "key:value",
			wantKey:   "key",
			wantValue: "value",
		},
		{
			name:    "invalid format k:v1:v2",
			str:     "k:v1:v2",
			wantErr: true,
		},
		{
			name:    "only key",
			str:     "key",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotValue, err := splitKV(tt.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("splitKV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotKey != tt.wantKey {
				t.Errorf("splitKV() gotKey = %v, want %v", gotKey, tt.wantKey)
			}
			if gotValue != tt.wantValue {
				t.Errorf("splitKV() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}

func Test_splitKVs(t *testing.T) {
	tests := []struct {
		name    string
		kvs     []string
		want    map[string][]string
		wantErr bool
	}{
		{
			name:    "contains invalid format",
			kvs:     []string{"a:b:c"},
			wantErr: true,
		},
		{
			name: "empty kvs",
			kvs:  []string{},
			want: map[string][]string{},
		},
		{
			kvs: []string{"a:b", "a:c", "b:d"},
			want: map[string][]string{
				"a": []string{"b", "c"},
				"b": []string{"d"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := splitKVs(tt.kvs)
			if (err != nil) != tt.wantErr {
				t.Errorf("splitKVs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitKVs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_basicAuth(t *testing.T) {
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				username: "gurl",
				password: "passw0rd",
			},
			want: "Z3VybDpwYXNzdzByZA==",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := basicAuth(tt.args.username, tt.args.password); got != tt.want {
				t.Errorf("basicAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}
