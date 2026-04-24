package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestRollbackCmd_RegisteredOnRoot(t *testing.T) {
	var found bool
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "rollback <mount> <path> <version>" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("rollback command not registered on root")
	}
}

func TestRollbackCmd_MissingArgs(t *testing.T) {
	cmd := &cobra.Command{Use: "root"}
	cmd.AddCommand(rollbackCmd)
	cmd.SetArgs([]string{"rollback"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error when args are missing")
	}
}

func TestRollbackCmd_InvalidVersion(t *testing.T) {
	cmd := &cobra.Command{Use: "root"}
	cmd.AddCommand(rollbackCmd)
	cmd.SetArgs([]string{"rollback", "secret", "myapp/db", "notanumber"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for non-numeric version")
	}
}

func TestRollbackCmd_ZeroVersion(t *testing.T) {
	cmd := &cobra.Command{Use: "root"}
	cmd.AddCommand(rollbackCmd)
	cmd.SetArgs([]string{"rollback", "secret", "myapp/db", "0"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for version 0")
	}
}

func TestRollbackCmd_DefaultUsage(t *testing.T) {
	if rollbackCmd.Use != "rollback <mount> <path> <version>" {
		t.Errorf("unexpected Use: %s", rollbackCmd.Use)
	}
	if rollbackCmd.Short == "" {
		t.Error("Short description should not be empty")
	}
}
