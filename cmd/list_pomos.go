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

var (
	listPomosAfter  string
	listPomosBefore string
)

var listPomosCmd = &cobra.Command{
	Use:   "list-pomos",
	Short: "List all pomodoro sessions",
	Long:  `List all pomodoro sessions in a table format, ordered from most recent to oldest.`,
	Run: func(_ *cobra.Command, _ []string) {
		queries, err := core.GetDB()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
			os.Exit(1)
		}

		ctx := context.Background()

		var pomodoros []db.Pomodori
		if listPomosAfter != "" || listPomosBefore != "" {
			afterTime := time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local)
			beforeTime := time.Date(2100, 1, 1, 0, 0, 0, 0, time.Local)

			if listPomosAfter != "" {
				parsed, err := time.ParseInLocation("2006-01-02", listPomosAfter, time.Local)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error parsing --after date: %v\n", err)
					os.Exit(1)
				}
				afterTime = parsed
			}

			if listPomosBefore != "" {
				parsed, err := time.ParseInLocation("2006-01-02", listPomosBefore, time.Local)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error parsing --before date: %v\n", err)
					os.Exit(1)
				}
				// Add one day to include the entire "before" date
				beforeTime = parsed.AddDate(0, 0, 1)
			}

			pomodoros, err = queries.ListPomodoriFiltered(ctx, db.ListPomodoriFilteredParams{
				StartTime:   afterTime,
				StartTime_2: beforeTime,
			})
		} else {
			pomodoros, err = queries.ListPomodori(ctx)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing pomodoros: %v\n", err)
			os.Exit(1)
		}

		if len(pomodoros) == 0 {
			fmt.Println("No pomodoro sessions recorded yet.")
			return
		}

		printPomodoroTable(pomodoros)
	},
}

func printPomodoroTable(pomodoros []db.Pomodori) {
	fmt.Println(strings.Repeat("-", 100))
	fmt.Printf("%-5s %-12s %-10s %-12s %-12s %-30s\n", "ID", "Date", "Duration", "Status", "Start Time", "Tags")
	fmt.Println(strings.Repeat("-", 100))

	now := time.Now()

	for _, p := range pomodoros {
		status := getPomoStatus(p, now)

		date := p.StartTime.Format("2006-01-02")
		startTime := p.StartTime.Format("15:04:05")

		duration := fmt.Sprintf("%dm", p.DurationMinutes)

		tags := ""
		if p.Tags.Valid && p.Tags.String != "" {
			tags = p.Tags.String
		}

		fmt.Printf("%-5d %-12s %-10s %-12s %-12s %-30s\n",
			p.ID,
			date,
			duration,
			status,
			startTime,
			tags,
		)
	}
	fmt.Println(strings.Repeat("-", 100))
	fmt.Printf("Total: %d pomodoro sessions\n", len(pomodoros))
}

func getPomoStatus(p db.Pomodori, now time.Time) string {
	expectedEndTime := p.StartTime.Add(time.Duration(p.DurationMinutes) * time.Minute)

	if now.Before(expectedEndTime) {
		return "running"
	}

	if p.Completed.Valid && p.Completed.Bool {
		return "completed"
	}

	return "cancelled"
}

func init() {
	rootCmd.AddCommand(listPomosCmd)
	listPomosCmd.Flags().StringVar(&listPomosAfter, "after", "", "Show pomodoros on or after this date (format: 2006-01-02)")
	listPomosCmd.Flags().StringVar(&listPomosBefore, "before", "", "Show pomodoros on or before this date (format: 2006-01-02)")
}
