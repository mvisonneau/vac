package client

import (
	"os"
	"testing"

	vault "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
)

func TestGetVaultClient(t *testing.T) {
	// With sufficient configuration
	os.Setenv("VAULT_ADDR", "http://localhost:8200")
	os.Setenv("VAULT_TOKEN", "s.xxxxxx")
	c, err := getVaultClient()
	assert.Nil(t, err)
	assert.IsType(t, &vault.Client{}, c)
	assert.NotNil(t, c)

	// Without VAULT_ADDR defined
	os.Unsetenv("VAULT_ADDR")
	c, err = getVaultClient()
	assert.Error(t, err, "initializing vault client: VAULT_ADDR env is not defined")
	assert.IsType(t, &vault.Client{}, c)
	assert.Nil(t, c)
}
