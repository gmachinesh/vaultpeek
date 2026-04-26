package vault

import (
	"errors"
	"fmt"
	"time"
)

const lockMetaKey = "__vaultpeek_locked"

// LockSecret marks a secret as locked by writing a metadata tag that
// prevents accidental writes. It stores the lock timestamp and optional reason.
func LockSecret(client *Client, mount, path, reason string) error {
	if client == nil {
		return errors.New("client is nil")
	}
	if mount == "" {
		return errors.New("mount is required")
	}
	if path == "" {
		return errors.New("path is required")
	}

	tags := map[string]string{
		lockMetaKey: fmt.Sprintf("%s|%s", time.Now().UTC().Format(time.RFC3339), reason),
	}

	return TagSecret(client, mount, path, tags)
}

// UnlockSecret removes the lock tag from a secret.
func UnlockSecret(client *Client, mount, path string) error {
	if client == nil {
		return errors.New("client is nil")
	}
	if mount == "" {
		return errors.New("mount is required")
	}
	if path == "" {
		return errors.New("path is required")
	}

	tags, err := GetTags(client, mount, path)
	if err != nil {
		return fmt.Errorf("fetch tags: %w", err)
	}

	if _, locked := tags[lockMetaKey]; !locked {
		return errors.New("secret is not locked")
	}

	delete(tags, lockMetaKey)
	return TagSecret(client, mount, path, tags)
}

// IsLocked reports whether a secret currently carries a lock tag.
func IsLocked(client *Client, mount, path string) (bool, error) {
	if client == nil {
		return false, errors.New("client is nil")
	}
	if mount == "" {
		return false, errors.New("mount is required")
	}
	if path == "" {
		return false, errors.New("path is required")
	}

	tags, err := GetTags(client, mount, path)
	if err != nil {
		return false, fmt.Errorf("fetch tags: %w", err)
	}

	_, locked := tags[lockMetaKey]
	return locked, nil
}
