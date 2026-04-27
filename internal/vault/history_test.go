package vault

import (
	"testing"
	"time"
)

func TestGetSecretHistory_NilClient(t *testing.T) {
	_, err := GetSecretHistory(nil, "secret", "myapp/config")
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestGetSecretHistory_EmptyMount(t *testing.T) {
	c := &Client{}
	_, err := GetSecretHistory(c, "", "myapp/config")
	if err == nil {
		t.Fatal("expected error for empty mount")
	}
}

func TestGetSecretHistory_EmptyPath(t *testing.T) {
	c := &Client{}
	_, err := GetSecretHistory(c, "secret", "")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestParseVersionMeta_FullData(t *testing.T) {
	m := map[string]interface{}{
		"created_time":  "2024-01-15T10:00:00Z",
		"deletion_time": "2024-02-01T00:00:00Z",
		"destroyed":     false,
		"deleted":       true,
	}
	meta := parseVersionMeta("3", m)

	if meta.Version != 3 {
		t.Errorf("expected version 3, got %d", meta.Version)
	}
	if meta.CreatedTime.IsZero() {
		t.Error("expected non-zero created_time")
	}
	if meta.DeletionTime.IsZero() {
		t.Error("expected non-zero deletion_time")
	}
	if meta.Destroyed {
		t.Error("expected destroyed=false")
	}
	if !meta.Deleted {
		t.Error("expected deleted=true")
	}
}

func TestParseVersionMeta_EmptyDeletionTime(t *testing.T) {
	m := map[string]interface{}{
		"created_time":  "2024-03-10T08:30:00Z",
		"deletion_time": "",
		"destroyed":     false,
		"deleted":       false,
	}
	meta := parseVersionMeta("1", m)

	if !meta.DeletionTime.IsZero() {
		t.Errorf("expected zero DeletionTime for empty string, got %v", meta.DeletionTime)
	}
}

func TestParseVersionMeta_InvalidVersion(t *testing.T) {
	m := map[string]interface{}{}
	meta := parseVersionMeta("notanumber", m)
	if meta.Version != 0 {
		t.Errorf("expected version 0 for invalid string, got %d", meta.Version)
	}
}

func TestVersionMeta_Fields(t *testing.T) {
	now := time.Now()
	v := VersionMeta{
		Version:     2,
		CreatedTime: now,
		Destroyed:   false,
		Deleted:     false,
	}
	if v.Version != 2 {
		t.Errorf("unexpected version: %d", v.Version)
	}
	if v.CreatedTime != now {
		t.Error("unexpected created time")
	}
}
