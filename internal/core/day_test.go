package core

import (
	"context"
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
	origDBPath := os.Getenv("TD_DB_PATH")
	defer func() {
		os.Setenv("TD_DB_PATH", origDBPath)
		CloseDB()
	}()

	os.Setenv("TD_DB_PATH", dir+"/td.db")
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
	origDBPath := os.Getenv("TD_DB_PATH")
	defer func() {
		os.Setenv("TD_DB_PATH", origDBPath)
		CloseDB()
	}()

	os.Setenv("TD_DB_PATH", dir+"/td.db")
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
	origDBPath := os.Getenv("TD_DB_PATH")
	defer func() {
		os.Setenv("TD_DB_PATH", origDBPath)
		CloseDB()
	}()

	os.Setenv("TD_DB_PATH", dir+"/td.db")
	os.Setenv("TD_INTERVAL_MODE", "weekly") // Should be ignored for day log
	CloseDB()

	now := time.Date(2025, 5, 12, 19, 44, 0, 0, time.UTC)

	// Record a pomodoro session
	dur := 1
	status := "completed"
	if err := SavePomodoroLog(dur, status, []string{}, now); err != nil {
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

func TestSavePomodoroLogWithTags(t *testing.T) {
	dir := t.TempDir()
	origDBPath := os.Getenv("TD_DB_PATH")
	defer func() {
		os.Setenv("TD_DB_PATH", origDBPath)
		CloseDB()
	}()

	os.Setenv("TD_DB_PATH", dir+"/td.db")
	CloseDB()

	now := time.Date(2025, 5, 13, 14, 30, 0, 0, time.UTC)
	tags := []string{"work", "planning", "feature-x"}

	// Save pomodoro with tags
	err := SavePomodoroLog(25, "completed", tags, now)
	if err != nil {
		t.Fatalf("SavePomodoroLog failed: %v", err)
	}

	// Verify tags were saved by querying database directly
	queries, err := GetDB()
	if err != nil {
		t.Fatalf("GetDB failed: %v", err)
	}

	pomodoros, err := queries.ListPomodori(ctx())
	if err != nil {
		t.Fatalf("ListPomodori failed: %v", err)
	}

	if len(pomodoros) == 0 {
		t.Fatal("no pomodoros found")
	}

	pomo := pomodoros[0]
	if !pomo.Tags.Valid {
		t.Fatal("tags should be valid")
	}

	expectedTags := "work,planning,feature-x"
	if pomo.Tags.String != expectedTags {
		t.Errorf("expected tags %q, got %q", expectedTags, pomo.Tags.String)
	}
}

func TestSavePomodoroLogWithEmptyTags(t *testing.T) {
	dir := t.TempDir()
	origDBPath := os.Getenv("TD_DB_PATH")
	defer func() {
		os.Setenv("TD_DB_PATH", origDBPath)
		CloseDB()
	}()

	os.Setenv("TD_DB_PATH", dir+"/td.db")
	CloseDB()

	now := time.Date(2025, 5, 13, 15, 0, 0, 0, time.UTC)

	// Save pomodoro with no tags
	err := SavePomodoroLog(25, "completed", []string{}, now)
	if err != nil {
		t.Fatalf("SavePomodoroLog failed: %v", err)
	}

	// Verify tags field is null/empty
	queries, err := GetDB()
	if err != nil {
		t.Fatalf("GetDB failed: %v", err)
	}

	pomodoros, err := queries.ListPomodori(ctx())
	if err != nil {
		t.Fatalf("ListPomodori failed: %v", err)
	}

	if len(pomodoros) == 0 {
		t.Fatal("no pomodoros found")
	}

	pomo := pomodoros[0]
	if pomo.Tags.Valid && pomo.Tags.String != "" {
		t.Errorf("expected empty/null tags, got %q", pomo.Tags.String)
	}
}

func TestSaveDayEnd(t *testing.T) {
	dir := t.TempDir()
	origDBPath := os.Getenv("TD_DB_PATH")
	defer func() {
		os.Setenv("TD_DB_PATH", origDBPath)
		CloseDB()
	}()

	os.Setenv("TD_DB_PATH", dir+"/td.db")
	CloseDB()

	now := time.Date(2025, 5, 14, 12, 0, 0, 0, time.UTC)

	// First create a day start
	start := DayStart{
		ShutdownTime: now.Add(8 * time.Hour),
		DayGoal:      "Test goal",
		StartedAt:    now,
	}
	err := SaveDayStart(now, start)
	if err != nil {
		t.Fatalf("SaveDayStart failed: %v", err)
	}

	// Save day end
	endTime := now.Add(9 * time.Hour)
	end := DayEnd{
		Rating:     2,
		Reason:     "Very productive day",
		FocusHours: 6.5,
		FinishedAt: endTime,
	}
	err = SaveDayEnd(now, end)
	if err != nil {
		t.Fatalf("SaveDayEnd failed: %v", err)
	}

	// Verify day end was saved
	loaded, err := LoadDayLog(now)
	if err != nil {
		t.Fatalf("LoadDayLog failed: %v", err)
	}

	if loaded.End.Rating != 2 {
		t.Errorf("expected rating 2, got %d", loaded.End.Rating)
	}
	if loaded.End.Reason != "Very productive day" {
		t.Errorf("expected reason %q, got %q", "Very productive day", loaded.End.Reason)
	}
	if loaded.End.FocusHours != 6.5 {
		t.Errorf("expected focus hours 6.5, got %f", loaded.End.FocusHours)
	}
}

func ctx() context.Context {
	return context.Background()
}
