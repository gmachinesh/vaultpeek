package vault

import (
	"fmt"
)

// RenameSecret copies a secret from srcPath to dstPath within the same mount,
// then deletes the original.
func RenameSecret(client *Client, mount, srcPath, dstPath string) error {
	if mount == "" {
		return fmt.Errorf("mount must not be empty")
	}
	if srcPath == "" {
		return fmt.Errorf("source path must not be empty")
	}
	if dstPath == "" {
		return fmt.Errorf("destination path must not be empty")
	}
	if srcPath == dstPath {
		return fmt.Errorf("source and destination paths must differ")
	}

	if err := CopySecret(client, mount, srcPath, mount, dstPath); err != nil {
		return fmt.Errorf("rename copy step failed: %w", err)
	}

	if err := DeleteSecret(client, mount, srcPath); err != nil {
		return fmt.Errorf("rename delete step failed: %w", err)
	}

	return nil
}
