package vault

import (
	"errors"
	"fmt"
)

// PromoteSecret copies a secret from a source mount/path to a destination
// mount/path, optionally overwriting an existing secret at the destination.
// It returns an error if the source secret cannot be fetched or if the
// destination already exists and overwrite is false.
func PromoteSecret(src, dst Client, srcMount, srcPath, dstMount, dstPath string, overwrite bool) error {
	if srcMount == "" {
		return errors.New("source mount must not be empty")
	}
	if srcPath == "" {
		return errors.New("source path must not be empty")
	}
	if dstMount == "" {
		return errors.New("destination mount must not be empty")
	}
	if dstPath == "" {
		return errors.New("destination path must not be empty")
	}

	data, err := FetchSecret(src, srcMount, srcPath)
	if err != nil {
		return fmt.Errorf("promote: fetch source %s/%s: %w", srcMount, srcPath, err)
	}

	if !overwrite {
		existing, err := FetchSecret(dst, dstMount, dstPath)
		if err == nil && len(existing) > 0 {
			return fmt.Errorf("promote: destination %s/%s already exists; use --overwrite to replace", dstMount, dstPath)
		}
	}

	writePath := kvV2DataPath(dstMount, dstPath)
	_, err = dst.Logical().Write(writePath, map[string]interface{}{
		"data": data,
	})
	if err != nil {
		return fmt.Errorf("promote: write destination %s/%s: %w", dstMount, dstPath, err)
	}

	return nil
}
