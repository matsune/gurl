package gurl

import (
	"fmt"
	"net/http"
	"strings"
)

type (
	Basic struct {
		User     string
		Password string
	}
)
type Options struct {
	Method       string
	URL          string
	Basic        *Basic
	CustomHeader http.Header
	Body         BodyData
}

func (opts *Options) buildHeader() *http.Header {
	h := http.Header{}
	for k, arr := range opts.CustomHeader {
		for _, v := range arr {
			h.Add(k, v)
		}
	}
	if opts.Body != nil {
		h.Set("Content-Type", opts.Body.ContentType())
	}
	if opts.Basic != nil {
		h.Set("Authorization", fmt.Sprintf("Basic %s", basicAuth(opts.Basic.User, opts.Basic.Password)))
	}
	return &h
}

func (opts *Options) buildRequest() (req *http.Request, err error) {
	if opts.Body != nil {
		req, err = http.NewRequest(opts.Method, opts.URL, strings.NewReader(opts.Body.Raw()))
	} else {
		req, err = http.NewRequest(opts.Method, opts.URL, nil)
	}
	req.Header = *opts.buildHeader()
	return
}
