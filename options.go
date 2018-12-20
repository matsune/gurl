package gurl

import (
	"net/http"
	"net/url"
	"strings"
)

type Options struct {
	Method string
	URL    string
	Header http.Header
	Body   BodyData
}

type (
	BodyData interface {
		ContentType() string
		Raw() string
	}

	JSONData    string
	XMLData     string
	EncodedData url.Values
)

func (JSONData) ContentType() string {
	return "application/json"
}
func (XMLData) ContentType() string {
	return "application/xml"
}
func (EncodedData) ContentType() string {
	return "application/x-www-form-urlencoded"
}

func (j JSONData) Raw() string {
	return string(j)
}
func (x XMLData) Raw() string {
	return string(x)
}
func (f EncodedData) Raw() string {
	return url.Values(f).Encode()
}

func (opts *Options) buildRequest() (req *http.Request, err error) {
	if opts.Body != nil {
		req, err = http.NewRequest(opts.Method, opts.URL, strings.NewReader(opts.Body.Raw()))
		req.Header = opts.Header
		req.Header.Set("Content-Type", opts.Body.ContentType())
	} else {
		req, err = http.NewRequest(opts.Method, opts.URL, nil)
		req.Header = opts.Header
	}
	return
}
