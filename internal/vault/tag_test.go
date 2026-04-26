package vault

import (
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

func TestTagSecret_NilClient(t *testing.T) {
	err := TagSecret(nil, "secret", "myapp/config", map[string]string{"env": "prod"})
	if err == nil || err.Error() != "client is nil" {
		t.Fatalf("expected 'client is nil', got %v", err)
	}
}

func TestTagSecret_EmptyMount(t *testing.T) {
	c := &Client{client: &vaultapi.Client{}}
	err := TagSecret(c, "", "myapp/config", map[string]string{"env": "prod"})
	if err == nil || err.Error() != "mount is required" {
		t.Fatalf("expected 'mount is required', got %v", err)
	}
}

func TestTagSecret_EmptyPath(t *testing.T) {
	c := &Client{client: &vaultapi.Client{}}
	err := TagSecret(c, "secret", "", map[string]string{"env": "prod"})
	if err == nil || err.Error() != "path is required" {
		t.Fatalf("expected 'path is required', got %v", err)
	}
}

func TestTagSecret_EmptyTags(t *testing.T) {
	c := &Client{client: &vaultapi.Client{}}
	err := TagSecret(c, "secret", "myapp/config", map[string]string{})
	if err == nil || err.Error() != "tags must not be empty" {
		t.Fatalf("expected 'tags must not be empty', got %v", err)
	}
}

func TestGetTags_NilClient(t *testing.T) {
	_, err := GetTags(nil, "secret", "myapp/config")
	if err == nil || err.Error() != "client is nil" {
		t.Fatalf("expected 'client is nil', got %v", err)
	}
}

func TestGetTags_EmptyMount(t *testing.T) {
	c := &Client{client: &vaultapi.Client{}}
	_, err := GetTags(c, "", "myapp/config")
	if err == nil || err.Error() != "mount is required" {
		t.Fatalf("expected 'mount is required', got %v", err)
	}
}

func TestGetTags_EmptyPath(t *testing.T) {
	c := &Client{client: &vaultapi.Client{}}
	_, err := GetTags(c, "secret", "")
	if err == nil || err.Error() != "path is required" {
		t.Fatalf("expected 'path is required', got %v", err)
	}
}
