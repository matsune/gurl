package gurl

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey"
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

	// show prompt if Basic auth option doesn't have password
	if opts.Basic != nil && len(opts.Basic.Password) == 0 {
		p := ""
		s := &survey.Password{
			Message: fmt.Sprintf("Password for user %s:", opts.Basic.User),
		}
		if err := survey.AskOne(s, &p, nil); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return exitError
		}
		opts.Basic.Password = p
	}

	if cmdArgs.isInteractive {
		if err = runInteractive(opts); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return exitError
		}
	}

	if err := run(opts); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitError
	}

	if cmdArgs.flags.OutOneline {
		str, err := opts.outputOneline()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return exitError
		}
		fmt.Print(str)
	}

	return exitOK
}

func run(opts *Options) error {
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
