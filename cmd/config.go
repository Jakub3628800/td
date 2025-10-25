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
	Short: "Display all configuration options",
	Long:  `Display all configuration options currently stored in the database.`,
	Run: func(_ *cobra.Command, _ []string) {
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
	},
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
}
