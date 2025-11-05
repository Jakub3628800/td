package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/Jakub3628800/td/internal/core"
	"github.com/Jakub3628800/td/internal/db"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration options",
	Long:  `Manage configuration options stored in the database. Run without subcommands to list all configs.`,
	Run: func(_ *cobra.Command, _ []string) {
		listConfig()
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

func listConfig() {
	queries, err := core.GetDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	configs, err := queries.ListConfig(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing config: %v\n", err)
		os.Exit(1)
	}

	if len(configs) == 0 {
		fmt.Println("No configuration options set.")
		return
	}

	printConfigTable(configs)
}

func printConfigTable(configs []db.Config) {
	// Calculate column widths
	maxKeyLen := len("Key")
	maxValueLen := len("Value")
	maxUpdatedLen := len("Updated At")

	// Find max lengths from data
	for _, c := range configs {
		if len(c.Key) > maxKeyLen {
			maxKeyLen = len(c.Key)
		}

		value := ""
		if c.Value.Valid {
			value = c.Value.String
		}
		if len(value) > maxValueLen {
			maxValueLen = len(value)
		}

		updated := ""
		if c.UpdatedAt.Valid {
			updated = c.UpdatedAt.Time.Format(time.RFC3339)
		}
		if len(updated) > maxUpdatedLen {
			maxUpdatedLen = len(updated)
		}
	}

	// Ensure minimum widths
	if maxValueLen < 8 {
		maxValueLen = 8
	}

	// Print header
	totalWidth := maxKeyLen + maxValueLen + maxUpdatedLen + 7 // 7 for separators and spacing
	fmt.Println(strings.Repeat("-", totalWidth))
	fmt.Printf("%-"+fmt.Sprint(maxKeyLen)+"s | %-"+fmt.Sprint(maxValueLen)+"s | %-"+fmt.Sprint(maxUpdatedLen)+"s\n",
		"Key", "Value", "Updated At")
	fmt.Println(strings.Repeat("-", totalWidth))

	// Print rows
	for _, c := range configs {
		value := "<null>"
		if c.Value.Valid {
			value = c.Value.String
		}

		updated := ""
		if c.UpdatedAt.Valid {
			updated = c.UpdatedAt.Time.Format(time.RFC3339)
		}

		fmt.Printf("%-"+fmt.Sprint(maxKeyLen)+"s | %-"+fmt.Sprint(maxValueLen)+"s | %-"+fmt.Sprint(maxUpdatedLen)+"s\n",
			c.Key, value, updated)
	}
	fmt.Println(strings.Repeat("-", totalWidth))
	fmt.Printf("Total: %d configuration option(s)\n", len(configs))
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
}
