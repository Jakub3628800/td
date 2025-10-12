package cmd

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Jakub3628800/td/internal/core"
)

func TestPomoIntegration(t *testing.T) {
	// Skip this test as it requires user interaction
	t.Skip("Skipping integration test that requires user interaction")

	// The test body is skipped completely
}

// TestRootCommandPomoKey tests that pressing 'p' in the root command
// launches a pomodoro session and adds a record when completed
func TestRootCommandPomoKey(t *testing.T) {
	// This is a more complex test that would require mocking the tea.Program
	// and simulating key presses, which is beyond the scope of this implementation.
	// In a real-world scenario, you would use a testing framework that can
	// simulate user input and verify the resulting state changes.
	t.Skip("Integration test requiring user interaction - skipped")
}

func TestRecordPomoSession(t *testing.T) {
	dir := t.TempDir()
	origDBPath := os.Getenv("TD_DB_PATH")
	defer func() {
		os.Setenv("TD_DB_PATH", origDBPath)
		core.CloseDB()
	}()

	os.Setenv("TD_DB_PATH", dir+"/td.db")
	os.Setenv("TD_INTERVAL_MODE", "daily")
	core.CloseDB()

	now := time.Now()
	duration := 5
	recordPomoSession(duration, "completed", []string{})

	log, err := core.LoadDayLog(now)
	if err != nil {
		t.Fatalf("LoadDayLog failed: %v", err)
	}
	if len(log.Pomodoros) == 0 {
		t.Fatal("no pomodoro sessions recorded")
	}
	p := log.Pomodoros[0]
	if p.Duration != duration || p.Status != "completed" {
		t.Errorf("unexpected pomodoro entry: %+v", p)
	}
}

func TestRecordPomoSessionWithTags(t *testing.T) {
	dir := t.TempDir()
	origDBPath := os.Getenv("TD_DB_PATH")
	defer func() {
		os.Setenv("TD_DB_PATH", origDBPath)
		core.CloseDB()
	}()

	os.Setenv("TD_DB_PATH", dir+"/td.db")
	core.CloseDB()

	duration := 25
	tags := []string{"work", "planning", "feature-x"}

	recordPomoSession(duration, "completed", tags)

	// Verify pomodoro was saved with tags
	queries, err := core.GetDB()
	if err != nil {
		t.Fatalf("GetDB failed: %v", err)
	}

	ctx := context.Background()
	pomodoros, err := queries.ListPomodori(ctx)
	if err != nil {
		t.Fatalf("ListPomodori failed: %v", err)
	}

	if len(pomodoros) == 0 {
		t.Fatal("no pomodoro sessions recorded")
	}

	p := pomodoros[0]
	if p.DurationMinutes != int64(duration) {
		t.Errorf("expected duration %d, got %d", duration, p.DurationMinutes)
	}
	if !p.Completed.Valid || !p.Completed.Bool {
		t.Error("expected completed status")
	}
	if !p.Tags.Valid || p.Tags.String != "work,planning,feature-x" {
		t.Errorf("expected tags 'work,planning,feature-x', got %v", p.Tags)
	}
}

func TestRecordPomoSessionCancelled(t *testing.T) {
	dir := t.TempDir()
	origDBPath := os.Getenv("TD_DB_PATH")
	defer func() {
		os.Setenv("TD_DB_PATH", origDBPath)
		core.CloseDB()
	}()

	os.Setenv("TD_DB_PATH", dir+"/td.db")
	core.CloseDB()

	now := time.Now()
	duration := 25

	recordPomoSession(duration, "cancelled", []string{})

	log, err := core.LoadDayLog(now)
	if err != nil {
		t.Fatalf("LoadDayLog failed: %v", err)
	}
	if len(log.Pomodoros) == 0 {
		t.Fatal("no pomodoro sessions recorded")
	}

	p := log.Pomodoros[0]
	if p.Status != "cancelled" {
		t.Errorf("expected status 'cancelled', got %s", p.Status)
	}
}
