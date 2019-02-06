package gurl

import (
	"fmt"

	"github.com/AlecAivazis/survey"
)

const gurlVersion = "1.0"

func printVersion() {
	fmt.Printf("gurl version %s\n", gurlVersion)
}

func Run(osArgs []string) error {
	// parse flags and fields
	cmdArgs, err := parseArgs(osArgs)
	if err != nil {
		return err
	}

	if cmdArgs.flags.Version {
		printVersion()
		return nil
	}

	opts, err := cmdArgs.buildOptions()
	if err != nil {
		return err
	}

	// show prompt if Basic auth option doesn't have password
	if opts.Basic != nil && len(opts.Basic.Password) == 0 {
		p := ""
		s := &survey.Password{
			Message: fmt.Sprintf("Password for user %s:", opts.Basic.User),
		}
		if err := survey.AskOne(s, &p, nil); err != nil {
			return err
		}
		opts.Basic.Password = p
	}

	if cmdArgs.isInteractive() {
		if err = runInteractive(opts); err != nil {
			return err
		}
	}

	if err := run(opts); err != nil {
		return err
	}

	if cmdArgs.flags.Oneline {
		str, err := opts.outputOneline()
		if err != nil {
			return err
		}
		fmt.Print(str)
	}

	return nil
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
