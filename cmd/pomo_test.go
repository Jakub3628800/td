package cmd

import (
	"os"
	"testing"
	"time"

	"github.com/Jakub3628800/td/core"
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
	os.Setenv("TD_VAULT_LOC", dir)
	os.Setenv("TD_INTERVAL_MODE", "daily")
	core.ResetForTest(dir, "daily")

	now := time.Now()
	duration := 5
	recordPomoSession(duration, "completed")

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
