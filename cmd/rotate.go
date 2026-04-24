package cmd

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpeek/internal/vault"
)

var (
	rotateMount  string
	rotateLength int
)

var rotateCmd = &cobra.Command{
	Use:   "rotate <path>",
	Short: "Rotate all values in a secret with randomly generated strings",
	Args:  cobra.ExactArgs(1),
	RunE:  runRotate,
}

func init() {
	rotateCmd.Flags().StringVar(&rotateMount, "mount", "secret", "KV v2 mount name")
	rotateCmd.Flags().IntVar(&rotateLength, "length", 32, "Length of generated secret values")
	rootCmd.AddCommand(rotateCmd)
}

func runRotate(cmd *cobra.Command, args []string) error {
	path := args[0]

	client, err := vault.NewClient()
	if err != nil {
		return fmt.Errorf("rotate: %w", err)
	}

	result, err := vault.RotateSecret(client, rotateMount, path, func(key string) string {
		return randomString(rotateLength)
	})
	if err != nil {
		return fmt.Errorf("rotate: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Rotated %s/%s\n", result.Mount, result.Path)
	fmt.Fprintf(cmd.OutOrStdout(), "Keys rotated: %s\n", strings.Join(result.NewKeys, ", "))
	return nil
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
