package vault

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// ExportFormat defines the output format for exported secrets.
type ExportFormat string

const (
	FormatJSON ExportFormat = "json"
	FormatEnv  ExportFormat = "env"
)

// ExportSecret writes the key/value pairs of a secret at the given mount and
// path to w in the requested format.
func ExportSecret(client *Client, mount, path string, format ExportFormat, w io.Writer) error {
	if client == nil {
		return fmt.Errorf("client must not be nil")
	}
	if mount == "" {
		return fmt.Errorf("mount must not be empty")
	}
	if path == "" {
		return fmt.Errorf("path must not be empty")
	}

	sm, err := FetchSecret(client, mount, path)
	if err != nil {
		return fmt.Errorf("fetch secret: %w", err)
	}

	keys := sm.Keys()
	sort.Strings(keys)

	switch format {
	case FormatJSON:
		return exportJSON(sm, keys, w)
	case FormatEnv:
		return exportEnv(sm, keys, w)
	default:
		return fmt.Errorf("unsupported format: %q", format)
	}
}

func exportJSON(sm SecretMap, keys []string, w io.Writer) error {
	m := make(map[string]string, len(keys))
	for _, k := range keys {
		m[k] = sm.ValueFor(k)
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(m)
}

func exportEnv(sm SecretMap, keys []string, w io.Writer) error {
	for _, k := range keys {
		if _, err := fmt.Fprintf(w, "%s=%q\n", k, sm.ValueFor(k)); err != nil {
			return err
		}
	}
	return nil
}
