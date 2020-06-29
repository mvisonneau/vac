package client

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// Without initialization error
	os.Setenv("VAULT_ADDR", "http://localhost:8200")
	os.Setenv("VAULT_TOKEN", "s.xxxxxx")
	c, err := New()
	assert.Nil(t, err)
	assert.IsType(t, &Client{}, c)
	assert.NotNil(t, c)

	// With an initialization error
	os.Unsetenv("VAULT_ADDR")
	c, err = New()
	assert.Error(t, err, "initializing vault client: VAULT_ADDR env is not defined")
	assert.IsType(t, &Client{}, c)
	assert.Nil(t, c)
}
