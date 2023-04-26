package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// TODO: login tests with mocked server

	config := &AuthConfig{
		AuthMethod:     "",
		AuthPath:       "",
		AuthNoStore:    true,
		AuthMethodArgs: map[string]string{},
	}

	c, err := New(config)
	assert.Error(t, err, "initializing vault client")
	assert.IsType(t, &Client{}, c)
	assert.Nil(t, c)
}
