package gurl

import (
	"fmt"

	flags "github.com/jessevdk/go-flags"
)

type cmdFlags struct {
	Version     bool     `short:"v" long:"version" description:"Show version"`
	Interactive bool     `short:"i" long:"interactive" description:"Interactive mode"`
	OutOneline  bool     `short:"o" long:"out-oneline" description:"Output oneline command (only interactive mode)"`
	Basic       string   `short:"u" long:"user" description:"Basic auth <user[:password]>"`
	Headers     []string `short:"H" long:"header" description:"Extra header <key:value>"`
	JSON        string   `short:"j" long:"json" description:"JSON data"`
	XML         string   `short:"x" long:"xml" description:"XML data"`
	Form        []string `short:"f" long:"form" description:"Form URL Encoded data <key:value>"`
}

type cmdArgs struct {
	cmdName       string
	flags         cmdFlags
	rest          []string
	isInteractive bool
}

func parseArgs(args []string) (*cmdArgs, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("args is empty")
	}

	var f cmdFlags
	p := flags.NewParser(&f, flags.Default)
	p.Usage = "[METHOD] URL [OPTIONS]"
	args, err := p.ParseArgs(args)
	if err != nil {
		return nil, err
	}
	return &cmdArgs{
		cmdName:       args[0],
		flags:         f,
		rest:          args[1:],
		isInteractive: f.Interactive || len(args[1:]) == 0,
	}, nil
}
