package cmd

import (
	"fmt"

	"github.com/hashicorp/go-secure-stdlib/mlock"
	log "github.com/sirupsen/logrus"
	"github.com/syndtr/gocapability/capability"
	cli "github.com/urfave/cli/v2"
)

// ExecWrapper mlocks the process memory (if supported) before our `run` functions,
// and gracefully logs and exits afterwards.
func ExecWrapper(f func(ctx *cli.Context) (int, error)) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		caps, err := capability.NewPid2(0)
		if err != nil {
			return exit(1, fmt.Errorf("error getting capabilities: %w", err))
		}

		if err = caps.Load(); err != nil {
			return exit(1, fmt.Errorf("error loading capabilities: %w", err))
		}

		if caps.Get(capability.EFFECTIVE, capability.CAP_IPC_LOCK) { // mlock.Supported() is assumed
			if err = mlock.LockMemory(); err != nil {
				return exit(1, fmt.Errorf("error locking vac memory: %w", err))
			}
		} else {
			log.Warn("unable to lock memory, missing CAP_IPC_LOCK capability")
		}

		return exit(f(ctx))
	}
}
