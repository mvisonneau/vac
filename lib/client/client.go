package client

import (
	vault "github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
)

// Client ..
type Client struct {
	*vault.Client
}

// New ..
func New() (*Client, error) {
	vault, err := getVaultClient()
	if err != nil {
		return nil, errors.Wrap(err, "initializing vault client")
	}
	return &Client{vault}, nil
}
