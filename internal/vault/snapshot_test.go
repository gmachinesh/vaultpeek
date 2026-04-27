package vault

import (
	"testing"
)

func TestTakeSnapshot_NilClient(t *testing.T) {
	_, err := TakeSnapshot(nil, "secret", "myapp")
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestTakeSnapshot_EmptyMount(t *testing.T) {
	client := &Client{}
	_, err := TakeSnapshot(client, "", "myapp")
	if err == nil {
		t.Fatal("expected error for empty mount")
	}
}

func TestTakeSnapshot_EmptyPrefix(t *testing.T) {
	client := &Client{}
	_, err := TakeSnapshot(client, "secret", "")
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestSnapshot_Fields(t *testing.T) {
	sm := SecretMap{"data": map[string]interface{}{"key": "val"}}
	snap := &Snapshot{
		Mount:   "secret",
		Prefix:  "myapp",
		Secrets: map[string]SecretMap{"db": sm},
	}

	if snap.Mount != "secret" {
		t.Errorf("expected mount 'secret', got %q", snap.Mount)
	}
	if snap.Prefix != "myapp" {
		t.Errorf("expected prefix 'myapp', got %q", snap.Prefix)
	}
	if len(snap.Secrets) != 1 {
		t.Errorf("expected 1 secret, got %d", len(snap.Secrets))
	}
}

func TestSnapshotResult_Fields(t *testing.T) {
	result := &SnapshotResult{
		Snapshot: &Snapshot{
			Mount:   "kv",
			Prefix:  "svc",
			Secrets: map[string]SecretMap{},
		},
		Count: 0,
	}

	if result.Count != 0 {
		t.Errorf("expected count 0, got %d", result.Count)
	}
	if result.Snapshot == nil {
		t.Fatal("expected non-nil snapshot")
	}
}
