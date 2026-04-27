package vault

import (
	"errors"
	"fmt"
	"time"
)

// TouchResult holds metadata about a secret touch operation.
type TouchResult struct {
	Mount   string
	Path    string
	Version int
	TouchedAt time.Time
}

// TouchSecret re-writes a secret with its existing data to create a new
// version, effectively "touching" it. This is useful for resetting TTLs
// or triggering audit events without changing secret values.
func TouchSecret(client *Client, mount, path string) (*TouchResult, error) {
	if client == nil {
		return nil, errors.New("vault: client is nil")
	}
	if mount == "" {
		return nil, errors.New("vault: mount is required")
	}
	if path == "" {
		return nil, errors.New("vault: path is required")
	}

	sm, err := FetchSecret(client, mount, path)
	if err != nil {
		return nil, fmt.Errorf("vault: touch fetch error: %w", err)
	}

	data := make(map[string]interface{}, len(sm))
	for _, k := range sm.Keys() {
		v, _ := sm.ValueFor(k)
		data[k] = v
	}

	dataPath := kvV2DataPath(mount, path)
	secret, err := client.Logical().Write(dataPath, map[string]interface{}{"data": data})
	if err != nil {
		return nil, fmt.Errorf("vault: touch write error: %w", err)
	}

	version := 0
	if secret != nil && secret.Data != nil {
		if meta, ok := secret.Data["version"]; ok {
			if v, ok := meta.(float64); ok {
				version = int(v)
			}
		}
	}

	return &TouchResult{
		Mount:     mount,
		Path:      path,
		Version:   version,
		TouchedAt: time.Now().UTC(),
	}, nil
}
