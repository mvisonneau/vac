package cli

import (
	"os"
	"time"

	"github.com/mvisonneau/vac/cmd"
	"github.com/urfave/cli"
)

// Run handles the instanciation of the CLI application
func Run(version string) {
	NewApp(version, time.Now()).Run(os.Args)
}

// NewApp configures the CLI application
func NewApp(version string, start time.Time) (app *cli.App) {
	app = cli.NewApp()
	app.Name = "vac"
	app.Version = version
	app.Usage = "Manage AWS credentials dynamically using Vault"
	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "engine, e",
			EnvVar: "VAC_ENGINE",
			Usage:  "engine `path`",
		},
		cli.StringFlag{
			Name:   "role, r",
			EnvVar: "VAC_ROLE",
			Usage:  "role `name`",
		},
		cli.StringFlag{
			Name:   "state, s",
			EnvVar: "VAC_STATE_PATH",
			Usage:  "state `path`",
			Value:  "~/.vac_state",
		},
		cli.StringFlag{
			Name:   "log-level",
			EnvVar: "VAC_LOG_LEVEL",
			Usage:  "log `level` (debug,info,warn,fatal,panic)",
			Value:  "info",
		},
		cli.StringFlag{
			Name:   "log-format",
			EnvVar: "VAC_LOG_FORMAT",
			Usage:  "log `format` (json,text)",
			Value:  "text",
		},
	}

	app.Action = cmd.ExecWrapper(cmd.Switch)

	app.Commands = []cli.Command{
		{
			Name:   "get",
			Usage:  "get the creds in credential_process format (json)",
			Action: cmd.ExecWrapper(cmd.Get),
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
