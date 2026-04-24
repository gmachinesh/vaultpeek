package vault_test

import (
	"context"
	"testing"

	"github.com/user/vaultpeek/internal/vault"
)

func TestGetSecretVersion_NilClient(t *testing.T) {
	_, err := vault.GetSecretVersion(context.Background(), nil, "secret", "myapp/db", 1)
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestGetSecretVersion_EmptyMount(t *testing.T) {
	client, _ := vault.NewClient("http://127.0.0.1:8200", "test-token")
	_, err := vault.GetSecretVersion(context.Background(), client, "", "myapp/db", 1)
	if err == nil {
		t.Fatal("expected error for empty mount")
	}
}

func TestGetSecretVersion_EmptyPath(t *testing.T) {
	client, _ := vault.NewClient("http://127.0.0.1:8200", "test-token")
	_, err := vault.GetSecretVersion(context.Background(), client, "secret", "", 1)
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestGetSecretVersion_InvalidVersion(t *testing.T) {
	client, _ := vault.NewClient("http://127.0.0.1:8200", "test-token")
	_, err := vault.GetSecretVersion(context.Background(), client, "secret", "myapp/db", 0)
	if err == nil {
		t.Fatal("expected error for version < 1")
	}
}

func TestRollbackSecret_NilClient(t *testing.T) {
	err := vault.RollbackSecret(context.Background(), nil, "secret", "myapp/db", 2)
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestRollbackSecret_EmptyMount(t *testing.T) {
	client, _ := vault.NewClient("http://127.0.0.1:8200", "test-token")
	err := vault.RollbackSecret(context.Background(), client, "", "myapp/db", 2)
	if err == nil {
		t.Fatal("expected error for empty mount")
	}
}

func TestRollbackSecret_EmptyPath(t *testing.T) {
	client, _ := vault.NewClient("http://127.0.0.1:8200", "test-token")
	err := vault.RollbackSecret(context.Background(), client, "secret", "", 2)
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestRollbackSecret_InvalidVersion(t *testing.T) {
	client, _ := vault.NewClient("http://127.0.0.1:8200", "test-token")
	err := vault.RollbackSecret(context.Background(), client, "secret", "myapp/db", -1)
	if err == nil {
		t.Fatal("expected error for negative version")
	}
}
