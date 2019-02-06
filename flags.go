package gurl

import (
	"strings"

	flags "github.com/jessevdk/go-flags"
)

type cmdFlags struct {
	Version     bool     `short:"v" long:"version" description:"Show version"`
	Interactive bool     `short:"i" long:"interactive" description:"Interactive mode"`
	Oneline     bool     `short:"o" long:"oneline" description:"Output oneline command (only interactive mode)"`
	Basic       string   `short:"u" long:"user" description:"Basic auth <user[:password]>"`
	Headers     []string `short:"H" long:"header" description:"Extra header <key:value>"`
	JSON        string   `short:"j" long:"json" description:"JSON data"`
	XML         string   `short:"x" long:"xml" description:"XML data"`
	Form        []string `short:"f" long:"form" description:"Form URL Encoded data <key:value>"`
}

func parseFlags(osArgs []string) (*cmdFlags, []string, error) {
	var f cmdFlags
	p := flags.NewParser(&f, flags.Default)
	p.Usage = "[METHOD] URL [OPTIONS]"

	args, err := p.ParseArgs(osArgs)
	if err != nil {
		return nil, nil, err
	}

	return &f, args[1:], nil
}

func (f cmdFlags) headers() (map[string][]string, error) {
	return splitKVs(f.Headers)
}

func (f cmdFlags) bodyData() (BodyData, error) {
	var body BodyData
	if len(f.JSON) > 0 {
		json := JSONData(f.JSON)
		body = &json
	} else if len(f.XML) > 0 {
		xml := XMLData(f.XML)
		body = &xml
	} else if len(f.Form) > 0 {
		v, err := splitKVs(f.Form)
		if err != nil {
			return nil, err
		}
		b := EncodedData(v)
		body = &b
	}
	return body, nil
}

func (f cmdFlags) basic() *Basic {
	var b *Basic
	if len(f.Basic) > 0 {
		var user, pass string
		kvs := strings.Split(f.Basic, ":")
		user = kvs[0]
		if len(kvs) > 1 {
			pass = kvs[1]
		}
		b = &Basic{
			User:     user,
			Password: pass,
		}
	}
	return b
}
