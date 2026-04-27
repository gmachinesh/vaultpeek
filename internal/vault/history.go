package vault

import (
	"fmt"
	"time"
)

// VersionMeta holds metadata for a single secret version.
type VersionMeta struct {
	Version      int
	CreatedTime  time.Time
	DeletionTime time.Time
	Destroyed    bool
	Deleted      bool
}

// SecretHistory contains all version metadata for a secret.
type SecretHistory struct {
	Mount    string
	Path     string
	Versions []VersionMeta
}

// GetSecretHistory returns version metadata for all versions of a KV v2 secret.
func GetSecretHistory(client *Client, mount, path string) (*SecretHistory, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}
	if mount == "" {
		return nil, fmt.Errorf("mount is required")
	}
	if path == "" {
		return nil, fmt.Errorf("path is required")
	}

	metaPath := fmt.Sprintf("%s/metadata/%s", mount, path)
	secret, err := client.Logical().Read(metaPath)
	if err != nil {
		return nil, fmt.Errorf("reading metadata: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no metadata found for %s/%s", mount, path)
	}

	versionsRaw, ok := secret.Data["versions"]
	if !ok {
		return &SecretHistory{Mount: mount, Path: path}, nil
	}

	versionsMap, ok := versionsRaw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected versions format")
	}

	history := &SecretHistory{Mount: mount, Path: path}
	for vStr, vRaw := range versionsMap {
		vMap, ok := vRaw.(map[string]interface{})
		if !ok {
			continue
		}
		meta := parseVersionMeta(vStr, vMap)
		history.Versions = append(history.Versions, meta)
	}
	return history, nil
}

func parseVersionMeta(vStr string, m map[string]interface{}) VersionMeta {
	meta := VersionMeta{}
	fmt.Sscanf(vStr, "%d", &meta.Version)

	if v, ok := m["created_time"].(string); ok {
		meta.CreatedTime, _ = time.Parse(time.RFC3339Nano, v)
	}
	if v, ok := m["deletion_time"].(string); ok && v != "" {
		meta.DeletionTime, _ = time.Parse(time.RFC3339Nano, v)
	}
	if v, ok := m["destroyed"].(bool); ok {
		meta.Destroyed = v
	}
	if v, ok := m["deleted"].(bool); ok {
		meta.Deleted = v
	}
	return meta
}
