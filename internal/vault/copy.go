package vault

import (
	"context"
	"fmt"
)

// CopySecret reads a secret from srcPath and writes all its key/value pairs to dstPath.
func CopySecret(ctx context.Context, client *Client, srcPath, dstPath string) error {
	secret, err := FetchSecret(ctx, client, srcPath)
	if err != nil {
		return fmt.Errorf("copy: fetch source %q: %w", srcPath, err)
	}

	data := make(map[string]interface{}, len(secret))
	for _, k := range secret.Keys() {
		v, _ := secret.ValueFor(k)
		data[k] = v
	}

	writePath := kvV2DataPath(dstPath)
	_, err = client.http.R().
		SetContext(ctx).
		SetBody(map[string]interface{}{"data": data}).
		Post(writePath)
	if err != nil {
		return fmt.Errorf("copy: write destination %q: %w", dstPath, err)
	}

	return nil
}
