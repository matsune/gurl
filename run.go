package gurl

import (
	"fmt"
	"os"
	"strings"
)

const gurlVersion = "1.0"

const (
	exitOK = iota
	exitError
)

var interactiveMode bool

func Run(args []string) int {
	f, args, err := parseFlags(args)
	if err != nil {
		return exitError
	}

	if f.Version {
		fmt.Printf("gurl version %s\n", gurlVersion)
		return exitOK
	}

	opts, err := parseOptions(f, args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitError
	}

	var run func(*Options) error
	if interactiveMode {
		run = runInteractive
	} else {
		run = runOneline
	}
	if err = run(opts); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitError
	}
	return exitOK
}

func runOneline(opts *Options) error {
	g, err := New(opts)
	if err != nil {
		return err
	}

	if err := g.DoRequest(); err != nil {
		return err
	}

	if err := g.Render(); err != nil {
		return err
	}

	return nil
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

	interactiveMode = len(args) == 0 || f.Interactive

	if !interactiveMode && len(_url) == 0 {
		return nil, fmt.Errorf("no url")
	}

	if !interactiveMode && len(method) == 0 {
		method = "GET"
	}

	header, err := splitKVs(f.Headers, ":")
	if err != nil {
		return nil, err
	}

	var body BodyData
	if f.JSON != nil {
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

	opts := Options{
		Method: method,
		URL:    _url,
		Header: header,
		Body:   body,
	}
	if len(f.Basic) > 0 {
		user, pass, err := split(f.Basic, ":")
		if err != nil {
			return nil, err
		}
		opts.SetBasic(user, pass)
	}
	return &opts, nil
}
