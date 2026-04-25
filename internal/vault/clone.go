package vault

import (
	"errors"
	"fmt"
)

// CloneOptions controls the behaviour of CloneSecret.
type CloneOptions struct {
	// Overwrite allows the destination path to be overwritten if it already
	// contains data. When false, CloneSecret returns an error if the
	// destination already exists.
	Overwrite bool
}

// CloneSecret copies every key/value pair from srcMount/srcPath into
// dstMount/dstPath, optionally merging with existing data at the
// destination when opts.Overwrite is true.
//
// Unlike CopySecret, CloneSecret accepts an explicit options struct and
// performs an existence check before writing.
func CloneSecret(c *Client, srcMount, srcPath, dstMount, dstPath string, opts CloneOptions) error {
	if c == nil {
		return errors.New("vault: client must not be nil")
	}
	if srcMount == "" {
		return errors.New("vault: source mount must not be empty")
	}
	if srcPath == "" {
		return errors.New("vault: source path must not be empty")
	}
	if dstMount == "" {
		return errors.New("vault: destination mount must not be empty")
	}
	if dstPath == "" {
		return errors.New("vault: destination path must not be empty")
	}
	if srcMount == dstMount && srcPath == dstPath {
		return errors.New("vault: source and destination paths must differ")
	}

	src, err := FetchSecret(c, srcMount, srcPath)
	if err != nil {
		return fmt.Errorf("vault: fetch source secret: %w", err)
	}

	if !opts.Overwrite {
		existing, err := FetchSecret(c, dstMount, dstPath)
		if err == nil && len(existing) > 0 {
			return fmt.Errorf("vault: destination %s/%s already exists; use --overwrite to replace", dstMount, dstPath)
		}
	}

	dataPath := kvV2DataPath(dstMount, dstPath)
	_, err = c.Logical().Write(dataPath, map[string]interface{}{"data": src})
	if err != nil {
		return fmt.Errorf("vault: write destination secret: %w", err)
	}
	return nil
}
