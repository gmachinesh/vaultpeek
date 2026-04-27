package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultpeek/internal/vault"
)

var lockCmd = &cobra.Command{
	Use:   "lock <mount> <path>",
	Short: "Lock a secret to prevent modifications",
	Long: `Lock a secret at the given mount and path.

A locked secret cannot be modified or deleted until it is explicitly
unlocked. This is useful for protecting critical secrets in production
environments.

Example:
  vaultpeek lock secret prod/database/password`,
	Args: cobra.ExactArgs(2),
	RunE: runLock,
}

var unlockCmd = &cobra.Command{
	Use:   "unlock <mount> <path>",
	Short: "Unlock a previously locked secret",
	Long: `Unlock a secret at the given mount and path, restoring write access.

Example:
  vaultpeek unlock secret prod/database/password`,
	Args: cobra.ExactArgs(2),
	RunE: runUnlock,
}

var statusCmd = &cobra.Command{
	Use:   "lock-status <mount> <path>",
	Short: "Check whether a secret is locked",
	Long: `Check the lock status of a secret at the given mount and path.

Example:
  vaultpeek lock-status secret prod/database/password`,
	Args: cobra.ExactArgs(2),
	RunE: runLockStatus,
}

func init() {
	rootCmd.AddCommand(lockCmd)
	rootCmd.AddCommand(unlockCmd)
	rootCmd.AddCommand(statusCmd)
}

func runLock(cmd *cobra.Command, args []string) error {
	mount := args[0]
	path := args[1]

	client, err := vault.NewClient(os.Getenv("VAULT_ADDR"), os.Getenv("VAULT_TOKEN"))
	if err != nil {
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	if err := vault.LockSecret(client, mount, path); err != nil {
		return fmt.Errorf("failed to lock secret %s/%s: %w", mount, path, err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Secret %s/%s locked successfully.\n", mount, path)
	return nil
}

func runUnlock(cmd *cobra.Command, args []string) error {
	mount := args[0]
	path := args[1]

	client, err := vault.NewClient(os.Getenv("VAULT_ADDR"), os.Getenv("VAULT_TOKEN"))
	if err != nil {
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	if err := vault.UnlockSecret(client, mount, path); err != nil {
		return fmt.Errorf("failed to unlock secret %s/%s: %w", mount, path, err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Secret %s/%s unlocked successfully.\n", mount, path)
	return nil
}

func runLockStatus(cmd *cobra.Command, args []string) error {
	mount := args[0]
	path := args[1]

	client, err := vault.NewClient(os.Getenv("VAULT_ADDR"), os.Getenv("VAULT_TOKEN"))
	if err != nil {
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	locked, err := vault.IsLocked(client, mount, path)
	if err != nil {
		return fmt.Errorf("failed to check lock status for %s/%s: %w", mount, path, err)
	}

	if locked {
		fmt.Fprintf(cmd.OutOrStdout(), "Secret %s/%s is LOCKED.\n", mount, path)
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "Secret %s/%s is UNLOCKED.\n", mount, path)
	}
	return nil
}
