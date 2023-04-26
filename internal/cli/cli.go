package cli

import (
	"fmt"
	"os"
	"time"

	cli "github.com/urfave/cli/v2"

	"github.com/mvisonneau/vac/internal/cli/flags"
	"github.com/mvisonneau/vac/internal/cmd"
)

// Run handles the instanciation of the CLI application
func Run(version string, args []string) {
	err := NewApp(version, time.Now()).Run(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// NewApp configures the CLI application
func NewApp(version string, start time.Time) (app *cli.App) {
	app = cli.NewApp()
	app.Name = "vac"
	app.Version = version
	app.Usage = "Manage AWS credentials dynamically using Vault"
	app.EnableBashCompletion = true

	app.Flags = cli.FlagsByName{
		flags.Engine,
		flags.LogFormat,
		flags.LogLevel,
		flags.Role,
		flags.State,
	}

	app.Action = cmd.ExecWrapper(cmd.Switch)

	app.Commands = cli.CommandsByName{
		{
			Name:   "get",
			Usage:  "get the creds in credential_process format (json)",
			Action: cmd.ExecWrapper(cmd.Get),
			Flags: cli.FlagsByName{
				flags.ForceGenerate,
				flags.MinTTL,
				flags.TTL,
			},
		},
		{
			Name:   "status",
			Usage:  "returns some info about the current context, cached credentials and Vault server connectivity",
			Action: cmd.ExecWrapper(cmd.Status),
		},
	}

	app.Metadata = map[string]interface{}{
		"startTime": start,
	}

	return
}
