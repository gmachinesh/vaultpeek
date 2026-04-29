package vault

import (
	"testing"
)

func makeCompareSecretMap(data map[string]interface{}) SecretMap {
	return SecretMap{data: data}
}

func TestCompareResult_IsDrifted_NoDrift(t *testing.T) {
	r := &CompareResult{Match: []string{"key1", "key2"}}
	if r.IsDrifted() {
		t.Error("expected no drift")
	}
}

func TestCompareResult_IsDrifted_OnlyInA(t *testing.T) {
	r := &CompareResult{OnlyInA: []string{"key1"}}
	if !r.IsDrifted() {
		t.Error("expected drift")
	}
}

func TestCompareResult_IsDrifted_Differ(t *testing.T) {
	r := &CompareResult{Differ: []string{"key1"}}
	if !r.IsDrifted() {
		t.Error("expected drift")
	}
}

func TestCompareResult_Summary_Identical(t *testing.T) {
	r := &CompareResult{
		Mount:  "secret",
		Path:   "app/config",
		SrcEnv: "staging",
		DstEnv: "prod",
		Match:  []string{"key1"},
	}
	got := r.Summary()
	want := "secret/app/config: identical across staging and prod"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestCompareResult_Summary_Drifted(t *testing.T) {
	r := &CompareResult{
		Mount:   "secret",
		Path:    "app/config",
		SrcEnv:  "staging",
		DstEnv:  "prod",
		OnlyInA: []string{"a"},
		Differ:  []string{"b"},
	}
	got := r.Summary()
	if got == "" {
		t.Error("expected non-empty summary")
	}
}

func TestCompareAcrossEnvs_NilSrcClient(t *testing.T) {
	_, err := CompareAcrossEnvs(nil, "staging", &Client{}, "prod", "secret", "app")
	if err == nil {
		t.Error("expected error for nil src client")
	}
}

func TestCompareAcrossEnvs_NilDstClient(t *testing.T) {
	_, err := CompareAcrossEnvs(&Client{}, "staging", nil, "prod", "secret", "app")
	if err == nil {
		t.Error("expected error for nil dst client")
	}
}

func TestCompareAcrossEnvs_EmptyMount(t *testing.T) {
	_, err := CompareAcrossEnvs(&Client{}, "staging", &Client{}, "prod", "", "app")
	if err == nil {
		t.Error("expected error for empty mount")
	}
}

func TestCompareAcrossEnvs_EmptyPath(t *testing.T) {
	_, err := CompareAcrossEnvs(&Client{}, "staging", &Client{}, "prod", "secret", "")
	if err == nil {
		t.Error("expected error for empty path")
	}
}

func TestToStringSet(t *testing.T) {
	keys := []string{"a", "b", "c"}
	s := toStringSet(keys)
	if len(s) != 3 {
		t.Errorf("expected 3 keys, got %d", len(s))
	}
	if _, ok := s["b"]; !ok {
		t.Error("expected key 'b' in set")
	}
}
