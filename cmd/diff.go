package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultpeek/vaultpeek/internal/diff"
	"github.com/vaultpeek/vaultpeek/internal/vault"
)

var diffCmd = &cobra.Command{
	Use:   "diff <path-a> <path-b>",
	Short: "Diff secrets between two Vault paths",
	Args:  cobra.ExactArgs(2),
	RunE:  runDiff,
}

func init() {
	rootCmd.AddCommand(diffCmd)
}

func runDiff(cmd *cobra.Command, args []string) error {
	client, err := vault.NewClient(vaultAddr, vaultToken)
	if err != nil {
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	pathA := args[0]
	pathB := args[1]

	secretA, err := vault.FetchSecret(client, pathA)
	if err != nil {
		return fmt.Errorf("failed to fetch secret at %q: %w", pathA, err)
	}

	secretB, err := vault.FetchSecret(client, pathB)
	if err != nil {
		return fmt.Errorf("failed to fetch secret at %q: %w", pathB, err)
	}

	result := diff.Compare(secretA, secretB)

	if result.IsEmpty() {
		fmt.Println("No differences found.")
		return nil
	}

	output := diff.Render(result, pathA, pathB)
	fmt.Fprint(os.Stdout, output)
	return nil
}
