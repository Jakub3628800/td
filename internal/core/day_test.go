package core

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestDayLog_MarshalUnmarshal(t *testing.T) {
	shutdown := time.Date(2025, 5, 10, 17, 0, 0, 0, time.UTC)
	finished := time.Date(2025, 5, 10, 19, 0, 0, 0, time.UTC)
	log := DayLog{
		Date: "2025-05-10",
		Start: DayStart{
			ShutdownTime: shutdown,
			DayGoal:      "Test goal",
			StartedAt:    shutdown,
		},
		End: DayEnd{
			Rating:     1,
			Reason:     "Good day",
			FocusHours: 3.5,
			FinishedAt: finished,
		},
		Pomodoros: []PomodoroLog{{Duration: 25, Status: "completed", Timestamp: shutdown}},
	}
	b, err := MarshalDayLog(log)
	if err != nil {
		t.Fatalf("MarshalDayLog failed: %v", err)
	}
	var out DayLog
	if err := UnmarshalDayLog(b, &out); err != nil {
		t.Fatalf("UnmarshalDayLog failed: %v", err)
	}
	if out.Date != log.Date || out.Start.DayGoal != log.Start.DayGoal || out.End.Rating != log.End.Rating {
		t.Errorf("Unmarshalled log does not match original")
	}
}

func MarshalDayLog(log DayLog) ([]byte, error) {
	return json.MarshalIndent(log, "", "  ")
}

func UnmarshalDayLog(b []byte, log *DayLog) error {
	return json.Unmarshal(b, log)
}

func TestDayLog_FileOps(t *testing.T) {
	dir := t.TempDir()
	origVaultLoc := os.Getenv("TD_VAULT_LOC")
	defer func() {
		os.Setenv("TD_VAULT_LOC", origVaultLoc)
		CloseDB()
	}()

	os.Setenv("TD_VAULT_LOC", dir)
	os.Setenv("TD_INTERVAL_MODE", "daily")
	CloseDB()

	now := time.Date(2025, 5, 10, 12, 0, 0, 0, time.UTC)
	start := DayStart{
		ShutdownTime: now,
		DayGoal:      "Goal",
		StartedAt:    now,
	}
	if err := SaveDayStart(now, start); err != nil {
		t.Fatalf("SaveDayStart failed: %v", err)
	}
	loaded, err := LoadDayLog(now)
	if err != nil {
		t.Fatalf("LoadDayLog failed: %v", err)
	}
	if loaded.Start.DayGoal != "Goal" {
		t.Errorf("Loaded log mismatch: got %v", loaded.Start.DayGoal)
	}
}

func TestDayLog_Update(t *testing.T) {
	dir := t.TempDir()
	origVaultLoc := os.Getenv("TD_VAULT_LOC")
	defer func() {
		os.Setenv("TD_VAULT_LOC", origVaultLoc)
		CloseDB()
	}()

	os.Setenv("TD_VAULT_LOC", dir)
	os.Setenv("TD_INTERVAL_MODE", "daily")
	CloseDB()

	now := time.Date(2025, 5, 11, 12, 0, 0, 0, time.UTC)
	start := DayStart{
		ShutdownTime: now,
		DayGoal:      "Original",
		StartedAt:    now,
	}
	if err := SaveDayStart(now, start); err != nil {
		t.Fatalf("SaveDayStart failed: %v", err)
	}
	err := UpdateDayLog(now, func(log *DayLog) {
		log.Start.DayGoal = "Updated"
		log.Start.StartedAt = now
	})
	if err != nil {
		t.Fatalf("UpdateDayLog failed: %v", err)
	}
	loaded, _ := LoadDayLog(now)
	if loaded.Start.DayGoal != "Updated" {
		t.Errorf("UpdateDayLog did not update: got %v", loaded.Start.DayGoal)
	}
}

func TestPomoSessionWritesToDatabase(t *testing.T) {
	dir := t.TempDir()
	origVaultLoc := os.Getenv("TD_VAULT_LOC")
	defer func() {
		os.Setenv("TD_VAULT_LOC", origVaultLoc)
		CloseDB()
	}()

	os.Setenv("TD_VAULT_LOC", dir)
	os.Setenv("TD_INTERVAL_MODE", "weekly") // Should be ignored for day log
	CloseDB()

	now := time.Date(2025, 5, 12, 19, 44, 0, 0, time.UTC)

	// Record a pomodoro session
	dur := 1
	status := "completed"
	if err := SavePomodoroLog(dur, status, now); err != nil {
		t.Fatalf("SavePomodoroLog failed: %v", err)
	}

	// Check that the pomodoro was recorded in the database
	loaded, err := LoadDayLog(now)
	if err != nil {
		t.Fatalf("LoadDayLog failed: %v", err)
	}
	if len(loaded.Pomodoros) == 0 || loaded.Pomodoros[0].Duration != dur || loaded.Pomodoros[0].Status != status {
		t.Errorf("Pomodoro session not recorded correctly in database: %+v", loaded.Pomodoros)
	}
}
