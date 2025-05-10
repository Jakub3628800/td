package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "td",
	Short: "A simple, efficient TUI app for tracking tasks",
	Long:  `To-Do ToDay (td) is a simple, efficient TUI app for tracking tasks with a focus on daily workflow.`,
	Run: func(cmd *cobra.Command, args []string) {
		dayCmd.Run(cmd, args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Remove toggle flag
}
