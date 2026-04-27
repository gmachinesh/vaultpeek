package vault

import (
	"testing"
	"time"
)

func TestWatchEvent_Fields(t *testing.T) {
	now := time.Now()
	old := SecretMap{"key": "old_value"}
	new_ := SecretMap{"key": "new_value"}

	ev := WatchEvent{
		Mount:   "secret",
		Path:    "app/db",
		OldData: old,
		NewData: new_,
		At:      now,
	}

	if ev.Mount != "secret" {
		t.Errorf("expected mount 'secret', got %q", ev.Mount)
	}
	if ev.Path != "app/db" {
		t.Errorf("expected path 'app/db', got %q", ev.Path)
	}
	if ev.OldData["key"] != "old_value" {
		t.Errorf("unexpected OldData: %v", ev.OldData)
	}
	if ev.NewData["key"] != "new_value" {
		t.Errorf("unexpected NewData: %v", ev.NewData)
	}
	if !ev.At.Equal(now) {
		t.Errorf("unexpected timestamp: %v", ev.At)
	}
}

func TestWatchEvent_NoChange(t *testing.T) {
	a := SecretMap{"x": "1", "y": "2"}
	b := SecretMap{"x": "1", "y": "2"}
	if !secretMapsEqual(a, b) {
		t.Fatal("identical maps should be equal")
	}
}

func TestWatchEvent_AddedKey(t *testing.T) {
	a := SecretMap{"x": "1"}
	b := SecretMap{"x": "1", "y": "2"}
	if secretMapsEqual(a, b) {
		t.Fatal("maps with different keys should not be equal")
	}
}

func TestWatchEvent_RemovedKey(t *testing.T) {
	a := SecretMap{"x": "1", "y": "2"}
	b := SecretMap{"x": "1"}
	if secretMapsEqual(a, b) {
		t.Fatal("maps with removed key should not be equal")
	}
}

func TestWatchEvent_NilMaps(t *testing.T) {
	if !secretMapsEqual(nil, nil) {
		t.Fatal("two nil maps should be equal")
	}
	if secretMapsEqual(nil, SecretMap{"k": "v"}) {
		t.Fatal("nil vs non-nil maps should not be equal")
	}
}
