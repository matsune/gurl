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

	assert.Equal(t, `[Header]
Content-Type: application/json, charset=utf-8
Content-Length: 35
`, DefaultHeaderRender(header))
}
