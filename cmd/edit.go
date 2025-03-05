package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"td/core"
)

var copyPrevious bool

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit today's task file",
	Long:  `Open today's task file in your default editor.`,
	Run: func(_ *cobra.Command, _ []string) {
		date := time.Now()
		err := core.OpenEditor(date, 1, copyPrevious) // Start at line 1
		if err != nil {
			fmt.Printf("Error opening editor: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
	editCmd.Flags().BoolVarP(&copyPrevious, "copy-previous", "c", false, "Copy content from the previous day's file")
}
