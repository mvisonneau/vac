package cmd

import (
	"github.com/mvisonneau/vac/pkg/client"
	"time"

	"github.com/hashicorp/go-secure-stdlib/mlock"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"

	"github.com/mvisonneau/go-helpers/logger"
)

var start time.Time

// Config ..
type Config struct {
	Engine    string
	Role      string
	StatePath string

	AuthInfo client.AuthInfo
}

func configure(ctx *cli.Context) (*Config, error) {
	start = ctx.App.Metadata["startTime"].(time.Time)

	if err := logger.Configure(logger.Config{
		Level:  ctx.String("log-level"),
		Format: ctx.String("log-format"),
	}); err != nil {
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

		AuthInfo: client.AuthInfo{
			Method:   ctx.String("auth"),
			RoleName: ctx.String("auth-k8s-role"),
		},
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

// ExecWrapper mlocks the process memory (if supported) before our `run` functions,
// and gracefully logs and exits afterwards.
func ExecWrapper(f func(ctx *cli.Context) (int, error)) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		if err := mlock.LockMemory(); err != nil {
			return exit(1, fmt.Errorf("error locking vac memory: %w", err))
		}
		return exit(f(ctx))
	}
}
