package gurl

import (
	"fmt"
	"net/http"
	"strings"
)

type cmdArgs struct {
	flags  *cmdFlags
	fields []string
}

func parseArgs(args []string) (*cmdArgs, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("args is empty")
	}

	f, rest, err := parseFlags(args)
	if err != nil {
		return nil, err
	}

	return &cmdArgs{
		flags:  f,
		fields: rest,
	}, nil
}

// command becomes interactive mode if args has -i flag or no args
func (c cmdArgs) isInteractive() bool {
	return c.flags.Interactive || len(c.fields) == 0
}

func (c cmdArgs) buildOptions() (*Options, error) {
	// validate fields
	var url, method string
	for _, field := range c.fields {
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
	if !c.isInteractive() && len(url) == 0 {
		return nil, fmt.Errorf("no URL")
	}

	header, err := c.flags.headers()
	if err != nil {
		return nil, err
	}

	body, err := c.flags.bodyData()
	if err != nil {
		return nil, err
	}

	if !c.isInteractive() && len(method) == 0 {
		if body == nil {
			method = http.MethodGet
		} else {
			method = http.MethodPost
		}
	}

	b := c.flags.basic()

	opts := Options{
		Method:       method,
		URL:          url,
		Basic:        b,
		CustomHeader: header,
		Body:         body,
	}

	return &opts, nil
}
