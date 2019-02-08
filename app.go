package gurl

import (
	"fmt"
	"net/http"
)

type App struct {
	version  string
	client   *http.Client
	renderer *Renderer
}

func New() *App {
	return &App{
		client:   new(http.Client),
		renderer: NewRenderer(),
	}
}

func (a *App) SetVersion(ver string) {
	a.version = ver
}

func (a *App) printVersion() {
	fmt.Printf("gurl version %s\n", a.version)
}

func (a *App) Run(osArgs []string) error {
	// parse flag options and others
	flags, fields, err := parseFlags(osArgs)
	if err != nil {
		return err
	}

	if flags.Version {
		a.printVersion()
		return nil
	}

	// becomes interactive mode if args has -i flag or no args
	isInteractive := flags.Interactive || len(fields) == 0

	opts, err := buildOptions(flags, fields, isInteractive)
	if err != nil {
		return err
	}

	// show prompt if Basic auth option doesn't have password
	if opts.Basic != nil && isEmpty(opts.Basic.Password) {
		msg := fmt.Sprintf("Password for user %s:", opts.Basic.User)
		p, err := prompt.askPassword(msg)
		if err != nil {
			return err
		}
		opts.Basic.Password = p
	}

	if isInteractive {
		// start interactive prompt
		if err = goInteractive(opts); err != nil {
			return err
		}
	}

	req, err := opts.buildRequest()
	if err != nil {
		return err
	}

	res, err := a.client.Do(req)
	if err != nil {
		return err
	}

	if err := a.renderer.render(res); err != nil {
		return err
	}

	return nil
}
