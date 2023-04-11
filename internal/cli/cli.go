package cli

import (
	"fmt"
	"os"
	"time"

	cli "github.com/urfave/cli/v2"

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
		&cli.StringFlag{
			Name:    "engine",
			Aliases: []string{"e"},
			EnvVars: []string{"VAC_ENGINE"},
			Usage:   "engine `path`",
		},
		&cli.StringFlag{
			Name:    "role",
			Aliases: []string{"r"},
			EnvVars: []string{"VAC_ROLE"},
			Usage:   "role `name`",
		},
		&cli.StringFlag{
			Name:    "state",
			Aliases: []string{"s"},
			EnvVars: []string{"VAC_STATE_PATH"},
			Usage:   "state `path`",
			Value:   "~/.vac_state",
		},
		&cli.StringFlag{
			Name:    "log-level",
			EnvVars: []string{"VAC_LOG_LEVEL"},
			Usage:   "log `level` (debug,info,warn,fatal,panic)",
			Value:   "info",
		},
		&cli.StringFlag{
			Name:    "log-format",
			EnvVars: []string{"VAC_LOG_FORMAT"},
			Usage:   "log `format` (json,text)",
			Value:   "text",
		},
		&cli.StringFlag{
			Name:    "auth",
			EnvVars: []string{"VAC_AUTH"},
			Usage:   "auth method (token, kubernetes)",
			Value:   "token",
		},
		&cli.StringFlag{
			Name:    "auth-k8s-role",
			EnvVars: []string{"VAC_AUTH_K8S_ROLE"},
			Usage:   "Kubernetes role to authenticate to (for --auth kubernetes)",
		},
		&cli.StringFlag{
			Name:    "auth-k8s-mount",
			EnvVars: []string{"VAC_AUTH_K8S_MOUNT"},
			Usage:   "Kubernetes auth mount path (for --auth kubernetes)",
			Value:   "kubernetes",
		},
	}

	app.Action = cmd.ExecWrapper(cmd.Switch)

	app.Commands = cli.CommandsByName{
		{
			Name:   "get",
			Usage:  "get the creds in credential_process format (json)",
			Action: cmd.ExecWrapper(cmd.Get),
			Flags: cli.FlagsByName{
				&cli.DurationFlag{
					Name:    "min-ttl",
					EnvVars: []string{"VAC_MIN_TTL"},
					Usage:   "min-ttl `duration`",
					Value:   0,
				},
				&cli.DurationFlag{
					Name:    "ttl",
					Aliases: []string{"t"},
					EnvVars: []string{"VAC_TTL"},
					Usage:   "ttl `duration`",
					Value:   0,
				},
				&cli.BoolFlag{
					Name:    "force-generate",
					Aliases: []string{"f"},
					EnvVars: []string{"VAC_FORCE_GENERATE"},
					Usage:   "bypass currently cached creds and generate new ones",
				},
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
