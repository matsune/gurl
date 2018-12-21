package gurl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFlagsVersion(t *testing.T) {
	args := []string{"./gurl", "-v"}
	f, rest, err := parseFlags(args)
	if assert.NoError(t, err) {
		assert.Empty(t, rest)
		assert.True(t, f.Version)
	}
}

func TestParseFlagsBasic(t *testing.T) {
	args := []string{"gurl", "-u=user:pass"}
	f, rest, err := parseFlags(args)
	if assert.NoError(t, err) {
		assert.Empty(t, rest)
		assert.Equal(t, "user:pass", f.Basic)
	}
}

func TestParseFlagsHeader(t *testing.T) {
	args := []string{"gurl", "-H", "A:B", "-H", "c:d", "get", "http://localhost"}
	f, rest, err := parseFlags(args)
	if assert.NoError(t, err) {
		assert.ElementsMatch(t, []string{"get", "http://localhost"}, rest)
		assert.ElementsMatch(t, []string{"A:B", "c:d"}, f.Headers)
	}
}

func TestParseFlagsJSON(t *testing.T) {
	args := []string{"gurl", "-j", `{"user": "u", "password": "p"}`}
	f, rest, err := parseFlags(args)
	if assert.NoError(t, err) {
		assert.Empty(t, rest)
		assert.NotNil(t, f.JSON)
		assert.Equal(t, `{"user": "u", "password": "p"}`, *f.JSON)
	}
}

func TestParseFlagsXML(t *testing.T) {
	args := []string{"gurl", "-x", `<user>u</user><password>p</password>`}
	f, rest, err := parseFlags(args)
	if assert.NoError(t, err) {
		assert.Empty(t, rest)
		assert.NotNil(t, f.XML)
		assert.Equal(t, `<user>u</user><password>p</password>`, *f.XML)
	}
}

func TestParseFlagsEncoded(t *testing.T) {
	args := []string{"gurl", "-d", "user=u", "-d", "password=p"}
	f, rest, err := parseFlags(args)
	if assert.NoError(t, err) {
		assert.Empty(t, rest)
		assert.ElementsMatch(t, []string{"user=u", "password=p"}, f.Encoded)
	}
}

// Fail cases

func TestParseFlagsBasicFail(t *testing.T) {
	args := []string{"gurl", "--user"}
	_, _, err := parseFlags(args)
	assert.Error(t, err)
}
func TestParseFlagsHeaderFail(t *testing.T) {
	args := []string{"gurl", "-H"}
	_, _, err := parseFlags(args)
	assert.Error(t, err)
}

func TestParseFlagsJSONFail(t *testing.T) {
	args := []string{"gurl", "--json"}
	_, _, err := parseFlags(args)
	assert.Error(t, err)
}

func TestParseFlagsXMLFail(t *testing.T) {
	args := []string{"gurl", "--xml"}
	_, _, err := parseFlags(args)
	assert.Error(t, err)
}

func TestParseFlagsEncodedFail(t *testing.T) {
	args := []string{"gurl", "--data"}
	_, _, err := parseFlags(args)
	assert.Error(t, err)
}
