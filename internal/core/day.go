package core

import (
	"context"
	"database/sql"
	"time"

	"github.com/Jakub3628800/td/internal/db"
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

func LoadDayLog(date time.Time) (DayLog, error) {
	queries, err := GetDB()
	if err != nil {
		return DayLog{}, err
	}

	ctx := context.Background()
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	var log DayLog
	log.Date = date.Format("2006-01-02")

	dayRecord, err := queries.GetDayByDate(ctx, dateOnly)
	if err != nil && err != sql.ErrNoRows {
		return log, err
	}

	if err != sql.ErrNoRows {
		log.Start = DayStart{
			ShutdownTime: dayRecord.ShutdownTime,
			DayGoal:      dayRecord.DayGoal,
			StartedAt:    dayRecord.StartedAt,
		}

		if dayRecord.Rating.Valid {
			log.End.Rating = int(dayRecord.Rating.Int64)
		}
		if dayRecord.Reason.Valid {
			log.End.Reason = dayRecord.Reason.String
		}
		if dayRecord.FocusHours.Valid {
			log.End.FocusHours = dayRecord.FocusHours.Float64
		}
		if dayRecord.FinishedAt.Valid {
			log.End.FinishedAt = dayRecord.FinishedAt.Time
		}
	}

	pomodoros, err := queries.ListPomodori(ctx)
	if err != nil {
		return log, err
	}

	for _, pomo := range pomodoros {
		if pomo.StartTime.Year() == date.Year() && pomo.StartTime.YearDay() == date.YearDay() {
			status := "cancelled"
			if pomo.Completed.Valid && pomo.Completed.Bool {
				status = "completed"
			}
			log.Pomodoros = append(log.Pomodoros, PomodoroLog{
				Duration:  int(pomo.DurationMinutes),
				Status:    status,
				Timestamp: pomo.StartTime,
			})
		}
	}

	return log, nil
}

func SaveDayStart(date time.Time, start DayStart) error {
	queries, err := GetDB()
	if err != nil {
		return err
	}

	ctx := context.Background()
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	_, err = queries.GetDayByDate(ctx, dateOnly)
	if err == sql.ErrNoRows {
		_, err = queries.CreateDay(ctx, db.CreateDayParams{
			Date:         dateOnly,
			ShutdownTime: start.ShutdownTime,
			DayGoal:      start.DayGoal,
			StartedAt:    start.StartedAt,
		})
		return err
	} else if err != nil {
		return err
	}

	return queries.UpdateDayStart(ctx, db.UpdateDayStartParams{
		ShutdownTime: start.ShutdownTime,
		DayGoal:      start.DayGoal,
		StartedAt:    start.StartedAt,
		Date:         dateOnly,
	})
}

func SaveDayEnd(date time.Time, end DayEnd) error {
	queries, err := GetDB()
	if err != nil {
		return err
	}

	ctx := context.Background()
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	return queries.UpdateDayEnd(ctx, db.UpdateDayEndParams{
		Rating:     sql.NullInt64{Int64: int64(end.Rating), Valid: true},
		Reason:     sql.NullString{String: end.Reason, Valid: end.Reason != ""},
		FocusHours: sql.NullFloat64{Float64: end.FocusHours, Valid: true},
		FinishedAt: sql.NullTime{Time: end.FinishedAt, Valid: !end.FinishedAt.IsZero()},
		Date:       dateOnly,
	})
}

func UpdateDayLog(date time.Time, updateFn func(*DayLog)) error {
	queries, err := GetDB()
	if err != nil {
		return err
	}

	ctx := context.Background()
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	log, err := LoadDayLog(date)
	if err != nil {
		return err
	}
	updateFn(&log)

	_, err = queries.GetDayByDate(ctx, dateOnly)
	dayExists := err == nil

	if !log.Start.StartedAt.IsZero() {
		if !dayExists {
			if err := SaveDayStart(date, log.Start); err != nil {
				return err
			}
		} else {
			err := queries.UpdateDayStart(ctx, db.UpdateDayStartParams{
				ShutdownTime: log.Start.ShutdownTime,
				DayGoal:      log.Start.DayGoal,
				StartedAt:    log.Start.StartedAt,
				Date:         dateOnly,
			})
			if err != nil {
				return err
			}
		}
	}

	if !log.End.FinishedAt.IsZero() {
		if dayExists {
			if err := SaveDayEnd(date, log.End); err != nil {
				return err
			}
		}
	}

	for _, pomo := range log.Pomodoros {
		if err := SavePomodoroLog(pomo.Duration, pomo.Status, pomo.Timestamp); err != nil {
			return err
		}
	}

	return nil
}

func SavePomodoroLog(duration int, status string, timestamp time.Time) error {
	queries, err := GetDB()
	if err != nil {
		return err
	}

	ctx := context.Background()
	completed := status == "completed"

	_, err = queries.InsertPomodoro(ctx, db.InsertPomodoroParams{
		StartTime:       timestamp,
		EndTime:         sql.NullTime{},
		DurationMinutes: int64(duration),
		Completed:       sql.NullBool{Bool: completed, Valid: true},
	})
	return err
}
