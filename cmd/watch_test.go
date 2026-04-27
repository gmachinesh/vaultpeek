package cmd

import (
	"strings"
	"testing"
	"time"
)

func TestWatchCmd_RegisteredOnRoot(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Name() == "watch" {
			return
		}
	}
	t.Fatal("watch command not registered on root")
}

func TestWatchCmd_MissingArgs(t *testing.T) {
	rootCmd.SetArgs([]string{"watch"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when no args provided")
	}
}

func TestWatchCmd_DefaultFlags(t *testing.T) {
	cmd := watchCmd

	intervalFlag := cmd.Flags().Lookup("interval")
	if intervalFlag == nil {
		t.Fatal("expected --interval flag")
	}
	if intervalFlag.DefValue != "10s" {
		t.Fatalf("expected default interval 10s, got %s", intervalFlag.DefValue)
	}

	if cmd.Flags().Lookup("token") == nil {
		t.Fatal("expected --token flag")
	}
	if cmd.Flags().Lookup("address") == nil {
		t.Fatal("expected --address flag")
	}
}

func TestMaskValue_Short(t *testing.T) {
	result := maskValue("ab")
	if result != "**" {
		t.Fatalf("expected '**', got %q", result)
	}
}

func TestMaskValue_Long(t *testing.T) {
	result := maskValue("supersecret")
	if !strings.HasPrefix(result, "su") {
		t.Fatalf("expected result to start with 'su', got %q", result)
	}
	if !strings.HasSuffix(result, "et") {
		t.Fatalf("expected result to end with 'et', got %q", result)
	}
	if strings.Contains(result, "persecr") {
		t.Fatalf("expected middle to be masked, got %q", result)
	}
}

func TestWatchCmd_IntervalFlag(t *testing.T) {
	f := watchCmd.Flags().Lookup("interval")
	if f == nil {
		t.Fatal("interval flag not found")
	}
	var d time.Duration
	if err := watchCmd.Flags().Set("interval", "30s"); err != nil {
		t.Fatalf("failed to set interval: %v", err)
	}
	d, _ = watchCmd.Flags().GetDuration("interval")
	if d != 30*time.Second {
		t.Fatalf("expected 30s, got %v", d)
	}
	// reset
	_ = watchCmd.Flags().Set("interval", "10s")
}
