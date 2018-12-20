package gurl

import (
	"fmt"
	"os"
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
