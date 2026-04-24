package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpeek/internal/vault"
)

var (
	exportMount  string
	exportFormat string
)

var exportCmd = &cobra.Command{
	Use:   "export <path>",
	Short: "Export a secret to stdout in JSON or env format",
	Args:  cobra.ExactArgs(1),
	RunE:  runExport,
}

func init() {
	exportCmd.Flags().StringVar(&exportMount, "mount", "secret", "KV v2 mount path")
	exportCmd.Flags().StringVar(&exportFormat, "format", "json", "Output format: json or env")
	rootCmd.AddCommand(exportCmd)
}

func runExport(cmd *cobra.Command, args []string) error {
	path := args[0]

	client, err := vault.NewClient("", "")
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	format := vault.ExportFormat(exportFormat)
	if format != vault.FormatJSON && format != vault.FormatEnv {
		return fmt.Errorf("unsupported format %q: choose 'json' or 'env'", exportFormat)
	}

	if err := vault.ExportSecret(client, exportMount, path, format, os.Stdout); err != nil {
		return fmt.Errorf("export: %w", err)
	}
	return nil
}
