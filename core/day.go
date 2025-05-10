package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type DayLog struct {
	Date      string        `json:"date"`
	Start     DayStart      `json:"start"`
	End       DayEnd        `json:"end"`
	Pomodoros []PomodoroLog `json:"pomodoros"`
}

type DayStart struct {
	ShutdownTime time.Time `json:"shutdown_time"`
	DayGoal      string    `json:"day_goal"`
	StartedAt    time.Time `json:"started_at"`
}

type DayEnd struct {
	Rating     int       `json:"rating"`
	Reason     string    `json:"reason"`
	FocusHours float64   `json:"focus_hours"`
	FinishedAt time.Time `json:"finished_at"`
}

type PomodoroLog struct {
	Duration  int       `json:"duration"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

func dayLogFilename(date time.Time) string {
	year, month, day := date.Date()
	vaultLoc := GetVaultLocation()
	return filepath.Join(vaultLoc, fmt.Sprintf("%d/%s/%02d.json", year, month.String(), day))
}

func LoadDayLog(date time.Time) (DayLog, error) {
	filename := dayLogFilename(date)
	var log DayLog
	log.Date = date.Format("2006-01-02")
	if b, err := os.ReadFile(filename); err == nil {
		err = json.Unmarshal(b, &log)
		return log, err
	}
	return log, nil
}

func SaveDayLog(date time.Time, log DayLog) error {
	filename := dayLogFilename(date)
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, b, 0600)
}

func UpdateDayLog(date time.Time, updateFn func(*DayLog)) error {
	log, err := LoadDayLog(date)
	if err != nil {
		return err
	}
	updateFn(&log)
	return SaveDayLog(date, log)
}
