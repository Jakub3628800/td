package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "td",
	Short: "A simple, efficient tool for pomodoro tracking",
	Long:  `To-Do ToDay (td) is a simple tool for tracking pomodoro sessions with configuration management.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func SetVersion(version string) {
	Version = version
}

func init() {
}
