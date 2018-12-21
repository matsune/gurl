package gurl

import (
	"fmt"
	"os"
	"strings"
)

const VERSION = "1.0"

const (
	exitOK = iota
	exitError
)

func Run(args []string) int {
	f, args, err := parseFlags(args)
	if err != nil {
		return exitError
	}

	if f.Version {
		fmt.Printf("gurl version %s\n", VERSION)
		return exitOK
	}

	if len(args) == 0 {
		return runInteractive()
	} else {
		return runOneline(f, args)
	}
}

func runInteractive() int {
	fmt.Fprintln(os.Stderr, "unimplemented interactive mode")
	return exitError
}

func runOneline(f *Flags, args []string) int {
	opts, err := parseOptions(f, args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitError
	}

	g, err := New(opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitError
	}

	if err := g.DoRequest(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitError
	}

	if err := g.Render(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitError
	}

	return exitOK
}

func parseOptions(f *Flags, args []string) (*Options, error) {
	var _url, method string
	var err error

	for _, arg := range args {
		if isURL(arg) {
			if len(_url) > 0 {
				err = fmt.Errorf("has multiple URLs")
				return nil, err
			}
			_url = arg
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

	if len(_url) == 0 {
		return nil, fmt.Errorf("no url")
	}

	if len(method) == 0 {
		method = "GET"
	}

	header, err := splitKVs(f.Headers, ":")
	if err != nil {
		return nil, err
	}

	var body BodyData
	if len(f.Basic) > 0 {
		user, pass, err := split(f.Basic, ":")
		if err != nil {
			return nil, err
		}
		header["Authorization"] = []string{fmt.Sprintf("Basic %s", basicAuth(user, pass))}
	} else if f.JSON != nil {
		json := JSONData(*f.JSON)
		body = &json
	} else if f.XML != nil {
		xml := XMLData(*f.XML)
		body = &xml
	} else if f.Encoded != nil {
		v, err := splitKVs(f.Encoded, "=")
		if err != nil {
			return nil, err
		}
		b := EncodedData(v)
		body = &b
	}

	return &Options{
		Method: method,
		URL:    _url,
		Header: header,
		Body:   body,
	}, nil
}
