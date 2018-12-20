package gurl

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	flags "github.com/jessevdk/go-flags"
)

type Flags struct {
	Version bool    `short:"v" long:"version" description:"Show version"`
	JSON    *string `short:"j" long:"json" description:"JSON data"`
}

func parseFlags(args []string) (*Flags, []string, error) {
	var f Flags
	parser := flags.NewParser(&f, flags.Default)
	parser.Usage = "METHOD URL [OPTIONS]"
	args, err := parser.ParseArgs(os.Args)
	if err != nil {
		return nil, nil, err
	}
	return &f, args, nil
}

type Options struct {
	Method string
	URL    string
	Body   Body
}

type Body interface {
	Build() io.Reader
	ContentType() string
}

type JSON struct {
	Raw string
}

func (j *JSON) Build() io.Reader {
	return bytes.NewBuffer([]byte(j.Raw))
}
func (JSON) ContentType() string {
	return "application/json"
}

func parseOptions(f *Flags, args []string) (*Options, error) {
	var url, method string
	var err error

	for _, arg := range args[1:] {
		if isURL(arg) {
			if len(url) > 0 {
				err = fmt.Errorf("has multiple URLs")
				return nil, err
			}
			url = arg
		} else if isMethod(arg) {
			if len(method) > 0 {
				err = fmt.Errorf("has multiple methods")
				return nil, err
			}
			method = strings.ToUpper(arg)
		} else {
			err = fmt.Errorf("unknown argument: %s", arg)
			return nil, err
		}
	}

	if len(url) == 0 {
		return nil, fmt.Errorf("no url")
	}

	if len(method) == 0 {
		method = "GET"
	}

	var body Body
	if f.JSON != nil {
		body = &JSON{
			Raw: *f.JSON,
		}
	}

	return &Options{
		Method: method,
		URL:    url,
		Body:   body,
	}, nil
}

func buildRequest(opts *Options) (req *http.Request, err error) {
	if opts.Body != nil {
		req, err = http.NewRequest(opts.Method, opts.URL, opts.Body.Build())
		req.Header.Set("Content-Type", opts.Body.ContentType())
	} else {
		req, err = http.NewRequest(opts.Method, opts.URL, nil)
	}
	return
}
