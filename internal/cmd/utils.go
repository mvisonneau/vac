package cmd

import (
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"

	"github.com/mvisonneau/go-helpers/logger"
	"github.com/mvisonneau/vac/internal/cli/flags"
)

var start time.Time

// Config ..
type Config struct {
	Engine    string
	Role      string
	StatePath string
}

func configure(ctx *cli.Context) (*Config, error) {
	start = ctx.App.Metadata["startTime"].(time.Time)

	if err := logger.Configure(logger.Config{
		Format: flags.LogFormat.Get(ctx),
		Level:  flags.LogLevel.Get(ctx),
	}); err != nil {
		return nil, errors.Wrap(err, "configuring logger")
	}

	statePath, err := homedir.Expand(flags.State.Get(ctx))
	if err != nil {
		return nil, errors.Wrap(err, "expanding cache path value (go-homedir)")
	}

	return &Config{
		Engine:    flags.Engine.Get(ctx),
		Role:      flags.Role.Get(ctx),
		StatePath: statePath,
	}, nil
}

func exit(exitCode int, err error) cli.ExitCoder {
	defer log.WithFields(
		log.Fields{
			"execution-time": time.Since(start),
		},
	).Debug("exited..")

	if err != nil {
		log.Error(err.Error())
	}

	return cli.NewExitError("", exitCode)
}
