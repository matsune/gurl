package gurl

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultHeaderRenderer(t *testing.T) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	header.Add("Content-Type", "charset=utf-8")
	header.Set("Content-Length", "35")

	res := DefaultHeaderRender(header)
	assert.Contains(t, res, "--- Header ---\n\n")
	assert.Contains(t, res, "Content-Type: application/json, charset=utf-8")
	assert.Contains(t, res, "Content-Length: 35")
}

func TestJSONRenderer(t *testing.T) {
	str := `{"key1": "value1","key2": ["value2"]}`
	res := JSONRender(str)
	if assert.NotEmpty(t, res) {
		assert.Equal(t, `---  Body  ---

{
  "key1": "value1",
  "key2": [
    "value2"
  ]
}`, res)
	}
}
