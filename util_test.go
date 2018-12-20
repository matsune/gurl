package gurl

import "testing"

func TestIsURLTrue(t *testing.T) {
	cases := []string{
		"http://localhost:3000",
		"https://golang.org/",
	}
	for _, testCase := range cases {
		if !isURL(testCase) {
			t.Errorf("isURL(\"%s\") should be true", testCase)
		}
	}
}

func TestIsURLFalse(t *testing.T) {
	cases := []string{
		"htt//aaa.com",
	}
	for _, testCase := range cases {
		if isURL(testCase) {
			t.Errorf("isURL(\"%s\") should be false", testCase)
		}
	}
}

func TestIsMethodTrue(t *testing.T) {
	cases := []string{
		"get",
		"GET",
		"Get",
		"head",
		"HEAD",
		"post",
		"POST",
		"put",
		"PUT",
		"patch",
		"PATCH",
		"delete",
		"DELETE",
		"connect",
		"CONNECT",
		"options",
		"OPTIONS",
		"trace",
		"TRACE",
	}
	for _, testCase := range cases {
		if !isMethod(testCase) {
			t.Errorf("isMethod(\"%s\") should be true", testCase)
		}
	}
}

func TestIsMethodFalse(t *testing.T) {
	cases := []string{
		"ge",
		"GETT",
		"_get",
	}
	for _, testCase := range cases {
		if isMethod(testCase) {
			t.Errorf("isMethod(\"%s\") should be false", testCase)
		}
	}
}
