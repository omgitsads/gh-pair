package main

import (
	"os"

	"github.com/omgitsads/gh-pair/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
