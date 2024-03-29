package flags

import "github.com/urfave/cli/v2"

var (
	Engine = &cli.StringFlag{
		Name:    "engine",
		Aliases: []string{"e"},
		EnvVars: []string{"VAC_ENGINE"},
		Usage:   "engine `path`",
	}

	ForceGenerate = &cli.BoolFlag{
		Name:    "force-generate",
		Aliases: []string{"f"},
		EnvVars: []string{"VAC_FORCE_GENERATE"},
		Usage:   "bypass currently cached creds and generate new ones",
	}

	LogFormat = &cli.StringFlag{
		Name:    "log-format",
		EnvVars: []string{"VAC_LOG_FORMAT"},
		Usage:   "log `format` (json,text)",
		Value:   "text",
	}

	LogLevel = &cli.StringFlag{
		Name:    "log-level",
		EnvVars: []string{"VAC_LOG_LEVEL"},
		Usage:   "log `level` (debug,info,warn,fatal,panic)",
		Value:   "info",
	}

	MinTTL = &cli.DurationFlag{
		Name:    "min-ttl",
		EnvVars: []string{"VAC_MIN_TTL"},
		Usage:   "min-ttl `duration`",
		Value:   0,
	}

	Role = &cli.StringFlag{
		Name:    "role",
		Aliases: []string{"r"},
		EnvVars: []string{"VAC_ROLE"},
		Usage:   "role `name`",
	}

	State = &cli.StringFlag{
		Name:    "state",
		Aliases: []string{"s"},
		EnvVars: []string{"VAC_STATE_PATH"},
		Usage:   "state `path`",
		Value:   "~/.vac_state",
	}

	TTL = &cli.DurationFlag{
		Name:    "ttl",
		Aliases: []string{"t"},
		EnvVars: []string{"VAC_TTL"},
		Usage:   "ttl `duration`",
		Value:   0,
	}
)
