package gurl

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeaderRenderer(t *testing.T) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	header.Add("Content-Type", "charset=utf-8")
	header.Set("Content-Length", "35")

	res := DefaultHeaderRender(header)
	assert.Contains(t, res, "[Header]\n")
	assert.Contains(t, res, "Content-Type: application/json, charset=utf-8")
	assert.Contains(t, res, "Content-Length: 35")
}
