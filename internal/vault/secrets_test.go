package vault

import (
	"testing"
)

func makeSecretMap(env, path string, data map[string]interface{}) *SecretMap {
	return &SecretMap{Env: env, Path: path, Data: data}
}

func TestSecretMap_ValueFor(t *testing.T) {
	sm := makeSecretMap("dev", "myapp/config", map[string]interface{}{
		"DB_HOST": "localhost",
		"PORT":    8080,
	})

	v, ok := sm.ValueFor("DB_HOST")
	if !ok || v != "localhost" {
		t.Errorf("expected localhost, got %q (ok=%v)", v, ok)
	}

	v, ok = sm.ValueFor("PORT")
	if !ok || v != "8080" {
		t.Errorf("expected 8080, got %q (ok=%v)", v, ok)
	}

	_, ok = sm.ValueFor("MISSING")
	if ok {
		t.Error("expected ok=false for missing key")
	}
}

func TestSecretMap_Keys(t *testing.T) {
	sm := makeSecretMap("prod", "myapp/config", map[string]interface{}{
		"A": "1",
		"B": "2",
		"C": "3",
	})
	keys := sm.Keys()
	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(keys))
	}
}

func TestSecretMap_EmptyData(t *testing.T) {
	sm := makeSecretMap("staging", "empty/path", map[string]interface{}{})
	if len(sm.Keys()) != 0 {
		t.Error("expected no keys for empty secret")
	}
}
