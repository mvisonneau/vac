package cmd

import (
	"github.com/ktr0731/go-fuzzyfinder"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/mvisonneau/vac/pkg/client"
	"github.com/mvisonneau/vac/pkg/state"
)

// Switch ..
func Switch(ctx *cli.Context) (int, error) {
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

	if cfg.Engine == "" {
		awsSecretEngines, err := vac.ListAWSSecretEngines()
		if err != nil {
			log.Fatal(err)
		}

		selectedAWSSecretEngineID, err := fuzzyfinder.Find(
			awsSecretEngines,
			func(i int) string {
				return awsSecretEngines[i]
			},
			fuzzyfinder.WithDefaultIndex(indexOf(s.Current.Engine, awsSecretEngines)),
		)
		if err != nil {
			return 1, err
		}
		s.SetCurrentEngine(awsSecretEngines[selectedAWSSecretEngineID])
	} else {
		s.SetCurrentEngine(cfg.Engine)
	}

	if cfg.Role == "" {
		roles, err := vac.ListAWSSecretEngineRoles(s.Current.Engine)
		if err != nil {
			return 1, err
		}

		selectedRoleID, err := fuzzyfinder.Find(
			roles,
			func(i int) string {
				return roles[i]
			},
			fuzzyfinder.WithDefaultIndex(indexOf(s.Current.Role, roles)),
		)
		if err != nil {
			return 1, err
		}
		s.SetCurrentRole(roles[selectedRoleID])
	} else {
		s.SetCurrentRole(cfg.Role)
	}

	// Write state on disk
	if err = state.Write(s, cfg.StatePath); err != nil {
		return 1, err
	}

	log.WithFields(log.Fields{
		"engine": s.Current.Engine,
		"role":   s.Current.Role,
	}).Debugf("updated current engine & role in statefile")

	return 0, nil
}

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1
}
