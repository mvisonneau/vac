package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"

	"github.com/mvisonneau/vac/lib/client"
	"github.com/mvisonneau/vac/lib/state"
)

// Output ..
type Output struct {
	Version         int
	AccessKeyID     string `json:"AccessKeyId"`
	SecretAccessKey string
	SessionToken    string
	Expiration      time.Time
}

// GetConfig ..
type GetConfig struct {
	*Config
	TTL           time.Duration
	MinTTL        time.Duration
	ForceGenerate bool
}

// Get the credentials from Vault or statefile
func Get(ctx *cli.Context) (int, error) {
	globalCfg, err := configure(ctx)
	if err != nil {
		return 1, err
	}

	cfg := &GetConfig{
		Config:        globalCfg,
		TTL:           ctx.Duration("ttl"),
		MinTTL:        ctx.Duration("min-ttl"),
		ForceGenerate: ctx.Bool("force-generate"),
	}

	if cfg.TTL > 0 && cfg.MinTTL > cfg.TTL {
		return 1, fmt.Errorf("'min-ttl' cannot be longer than 'ttl'")
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
	if (s.Current.Engine == "" && cfg.Engine == "") ||
		(s.Current.Role == "" && cfg.Role == "") {
		if code, err := Switch(ctx); code != 0 {
			return code, err
		}
		// Reload the state
		s, err = state.Read(cfg.StatePath)
		if err != nil {
			return 1, err
		}
	}

	var engine, role string
	if cfg.Engine != "" {
		engine = cfg.Engine
	} else {
		engine = s.Current.Engine
	}

	if cfg.Role != "" {
		role = cfg.Role
	} else {
		role = s.Current.Role
	}

	creds := s.GetAWSCredentials(engine, role)
	if shouldRegenerateCreds(creds, cfg.MinTTL, cfg.ForceGenerate) {
		creds, err = vac.GenerateAWSCredentials(engine, role, cfg.TTL)
		if err != nil {
			return 1, err
		}

		s.SetAWSCredentials(engine, role, creds)
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

func shouldRegenerateCreds(creds *client.AWSCredentials, minTTL time.Duration, forceGenerate bool) bool {
	if creds == nil {
		log.Debug("no cached credentials found, generating new ones")
		return true
	}

	if time.Now().After(creds.Metadata.ExpireAt) {
		log.Debug("cached credentials expired, generating new ones")
		return true
	}

	if minTTL > 0 && time.Now().Add(minTTL).After(creds.Metadata.ExpireAt) {
		log.Debug("cached credentials expiring before the defined minimum TTL, generating new ones")
		return true
	}

	if forceGenerate {
		log.Debug("valid creds found but force-generate flag is set, generating new ones")
		return true
	}

	log.Debug("using cached credentials")
	return false
}
