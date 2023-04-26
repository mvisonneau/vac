package write

import (
	"fmt"
	"io"
	"strconv"

	"github.com/hashicorp/vault/api"

	"github.com/mvisonneau/vac/internal/base"
)

type CLIHandler struct {
	// for tests
	testStdin  io.Reader
	testStdout io.Writer
}

func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (*api.Secret, error) {
	pathArg := m["mount"]

	path := base.SanitizePath(pathArg)

	var data map[string]interface{}

	for k, v := range m {
		if k != "mount" && k != "lookup" {
			data[k] = v
		}
	}

	secret, err := c.Logical().Write(path, data)
	if err != nil {
		return nil, fmt.Errorf("failed to write data to path: %w", err)
	}

	token := secret.Auth.ClientToken

	// Parse "lookup" first - we want to return an early error if the user
	// supplied an invalid value here before we prompt them for a token. It would
	// be annoying to type your token and then be told you supplied an invalid
	// value that we could have known in advance.
	lookup := true
	if x, ok := m["lookup"]; ok {
		parsed, err := strconv.ParseBool(x)
		if err != nil {
			return nil, fmt.Errorf("failed to parse \"lookup\" as boolean: %w", err)
		}
		lookup = parsed
	}

	// If the user declined verification, return now. Note that we will not have
	// a lot of information about the token.
	if !lookup {
		return &api.Secret{
			Auth: &api.SecretAuth{
				ClientToken: token,
			},
		}, nil
	}

	// If we got this far, we want to lookup and lookup the token and pull it's
	// list of policies an metadata.
	c.SetToken(token)
	c.SetWrappingLookupFunc(func(string, string) string { return "" })

	secret, err = c.Auth().Token().LookupSelf()
	if err != nil {
		return nil, fmt.Errorf("error looking up token: %w", err)
	}
	if secret == nil {
		return nil, fmt.Errorf("empty response from lookup-self")
	}

	// Return an auth struct that "looks" like the response from an auth method.
	// lookup and lookup-self return their data in data, not auth. We try to
	// mirror that data here.
	id, err := secret.TokenID()
	if err != nil {
		return nil, fmt.Errorf("error accessing token ID: %w", err)
	}
	accessor, err := secret.TokenAccessor()
	if err != nil {
		return nil, fmt.Errorf("error accessing token accessor: %w", err)
	}
	// This populates secret.Auth
	_, err = secret.TokenPolicies()
	if err != nil {
		return nil, fmt.Errorf("error accessing token policies: %w", err)
	}
	metadata, err := secret.TokenMetadata()
	if err != nil {
		return nil, fmt.Errorf("error accessing token metadata: %w", err)
	}
	dur, err := secret.TokenTTL()
	if err != nil {
		return nil, fmt.Errorf("error converting token TTL: %w", err)
	}
	renewable, err := secret.TokenIsRenewable()
	if err != nil {
		return nil, fmt.Errorf("error checking if token is renewable: %w", err)
	}
	return &api.Secret{
		Auth: &api.SecretAuth{
			ClientToken:      id,
			Accessor:         accessor,
			Policies:         secret.Auth.Policies,
			TokenPolicies:    secret.Auth.TokenPolicies,
			IdentityPolicies: secret.Auth.IdentityPolicies,
			Metadata:         metadata,

			LeaseDuration: int(dur.Seconds()),
			Renewable:     renewable,
		},
	}, nil
}

func (h *CLIHandler) Help() string {
	return ""
}
