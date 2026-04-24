package vault

import (
	"context"
	"errors"
	"fmt"
)

// VersionedSecret holds a specific version of a secret's data.
type VersionedSecret struct {
	Version int
	Data    map[string]interface{}
}

// GetSecretVersion fetches a specific version of a KV v2 secret.
func GetSecretVersion(ctx context.Context, client *Client, mount, path string, version int) (*VersionedSecret, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	if mount == "" {
		return nil, errors.New("mount is required")
	}
	if path == "" {
		return nil, errors.New("path is required")
	}
	if version < 1 {
		return nil, errors.New("version must be >= 1")
	}

	vaultPath := kvV2DataPath(mount, path)
	params := map[string][]string{
		"version": {fmt.Sprintf("%d", version)},
	}

	secret, err := client.Logical().ReadWithDataWithContext(ctx, vaultPath, params)
	if err != nil {
		return nil, fmt.Errorf("reading version %d of %s/%s: %w", version, mount, path, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("version %d of %s/%s not found", version, mount, path)
	}

	rawData, ok := secret.Data["data"]
	if !ok {
		return nil, fmt.Errorf("no data field in version %d of %s/%s", version, mount, path)
	}
	data, ok := rawData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected data format in version %d of %s/%s", version, mount, path)
	}

	return &VersionedSecret{Version: version, Data: data}, nil
}

// RollbackSecret restores a KV v2 secret to a previously stored version.
func RollbackSecret(ctx context.Context, client *Client, mount, path string, version int) error {
	versioned, err := GetSecretVersion(ctx, client, mount, path, version)
	if err != nil {
		return fmt.Errorf("fetching version for rollback: %w", err)
	}

	vaultPath := kvV2DataPath(mount, path)
	body := map[string]interface{}{
		"data": versioned.Data,
	}

	_, err = client.Logical().WriteWithContext(ctx, vaultPath, body)
	if err != nil {
		return fmt.Errorf("writing rollback data to %s/%s: %w", mount, path, err)
	}
	return nil
}
