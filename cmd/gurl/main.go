package main

import (
	"fmt"
	"os"

	"github.com/matsune/gurl"
)

const version = "1.0"

func main() {
	app := gurl.New()
	app.SetVersion(version)

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
