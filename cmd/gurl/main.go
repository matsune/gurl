package main

import (
	"fmt"
	"os"

	"github.com/matsune/gurl"
)

func main() {
	if err := gurl.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
