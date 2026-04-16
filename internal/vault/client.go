package vault

import (
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with environment metadata.
type Client struct {
	api  *vaultapi.Client
	Env  string
}

// Config holds connection parameters for a Vault instance.
type Config struct {
	Address string
	Token   string
	Env     string
}

// NewClient creates a new Vault client from the given config.
func NewClient(cfg Config) (*Client, error) {
	if cfg.Address == "" {
		cfg.Address = os.Getenv("VAULT_ADDR")
	}
	if cfg.Token == "" {
		cfg.Token = os.Getenv("VAULT_TOKEN")
	}
	if cfg.Address == "" {
		return nil, fmt.Errorf("vault address is required (set VAULT_ADDR or pass --address)")
	}
	if cfg.Token == "" {
		return nil, fmt.Errorf("vault token is required (set VAULT_TOKEN or pass --token)")
	}

	apiCfg := vaultapi.DefaultConfig()
	apiCfg.Address = cfg.Address

	c, err := vaultapi.NewClient(apiCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}
	c.SetToken(cfg.Token)

	return &Client{api: c, Env: cfg.Env}, nil
}

// ReadSecret reads a KV v2 secret at the given path.
func (c *Client) ReadSecret(path string) (map[string]interface{}, error) {
	secret, err := c.api.Logical().Read(kvV2DataPath(path))
	if err != nil {
		return nil, fmt.Errorf("reading secret at %q: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no secret found at path %q", path)
	}
	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected secret data format at %q", path)
	}
	return data, nil
}

// kvV2DataPath converts a mount-relative path to the KV v2 data endpoint.
func kvV2DataPath(path string) string {
	return "secret/data/" + path
}
