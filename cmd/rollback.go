package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/user/vaultpeek/internal/vault"
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback <mount> <path> <version>",
	Short: "Rollback a KV v2 secret to a previous version",
	Args:  cobra.ExactArgs(3),
	RunE:  runRollback,
}

func init() {
	rootCmd.AddCommand(rollbackCmd)
}

func runRollback(cmd *cobra.Command, args []string) error {
	mount := args[0]
	path := args[1]
	versionStr := args[2]

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return fmt.Errorf("invalid version %q: must be a positive integer", versionStr)
	}
	if version < 1 {
		return fmt.Errorf("version must be >= 1, got %d", version)
	}

	addr, token, err := resolveVaultEnv()
	if err != nil {
		return err
	}

	client, err := vault.NewClient(addr, token)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	if err := vault.RollbackSecret(context.Background(), client, mount, path, version); err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Successfully rolled back %s/%s to version %d\n", mount, path, version)
	return nil
}
