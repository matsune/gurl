package main

import (
	"os"

	"github.com/matsune/gurl"
)

func main() {
	exit := gurl.Run(os.Args)
	os.Exit(exit)
}
