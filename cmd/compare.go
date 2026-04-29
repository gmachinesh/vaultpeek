package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"vaultpeek/internal/vault"
)

var compareMount string

var compareCmd = &cobra.Command{
	Use:   "compare <mount/path>",
	Short: "Compare a secret across two Vault environments",
	Long: `Fetches the same secret path from two Vault instances (or namespaces)
and displays a side-by-side diff of their key/value pairs.

The source Vault is configured via VAULT_ADDR / VAULT_TOKEN.
The destination Vault is configured via VAULT_ADDR_DST / VAULT_TOKEN_DST.`,
	Args: cobra.ExactArgs(1),
	RunE: runCompare,
}

func init() {
	compareCmd.Flags().StringVar(&compareMount, "mount", "secret", "KV mount name")
	rootCmd.AddCommand(compareCmd)
}

func runCompare(cmd *cobra.Command, args []string) error {
	path := strings.TrimPrefix(args[0], "/")
	if path == "" {
		return fmt.Errorf("path must not be empty")
	}

	srcAddr := os.Getenv("VAULT_ADDR")
	srcToken := os.Getenv("VAULT_TOKEN")
	dstAddr := os.Getenv("VAULT_ADDR_DST")
	dstToken := os.Getenv("VAULT_TOKEN_DST")

	if dstAddr == "" {
		return fmt.Errorf("VAULT_ADDR_DST must be set for compare")
	}
	if dstToken == "" {
		return fmt.Errorf("VAULT_TOKEN_DST must be set for compare")
	}

	srcClient, err := vault.NewClient(srcAddr, srcToken)
	if err != nil {
		return fmt.Errorf("src client: %w", err)
	}
	dstClient, err := vault.NewClient(dstAddr, dstToken)
	if err != nil {
		return fmt.Errorf("dst client: %w", err)
	}

	result, err := vault.CompareAcrossEnvs(srcClient, srcAddr, dstClient, dstAddr, compareMount, path)
	if err != nil {
		return err
	}

	fmt.Println(result.Summary())
	if len(result.OnlyInA) > 0 {
		fmt.Printf("  Only in src (%s): %s\n", result.SrcEnv, strings.Join(result.OnlyInA, ", "))
	}
	if len(result.OnlyInB) > 0 {
		fmt.Printf("  Only in dst (%s): %s\n", result.DstEnv, strings.Join(result.OnlyInB, ", "))
	}
	for _, k := range result.Differ {
		fmt.Printf("  ~ %s\n", k)
	}
	for _, k := range result.Match {
		fmt.Printf("  = %s\n", k)
	}
	return nil
}
