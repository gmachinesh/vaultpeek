package vault

import (
	"testing"

	"github.com/hashicorp/vault/api"
)

func TestPatchSecret_EmptyMount(t *testing.T) {
	c := &Client{}
	if err := PatchSecret(c, "", "mypath", map[string]string{"k": "v"}); err == nil {
		t.Fatal("expected error for empty mount")
	}
}

func TestPatchSecret_EmptyPath(t *testing.T) {
	c := &Client{}
	if err := PatchSecret(c, "secret", "", map[string]string{"k": "v"}); err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestPatchSecret_EmptyUpdates(t *testing.T) {
	c := &Client{}
	if err := PatchSecret(c, "secret", "mypath", map[string]string{}); err == nil {
		t.Fatal("expected error for empty updates")
	}
}

func TestPatchSecret_FetchError(t *testing.T) {
	server, client := newMockVaultServer(t, func(path string) (*api.Secret, error) {
		return nil, fmt.Errorf("vault unavailable")
	})
	defer server.Close()

	err := PatchSecret(client, "secret", "app/config", map[string]string{"key": "val"})
	if err == nil {
		t.Fatal("expected error when fetch fails")
	}
}

func TestPatchSecret_MergesAndDeletes(t *testing.T) {
	existing := map[string]interface{}{
		"keep":   "original",
		"update": "old",
		"remove": "bye",
	}

	var written map[string]interface{}

	server, client := newMockVaultServerWithWrite(t,
		func(path string) (*api.Secret, error) {
			return &api.Secret{Data: map[string]interface{}{"data": existing}}, nil
		},
		func(path string, data map[string]interface{}) (*api.Secret, error) {
			written = data["data"].(map[string]interface{})
			return &api.Secret{}, nil
		},
	)
	defer server.Close()

	err := PatchSecret(client, "secret", "app/config", map[string]string{
		"update": "new",
		"remove": "",
		"added":  "fresh",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if written["keep"] != "original" {
		t.Errorf("expected keep=original, got %v", written["keep"])
	}
	if written["update"] != "new" {
		t.Errorf("expected update=new, got %v", written["update"])
	}
	if _, ok := written["remove"]; ok {
		t.Error("expected 'remove' key to be deleted")
	}
	if written["added"] != "fresh" {
		t.Errorf("expected added=fresh, got %v", written["added"])
	}
}
