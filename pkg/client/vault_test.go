package client

import (
	vault "github.com/hashicorp/vault/api"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVaultClient(t *testing.T) {
	// TODO: login tests with mocked server

	config := &AuthConfig{
		AuthMethod:     "",
		AuthPath:       "",
		AuthNoStore:    true,
		AuthMethodArgs: map[string]string{},
	}

	c, err := getVaultClient(config)
	assert.Error(t, err, "initializing vault client")
	assert.IsType(t, &vault.Client{}, c)
	assert.Nil(t, c)
}
