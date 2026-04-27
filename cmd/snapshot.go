package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"vaultpeek/internal/vault"
)

var snapshotMount string

func init() {
	snapshotCmd := &cobra.Command{
		Use:   "snapshot <prefix>",
		Short: "Capture a point-in-time snapshot of all secrets under a prefix",
		Args:  cobra.ExactArgs(1),
		RunE:  runSnapshot,
	}

	snapshotCmd.Flags().StringVar(&snapshotMount, "mount", "secret", "KV mount name")
	rootCmd.AddCommand(snapshotCmd)
}

func runSnapshot(cmd *cobra.Command, args []string) error {
	prefix := args[0]

	client, err := vault.NewClient("", "")
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	result, err := vault.TakeSnapshot(client, snapshotMount, prefix)
	if err != nil {
		return fmt.Errorf("taking snapshot: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Snapshot taken at %s\n", result.Snapshot.TakenAt.Format("2006-01-02T15:04:05Z"))
	fmt.Fprintf(cmd.OutOrStdout(), "Mount: %s  Prefix: %s  Secrets: %d\n\n",
		result.Snapshot.Mount, result.Snapshot.Prefix, result.Count)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PATH\tKEYS")
	for path, sm := range result.Snapshot.Secrets {
		fmt.Fprintf(w, "%s\t%d\n", path, len(sm.Keys()))
	}
	w.Flush()

	return nil
}
