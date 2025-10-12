package core

import (
	"context"
	"time"
)

func HasRunningPomodoro() (bool, error) {
	queries, err := GetDB()
	if err != nil {
		return false, err
	}

	ctx := context.Background()
	pomodoros, err := queries.ListPomodori(ctx)
	if err != nil {
		return false, err
	}

	now := time.Now()
	for _, p := range pomodoros {
		expectedEndTime := p.StartTime.Add(time.Duration(p.DurationMinutes) * time.Minute)
		if now.Before(expectedEndTime) {
			return true, nil
		}
	}

	return false, nil
}
