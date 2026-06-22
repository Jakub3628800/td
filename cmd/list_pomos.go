package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Jakub3628800/td/internal/core"
	"github.com/Jakub3628800/td/internal/db"
)

var (
	listPomosAfter  string
	listPomosBefore string
)

var listPomosHelp = `List all pomodoro sessions in a table format, ordered from most recent to oldest.

Usage:
  td list-pomos [--after YYYY-MM-DD] [--before YYYY-MM-DD]
`

func runListPomosCommand(args []string) error {
	listPomosAfter = ""
	listPomosBefore = ""

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "-h" || arg == "--help":
			fmt.Print(listPomosHelp)
			return nil
		case arg == "--after":
			if i+1 >= len(args) {
				return fmt.Errorf("--after requires a value")
			}
			i++
			listPomosAfter = args[i]
		case strings.HasPrefix(arg, "--after="):
			listPomosAfter = strings.TrimPrefix(arg, "--after=")
		case arg == "--before":
			if i+1 >= len(args) {
				return fmt.Errorf("--before requires a value")
			}
			i++
			listPomosBefore = args[i]
		case strings.HasPrefix(arg, "--before="):
			listPomosBefore = strings.TrimPrefix(arg, "--before=")
		default:
			return fmt.Errorf("unknown list-pomos argument %q\n\n%s", arg, listPomosHelp)
		}
	}

	queries, err := core.GetDB()
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	ctx := context.Background()

	var pomodoros []db.Pomodori
	if listPomosAfter != "" || listPomosBefore != "" {
		afterTime := time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local)
		beforeTime := time.Date(2100, 1, 1, 0, 0, 0, 0, time.Local)

		if listPomosAfter != "" {
			parsed, err := time.ParseInLocation("2006-01-02", listPomosAfter, time.Local)
			if err != nil {
				return fmt.Errorf("error parsing --after date: %w", err)
			}
			afterTime = parsed
		}

		if listPomosBefore != "" {
			parsed, err := time.ParseInLocation("2006-01-02", listPomosBefore, time.Local)
			if err != nil {
				return fmt.Errorf("error parsing --before date: %w", err)
			}
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
		return fmt.Errorf("error listing pomodoros: %w", err)
	}

	if len(pomodoros) == 0 {
		fmt.Println("No pomodoro sessions recorded yet.")
		return nil
	}

	printPomodoroTable(pomodoros)
	return nil
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
