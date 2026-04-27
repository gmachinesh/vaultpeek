package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"vaultpeek/internal/vault"
)

var watchCmd = &cobra.Command{
	Use:   "watch <mount> <path>",
	Short: "Watch a secret for changes and print diffs",
	Args:  cobra.ExactArgs(2),
	RunE:  runWatch,
}

func init() {
	watchCmd.Flags().DurationP("interval", "i", 10*time.Second, "polling interval (e.g. 5s, 1m)")
	watchCmd.Flags().StringP("token", "t", "", "Vault token (overrides VAULT_TOKEN)")
	watchCmd.Flags().StringP("address", "a", "", "Vault address (overrides VAULT_ADDR)")
	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, args []string) error {
	mount := args[0]
	path := args[1]

	interval, _ := cmd.Flags().GetDuration("interval")
	addr, _ := cmd.Flags().GetString("address")
	token, _ := cmd.Flags().GetString("token")

	client, err := vault.NewClient(addr, token)
	if err != nil {
		return fmt.Errorf("watch: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ch, err := vault.WatchSecret(ctx, client, mount, path, vault.WatchOptions{Interval: interval})
	if err != nil {
		return fmt.Errorf("watch: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Watching %s/%s every %s (Ctrl+C to stop)\n", mount, path, interval)

	for event := range ch {
		fmt.Fprintf(cmd.OutOrStdout(), "\n[%s] Change detected in %s/%s\n",
			event.At.Format(time.RFC3339), event.Mount, event.Path)
		printWatchDiff(cmd, event)
	}

	return nil
}

func printWatchDiff(cmd *cobra.Command, e vault.WatchEvent) {
	allKeys := map[string]struct{}{}
	for k := range e.OldData {
		allKeys[k] = struct{}{}
	}
	for k := range e.NewData {
		allKeys[k] = struct{}{}
	}
	for k := range allKeys {
		oldVal := e.OldData[k]
		newVal := e.NewData[k]
		switch {
		case oldVal == "" && newVal != "":
			fmt.Fprintf(cmd.OutOrStdout(), "  + %s = %s\n", k, newVal)
		case oldVal != "" && newVal == "":
			fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", k)
		case oldVal != newVal:
			fmt.Fprintf(cmd.OutOrStdout(), "  ~ %s: %s → %s\n", k, maskValue(oldVal), maskValue(newVal))
		}
	}
}

func maskValue(v string) string {
	if len(v) <= 4 {
		return strings.Repeat("*", len(v))
	}
	return v[:2] + strings.Repeat("*", len(v)-4) + v[len(v)-2:]
}
