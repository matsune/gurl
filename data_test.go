package gurl

import "testing"

func TestJSONData_ContentType(t *testing.T) {
	tests := []struct {
		name string
		j    JSONData
		want string
	}{
		{
			name: "empty string",
			j:    JSONData(""),
			want: "application/json",
		},
		{
			j:    JSONData("anystring"),
			want: "application/json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.j.ContentType(); got != tt.want {
				t.Errorf("JSONData.ContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestXMLData_ContentType(t *testing.T) {
	tests := []struct {
		name string
		x    XMLData
		want string
	}{
		{
			name: "empty string",
			x:    XMLData(""),
			want: "application/xml",
		},
		{
			x:    XMLData("anystring"),
			want: "application/xml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.x.ContentType(); got != tt.want {
				t.Errorf("XMLData.ContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodedData_ContentType(t *testing.T) {
	tests := []struct {
		name string
		e    EncodedData
		want string
	}{
		{
			name: "empty map",
			e:    EncodedData(map[string][]string{}),
			want: "application/x-www-form-urlencoded",
		},
		{
			e:    EncodedData(map[string][]string{"a": []string{"b"}}),
			want: "application/x-www-form-urlencoded",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.ContentType(); got != tt.want {
				t.Errorf("EncodedData.ContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONData_Raw(t *testing.T) {
	tests := []struct {
		name string
		j    JSONData
		want string
	}{
		{
			j:    JSONData("anystring"),
			want: "anystring",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.j.Raw(); got != tt.want {
				t.Errorf("JSONData.Raw() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestXMLData_Raw(t *testing.T) {
	tests := []struct {
		name string
		x    XMLData
		want string
	}{
		{
			x:    XMLData("anystring"),
			want: "anystring",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.x.Raw(); got != tt.want {
				t.Errorf("XMLData.Raw() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodedData_Raw(t *testing.T) {
	tests := []struct {
		name string
		f    EncodedData
		want string
	}{
		{
			f:    EncodedData(map[string][]string{"a": []string{"b", "c"}, "b": []string{"d"}}),
			want: "a=b&a=c&b=d",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Raw(); got != tt.want {
				t.Errorf("EncodedData.Raw() = %v, want %v", got, tt.want)
			}
		})
	}
}
