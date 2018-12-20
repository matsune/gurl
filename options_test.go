package gurl

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionsBuildRequest(t *testing.T) {
	h := http.Header{}
	h.Set("CustomHeader", "GURL")

	opt := Options{
		Method: "GET",
		URL:    "http://localhost",
		Header: h,
		Body:   JSONData(`{"key": "value"}`),
	}
	req, err := opt.buildRequest()
	if assert.NoError(t, err) {
		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, "http://localhost", req.URL.String())
		assert.Equal(t, "GURL", req.Header.Get("CustomHeader"))
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"))

		buf := new(bytes.Buffer)
		buf.ReadFrom(req.Body)
		b := buf.String()
		assert.Equal(t, `{"key": "value"}`, b)
	}
}
