package client

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"syscall"
	"time"

	credAliCloud "github.com/hashicorp/vault-plugin-auth-alicloud"
	credCentrify "github.com/hashicorp/vault-plugin-auth-centrify"
	credCF "github.com/hashicorp/vault-plugin-auth-cf"
	credGcp "github.com/hashicorp/vault-plugin-auth-gcp/plugin"
	credOIDC "github.com/hashicorp/vault-plugin-auth-jwt"
	credKerb "github.com/hashicorp/vault-plugin-auth-kerberos"
	credOCI "github.com/hashicorp/vault-plugin-auth-oci"
	vault "github.com/hashicorp/vault/api"
	credAws "github.com/hashicorp/vault/builtin/credential/aws"
	credCert "github.com/hashicorp/vault/builtin/credential/cert"
	credGitHub "github.com/hashicorp/vault/builtin/credential/github"
	credLdap "github.com/hashicorp/vault/builtin/credential/ldap"
	credOkta "github.com/hashicorp/vault/builtin/credential/okta"
	credToken "github.com/hashicorp/vault/builtin/credential/token"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	vaultCommand "github.com/hashicorp/vault/command"

	"github.com/mvisonneau/vac/internal/base"
	"github.com/mvisonneau/vac/pkg/auth/write"
)

type AuthConfig struct {
	AuthMethod     string
	AuthPath       string
	AuthNoStore    bool
	AuthMethodArgs map[string]string
}

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

func extractToken(client *vault.Client, secret *vault.Secret) (*vault.Secret, error) {
	switch {
	case secret == nil:
		return nil, fmt.Errorf("empty response from auth helper")

	case secret.Auth != nil:
		return secret, nil

	case secret.WrapInfo != nil:
		if secret.WrapInfo.WrappedAccessor == "" {
			return nil, fmt.Errorf("wrapped response does not contain a token")
		}

		client.SetToken(secret.WrapInfo.Token)
		secret, err := client.Logical().Unwrap("")
		if err != nil {
			return nil, err
		}
		return extractToken(client, secret)

	default:
		return nil, fmt.Errorf("no auth or wrapping info in response")
	}
}

