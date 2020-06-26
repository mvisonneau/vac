package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mvisonneau/vac/lib/client"
	"github.com/mvisonneau/vac/lib/state"
	"github.com/urfave/cli"
)

// Output ..
type Output struct {
	Version         int
	AccessKeyID     string `json:"AccessKeyId"`
	SecretAccessKey string
	SessionToken    string
	Expiration      time.Time
}

// Run ..
func Run(ctx *cli.Context) (int, error) {
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

	// Check if a engine is currently selected
	if s.Current.Engine == "" || s.Current.Role == "" {
		if code, err := Switch(ctx); code != 0 {
			return code, err
		}
	}

	creds := s.GetCurrentAWSCredentials()
	if creds == nil || time.Now().After(creds.Metadata.ExpireAt) {
		creds, err = vac.GenerateAWSCredentials(s.Current.Engine, s.Current.Role)
		if err != nil {
			return 1, err
		}

		s.SetCurrentAWSCredentials(creds)
		if err = state.Write(s, cfg.StatePath); err != nil {
			return 1, err
		}
	}

	o := Output{
		Version:         1,
		AccessKeyID:     creds.AccessKeyID,
		SecretAccessKey: creds.SecretAccessKey,
		SessionToken:    creds.SecurityToken,
		Expiration:      creds.Metadata.ExpireAt,
	}

	outputBytes, err := json.Marshal(o)
	if err != nil {
		return 1, err
	}

	fmt.Println(string(outputBytes))
	return 0, nil
}
