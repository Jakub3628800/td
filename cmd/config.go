package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Jakub3628800/td/internal/core"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration options",
	Long:  `Configure td settings interactively using a terminal UI.`,
	Run: func(_ *cobra.Command, _ []string) {
		runConfigTUI()
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration option",
	Long:  `Set a configuration option. Only predefined config keys are allowed.`,
	Args:  cobra.ExactArgs(2),
	Run: func(_ *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		// Validate the key
		if err := core.ValidateConfigKey(key); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			fmt.Fprintf(os.Stderr, "\nAllowed configuration keys:\n")
			for k, desc := range core.AllowedConfigKeys {
				fmt.Fprintf(os.Stderr, "  %-25s - %s\n", k, desc)
			}
			os.Exit(1)
		}

		// Set the config
		queries, err := core.GetDB()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
			os.Exit(1)
		}

		ctx := context.Background()
		if err := core.SetConfig(ctx, queries, key, value); err != nil {
			fmt.Fprintf(os.Stderr, "Error setting config: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Configuration updated: %s = %s\n", key, value)
	},
}


func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
}
