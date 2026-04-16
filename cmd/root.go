package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	vaultAddr  string
	vaultToken string
)

var rootCmd = &cobra.Command{
	Use:   "vaultpeek",
	Short: "Inspect and diff HashiCorp Vault secrets across environments",
	Long: `vaultpeek is a CLI tool for inspecting Vault secret paths
and comparing secrets between environments or paths.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&vaultAddr,
		"address",
		os.Getenv("VAULT_ADDR"),
		"Vault server address (env: VAULT_ADDR)",
	)
	rootCmd.PersistentFlags().StringVar(
		&vaultToken,
		"token",
		os.Getenv("VAULT_TOKEN"),
		"Vault token (env: VAULT_TOKEN)",
	)
}
