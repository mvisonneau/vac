package main

import (
	"os"

	"github.com/mvisonneau/vac/internal/cli"
)

var version = ""

func main() {
	cli.Run(version, os.Args)
}
