package client

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	vault "github.com/hashicorp/vault/api"
	"github.com/mitchellh/go-homedir"
)

// AWSCredentials ..
type AWSCredentials struct {
	Metadata struct {
		CreatedAt time.Time `json:"created_at"`
		ExpireAt  time.Time `json:"expire_at"`
	} `json:"metadata"`
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	SecurityToken   string `json:"security_token"`
}

// getVaultClient : Get a Vault client using Vault official params
func getVaultClient() (*vault.Client, error) {
	c, err := vault.NewClient(nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating Vault client: %s", err.Error())
	}

	if len(os.Getenv("VAULT_ADDR")) == 0 {
		return nil, fmt.Errorf("VAULT_ADDR env is not defined")
	}

	if err := c.SetAddress(os.Getenv("VAULT_ADDR")); err != nil {
		return nil, fmt.Errorf("error settings vault client addr: %w", err)
	}

	token := os.Getenv("VAULT_TOKEN")
	if len(token) == 0 {
		home, _ := homedir.Dir()
		f, err := ioutil.ReadFile(filepath.Clean(home + "/.vault-token"))
		if err != nil {
			return nil, fmt.Errorf("Vault token is not defined (VAULT_TOKEN or ~/.vault-token)")
		}

		token = string(f)
	}

	c.SetToken(token)

	return c, nil
}

// ListAWSSecretEngines ..
func (c *Client) ListAWSSecretEngines() (engines []string, err error) {
	mounts, err := c.Sys().ListMounts()
	if err != nil {
		return
	}

	for mountName, mountSpec := range mounts {
		if mountSpec.Type == "aws" {
			engines = append(engines, strings.TrimSuffix(mountName, "/"))
		}
	}
	return
}

// ListAWSSecretEngineRoles ..
func (c *Client) ListAWSSecretEngineRoles(awsSecretEngine string) (roles []string, err error) {
	var foundRoles *vault.Secret
	foundRoles, err = c.Logical().List(fmt.Sprintf("/%s/roles", awsSecretEngine))
	if err != nil {
		return
	}

	if foundRoles != nil && foundRoles.Data != nil {
		if _, ok := foundRoles.Data["keys"]; ok {
			for _, role := range foundRoles.Data["keys"].([]interface{}) {
				roles = append(roles, role.(string))
			}
		}
	}

	return
}

// GenerateAWSCredentials ..
func (c *Client) GenerateAWSCredentials(secretEngineName, secretEngineRole string, ttl time.Duration) (creds *AWSCredentials, err error) {
	var output *vault.Secret
	payload := make(map[string]interface{})
	if ttl > 0 {
		payload["ttl"] = ttl.Seconds()
	}
	output, err = c.Logical().Write(fmt.Sprintf("/%s/sts/%s", secretEngineName, secretEngineRole), payload)
	if err != nil {
		return
	}

	creds = &AWSCredentials{}
	creds.Metadata.CreatedAt = time.Now()

	if leaseDuration, err := time.ParseDuration(fmt.Sprintf("%ds", output.LeaseDuration)); err == nil {
		creds.Metadata.ExpireAt = creds.Metadata.CreatedAt.Add(leaseDuration)
	} else {
		return creds, err
	}

	if output != nil && output.Data != nil {
		if _, ok := output.Data["access_key"]; ok {
			creds.AccessKeyID = output.Data["access_key"].(string)
		}

		if _, ok := output.Data["secret_key"]; ok {
			creds.SecretAccessKey = output.Data["secret_key"].(string)
		}

		if _, ok := output.Data["security_token"]; ok {
			creds.SecurityToken = output.Data["security_token"].(string)
		}
	}

	return
}
