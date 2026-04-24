package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpeek/internal/vault"
)

var (
	listMount  string
	listPrefix string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List secret keys under a given path prefix",
	RunE:  runList,
}

func init() {
	listCmd.Flags().StringVar(&listMount, "mount", "secret", "KV v2 mount name")
	listCmd.Flags().StringVar(&listPrefix, "prefix", "", "Path prefix to list (required)")
	_ = listCmd.MarkFlagRequired("prefix")
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	address, token, err := resolveVaultEnv()
	if err != nil {
		return err
	}

	client, err := vault.NewClient(address, token)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	keys, err := vault.ListSecrets(context.Background(), client, listMount, listPrefix)
	if err != nil {
		return fmt.Errorf("listing secrets: %w", err)
	}

	if len(keys) == 0 {
		fmt.Println("No secrets found.")
		return nil
	}

	for _, k := range keys {
		fmt.Println(k)
	}
	return nil
}

// resolveVaultEnv reads and validates the VAULT_ADDR and VAULT_TOKEN
// environment variables required to communicate with a Vault server.
func resolveVaultEnv() (address, token string, err error) {
	address = os.Getenv("VAULT_ADDR")
	token = os.Getenv("VAULT_TOKEN")

	if address == "" {
		return "", "", fmt.Errorf("VAULT_ADDR environment variable is not set")
	}
	if token == "" {
		return "", "", fmt.Errorf("VAULT_TOKEN environment variable is not set")
	}
	return address, token, nil
}
