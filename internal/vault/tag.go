package vault

import (
	"errors"
	"fmt"
)

// TagSecret writes a set of metadata tags (key/value pairs) to a KV v2 secret's
// custom_metadata field using the Vault metadata endpoint.
func TagSecret(client *Client, mount, path string, tags map[string]string) error {
	if client == nil {
		return errors.New("client is nil")
	}
	if mount == "" {
		return errors.New("mount is required")
	}
	if path == "" {
		return errors.New("path is required")
	}
	if len(tags) == 0 {
		return errors.New("tags must not be empty")
	}

	metaPath := fmt.Sprintf("%s/metadata/%s", mount, path)

	body := map[string]interface{}{
		"custom_metadata": tags,
	}

	_, err := client.Logical().Write(metaPath, body)
	if err != nil {
		return fmt.Errorf("writing metadata for %s/%s: %w", mount, path, err)
	}
	return nil
}

// GetTags retrieves the custom_metadata tags from a KV v2 secret's metadata.
func GetTags(client *Client, mount, path string) (map[string]string, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	if mount == "" {
		return nil, errors.New("mount is required")
	}
	if path == "" {
		return nil, errors.New("path is required")
	}

	metaPath := fmt.Sprintf("%s/metadata/%s", mount, path)

	secret, err := client.Logical().Read(metaPath)
	if err != nil {
		return nil, fmt.Errorf("reading metadata for %s/%s: %w", mount, path, err)
	}
	if secret == nil || secret.Data == nil {
		return map[string]string{}, nil
	}

	raw, ok := secret.Data["custom_metadata"]
	if !ok || raw == nil {
		return map[string]string{}, nil
	}

	rawMap, ok := raw.(map[string]interface{})
	if !ok {
		return nil, errors.New("unexpected custom_metadata format")
	}

	tags := make(map[string]string, len(rawMap))
	for k, v := range rawMap {
		if s, ok := v.(string); ok {
			tags[k] = s
		}
	}
	return tags, nil
}
