package diff

import (
	"testing"

	"github.com/user/vaultpeek/internal/vault"
)

func makeMap(data map[string]interface{}) vault.SecretMap {
	return vault.SecretMap(data)
}

func TestCompare_OnlyInA(t *testing.T) {
	a := makeMap(map[string]interface{}{"foo": "bar", "baz": "qux"})
	b := makeMap(map[string]interface{}{"foo": "bar"})
	r := Compare(a, b)
	if len(r.OnlyInA) != 1 || r.OnlyInA[0] != "baz" {
		t.Errorf("expected OnlyInA=[baz], got %v", r.OnlyInA)
	}
}

func TestCompare_OnlyInB(t *testing.T) {
	a := makeMap(map[string]interface{}{"foo": "bar"})
	b := makeMap(map[string]interface{}{"foo": "bar", "new": "val"})
	r := Compare(a, b)
	if len(r.OnlyInB) != 1 || r.OnlyInB[0] != "new" {
		t.Errorf("expected OnlyInB=[new], got %v", r.OnlyInB)
	}
}

func TestCompare_Differing(t *testing.T) {
	a := makeMap(map[string]interface{}{"key": "val1"})
	b := makeMap(map[string]interface{}{"key": "val2"})
	r := Compare(a, b)
	if len(r.Differing) != 1 || r.Differing[0].Key != "key" {
		t.Errorf("expected one differing key, got %v", r.Differing)
	}
	if r.Differing[0].ValueA != "val1" || r.Differing[0].ValueB != "val2" {
		t.Errorf("unexpected values: %+v", r.Differing[0])
	}
}

func TestCompare_Matching(t *testing.T) {
	a := makeMap(map[string]interface{}{"x": "1", "y": "2"})
	b := makeMap(map[string]interface{}{"x": "1", "y": "2"})
	r := Compare(a, b)
	if len(r.Matching) != 2 {
		t.Errorf("expected 2 matching keys, got %v", r.Matching)
	}
	if len(r.Differing) != 0 || len(r.OnlyInA) != 0 || len(r.OnlyInB) != 0 {
		t.Errorf("unexpected diff results: %+v", r)
	}
}
