package gurl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsURLTrue(t *testing.T) {
	cases := []string{
		"http://localhost:3000",
		"https://golang.org/",
	}
	for _, testCase := range cases {
		assert.True(t, isURL(testCase))
	}
}

func TestIsURLFalse(t *testing.T) {
	cases := []string{
		"htt//aaa.com",
	}
	for _, testCase := range cases {
		assert.False(t, isURL(testCase))
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
		assert.True(t, isMethod(testCase))
	}
}

func TestIsMethodFalse(t *testing.T) {
	cases := []string{
		"ge",
		"GETT",
		"_get",
	}
	for _, testCase := range cases {
		assert.False(t, isMethod(testCase))
	}
}
