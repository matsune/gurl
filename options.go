package gurl

import (
	"fmt"
	"net/http"
	"strings"
)

type Options struct {
	Method string
	URL    string
	Header http.Header
	Body   BodyData
}

func (opts *Options) SetBasic(user, pass string) {
	opts.SetHeader("Authorization", fmt.Sprintf("Basic %s", basicAuth(user, pass)))
}

func (opts *Options) SetHeader(k, v string) {
	opts.Header.Set(k, v)
}

func (opts *Options) buildRequest() (req *http.Request, err error) {
	if opts.Body != nil {
		req, err = http.NewRequest(opts.Method, opts.URL, strings.NewReader(opts.Body.Raw()))
		opts.SetHeader("Content-Type", opts.Body.ContentType())
	} else {
		req, err = http.NewRequest(opts.Method, opts.URL, nil)
	}
	req.Header = opts.Header
	return
}
