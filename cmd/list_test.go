package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestListCmd_MissingPrefix(t *testing.T) {
	cmd := rootCmd
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"list", "--mount", "secret"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing --prefix flag")
	}
	if !strings.Contains(err.Error(), "prefix") {
		t.Errorf("expected error to mention 'prefix', got: %v", err)
	}
}

func TestListCmd_RegisteredOnRoot(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "list" {
			return
		}
	}
	t.Error("list command not registered on root")
}
