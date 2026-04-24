package vault_test

import (
	"errors"
	"testing"

	"github.com/hashicorp/vault/api"
	. "github.com/yourorg/vaultpeek/internal/vault"
)

func TestPromoteSecret_EmptySrcMount(t *testing.T) {
	err := PromoteSecret(nil, nil, "", "a/b", "dst", "a/b", false)
	if err == nil || err.Error() != "source mount must not be empty" {
		t.Fatalf("expected source mount error, got %v", err)
	}
}

func TestPromoteSecret_EmptySrcPath(t *testing.T) {
	err := PromoteSecret(nil, nil, "secret", "", "dst", "a/b", false)
	if err == nil || err.Error() != "source path must not be empty" {
		t.Fatalf("expected source path error, got %v", err)
	}
}

func TestPromoteSecret_EmptyDstMount(t *testing.T) {
	err := PromoteSecret(nil, nil, "secret", "a/b", "", "a/b", false)
	if err == nil || err.Error() != "destination mount must not be empty" {
		t.Fatalf("expected destination mount error, got %v", err)
	}
}

func TestPromoteSecret_EmptyDstPath(t *testing.T) {
	err := PromoteSecret(nil, nil, "secret", "a/b", "dst", "", false)
	if err == nil || err.Error() != "destination path must not be empty" {
		t.Fatalf("expected destination path error, got %v", err)
	}
}

func TestPromoteSecret_FetchSourceError(t *testing.T) {
	src := &mockClient{
		logical: &mockLogical{
			readFn: func(path string) (*api.Secret, error) {
				return nil, errors.New("permission denied")
			},
		},
	}
	err := PromoteSecret(src, src, "secret", "dev/app", "secret", "prod/app", false)
	if err == nil {
		t.Fatal("expected error fetching source, got nil")
	}
}

func TestPromoteSecret_DestinationExistsNoOverwrite(t *testing.T) {
	data := map[string]interface{}{"key": "value"}
	client := &mockClient{
		logical: &mockLogical{
			readFn: func(path string) (*api.Secret, error) {
				return makeSecretResponse(data), nil
			},
		},
	}
	err := PromoteSecret(client, client, "secret", "dev/app", "secret", "prod/app", false)
	if err == nil {
		t.Fatal("expected overwrite error, got nil")
	}
}

func TestPromoteSecret_Success(t *testing.T) {
	data := map[string]interface{}{"key": "value"}
	var writeCalled bool
	client := &mockClient{
		logical: &mockLogical{
			readFn: func(path string) (*api.Secret, error) {
				return makeSecretResponse(data), nil
			},
			writeFn: func(path string, d map[string]interface{}) (*api.Secret, error) {
				writeCalled = true
				return nil, nil
			},
		},
	}
	err := PromoteSecret(client, client, "secret", "dev/app", "secret", "prod/app", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !writeCalled {
		t.Fatal("expected write to be called on destination client")
	}
}
