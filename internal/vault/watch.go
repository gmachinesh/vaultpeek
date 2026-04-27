package vault

import (
	"context"
	"fmt"
	"time"
)

// WatchEvent represents a change detected on a secret.
type WatchEvent struct {
	Mount   string
	Path    string
	OldData SecretMap
	NewData SecretMap
	At      time.Time
}

// WatchOptions configures the Watch behaviour.
type WatchOptions struct {
	Interval time.Duration
}

// DefaultWatchInterval is used when no interval is specified.
const DefaultWatchInterval = 10 * time.Second

// WatchSecret polls a KV v2 secret at the given interval and emits a
// WatchEvent on the returned channel whenever the data changes.
// The caller must cancel ctx to stop watching.
func WatchSecret(ctx context.Context, client *Client, mount, path string, opts WatchOptions) (<-chan WatchEvent, error) {
	if client == nil {
		return nil, fmt.Errorf("watch: client must not be nil")
	}
	if mount == "" {
		return nil, fmt.Errorf("watch: mount must not be empty")
	}
	if path == "" {
		return nil, fmt.Errorf("watch: path must not be empty")
	}

	interval := opts.Interval
	if interval <= 0 {
		interval = DefaultWatchInterval
	}

	events := make(chan WatchEvent, 1)

	go func() {
		defer close(events)

		prev, _ := FetchSecret(client, mount, path)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				curr, err := FetchSecret(client, mount, path)
				if err != nil {
					continue
				}
				if !secretMapsEqual(prev, curr) {
					events <- WatchEvent{
						Mount:   mount,
						Path:    path,
						OldData: prev,
						NewData: curr,
						At:      t,
					}
					prev = curr
				}
			}
		}
	}()

	return events, nil
}

func secretMapsEqual(a, b SecretMap) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
