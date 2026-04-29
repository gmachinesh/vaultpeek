package vault

import "fmt"

// CompareResult holds the result of comparing a secret across two environments.
type CompareResult struct {
	Mount   string
	Path    string
	SrcEnv  string
	DstEnv  string
	OnlyInA []string
	OnlyInB []string
	Differ  []string
	Match   []string
}

// IsDrifted returns true if there are any differences between the two environments.
func (r *CompareResult) IsDrifted() bool {
	return len(r.OnlyInA) > 0 || len(r.OnlyInB) > 0 || len(r.Differ) > 0
}

// Summary returns a human-readable one-line summary of the comparison.
func (r *CompareResult) Summary() string {
	if !r.IsDrifted() {
		return fmt.Sprintf("%s/%s: identical across %s and %s", r.Mount, r.Path, r.SrcEnv, r.DstEnv)
	}
	return fmt.Sprintf("%s/%s: drift detected (%d only-in-%s, %d only-in-%s, %d differing)",
		r.Mount, r.Path,
		len(r.OnlyInA), r.SrcEnv,
		len(r.OnlyInB), r.DstEnv,
		len(r.Differ),
	)
}

// CompareAcrossEnvs fetches the same secret path from two Vault clients and
// returns a CompareResult describing the differences.
func CompareAcrossEnvs(
	srcClient *Client, srcEnv string,
	dstClient *Client, dstEnv string,
	mount, path string,
) (*CompareResult, error) {
	if srcClient == nil || dstClient == nil {
		return nil, fmt.Errorf("compare: both src and dst clients must be non-nil")
	}
	if mount == "" {
		return nil, fmt.Errorf("compare: mount must not be empty")
	}
	if path == "" {
		return nil, fmt.Errorf("compare: path must not be empty")
	}

	srcSecret, err := FetchSecret(srcClient, mount, path)
	if err != nil {
		return nil, fmt.Errorf("compare: fetch from %s: %w", srcEnv, err)
	}
	dstSecret, err := FetchSecret(dstClient, mount, path)
	if err != nil {
		return nil, fmt.Errorf("compare: fetch from %s: %w", dstEnv, err)
	}

	srcKeys := toStringSet(srcSecret.Keys())
	dstKeys := toStringSet(dstSecret.Keys())

	result := &CompareResult{
		Mount:  mount,
		Path:   path,
		SrcEnv: srcEnv,
		DstEnv: dstEnv,
	}

	for k := range srcKeys {
		if _, ok := dstKeys[k]; !ok {
			result.OnlyInA = append(result.OnlyInA, k)
		} else if srcSecret.ValueFor(k) != dstSecret.ValueFor(k) {
			result.Differ = append(result.Differ, k)
		} else {
			result.Match = append(result.Match, k)
		}
	}
	for k := range dstKeys {
		if _, ok := srcKeys[k]; !ok {
			result.OnlyInB = append(result.OnlyInB, k)
		}
	}
	return result, nil
}

func toStringSet(keys []string) map[string]struct{} {
	s := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		s[k] = struct{}{}
	}
	return s
}
