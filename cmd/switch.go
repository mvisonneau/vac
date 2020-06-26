package cmd

import (
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/mvisonneau/vac/lib/client"
	"github.com/mvisonneau/vac/lib/state"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// Switch ..
func Switch(ctx *cli.Context) (int, error) {
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

	roles, err := vac.ListAWSSecretEngineRoles(awsSecretEngines[selectedAWSSecretEngineID])
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

	s.SetCurrentEngine(awsSecretEngines[selectedAWSSecretEngineID])
	s.SetCurrentRole(roles[selectedRoleID])
	if err = state.Write(s, cfg.StatePath); err != nil {
		return 1, err
	}

	log.Infof("engine: %s\n", s.Current.Engine)
	log.Infof("role: %s\n", s.Current.Role)

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
