package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	cli "github.com/urfave/cli/v2"
	"github.com/xeonx/timeago"

	"github.com/mvisonneau/vac/pkg/client"
	"github.com/mvisonneau/vac/pkg/state"
)

// Status ..
func Status(ctx *cli.Context) (int, error) {
	cfg, err := configure(ctx)
	if err != nil {
		return 1, err
	}

	vac, err := client.New(cfg.AuthConfig)
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

	for _, engine := range s.GetCachedEngines() {
		for _, role := range s.GetCachedEngineRoles(engine) {
			creds := s.AWSCredentials[engine][role]
			var color int
			if creds.Metadata.ExpireAt.After(time.Now().Add(time.Minute * 5)) {
				color = tablewriter.FgGreenColor
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

	// Vault token status
	secret, err := vac.Auth().Token().LookupSelf()
	if err != nil {
		return 1, err
	}

	tokenId, _ := secret.TokenID()
	tokenAccessor, _ := secret.TokenAccessor()

	_, err = secret.TokenPolicies()
	if err != nil {
		return 1, fmt.Errorf("error accessing token policies: %w", err)
	}
	tokenPolicies := strings.Join(secret.Auth.Policies, ", ")
	tokenTokenPolicies := strings.Join(secret.Auth.TokenPolicies, ", ")
	tokenIdentityPolicies := strings.Join(secret.Auth.IdentityPolicies, ", ")

	tokenDuration, _ := secret.TokenTTL()
	tokenRenewable, _ := secret.TokenIsRenewable()

	vaultOutput := [][]string{
		{"ClusterID", health.ClusterID},
		{"ClusterName", health.ClusterName},
		{"Initialized", strconv.FormatBool(health.Initialized)},
		{"Sealed", strconv.FormatBool(health.Sealed)},
		{"Version", health.Version},
		{"ClientToken", tokenId},
		{"Accessor", tokenAccessor},
		{"Policies", tokenPolicies},
		{"TokenPolicies", tokenTokenPolicies},
		{"IdentityPolicies", tokenIdentityPolicies},
		{"LeaseDuration", tokenDuration.String()},
		{"Renewable", strconv.FormatBool(tokenRenewable)},
	}

	vaultTable := tablewriter.NewWriter(os.Stdout)
	vaultTable.SetHeader([]string{"VAULT"})
	vaultTable.AppendBulk(vaultOutput)
	vaultTable.Render()

	return 0, nil
}
