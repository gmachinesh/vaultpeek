package vault

import (
	"context"
	"fmt"
	"strings"
)

// ListSecrets returns all secret keys under the given path prefix.
func ListSecrets(ctx context.Context, client *Client, mount, prefix string) ([]string, error) {
	listPath := fmt.Sprintf("%s/metadata/%s", mount, strings.Trim(prefix, "/"))

	secret, err := client.Logical().ListWithContext(ctx, listPath)
	if err != nil {
		return nil, fmt.Errorf("listing secrets at %q: %w", listPath, err)
	}
	if secret == nil || secret.Data == nil {
		return []string{}, nil
	}

	raw, ok := secret.Data["keys"]
	if !ok {
		return []string{}, nil
	}

	ifaces, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type for keys field")
	}

	keys := make([]string, 0, len(ifaces))
	for _, v := range ifaces {
		s, ok := v.(string)
		if !ok {
			continue
		}
		keys = append(keys, s)
	}
	return keys, nil
}
