package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpeek/internal/vault"
)

var copyCmd = &cobra.Command{
	Use:   "copy <src-path> <dst-path>",
	Short: "Copy a secret from one path to another",
	Args:  cobra.ExactArgs(2),
	RunE:  runCopy,
}

func init() {
	rootCmd.AddCommand(copyCmd)
}

func runCopy(cmd *cobra.Command, args []string) error {
	address := os.Getenv("VAULT_ADDR")
	token := os.Getenv("VAULT_TOKEN")

	client, err := vault.NewClient(address, token)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	src, dst := args[0], args[1]
	if err := vault.CopySecret(context.Background(), client, src, dst); err != nil {
		return fmt.Errorf("copy secret: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Copied %q → %q\n", src, dst)
	return nil
}
