package gurl

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	flags "github.com/jessevdk/go-flags"
)

const VERSION = "1.0"

const (
	exitOK = iota
	exitError
)

type Flags struct {
	Version bool     `short:"v" long:"version" description:"Show version"`
	Headers []string `short:"H" long:"header" description:"Extra header (format key:value)"`
	JSON    *string  `short:"j" long:"json" description:"JSON data"`
	XML     *string  `short:"x" long:"xml" description:"XML data"`
	Encoded []string `short:"d" long:"data" description:"Form URL Encoded data (format key=value)"`
}

func parseFlags(args []string) (*Flags, []string, error) {
	var f Flags
	p := flags.NewParser(&f, flags.Default)
	p.Usage = "[METHOD] URL [OPTIONS]"
	args, err := p.ParseArgs(os.Args)
	if err != nil {
		return nil, nil, err
	}
	return &f, args, nil
}

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

	for _, arg := range args[1:] {
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

	header := make(http.Header)
	for _, kv := range f.Headers {
		kvs := strings.Split(kv, ":")
		if len(kvs) != 2 {
			return nil, fmt.Errorf("invalid key:value format %s", kv)
		}
		key := kvs[0]
		value := kvs[1]
		header.Set(key, value)
	}

	var body BodyData
	if f.JSON != nil {
		json := JSONData(*f.JSON)
		body = &json
	} else if f.XML != nil {
		xml := XMLData(*f.XML)
		body = &xml
	} else if f.Encoded != nil {
		v := url.Values{}
		for _, kv := range f.Encoded {
			kvs := strings.Split(kv, "=")
			if len(kvs) != 2 {
				return nil, fmt.Errorf("invalid key=value format %s", kv)
			}
			key := kvs[0]
			value := kvs[1]
			v.Set(key, value)
		}
		ff := EncodedData(v)
		body = &ff
	}

	return &Options{
		Method: method,
		URL:    _url,
		Header: header,
		Body:   body,
	}, nil
}
