package cmd

import (
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/mvisonneau/go-helpers/logger"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
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

	lc := &logger.Config{
		Level:  ctx.String("log-level"),
		Format: ctx.String("log-format"),
	}

	if err := lc.Configure(); err != nil {
		return nil, errors.Wrap(err, "configuring logger")
	}

	statePath, err := homedir.Expand(ctx.String("state"))
	if err != nil {
		return nil, errors.Wrap(err, "expanding cache path value (go-homedir)")
	}

	return &Config{
		Engine:    ctx.String("engine"),
		Role:      ctx.String("role"),
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

// ExecWrapper gracefully logs and exits our `run` functions
func ExecWrapper(f func(ctx *cli.Context) (int, error)) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		return exit(f(ctx))
	}
}
