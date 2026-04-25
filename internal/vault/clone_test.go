package vault

import (
	"testing"
)

func TestCloneSecret_NilClient(t *testing.T) {
	err := CloneSecret(nil, "secret", "src", "secret", "dst", CloneOptions{})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestCloneSecret_EmptySrcMount(t *testing.T) {
	err := CloneSecret(&Client{}, "", "src", "secret", "dst", CloneOptions{})
	if err == nil || err.Error() != "vault: source mount must not be empty" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCloneSecret_EmptySrcPath(t *testing.T) {
	err := CloneSecret(&Client{}, "secret", "", "secret", "dst", CloneOptions{})
	if err == nil || err.Error() != "vault: source path must not be empty" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCloneSecret_EmptyDstMount(t *testing.T) {
	err := CloneSecret(&Client{}, "secret", "src", "", "dst", CloneOptions{})
	if err == nil || err.Error() != "vault: destination mount must not be empty" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCloneSecret_EmptyDstPath(t *testing.T) {
	err := CloneSecret(&Client{}, "secret", "src", "secret", "", CloneOptions{})
	if err == nil || err.Error() != "vault: destination path must not be empty" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCloneSecret_SamePaths(t *testing.T) {
	err := CloneSecret(&Client{}, "secret", "app/config", "secret", "app/config", CloneOptions{})
	if err == nil {
		t.Fatal("expected error when source and destination are identical")
	}
}

func TestCloneOptions_DefaultOverwriteFalse(t *testing.T) {
	var opts CloneOptions
	if opts.Overwrite {
		t.Fatal("expected Overwrite to default to false")
	}
}
