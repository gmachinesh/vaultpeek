package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCopySecret_Success(t *testing.T) {
	var written map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"data": map[string]interface{}{"API_KEY": "abc123"},
				},
			})
		case http.MethodPost:
			if err := json.NewDecoder(r.Body).Decode(&written); err != nil {
				http.Error(w, "bad body", http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		}
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	if err := CopySecret(context.Background(), client, "secret/src", "secret/dst"); err != nil {
		t.Fatalf("CopySecret: %v", err)
	}

	data, ok := written["data"].(map[string]interface{})
	if !ok {
		t.Fatal("expected written data map")
	}
	if data["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %v", data["API_KEY"])
	}
}

func TestCopySecret_FetchError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer server.Close()

	client, _ := NewClient(server.URL, "test-token")
	err := CopySecret(context.Background(), client, "secret/missing", "secret/dst")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
