package vault

import (
	"testing"
)

func TestKvV2DataPath(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"myapp/config", "secret/data/myapp/config"},
		{"prod/db", "secret/data/prod/db"},
		{"", "secret/data/"},
	}
	for _, tt := range tests {
		got := kvV2DataPath(tt.input)
		if got != tt.want {
			t.Errorf("kvV2DataPath(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestNewClient_MissingAddress(t *testing.T) {
	t.Setenv("VAULT_ADDR", "")
	t.Setenv("VAULT_TOKEN", "")
	_, err := NewClient(Config{})
	if err == nil {
		t.Fatal("expected error for missing address, got nil")
	}
}

func TestNewClient_MissingToken(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://127.0.0.1:8200")
	t.Setenv("VAULT_TOKEN", "")
	_, err := NewClient(Config{Address: "http://127.0.0.1:8200"})
	if err == nil {
		t.Fatal("expected error for missing token, got nil")
	}
}

func TestNewClient_Success(t *testing.T) {
	client, err := NewClient(Config{
		Address: "http://127.0.0.1:8200",
		Token:   "root",
		Env:     "dev",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.Env != "dev" {
		t.Errorf("expected env %q, got %q", "dev", client.Env)
	}
}
