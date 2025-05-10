/*
Copyright © 2024 Jakub Kriz
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/Jakub3628800/td/core"
)

var dateFlag string
var debugFlag bool

var captureCmd = &cobra.Command{
	Use:   "capture",
	Short: "Capture an action to today's list.",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		date, err := parseDate(dateFlag)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			return
		}
		contains, _ := core.ContainsLine(date, args[0])
		if contains != 0 {
			fmt.Println("This item already exists. Skipping")
		} else {
			err := core.AddTask(date, args[0])
			if err != nil {
				fmt.Println("Error appending line to file.")
			}
		}
	},
}

func printDebugInfo() {
	now := time.Now()
	vaultLoc := core.GetVaultLocation()
	intervalMode := core.GetIntervalMode()
	year, month, day := now.Date()
	// Always show the daily JSON file for day logs and Pomodoros
	file := fmt.Sprintf("%s/%d/%s/%02d.json", vaultLoc, year, month.String(), day)
	yellow := lipgloss.NewStyle().Foreground(lipgloss.Color("#fbc02d")).Bold(true).Render
	fmt.Printf("Today is %s\n", yellow(now.Format("2006-01-02")))
	fmt.Printf("currently set vault location: %s\n", yellow(vaultLoc))
	fmt.Printf("currently set vault mode: %s\n", yellow(intervalMode))
	fmt.Printf("file: %s\n", yellow(file))
	vars := []string{"TD_VAULT_LOC", "TD_INTERVAL_MODE", "TD_TEMPLATE_PATH", "TD_SKIP_WEEKEND", "TD_COPY_PREVIOUS", "TD_TEST_MODE"}
	fmt.Println("\nEnvironment variables:")
	for _, v := range vars {
		val := os.Getenv(v)
		fmt.Printf("%s=%s\n", yellow(v), yellow(val))
	}
}

func init() {
	rootCmd.AddCommand(captureCmd)
	captureCmd.Flags().StringVar(&dateFlag, "date", "today", "Date (today, tomorrow, yesterday, or YYYY-MM-DD)")
	rootCmd.PersistentFlags().BoolVarP(&debugFlag, "debug", "", false, "Show debug info about vault and today.")
	cobra.OnInitialize(func() {
		if debugFlag {
			printDebugInfo()
			os.Exit(0)
		}
	})
}

func parseDate(input string) (time.Time, error) {
	now := time.Now()
	switch input {
	case "today":
		return now, nil
	case "tomorrow":
		return now.AddDate(0, 0, 1), nil
	case "yesterday":
		return now.AddDate(0, 0, -1), nil
	default:
		return time.Parse("2006-01-02", input)
	}
}
