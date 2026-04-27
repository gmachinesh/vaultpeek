package vault

import (
	"context"
	"testing"
	"time"
)

func TestWatchSecret_NilClient(t *testing.T) {
	_, err := WatchSecret(context.Background(), nil, "secret", "myapp/config", WatchOptions{})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestWatchSecret_EmptyMount(t *testing.T) {
	c := &Client{}
	_, err := WatchSecret(context.Background(), c, "", "myapp/config", WatchOptions{})
	if err == nil {
		t.Fatal("expected error for empty mount")
	}
}

func TestWatchSecret_EmptyPath(t *testing.T) {
	c := &Client{}
	_, err := WatchSecret(context.Background(), c, "secret", "", WatchOptions{})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestWatchSecret_CancelStopsChannel(t *testing.T) {
	c := &Client{}
	ctx, cancel := context.WithCancel(context.Background())

	ch, err := WatchSecret(ctx, c, "secret", "app/cfg", WatchOptions{Interval: 50 * time.Millisecond})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cancel()

	select {
	case _, ok := <-ch:
		if ok {
			t.Fatal("expected channel to be closed after cancel")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("channel was not closed after context cancellation")
	}
}

func TestSecretMapsEqual_Equal(t *testing.T) {
	a := SecretMap{"key": "val"}
	b := SecretMap{"key": "val"}
	if !secretMapsEqual(a, b) {
		t.Fatal("expected maps to be equal")
	}
}

func TestSecretMapsEqual_Different(t *testing.T) {
	a := SecretMap{"key": "old"}
	b := SecretMap{"key": "new"}
	if secretMapsEqual(a, b) {
		t.Fatal("expected maps to differ")
	}
}

func TestSecretMapsEqual_DifferentLengths(t *testing.T) {
	a := SecretMap{"key": "val"}
	b := SecretMap{"key": "val", "extra": "x"}
	if secretMapsEqual(a, b) {
		t.Fatal("expected maps to differ due to length")
	}
}

func TestWatchOptions_DefaultInterval(t *testing.T) {
	if DefaultWatchInterval != 10*time.Second {
		t.Fatalf("expected default interval 10s, got %v", DefaultWatchInterval)
	}
}
