package vault

import (
	"errors"
	"fmt"
)

// PatchSecret merges the provided updates into an existing secret at mount/path.
// Existing keys not present in updates are preserved. Keys with empty string
// values in updates are removed from the resulting secret.
func PatchSecret(client *Client, mount, path string, updates map[string]string) error {
	if mount == "" {
		return errors.New("mount must not be empty")
	}
	if path == "" {
		return errors.New("path must not be empty")
	}
	if len(updates) == 0 {
		return errors.New("updates must not be empty")
	}

	existing, err := FetchSecret(client, mount, path)
	if err != nil {
		return fmt.Errorf("fetch existing secret: %w", err)
	}

	merged := make(map[string]interface{})

	// Copy existing keys.
	for _, k := range existing.Keys() {
		if v, ok := existing.ValueFor(k); ok {
			merged[k] = v
		}
	}

	// Apply updates: empty value means delete the key.
	for k, v := range updates {
		if v == "" {
			delete(merged, k)
		} else {
			merged[k] = v
		}
	}

	dataPath := kvV2DataPath(mount, path)
	body := map[string]interface{}{"data": merged}

	_, err = client.Logical().Write(dataPath, body)
	if err != nil {
		return fmt.Errorf("write patched secret: %w", err)
	}

	return nil
}
