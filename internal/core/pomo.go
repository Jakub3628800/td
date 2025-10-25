package core

import (
	"context"
	"database/sql"
)

func HasRunningPomodoro() (bool, error) {
	queries, err := GetDB()
	if err != nil {
		return false, err
	}

	ctx := context.Background()
	// Use GetActivePomodoro to check for sessions with end_time IS NULL
	// This properly identifies sessions that are still running
	_, err = queries.GetActivePomodoro(ctx)
	if err != nil {
		// sql.ErrNoRows means no active sessions
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
