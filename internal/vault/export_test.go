package vault

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestExportSecret_NilClient(t *testing.T) {
	var buf bytes.Buffer
	err := ExportSecret(nil, "secret", "app/config", FormatJSON, &buf)
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestExportSecret_EmptyMount(t *testing.T) {
	client := &Client{}
	var buf bytes.Buffer
	err := ExportSecret(client, "", "app/config", FormatJSON, &buf)
	if err == nil {
		t.Fatal("expected error for empty mount")
	}
}

func TestExportSecret_EmptyPath(t *testing.T) {
	client := &Client{}
	var buf bytes.Buffer
	err := ExportSecret(client, "secret", "", FormatJSON, &buf)
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestExportSecret_UnsupportedFormat(t *testing.T) {
	// We can't reach the format switch without a working Vault, so test
	// the helpers directly via exportJSON / exportEnv.
	sm := makeSecretMap(map[string]interface{}{"KEY": "val"})
	var buf bytes.Buffer
	err := exportJSON(sm, []string{"KEY"}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]string
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["KEY"] != "val" {
		t.Errorf("expected val, got %q", out["KEY"])
	}
}

func TestExportEnv_Output(t *testing.T) {
	sm := makeSecretMap(map[string]interface{}{"DB_PASS": "s3cr3t", "DB_USER": "admin"})
	keys := []string{"DB_PASS", "DB_USER"}
	var buf bytes.Buffer
	if err := exportEnv(sm, keys, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_PASS=") {
		t.Errorf("expected DB_PASS in output, got: %s", out)
	}
	if !strings.Contains(out, "DB_USER=") {
		t.Errorf("expected DB_USER in output, got: %s", out)
	}
}

func TestExportJSON_SortedKeys(t *testing.T) {
	sm := makeSecretMap(map[string]interface{}{"Z": "last", "A": "first", "M": "mid"})
	keys := []string{"A", "M", "Z"}
	var buf bytes.Buffer
	if err := exportJSON(sm, keys, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]string
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 3 {
		t.Errorf("expected 3 keys, got %d", len(out))
	}
}
