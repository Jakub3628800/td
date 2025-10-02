package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "td",
	Short: "A simple, efficient TUI app for tracking tasks",
	Long:  `To-Do ToDay (td) is a simple, efficient TUI app for tracking tasks with a focus on daily workflow.`,
	Run: func(cmd *cobra.Command, args []string) {
		dayCmd.Run(cmd, args)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
