package vault

import (
	"errors"
	"fmt"
)

// RotateResult holds the outcome of a secret rotation.
type RotateResult struct {
	Mount   string
	Path    string
	OldKeys []string
	NewKeys []string
}

// RotateSecret replaces the values of all keys in a secret with the provided
// generator function. The generator receives the key name and returns the new
// value. The old secret is preserved in memory and returned via RotateResult.
func RotateSecret(client *Client, mount, path string, generator func(key string) string) (*RotateResult, error) {
	if client == nil {
		return nil, errors.New("rotate: client must not be nil")
	}
	if mount == "" {
		return nil, errors.New("rotate: mount must not be empty")
	}
	if path == "" {
		return nil, errors.New("rotate: path must not be empty")
	}
	if generator == nil {
		return nil, errors.New("rotate: generator must not be nil")
	}

	secret, err := FetchSecret(client, mount, path)
	if err != nil {
		return nil, fmt.Errorf("rotate: fetch %s/%s: %w", mount, path, err)
	}

	oldKeys := secret.Keys()
	newData := make(map[string]interface{}, len(oldKeys))
	for _, k := range oldKeys {
		newData[k] = generator(k)
	}

	writePath := kvV2DataPath(mount, path)
	_, err = client.Logical().Write(writePath, map[string]interface{}{
		"data": newData,
	})
	if err != nil {
		return nil, fmt.Errorf("rotate: write %s/%s: %w", mount, path, err)
	}

	newKeys := make([]string, 0, len(newData))
	for k := range newData {
		newKeys = append(newKeys, k)
	}

	return &RotateResult{
		Mount:   mount,
		Path:    path,
		OldKeys: oldKeys,
		NewKeys: newKeys,
	}, nil
}
