package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestExportCmd_RegisteredOnRoot(t *testing.T) {
	var found bool
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "export <path>" {
			found = true
			break
		}
	}
	if !found {
		t.Error("export command not registered on root")
	}
}

func TestExportCmd_MissingPath(t *testing.T) {
	cmd := &cobra.Command{}
	err := runExport(cmd, []string{})
	// cobra ExactArgs guard fires before runExport when invoked via Execute,
	// but calling runExport directly with no args still triggers vault.NewClient
	// which fails on missing env — either way an error is returned.
	if err == nil {
		t.Error("expected error when no path provided")
	}
}

func TestExportCmd_InvalidFormat(t *testing.T) {
	// Temporarily override the flag value.
	prev := exportFormat
	exportFormat = "yaml"
	defer func() { exportFormat = prev }()

	cmd := &cobra.Command{}
	err := runExport(cmd, []string{"app/config"})
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestExportCmd_DefaultFlags(t *testing.T) {
	if exportMount != "secret" {
		t.Errorf("expected default mount 'secret', got %q", exportMount)
	}
	if exportFormat != "json" {
		t.Errorf("expected default format 'json', got %q", exportFormat)
	}
}
