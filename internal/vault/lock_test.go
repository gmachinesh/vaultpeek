package vault

import (
	"testing"
)

func TestLockSecret_NilClient(t *testing.T) {
	err := LockSecret(nil, "secret", "app/config", "")
	if err == nil || err.Error() != "client is nil" {
		t.Fatalf("expected 'client is nil', got %v", err)
	}
}

func TestLockSecret_EmptyMount(t *testing.T) {
	c := &Client{}
	err := LockSecret(c, "", "app/config", "")
	if err == nil || err.Error() != "mount is required" {
		t.Fatalf("expected 'mount is required', got %v", err)
	}
}

func TestLockSecret_EmptyPath(t *testing.T) {
	c := &Client{}
	err := LockSecret(c, "secret", "", "")
	if err == nil || err.Error() != "path is required" {
		t.Fatalf("expected 'path is required', got %v", err)
	}
}

func TestUnlockSecret_NilClient(t *testing.T) {
	err := UnlockSecret(nil, "secret", "app/config")
	if err == nil || err.Error() != "client is nil" {
		t.Fatalf("expected 'client is nil', got %v", err)
	}
}

func TestUnlockSecret_EmptyMount(t *testing.T) {
	c := &Client{}
	err := UnlockSecret(c, "", "app/config")
	if err == nil || err.Error() != "mount is required" {
		t.Fatalf("expected 'mount is required', got %v", err)
	}
}

func TestUnlockSecret_EmptyPath(t *testing.T) {
	c := &Client{}
	err := UnlockSecret(c, "secret", "")
	if err == nil || err.Error() != "path is required" {
		t.Fatalf("expected 'path is required', got %v", err)
	}
}

func TestIsLocked_NilClient(t *testing.T) {
	_, err := IsLocked(nil, "secret", "app/config")
	if err == nil || err.Error() != "client is nil" {
		t.Fatalf("expected 'client is nil', got %v", err)
	}
}

func TestIsLocked_EmptyMount(t *testing.T) {
	c := &Client{}
	_, err := IsLocked(c, "", "app/config")
	if err == nil || err.Error() != "mount is required" {
		t.Fatalf("expected 'mount is required', got %v", err)
	}
}

func TestIsLocked_EmptyPath(t *testing.T) {
	c := &Client{}
	_, err := IsLocked(c, "secret", "")
	if err == nil || err.Error() != "path is required" {
		t.Fatalf("expected 'path is required', got %v", err)
	}
}
