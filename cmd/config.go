package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/Jakub3628800/td/internal/core"
)

var configHelp = `Manage configuration options.

Usage:
  td config
  td config set <key> <value>
`

func runConfigCommand(args []string) error {
	if len(args) == 0 {
		runConfigTUI()
		return nil
	}

	switch args[0] {
	case "-h", "--help", "help":
		fmt.Print(configHelp)
		return nil
	case "set":
		return runConfigSet(args[1:])
	default:
		return fmt.Errorf("unknown config command %q\n\n%s", args[0], configHelp)
	}
}

func runConfigSet(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: td config set <key> <value>")
	}

	key := args[0]
	value := args[1]

	if err := core.ValidateConfigKey(key); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "\nAllowed configuration keys:\n")
		for k, desc := range core.AllowedConfigKeys {
			fmt.Fprintf(os.Stderr, "  %-25s - %s\n", k, desc)
		}
		return err
	}

	queries, err := core.GetDB()
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	ctx := context.Background()
	if err := core.SetConfig(ctx, queries, key, value); err != nil {
		return fmt.Errorf("error setting config: %w", err)
	}

	fmt.Printf("✓ Configuration updated: %s = %s\n", key, value)
	return nil
}
