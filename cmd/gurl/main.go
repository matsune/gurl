package main

import (
	"fmt"
	"os"

	"github.com/matsune/gurl"
)

func main() {
	app := gurl.New()
	app.SetVersion(gurl.Version)

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
