package vault_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultpeek/internal/vault"
)

func TestRenameSecret_EmptyMount(t *testing.T) {
	c := &vault.Client{}
	if err := vault.RenameSecret(c, "", "src", "dst"); err == nil {
		t.Fatal("expected error for empty mount")
	}
}

func TestRenameSecret_EmptySrc(t *testing.T) {
	c := &vault.Client{}
	if err := vault.RenameSecret(c, "secret", "", "dst"); err == nil {
		t.Fatal("expected error for empty source path")
	}
}

func TestRenameSecret_EmptyDst(t *testing.T) {
	c := &vault.Client{}
	if err := vault.RenameSecret(c, "secret", "src", ""); err == nil {
		t.Fatal("expected error for empty destination path")
	}
}

func TestRenameSecret_SamePaths(t *testing.T) {
	c := &vault.Client{}
	if err := vault.RenameSecret(c, "secret", "same", "same"); err == nil {
		t.Fatal("expected error when src == dst")
	}
}

func TestRenameSecret_Success(t *testing.T) {
	handled := map[string]int{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handled[r.Method+" "+r.URL.Path]++
		switch {
		case r.Method == http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data":{"data":{"key":"val"},"metadata":{}}}`)) 
		case r.Method == http.MethodPost || r.Method == http.MethodPut:
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{}`))
		case r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		}
	}))
	defer ts.Close()

	c, err := vault.NewClient(ts.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := vault.RenameSecret(c, "secret", "old/path", "new/path"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
