package cmd

import (
	"fmt"
	kvbuilder "github.com/hashicorp/go-secure-stdlib/kv-builder"
	"github.com/mitchellh/mapstructure"
	"github.com/mvisonneau/vac/pkg/client"
	"strings"
	"time"

	"github.com/hashicorp/vault/sdk/helper/mlock"
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

	*client.AuthConfig
}

// parseArgsData parses the given args in the format key=value into a map of
// the provided arguments. The given reader can also supply key=value pairs.
func parseArgsData(args []string) (map[string]interface{}, error) {
	builder := &kvbuilder.Builder{}
	if err := builder.Add(args...); err != nil {
		return nil, err
	}

	return builder.Map(), nil
}

// parseArgsDataString parses the args data and returns the values as strings.
// If the values cannot be represented as strings, an error is returned.
func parseArgsDataString(args []string) (map[string]string, error) {
	raw, err := parseArgsData(args)
	if err != nil {
		return nil, err
	}

	var result map[string]string
	if err := mapstructure.WeakDecode(raw, &result); err != nil {
		return nil, errors.Wrap(err, "failed to convert values to strings")
	}
	if result == nil {
		result = make(map[string]string)
	}
	return result, nil
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

	var authMethodArgs []string

	for _, v := range ctx.Args().Slice() {
		if strings.Contains(v, "=") {
			authMethodArgs = append(authMethodArgs, v)
		}
	}

	authMethodConfig, err := parseArgsDataString(authMethodArgs)
	if err != nil {
		return nil, fmt.Errorf("error parsing configuration: %s", err)
	}

	return &Config{
		Engine:    flags.Engine.Get(ctx),
		Role:      flags.Role.Get(ctx),
		StatePath: statePath,

		AuthConfig: &client.AuthConfig{
			AuthMethod:     ctx.String("auth-method"),
			AuthPath:       ctx.String("auth-path"),
			AuthNoStore:    ctx.Bool("auth-no-store"),
			AuthMethodArgs: authMethodConfig,
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
