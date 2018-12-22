package gurl

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionsBuildHeader(t *testing.T) {
	opts := Options{
		CustomHeader: http.Header{
			"a": []string{"b"},
		},
		Body: JSONData(`{"user":"u","pass":"p"}`),
	}
	h := opts.buildHeader()
	assert.Equal(t, "application/json", h.Get("Content-Type"))
	assert.Equal(t, "b", h.Get("A"))
}

func TestOptionsBuildRequest(t *testing.T) {
	opt := Options{
		Method: "post",
		URL:    "http://localhost",
		CustomHeader: http.Header{
			"a": []string{"b"},
		},
		Body: XMLData(`<user>u</user><password>p</password>`),
	}
	req, err := opt.buildRequest()
	if assert.NoError(t, err) {
		assert.Equal(t, http.MethodPost, req.Method)
		assert.Equal(t, "http://localhost", req.URL.String())
		assert.Equal(t, "b", req.Header.Get("a"))
		assert.Equal(t, "application/xml", req.Header.Get("Content-Type"))
		assert.Nil(t, opt.Basic)

		buf := new(bytes.Buffer)
		buf.ReadFrom(req.Body)
		b := buf.String()
		assert.Equal(t, `<user>u</user><password>p</password>`, b)
	}
}

func TestParseOptions(t *testing.T) {
	json := `{"user": "gurl","pass": "pass"}`
	c := cmdArgs{
		cmdName: "gurl",
		flags: cmdFlags{
			Basic:   "user:password",
			Headers: []string{"a:b"},
			JSON:    json,
		},
		rest:          []string{"POST", "http://localhost"},
		isInteractive: false,
	}
	opts, err := parseOptions(&c)
	if assert.NoError(t, err) {
		if assert.NotNil(t, opts.Basic) {
			assert.Equal(t, "user", opts.Basic.User)
			assert.Equal(t, "password", opts.Basic.Password)
		}
		assert.Equal(t, http.MethodPost, opts.Method)
		assert.Equal(t, "http://localhost", opts.URL)
		if assert.NotNil(t, opts.CustomHeader["a"]) {
			assert.Contains(t, opts.CustomHeader["a"], "b")
		}
		var jsonType *JSONData
		if assert.IsType(t, jsonType, opts.Body) {
			v, _ := opts.Body.(*JSONData)
			assert.Equal(t, json, v.Raw())
			assert.Equal(t, "application/json", v.ContentType())
		}
	}
}

func TestParseOptionsBasicNoPassword(t *testing.T) {
	c := cmdArgs{
		flags: cmdFlags{
			Basic: "user",
		},
		rest: []string{"http://localhost"},
	}
	opts, err := parseOptions(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, "user", opts.Basic.User)
		assert.Equal(t, "", opts.Basic.Password)
	}
}
