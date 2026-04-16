package vault

import (
	"context"
	"fmt"
	"net/http"
)

// DeleteSecret permanently deletes a secret at the given path in a KV v2 mount.
func DeleteSecret(ctx context.Context, client *Client, mount, secretPath string) error {
	if mount == "" {
		return fmt.Errorf("mount must not be empty")
	}
	if secretPath == "" {
		return fmt.Errorf("secret path must not be empty")
	}

	metaPath := fmt.Sprintf("%s/metadata/%s", mount, secretPath)

	resp, err := client.rawDelete(ctx, metaPath)
	if err != nil {
		return fmt.Errorf("delete request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("secret %q not found at mount %q", secretPath, mount)
	}
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d when deleting secret", resp.StatusCode)
	}

	return nil
}
