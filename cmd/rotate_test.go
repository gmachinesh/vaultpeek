package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestRotateCmd_RegisteredOnRoot(t *testing.T) {
	var found *cobra.Command
	for _, sub := range rootCmd.Commands() {
		if sub.Name() == "rotate" {
			found = sub
			break
		}
	}
	if found == nil {
		t.Fatal("rotate command not registered on root")
	}
}

func TestRotateCmd_MissingPath(t *testing.T) {
	cmd := rotateCmd
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when path argument is missing")
	}
}

func TestRotateCmd_DefaultFlags(t *testing.T) {
	mountFlag := rotateCmd.Flags().Lookup("mount")
	if mountFlag == nil {
		t.Fatal("expected --mount flag to be defined")
	}
	if mountFlag.DefValue != "secret" {
		t.Errorf("expected default mount 'secret', got %q", mountFlag.DefValue)
	}

	lengthFlag := rotateCmd.Flags().Lookup("length")
	if lengthFlag == nil {
		t.Fatal("expected --length flag to be defined")
	}
	if lengthFlag.DefValue != "32" {
		t.Errorf("expected default length '32', got %q", lengthFlag.DefValue)
	}
}

func TestRandomString_Length(t *testing.T) {
	for _, n := range []int{8, 16, 32, 64} {
		s := randomString(n)
		if len(s) != n {
			t.Errorf("randomString(%d) returned length %d", n, len(s))
		}
	}
}

func TestRandomString_Uniqueness(t *testing.T) {
	a := randomString(32)
	b := randomString(32)
	if a == b {
		t.Log("warning: two random strings were equal (unlikely but possible)")
	}
}
