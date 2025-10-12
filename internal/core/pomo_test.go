package core

import (
	"os"
	"testing"
	"time"
)

func TestHasRunningPomodoro_NoPomodoros(t *testing.T) {
	dir := t.TempDir()
	origDBPath := os.Getenv("TD_DB_PATH")
	defer func() {
		os.Setenv("TD_DB_PATH", origDBPath)
		CloseDB()
	}()

	os.Setenv("TD_DB_PATH", dir+"/td.db")
	CloseDB()

	hasRunning, err := HasRunningPomodoro()
	if err != nil {
		t.Fatalf("HasRunningPomodoro failed: %v", err)
	}
	if hasRunning {
		t.Error("expected no running pomodoro, but got true")
	}
}

func TestHasRunningPomodoro_WithRunningPomo(t *testing.T) {
	dir := t.TempDir()
	origDBPath := os.Getenv("TD_DB_PATH")
	defer func() {
		os.Setenv("TD_DB_PATH", origDBPath)
		CloseDB()
	}()

	os.Setenv("TD_DB_PATH", dir+"/td.db")
	CloseDB()

	// Create a pomodoro that started 5 minutes ago with 25 minute duration
	now := time.Now()
	startTime := now.Add(-5 * time.Minute)

	err := SavePomodoroLog(25, "running", []string{}, startTime)
	if err != nil {
		t.Fatalf("SavePomodoroLog failed: %v", err)
	}

	hasRunning, err := HasRunningPomodoro()
	if err != nil {
		t.Fatalf("HasRunningPomodoro failed: %v", err)
	}
	if !hasRunning {
		t.Error("expected running pomodoro, but got false")
	}
}

func TestHasRunningPomodoro_WithCompletedPomo(t *testing.T) {
	dir := t.TempDir()
	origDBPath := os.Getenv("TD_DB_PATH")
	defer func() {
		os.Setenv("TD_DB_PATH", origDBPath)
		CloseDB()
	}()

	os.Setenv("TD_DB_PATH", dir+"/td.db")
	CloseDB()

	// Create a pomodoro that started 30 minutes ago with 25 minute duration (expired)
	now := time.Now()
	startTime := now.Add(-30 * time.Minute)

	err := SavePomodoroLog(25, "completed", []string{}, startTime)
	if err != nil {
		t.Fatalf("SavePomodoroLog failed: %v", err)
	}

	hasRunning, err := HasRunningPomodoro()
	if err != nil {
		t.Fatalf("HasRunningPomodoro failed: %v", err)
	}
	if hasRunning {
		t.Error("expected no running pomodoro (completed), but got true")
	}
}

func TestHasRunningPomodoro_MultiplePomos(t *testing.T) {
	dir := t.TempDir()
	origDBPath := os.Getenv("TD_DB_PATH")
	defer func() {
		os.Setenv("TD_DB_PATH", origDBPath)
		CloseDB()
	}()

	os.Setenv("TD_DB_PATH", dir+"/td.db")
	CloseDB()

	now := time.Now()

	// Add old completed pomodoro
	err := SavePomodoroLog(25, "completed", []string{}, now.Add(-2*time.Hour))
	if err != nil {
		t.Fatalf("SavePomodoroLog failed: %v", err)
	}

	// Add currently running pomodoro
	err = SavePomodoroLog(25, "running", []string{"work"}, now.Add(-10*time.Minute))
	if err != nil {
		t.Fatalf("SavePomodoroLog failed: %v", err)
	}

	hasRunning, err := HasRunningPomodoro()
	if err != nil {
		t.Fatalf("HasRunningPomodoro failed: %v", err)
	}
	if !hasRunning {
		t.Error("expected running pomodoro among multiple, but got false")
	}
}
