package gurl

import (
	"fmt"
	"net/http"
)

const gurlVersion = "1.0"

func printVersion() {
	fmt.Printf("gurl version %s\n", gurlVersion)
}

func Run(osArgs []string) error {
	flags, fields, err := parseFlags(osArgs)
	if err != nil {
		return err
	}

	if flags.Version {
		printVersion()
		return nil
	}

	args := cmdArgs{
		flags:  flags,
		fields: fields,
	}

	opts, err := args.buildOptions()
	if err != nil {
		return err
	}

	// show prompt if Basic auth option doesn't have password
	if opts.Basic != nil && len(opts.Basic.Password) == 0 {
		p, err := askBasicPassword(opts.Basic.User)
		if err != nil {
			return err
		}
		opts.Basic.Password = p
	}

	if args.isInteractive() {
		if err = runInteractive(opts); err != nil {
			return err
		}
	}

	req, err := opts.buildRequest()
	if err != nil {
		return err
	}

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	renderer := NewRenderer()
	if err := renderer.render(res); err != nil {
		return err
	}

	return nil
}
