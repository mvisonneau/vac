//go:build !linux
// +build !linux

package cmd

import (
	"fmt"

	"github.com/hashicorp/go-secure-stdlib/mlock"
	cli "github.com/urfave/cli/v2"
)

// ExecWrapper mlocks the process memory (if supported) before our `run` functions,
// and gracefully logs and exits afterwards.
func ExecWrapper(f func(ctx *cli.Context) (int, error)) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		if mlock.Supported() {
			if err := mlock.LockMemory(); err != nil {
				return exit(1, fmt.Errorf("error locking vac memory: %w", err))
			}
		}

		return exit(f(ctx))
	}
}
