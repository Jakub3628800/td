package core

import (
	"encoding/json"
	"os"
	"path/filepath"
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
	os.Setenv("TD_VAULT_LOC", dir)
	os.Setenv("TD_INTERVAL_MODE", "daily")
	now := time.Date(2025, 5, 10, 12, 0, 0, 0, time.UTC)
	log := DayLog{
		Date: "2025-05-10",
		Start: DayStart{
			ShutdownTime: now,
			DayGoal:      "Goal",
			StartedAt:    now,
		},
	}
	if err := SaveDayLog(now, log); err != nil {
		t.Fatalf("SaveDayLog failed: %v", err)
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
	os.Setenv("TD_VAULT_LOC", dir)
	os.Setenv("TD_INTERVAL_MODE", "daily")
	now := time.Date(2025, 5, 10, 12, 0, 0, 0, time.UTC)
	if err := SaveDayLog(now, DayLog{Date: "2025-05-10"}); err != nil {
		t.Fatalf("SaveDayLog failed: %v", err)
	}
	err := UpdateDayLog(now, func(log *DayLog) {
		log.Start.DayGoal = "Updated"
	})
	if err != nil {
		t.Fatalf("UpdateDayLog failed: %v", err)
	}
	loaded, _ := LoadDayLog(now)
	if loaded.Start.DayGoal != "Updated" {
		t.Errorf("UpdateDayLog did not update: got %v", loaded.Start.DayGoal)
	}
}

func TestPomoSessionWritesToDailyJson(t *testing.T) {
	dir := t.TempDir()
	os.Setenv("TD_VAULT_LOC", dir)
	os.Setenv("TD_INTERVAL_MODE", "weekly") // Should be ignored for day log
	now := time.Date(2025, 5, 10, 19, 44, 0, 0, time.UTC)
	ResetForTest(dir, "weekly")

	// Record a pomodoro session
	dur := 1
	status := "completed"
	err := UpdateDayLog(now, func(log *DayLog) {
		pomo := PomodoroLog{
			Duration:  dur,
			Status:    status,
			Timestamp: now,
		}
		log.Pomodoros = append(log.Pomodoros, pomo)
	})
	if err != nil {
		t.Fatalf("UpdateDayLog failed: %v", err)
	}

	// Check that the daily JSON file exists and contains the pomodoro
	file := filepath.Join(dir, "2025", "May", "10.json")
	b, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("Expected daily JSON file not found: %v", err)
	}
	var log DayLog
	if err := json.Unmarshal(b, &log); err != nil {
		t.Fatalf("Failed to unmarshal daily log: %v", err)
	}
	if len(log.Pomodoros) == 0 || log.Pomodoros[0].Duration != dur || log.Pomodoros[0].Status != status {
		t.Errorf("Pomodoro session not recorded correctly in daily JSON file: %+v", log.Pomodoros)
	}
}
