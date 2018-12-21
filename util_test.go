package gurl

import (
	"encoding/base64"
	"fmt"
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

func TestSplit(t *testing.T) {
	k, v, err := split("a:b", ":")
	if assert.NoError(t, err) {
		assert.Equal(t, "a", k)
		assert.Equal(t, "b", v)
	}
}

func TestSplitKVs(t *testing.T) {
	kvs := []string{"a:b", "c:d"}
	m, err := splitKVs(kvs, ":")
	if assert.NoError(t, err) {
		if assert.Contains(t, m, "a") {
			assert.Equal(t, []string{"b"}, m["a"])
		}
		if assert.Contains(t, m, "c") {
			assert.Equal(t, []string{"d"}, m["c"])
		}
	}
}

func TestBasicAuth(t *testing.T) {
	user := "gurl"
	pass := "passw0rd"
	encoded := basicAuth(user, pass)

	payload, err := base64.StdEncoding.DecodeString(encoded)
	if assert.NoError(t, err) {
		assert.Equal(t, fmt.Sprintf("%s:%s", user, pass), string(payload))
	}
}
