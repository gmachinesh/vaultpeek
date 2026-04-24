package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpeek/internal/vault"
)

var renameMount string

var renameCmd = &cobra.Command{
	Use:   "rename <src-path> <dst-path>",
	Short: "Rename a secret by copying it to a new path and deleting the original",
	Args:  cobra.ExactArgs(2),
	RunE:  runRename,
}

func init() {
	renameCmd.Flags().StringVar(&renameMount, "mount", "secret", "KV v2 mount name")
	rootCmd.AddCommand(renameCmd)
}

func runRename(cmd *cobra.Command, args []string) error {
	srcPath := args[0]
	dstPath := args[1]

	if srcPath == dstPath {
		return fmt.Errorf("source and destination paths are the same: %q", srcPath)
	}

	addr := os.Getenv("VAULT_ADDR")
	token := os.Getenv("VAULT_TOKEN")

	client, err := vault.NewClient(addr, token)
	if err != nil {
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	if err := vault.RenameSecret(client, renameMount, srcPath, dstPath); err != nil {
		return fmt.Errorf("rename failed: %w", err)
	}

	fmt.Printf("Renamed %s/%s → %s/%s\n", renameMount, srcPath, renameMount, dstPath)
	return nil
}
