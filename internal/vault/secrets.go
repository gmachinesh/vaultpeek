package vault

import "fmt"

// SecretMap represents a flat map of secret key-value pairs for a path.
type SecretMap struct {
	Env  string
	Path string
	Data map[string]interface{}
}

// FetchSecret retrieves a secret from Vault and wraps it in a SecretMap.
func FetchSecret(client *Client, path string) (*SecretMap, error) {
	data, err := client.ReadSecret(path)
	if err != nil {
		return nil, err
	}
	return &SecretMap{
		Env:  client.Env,
		Path: path,
		Data: data,
	}, nil
}

// Keys returns a sorted list of keys present in the secret.
func (s *SecretMap) Keys() []string {
	keys := make([]string, 0, len(s.Data))
	for k := range s.Data {
		keys = append(keys, k)
	}
	return keys
}

// ValueFor returns the string representation of a secret key's value.
func (s *SecretMap) ValueFor(key string) (string, bool) {
	v, ok := s.Data[key]
	if !ok {
		return "", false
	}
	return fmt.Sprintf("%v", v), true
}
