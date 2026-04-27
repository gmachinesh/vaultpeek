package vault

import (
	"fmt"
	"strings"
)

// SearchResult holds a matched secret path and the matching key within it.
type SearchResult struct {
	Mount string
	Path  string
	Key   string
	Value string
}

// SearchSecrets searches all secrets under a mount/prefix for keys or values
// matching the given query string (case-insensitive substring match).
func SearchSecrets(client *Client, mount, prefix, query string, searchValues bool) ([]SearchResult, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}
	if mount == "" {
		return nil, fmt.Errorf("mount is required")
	}
	if query == "" {
		return nil, fmt.Errorf("query is required")
	}

	paths, err := ListSecrets(client, mount, prefix)
	if err != nil {
		return nil, fmt.Errorf("listing secrets: %w", err)
	}

	lower := strings.ToLower(query)
	var results []SearchResult

	for _, p := range paths {
		sm, err := FetchSecret(client, mount, p)
		if err != nil {
			continue
		}
		for _, k := range sm.Keys() {
			keyMatch := strings.Contains(strings.ToLower(k), lower)
			valMatch := false
			if searchValues {
				v, _ := sm.ValueFor(k)
				valMatch = strings.Contains(strings.ToLower(v), lower)
			}
			if keyMatch || valMatch {
				v, _ := sm.ValueFor(k)
				results = append(results, SearchResult{
					Mount: mount,
					Path:  p,
					Key:   k,
					Value: v,
				})
			}
		}
	}

	return results, nil
}
