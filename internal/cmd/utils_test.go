package cmd

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	cli "github.com/urfave/cli/v2"
)

func NewTestContext() (ctx *cli.Context, flags, globalFlags *flag.FlagSet) {
	app := cli.NewApp()
	app.Name = "vac"

	app.Metadata = map[string]interface{}{
		"startTime": time.Now(),
	}

	globalFlags = flag.NewFlagSet("test", flag.ContinueOnError)
	globalCtx := cli.NewContext(app, globalFlags, nil)

	flags = flag.NewFlagSet("test", flag.ContinueOnError)
	ctx = cli.NewContext(app, flags, globalCtx)

	globalFlags.String("log-level", "fatal", "")
	globalFlags.String("log-format", "text", "")

	return
}

func TestExit(t *testing.T) {
	err := exit(20, fmt.Errorf("test"))
	assert.Equal(t, "", err.Error())
	assert.Equal(t, 20, err.ExitCode())
}

const charSet = "abcdefghijklmnopqrstuvwxyz"

func randString(strlen int) string {
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = charSet[rand.Intn(len(charSet))] //nolint:gosec
	}
	return string(result)
}

func TestFileLock(t *testing.T) {
	fp := filepath.Join(os.TempDir(), randString(8))
	ok, unlock, err := fileLock(fp)
	assert.Nil(t, err)
	assert.True(t, ok)
	assert.FileExists(t, fp)
	defer os.Remove(fp)
	defer unlock()
}
