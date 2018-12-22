package gurl

import (
	flags "github.com/jessevdk/go-flags"
)

type Flags struct {
	Version     bool     `short:"v" long:"version" description:"Show version"`
	Interactive bool     `short:"i" long:"interactive" description:"Interactive mode"`
	OutOneline  bool     `short:"o" long:"out-oneline" description:"Output oneline command (only interactive mode)"`
	Basic       string   `short:"u" long:"user" description:"Basic auth <user:password>"`
	Headers     []string `short:"H" long:"header" description:"Extra header <key:value>"`
	JSON        *string  `short:"j" long:"json" description:"JSON data"`
	XML         *string  `short:"x" long:"xml" description:"XML data"`
	Encoded     []string `short:"d" long:"data" description:"Form URL Encoded data <key=value>"`
}

func parseFlags(args []string) (*Flags, []string, error) {
	var f Flags
	p := flags.NewParser(&f, flags.Default)
	p.Usage = "[METHOD] URL [OPTIONS]"
	rest, err := p.ParseArgs(args)
	if err != nil {
		return nil, nil, err
	}
	return &f, rest[1:], nil
}
