package gurl

import (
	"fmt"
	"net/http"
)

type App struct {
	version string
	Prompt
	http.Client
	Renderer
}

func New() *App {
	return &App{
		version: appVersion,
		Client:  http.Client{},
		Prompt: &surveyPrompt{
			editor: "vim",
		},
		Renderer: NewRenderer(),
	}
}

func (a *App) printVersion() {
	fmt.Printf("gurl version %s\n", a.version)
}

func Run(osArgs []string) error {
	// parse flag options and others
	flags, fields, err := parseFlags(osArgs)
	if err != nil {
		return err
	}
	a := New()

	if flags.Version {
		a.printVersion()
		return nil
	}

	// becomes interactive mode if args has -i flag or no args
	isInteractive := flags.Interactive || len(fields) == 0

	opts, err := makeOptions(flags, fields, isInteractive)
	if err != nil {
		return err
	}

	// show prompt if Basic auth option doesn't have password
	if opts.Basic != nil && isEmpty(opts.Basic.Password) {
		msg := fmt.Sprintf("Password for user %s:", opts.Basic.User)
		p, err := a.InputPassword(msg)
		if err != nil {
			return err
		}
		opts.Basic.Password = p
	}

	if isInteractive {
		// start interactive prompt
		if err = a.interactive(opts); err != nil {
			return err
		}
	}

	req, err := opts.httpRequest()
	if err != nil {
		return err
	}

	res, err := a.Do(req)
	if err != nil {
		return err
	}

	if err := render(a.Renderer, res); err != nil {
		return err
	}

	if isInteractive {
		str, err := opts.oneliner(osArgs[0])
		if err != nil {
			return err
		}
		fmt.Print(a.Oneliner(str))
	}

	return nil
}
