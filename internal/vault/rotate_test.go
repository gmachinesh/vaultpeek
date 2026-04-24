package vault

import (
	"errors"
	"testing"
)

func TestRotateSecret_NilClient(t *testing.T) {
	_, err := RotateSecret(nil, "secret", "myapp/config", func(k string) string { return "v" })
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestRotateSecret_EmptyMount(t *testing.T) {
	client := &Client{}
	_, err := RotateSecret(client, "", "myapp/config", func(k string) string { return "v" })
	if err == nil {
		t.Fatal("expected error for empty mount")
	}
}

func TestRotateSecret_EmptyPath(t *testing.T) {
	client := &Client{}
	_, err := RotateSecret(client, "secret", "", func(k string) string { return "v" })
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestRotateSecret_NilGenerator(t *testing.T) {
	client := &Client{}
	_, err := RotateSecret(client, "secret", "myapp/config", nil)
	if err == nil {
		t.Fatal("expected error for nil generator")
	}
}

func TestRotateSecret_FetchError(t *testing.T) {
	server, client := newMockVaultServer(t, map[string]mockResponse{
		"/v1/secret/data/myapp/config": {
			statusCode: 403,
			body:       `{"errors":["permission denied"]}`,
		},
	})
	defer server.Close()

	_, err := RotateSecret(client, "secret", "myapp/config", func(k string) string { return "new" })
	if err == nil {
		t.Fatal("expected error on fetch failure")
	}
	if !errors.Is(err, err) {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRotateSecret_Success(t *testing.T) {
	server, client := newMockVaultServer(t, map[string]mockResponse{
		"/v1/secret/data/myapp/config": {
			statusCode: 200,
			body: `{"data":{"data":{"DB_PASS":"old","API_KEY":"old2"},"metadata":{}}}`,
		},
	})
	defer server.Close()

	result, err := RotateSecret(client, "secret", "myapp/config", func(k string) string {
		return "rotated-" + k
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Mount != "secret" {
		t.Errorf("expected mount 'secret', got %q", result.Mount)
	}
	if result.Path != "myapp/config" {
		t.Errorf("expected path 'myapp/config', got %q", result.Path)
	}
	if len(result.OldKeys) != 2 {
		t.Errorf("expected 2 old keys, got %d", len(result.OldKeys))
	}
}
