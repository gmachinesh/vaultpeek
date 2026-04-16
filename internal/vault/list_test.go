package vault

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListSecrets_Empty(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":[]}`)) //nolint
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	keys, err := ListSecrets(context.Background(), client, "secret", "myapp")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(keys) != 0 {
		t.Errorf("expected empty keys, got %v", keys)
	}
}

func TestListSecrets_Keys(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"keys":["alpha","beta","gamma"]}}`)) //nolint
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	keys, err := ListSecrets(context.Background(), client, "secret", "myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	expected := map[string]bool{"alpha": true, "beta": true, "gamma": true}
	for _, k := range keys {
		if !expected[k] {
			t.Errorf("unexpected key: %s", k)
		}
	}
}
