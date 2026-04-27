package vault

import (
	"testing"
)

func TestSearchSecrets_NilClient(t *testing.T) {
	_, err := SearchSecrets(nil, "secret", "", "foo", false)
	if err == nil || err.Error() != "client is nil" {
		t.Fatalf("expected 'client is nil', got %v", err)
	}
}

func TestSearchSecrets_EmptyMount(t *testing.T) {
	c := &Client{}
	_, err := SearchSecrets(c, "", "", "foo", false)
	if err == nil || err.Error() != "mount is required" {
		t.Fatalf("expected 'mount is required', got %v", err)
	}
}

func TestSearchSecrets_EmptyQuery(t *testing.T) {
	c := &Client{}
	_, err := SearchSecrets(c, "secret", "", "", false)
	if err == nil || err.Error() != "query is required" {
		t.Fatalf("expected 'query is required', got %v", err)
	}
}

func TestSearchResult_Fields(t *testing.T) {
	r := SearchResult{
		Mount: "secret",
		Path:  "app/config",
		Key:   "db_password",
		Value: "s3cr3t",
	}
	if r.Mount != "secret" {
		t.Errorf("unexpected Mount: %s", r.Mount)
	}
	if r.Path != "app/config" {
		t.Errorf("unexpected Path: %s", r.Path)
	}
	if r.Key != "db_password" {
		t.Errorf("unexpected Key: %s", r.Key)
	}
	if r.Value != "s3cr3t" {
		t.Errorf("unexpected Value: %s", r.Value)
	}
}

func TestSearchSecrets_NoResults(t *testing.T) {
	// Verify empty slice (not nil error) when listing returns nothing.
	// We rely on ListSecrets returning an error for an unconfigured client,
	// so SearchSecrets should propagate that error.
	c := &Client{}
	_, err := SearchSecrets(c, "secret", "nonexistent/", "anything", true)
	if err == nil {
		t.Fatal("expected an error from unconfigured client, got nil")
	}
}
