package vault

import (
	"fmt"
	"time"
)

// Snapshot represents a point-in-time capture of all keys and values at a path prefix.
type Snapshot struct {
	Mount     string
	Prefix    string
	TakenAt   time.Time
	Secrets   map[string]SecretMap // key: relative path, value: secret data
}

// SnapshotResult holds the outcome of a snapshot operation.
type SnapshotResult struct {
	Snapshot *Snapshot
	Count    int
}

// TakeSnapshot recursively lists and fetches all secrets under prefix, returning a Snapshot.
func TakeSnapshot(client *Client, mount, prefix string) (*SnapshotResult, error) {
	if client == nil {
		return nil, fmt.Errorf("client must not be nil")
	}
	if mount == "" {
		return nil, fmt.Errorf("mount must not be empty")
	}
	if prefix == "" {
		return nil, fmt.Errorf("prefix must not be empty")
	}

	keys, err := ListSecrets(client, mount, prefix)
	if err != nil {
		return nil, fmt.Errorf("listing secrets: %w", err)
	}

	snapshot := &Snapshot{
		Mount:   mount,
		Prefix:  prefix,
		TakenAt: time.Now().UTC(),
		Secrets: make(map[string]SecretMap, len(keys)),
	}

	for _, key := range keys {
		path := prefix + "/" + key
		sm, err := FetchSecret(client, mount, path)
		if err != nil {
			return nil, fmt.Errorf("fetching secret %q: %w", path, err)
		}
		snapshot.Secrets[key] = sm
	}

	return &SnapshotResult{
		Snapshot: snapshot,
		Count:    len(snapshot.Secrets),
	}, nil
}
