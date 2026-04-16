package diff

import (
	"fmt"
	"sort"

	"github.com/user/vaultpeek/internal/vault"
)

// Result holds the diff outcome between two secret maps.
type Result struct {
	OnlyInA  []string
	OnlyInB  []string
	Differing []KeyDiff
	Matching  []string
}

// KeyDiff represents a key whose value differs between two environments.
type KeyDiff struct {
	Key    string
	ValueA string
	ValueB string
}

// Compare returns a Result describing differences between two SecretMaps.
func Compare(a, b vault.SecretMap) Result {
	result := Result{}

	keysA := a.Keys()
	keysB := b.Keys()

	setA := toSet(keysA)
	setB := toSet(keysB)

	for _, k := range keysA {
		if !setB[k] {
			result.OnlyInA = append(result.OnlyInA, k)
			continue
		}
		va, _ := a.ValueFor(k)
		vb, _ := b.ValueFor(k)
		if va != vb {
			result.Differing = append(result.Differing, KeyDiff{Key: k, ValueA: fmt.Sprintf("%v", va), ValueB: fmt.Sprintf("%v", vb)})
		} else {
			result.Matching = append(result.Matching, k)
		}
	}

	for _, k := range keysB {
		if !setA[k] {
			result.OnlyInB = append(result.OnlyInB, k)
		}
	}

	sort.Strings(result.OnlyInA)
	sort.Strings(result.OnlyInB)
	sort.Strings(result.Matching)

	return result
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