func lookupToken(c *vault.Client, token string) (*vault.Secret, error) {
	// If we got this far, we want to lookup and lookup the token and pull its
	// list of policies a metadata.
	c.SetToken(token)
	c.SetWrappingLookupFunc(func(string, string) string { return "" })

	secret, err := c.Auth().Token().LookupSelf()
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
	return &vault.Secret{
		Auth: &vault.SecretAuth{
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

// getVaultClient : Get a Vault client using Vault official params
func getVaultClient(authConfig *AuthConfig) (*vault.Client, error) {
	clientConfig := vault.DefaultConfig()
	if err := clientConfig.ReadEnvironment(); err != nil {
		return nil, fmt.Errorf("failed to read environment: %s", err)
	}

	client, err := vault.NewClient(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Vault client: %s", err)
	}

	tokenHelper, err := vaultCommand.DefaultTokenHelper()
	if err != nil {
		return nil, fmt.Errorf("failed to get token helper: %s", err)
	}

	// Get the token if it came in from the environment
	token := client.Token()

	if token == "" {
		token, err = tokenHelper.Get()
		if err != nil {
			return nil, fmt.Errorf("failed to get token from token helper: %s", err)
		}
	}

	config := authConfig.AuthMethodArgs

	var secret *vault.Secret

	lookup := true
	if x, ok := config["lookup"]; ok {
		parsed, err := strconv.ParseBool(x)
		if err != nil {
			return nil, fmt.Errorf("failed to parse \"lookup\" as boolean: %w", err)
		}
		lookup = parsed
	}

	if lookup {
		secret, err = lookupToken(client, token)
		if err != nil {
			if errors.Is(err, syscall.ECONNREFUSED) {
				return nil, fmt.Errorf("failed to connect to vault: %s", err)
			}
		}
	}

	if secret == nil || secret.Renewable {
		loginHandlers := map[string]vaultCommand.LoginHandler{
			"alicloud": &credAliCloud.CLIHandler{},
			"aws":      &credAws.CLIHandler{},
			"centrify": &credCentrify.CLIHandler{},
			"cert":     &credCert.CLIHandler{},
			"cf":       &credCF.CLIHandler{},
			"gcp":      &credGcp.CLIHandler{},
			"github":   &credGitHub.CLIHandler{},
			"kerberos": &credKerb.CLIHandler{},
			"ldap":     &credLdap.CLIHandler{},
			"oci":      &credOCI.CLIHandler{},
			"oidc":     &credOIDC.CLIHandler{},
			"okta":     &credOkta.CLIHandler{},
			"pcf":      &credCF.CLIHandler{}, // Deprecated.
			"radius": &credUserpass.CLIHandler{
				DefaultMount: "radius",
			},
			"token": &credToken.CLIHandler{},
			"userpass": &credUserpass.CLIHandler{
				DefaultMount: "userpass",
			},
			"write": &write.CLIHandler{}, // used for jwt (gitlab) and approle auth
		}

		// Get the auth method
		authMethod := base.SanitizePath(authConfig.AuthMethod)
		if authMethod == "" {
			authMethod = "token"
		}

		flagPath := authConfig.AuthPath

		// If no path is specified, we default the path to the method type
		// or use the plugin name if it's a plugin
		authPath := flagPath
		if authPath == "" && authMethod != "write" {
			authPath = base.EnsureTrailingSlash(authMethod)
		}

		// Get the handler function
		authHandler, ok := loginHandlers[authMethod]
		if !ok {
			return nil, fmt.Errorf(
				"unknown auth method: %s. Use \"vault auth list\" to see the "+
					"complete list of auth methods. Additionally, some "+
					"auth methods are only available via the HTTP API",
				authMethod)
		}

		// If the user did not specify a mount path, use the provided mount path.
		if config["mount"] == "" && authPath != "" {
			config["mount"] = authPath
		}

		// Evolving token formats across Vault versions have caused issues during CLI logins. Unless
		// token auth is being used, omit any token picked up from TokenHelper.
		if authMethod != "token" {
			client.SetToken("")
		}

		// Authenticate delegation to the auth handler
		secret, err = authHandler.Auth(client, config)
		if err != nil {
			return nil, fmt.Errorf("error authenticating: %s", err)
		}

		// Unset any previous token wrapping functionality. If the original request
		// was for a wrapped token, we don't want future requests to be wrapped.
		client.SetWrappingLookupFunc(func(string, string) string { return "" })

		flagNoStore := authConfig.AuthNoStore

		// Recursively extract the token, handling wrapping
		secret, err = extractToken(client, secret)
		if err != nil {
			return nil, fmt.Errorf("error extracting token: %s", err)
		}
		if secret == nil {
			return nil, fmt.Errorf("vault returned an empty secret")
		}

		if secret.Auth == nil {
			return nil, fmt.Errorf(base.WrapAtLength(
				"Vault returned a secret, but the secret has no authentication " +
					"information attached. This should never happen and is likely a " +
					"bug."))
		}

		// Pull the token itself out, since we don't need the rest of the auth
		// information anymore/.
		token = secret.Auth.ClientToken

		if !flagNoStore {
			// Grab the token helper so we can store
			if err != nil {
				return nil, fmt.Errorf(base.WrapAtLength(fmt.Sprintf(
					"Error initializing token helper. Please verify that the token "+
						"helper is available and properly configured for your system. The "+
						"error was: %s", err)))
			}

			// Store the token in the local client
			if err := tokenHelper.Store(token); err != nil {
				return nil, fmt.Errorf("error storing token: %s", err)
			}
		}

		//home, _ := homedir.Dir()
		//f, err := ioutil.ReadFile(filepath.Clean(home + "/.vault-token"))
		//if err != nil {
		//	return nil, fmt.Errorf("Vault token is not defined (VAULT_TOKEN or ~/.vault-token)")
		//}
		//
		//// The vault client does not handle a trailing newline, so we ensure it
		//// has been removed
		//token = strings.TrimSuffix(string(f), "\n")
	} else {
		token = secret.Auth.ClientToken
	}

	client.SetToken(token)

	return client, nil
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
