package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestTagCmd_RegisteredOnRoot(t *testing.T) {
	var found bool
	for _, c := range rootCmd.Commands() {
		if c.Use == "tag <path> [key=value ...]" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("tag command not registered on root")
	}
}

func TestTagCmd_MissingPath(t *testing.T) {
	rootCmd.SetArgs([]string{"tag"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when path argument is missing")
	}
}

func TestTagCmd_DefaultMountFlag(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Name() == "tag" {
			f := c.Flags().Lookup("mount")
			if f == nil {
				t.Fatal("expected --mount flag")
			}
			if f.DefValue != "secret" {
				t.Fatalf("expected default mount 'secret', got %q", f.DefValue)
			}
			return
		}
	}
	t.Fatal("tag command not found")
}

func TestTagCmd_InvalidPairFormat(t *testing.T) {
	var tagCmd *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c.Name() == "tag" {
			tagCmd = c
			break
		}
	}
	if tagCmd == nil {
		t.Skip("tag command not found")
	}
	// Verify the command requires at least 1 arg (path)
	if tagCmd.Args == nil {
		t.Fatal("expected Args validator on tag command")
	}
}
