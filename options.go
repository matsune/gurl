package gurl

import (
	"bytes"
	"encoding/json"
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

func makeOptions(flags *cmdFlags, fields []string, isInteractive bool) (*Options, error) {
	// validate fields
	var url, method string
	for _, field := range fields {
		if isMethod(field) {
			if len(method) > 0 {
				return nil, fmt.Errorf("multiple methods: '%s %s'", method, field)
			}
			method = strings.ToUpper(field)
		} else {
			if len(url) > 0 {
				return nil, fmt.Errorf("multiple URLs: '%s %s'", url, field)
			}
			url = field
		}
	}

	// URL is required if not interactive mode
	if !isInteractive && len(url) == 0 {
		return nil, fmt.Errorf("no URL")
	}

	header, err := flags.headers()
	if err != nil {
		return nil, err
	}

	body, err := flags.bodyData()
	if err != nil {
		return nil, err
	}

	if !isInteractive && len(method) == 0 {
		if body == nil {
			method = http.MethodGet
		} else {
			method = http.MethodPost
		}
	}

	b := flags.basic()

	opts := Options{
		Method:       method,
		URL:          url,
		Basic:        b,
		CustomHeader: header,
		Body:         body,
	}

	return &opts, nil
}

func (opts *Options) httpRequest() (*http.Request, error) {
	m := strings.ToUpper(opts.Method)
	u := opts.URL
	if strings.Index(u, "http") != 0 {
		u = "http://" + u
	}

	var req *http.Request
	var err error
	if opts.Body != nil {
		req, err = http.NewRequest(m, u, strings.NewReader(opts.Body.Raw()))
	} else {
		req, err = http.NewRequest(m, u, nil)
	}
	if err != nil {
		return nil, err
	}
	req.Header = opts.httpHeader()
	return req, nil
}

func (opts *Options) httpHeader() http.Header {
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

	h.Set("User-Agent", fmt.Sprintf("gurl/%s", Version))
	return h
}

func (opts *Options) oneliner(basename string) (string, error) {
	url := fmt.Sprintf(`"%s"`, opts.URL)
	args := []string{basename, opts.Method, url}

	if opts.Basic != nil {
		args = append(args, fmt.Sprintf("-u %s:", opts.Basic.User))
	}

	if len(opts.CustomHeader) > 0 {
		for k, arr := range opts.CustomHeader {
			for _, v := range arr {
				args = append(args, fmt.Sprintf("-H %s=%s", k, v))
			}
		}
	}

	if opts.Body != nil {
		var d string
		switch v := opts.Body.(type) {
		case JSONData:
			buf := new(bytes.Buffer)
			if err := json.Compact(buf, []byte(v)); err != nil {
				return "", err
			}
			d = fmt.Sprintf("-j '%s'", buf)
		case XMLData:
			d = fmt.Sprintf("-x '%s'", v)
		case EncodedData:
			for k, arr := range v {
				for _, v := range arr {
					d += fmt.Sprintf("-d %s=%s ", k, v)
				}
			}
		}
		args = append(args, d)
	}

	return strings.Join(args, " "), nil
}
