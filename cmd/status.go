package cmd

import (
	"os"
	"strconv"
	"time"

	"github.com/mvisonneau/vac/lib/client"
	"github.com/mvisonneau/vac/lib/state"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
	"github.com/xeonx/timeago"
)

// Status ..
func Status(ctx *cli.Context) (int, error) {
	cfg, err := configure(ctx)
	if err != nil {
		return 1, err
	}

	vac, err := client.New()
	if err != nil {
		return 1, err
	}

	s, err := state.Read(cfg.StatePath)
	if err != nil {
		return 1, err
	}

	stateOutput := [][]string{
		{"Current Engine", s.Current.Engine},
		{"Current Role", s.Current.Role},
	}

	stateTable := tablewriter.NewWriter(os.Stdout)
	stateTable.SetHeader([]string{"LOCAL STATE"})
	stateTable.AppendBulk(stateOutput)
	stateTable.Render()

	// Credentials info
	credsTable := tablewriter.NewWriter(os.Stdout)
	credsTable.SetHeader([]string{"ENGINE", "ROLE", "EXPIRATION"})

	for engine, mapRoles := range s.AWSCredentials {
		for role, creds := range mapRoles {
			var color int
			if creds.Metadata.ExpireAt.After(time.Now().Add(time.Minute * 5)) {
				color = tablewriter.FgRedColor
			} else if creds.Metadata.ExpireAt.After(time.Now()) {
				color = tablewriter.FgYellowColor
			} else {
				color = tablewriter.FgRedColor
			}

			credsTable.Rich(
				[]string{
					engine,
					role,
					timeago.English.Format(creds.Metadata.ExpireAt),
				},
				[]tablewriter.Colors{
					{},
					{},
					{color},
				},
			)
		}
	}

	if credsTable.NumLines() > 0 {
		credsTable.Render()
	}

	// Vault status
	health, err := vac.Sys().Health()
	if err != nil {
		return 1, err
	}

	vaultOutput := [][]string{
		{"ClusterID", health.ClusterID},
		{"ClusterName", health.ClusterName},
		{"Initialized", strconv.FormatBool(health.Initialized)},
		{"Sealed", strconv.FormatBool(health.Sealed)},
		{"Version", health.Version},
	}

	vaultTable := tablewriter.NewWriter(os.Stdout)
	vaultTable.SetHeader([]string{"VAULT"})
	vaultTable.AppendBulk(vaultOutput)
	vaultTable.Render()

	return 0, nil
}
