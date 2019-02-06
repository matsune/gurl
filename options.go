package gurl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	u := opts.URL
	if strings.Index(u, "http") != 0 {
		u = "http://" + u
	}
	if opts.Body != nil {
		req, err = http.NewRequest(m, u, strings.NewReader(opts.Body.Raw()))
	} else {
		req, err = http.NewRequest(m, u, nil)
	}
	req.Header = *opts.buildHeader()
	return
}

func (opts *Options) outputOneline() (string, error) {
	path := os.Args[0]
	m := opts.Method
	url := opts.URL

	args := []string{path, m, url}

	if opts.Basic != nil {
		u := fmt.Sprintf("-u %s:%s", opts.Basic.User, opts.Basic.Password)
		args = append(args, u)
	}

	if len(opts.CustomHeader) != 0 {
		var h string
		for k, arr := range opts.CustomHeader {
			for _, v := range arr {
				h += fmt.Sprintf("-H %s=%s ", k, v)
			}
		}
		args = append(args, h)
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

	return sectionStr("> one-liners") + strings.Join(args, " "), nil
}
