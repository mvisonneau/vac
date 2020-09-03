package main

import (
	"os"

	"github.com/mvisonneau/vac/cli"
)

var version = ""

func main() {
	cli.Run(version, os.Args)
}
