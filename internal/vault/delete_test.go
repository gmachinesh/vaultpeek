package vault

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteSecret_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	err = DeleteSecret(context.Background(), client, "secret", "myapp/config")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestDeleteSecret_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	err = DeleteSecret(context.Background(), client, "secret", "missing/path")
	if err == nil {
		t.Fatal("expected error for not found, got nil")
	}
}

func TestDeleteSecret_EmptyMount(t *testing.T) {
	client := &Client{}
	err := DeleteSecret(context.Background(), client, "", "some/path")
	if err == nil {
		t.Fatal("expected error for empty mount")
	}
}

func TestDeleteSecret_EmptyPath(t *testing.T) {
	client := &Client{}
	err := DeleteSecret(context.Background(), client, "secret", "")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}
