package gurl

import (
	"fmt"
	"os"
)

const gurlVersion = "1.0"

const (
	exitOK = iota
	exitError
)

func Run(args []string) int {
	cmdArgs, err := parseArgs(args)
	if err != nil {
		return exitError
	}

	if cmdArgs.flags.Version {
		fmt.Printf("gurl version %s\n", gurlVersion)
		return exitOK
	}

	opts, err := parseOptions(cmdArgs)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitError
	}

	if cmdArgs.isInteractive {
		err = runInteractive(opts, cmdArgs.flags.OutOneline)
	} else {
		err = runOneline(opts)
	}
	if err != nil {
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
