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
	address := os.Getenv("VAULT_ADDR")
	token := os.Getenv("VAULT_TOKEN")

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
