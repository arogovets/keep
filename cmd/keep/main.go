package main

import (
	"fmt"
	"os"

	"github.com/ieffai/stay-awake-cli/internal/cli"
)

var version = "dev"

func main() {
	code, err := cli.Run(os.Args[1:], os.Stdout, os.Stderr, version)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(code)
}
