package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"vaultpeek/internal/vault"
)

var tagMount string

func init() {
	tagCmd := &cobra.Command{
		Use:   "tag <path> [key=value ...]",
		Short: "Set or display custom metadata tags on a secret",
		Args:  cobra.MinimumNArgs(1),
		RunE:  runTag,
	}
	tagCmd.Flags().StringVar(&tagMount, "mount", "secret", "KV v2 mount path")
	rootCmd.AddCommand(tagCmd)
}

func runTag(cmd *cobra.Command, args []string) error {
	path := args[0]
	pairs := args[1:]

	client, err := vault.NewClient(vault.Config{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
	})
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	if len(pairs) == 0 {
		// Read mode: display existing tags
		tags, err := vault.GetTags(client, tagMount, path)
		if err != nil {
			return fmt.Errorf("fetching tags: %w", err)
		}
		if len(tags) == 0 {
			fmt.Println("(no tags set)")
			return nil
		}
		for k, v := range tags {
			fmt.Printf("%s=%s\n", k, v)
		}
		return nil
	}

	// Write mode: parse key=value pairs
	tags := make(map[string]string, len(pairs))
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid tag format %q, expected key=value", p)
		}
		tags[parts[0]] = parts[1]
	}

	if err := vault.TagSecret(client, tagMount, path, tags); err != nil {
		return fmt.Errorf("tagging secret: %w", err)
	}
	fmt.Printf("Tagged %s/%s with %d tag(s)\n", tagMount, path, len(tags))
	return nil
}
