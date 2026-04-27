package vault_test

import (
	"testing"

	"github.com/your-org/vaultpeek/internal/vault"
)

func TestTouchSecret_NilClient(t *testing.T) {
	_, err := vault.TouchSecret(nil, "secret", "myapp/config")
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestTouchSecret_EmptyMount(t *testing.T) {
	c, _ := vault.NewClient("http://127.0.0.1:8200", "token")
	_, err := vault.TouchSecret(c, "", "myapp/config")
	if err == nil {
		t.Fatal("expected error for empty mount")
	}
}

func TestTouchSecret_EmptyPath(t *testing.T) {
	c, _ := vault.NewClient("http://127.0.0.1:8200", "token")
	_, err := vault.TouchSecret(c, "secret", "")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestTouchResult_Fields(t *testing.T) {
	res := &vault.TouchResult{
		Mount: "secret",
		Path:  "app/db",
		Version: 3,
	}
	if res.Mount != "secret" {
		t.Errorf("expected mount 'secret', got %q", res.Mount)
	}
	if res.Path != "app/db" {
		t.Errorf("expected path 'app/db', got %q", res.Path)
	}
	if res.Version != 3 {
		t.Errorf("expected version 3, got %d", res.Version)
	}
	if res.TouchedAt.IsZero() {
		t.Error("expected non-zero TouchedAt")
	}
}

func TestTouchSecret_FetchError(t *testing.T) {
	// Using a non-existent server forces a fetch error.
	c, _ := vault.NewClient("http://127.0.0.1:19999", "fake-token")
	_, err := vault.TouchSecret(c, "secret", "nonexistent/path")
	if err == nil {
		t.Fatal("expected error when vault is unreachable")
	}
}
