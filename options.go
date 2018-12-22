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

	Options struct {
		Method       string
		URL          string
		Basic        *Basic
		CustomHeader http.Header
		Body         BodyData
	}
)

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
	m := strings.ToUpper(opts.Method)
	if opts.Body != nil {
		req, err = http.NewRequest(m, opts.URL, strings.NewReader(opts.Body.Raw()))
	} else {
		req, err = http.NewRequest(m, opts.URL, nil)
	}
	req.Header = *opts.buildHeader()
	return
}

func parseOptions(cmdArgs *cmdArgs) (*Options, error) {
	if cmdArgs == nil {
		return nil, fmt.Errorf("cmdArgs is nil")
	}

	var u, m string
	var err error

	for _, arg := range cmdArgs.rest {
		if isURL(arg) {
			if len(u) > 0 {
				err = fmt.Errorf("multiple URLs: '%s %s'", u, arg)
				return nil, err
			}
			u = arg
		} else if isMethod(arg) {
			if len(m) > 0 {
				err = fmt.Errorf("multiple methods: '%s %s'", m, arg)
				return nil, err
			}
			m = strings.ToUpper(arg)
		} else {
			err = fmt.Errorf("unknown argument: %s", arg)
			return nil, err
		}
	}

	if !cmdArgs.isInteractive && len(u) == 0 {
		return nil, fmt.Errorf("no URL")
	}

	header, err := splitKVs(cmdArgs.flags.Headers, ":")
	if err != nil {
		return nil, err
	}

	var body BodyData
	if cmdArgs.flags.JSON != nil {
		json := JSONData(*cmdArgs.flags.JSON)
		body = &json
	} else if cmdArgs.flags.XML != nil {
		xml := XMLData(*cmdArgs.flags.XML)
		body = &xml
	} else if cmdArgs.flags.Form != nil {
		v, err := splitKVs(cmdArgs.flags.Form, ":")
		if err != nil {
			return nil, err
		}
		b := EncodedData(v)
		body = &b
	}

	if !cmdArgs.isInteractive && len(m) == 0 {
		if body == nil {
			m = http.MethodGet
		} else {
			m = http.MethodPost
		}
	}

	var b *Basic
	if len(cmdArgs.flags.Basic) > 0 {
		user, pass, err := split(cmdArgs.flags.Basic, ":")
		if err != nil {
			return nil, err
		}
		b = &Basic{
			User:     user,
			Password: pass,
		}
	}

	opts := Options{
		Method:       m,
		URL:          u,
		Basic:        b,
		CustomHeader: header,
		Body:         body,
	}

	return &opts, nil
}
