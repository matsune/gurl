package gurl

import (
	"net/http"
	"strings"
)

type Options struct {
	Method string
	URL    string
	Header http.Header
	Body   BodyData
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
