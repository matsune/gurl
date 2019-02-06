package main

import (
	"fmt"
	"os"

	"github.com/matsune/gurl"
)

const version = "1.0"

func main() {
	app := gurl.New(os.Args, version)

	if err := app.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
