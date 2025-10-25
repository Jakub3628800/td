package core

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	"github.com/Jakub3628800/td/internal/db"
)

var globalDB *sql.DB
var globalQueries *db.Queries

func GetDB() (*db.Queries, error) {
	if globalQueries != nil {
		return globalQueries, nil
	}

	dbPath := os.Getenv("TD_DB_PATH")
	if dbPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		dbPath = filepath.Join(homeDir, ".local", "share", "td", "td.db")
	}

	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, err
	}

	database, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	globalDB = database
	globalQueries = db.New(database)

	if err := initializeDB(); err != nil {
		return nil, err
	}

	return globalQueries, nil
}

func migrateDefaultDeviceToConfig() error {
	// Check if the old spotify_tokens table has the default_device_id column
	var columnExists bool
	err := globalDB.QueryRow(`
		SELECT COUNT(*) > 0
		FROM pragma_table_info('spotify_tokens')
		WHERE name = 'default_device_id'
	`).Scan(&columnExists)
	if err != nil {
		return err
	}

	if !columnExists {
		return nil // Column doesn't exist, nothing to migrate
	}

	// Check if there's data in default_device_id
	var deviceID sql.NullString
	err = globalDB.QueryRow(`
		SELECT default_device_id FROM spotify_tokens WHERE id = 1
	`).Scan(&deviceID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// If there's data, migrate it to config table
	if err == nil && deviceID.Valid && deviceID.String != "" {
		queries := getQueries()
		ctx := context.Background()
		err = queries.SetConfig(ctx, db.SetConfigParams{
			Key:   "spotify_default_device",
			Value: deviceID,
		})
		if err != nil {
			return err
		}
	}

	// Drop the old column (SQLite doesn't support DROP COLUMN directly, so we recreate the table)
	_, err = globalDB.Exec(`
		CREATE TABLE IF NOT EXISTS spotify_tokens_new (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			access_token TEXT NOT NULL,
			refresh_token TEXT NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		INSERT OR IGNORE INTO spotify_tokens_new
		SELECT id, access_token, refresh_token, expires_at, updated_at FROM spotify_tokens;

		DROP TABLE IF EXISTS spotify_tokens;
		ALTER TABLE spotify_tokens_new RENAME TO spotify_tokens;
	`)

	return err
}

func initializeDB() error {
	schema := `
CREATE TABLE IF NOT EXISTS metadata (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    pomo_active BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS pomodori (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP,
    duration_minutes INTEGER NOT NULL,
    completed BOOLEAN DEFAULT FALSE,
    tags TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS days (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date DATE NOT NULL UNIQUE,
    shutdown_time TIMESTAMP NOT NULL,
    day_goal TEXT NOT NULL,
    started_at TIMESTAMP NOT NULL,
    rating INTEGER,
    reason TEXT,
    focus_hours REAL,
    finished_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS spotify_tokens (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS config (
    key TEXT PRIMARY KEY,
    value TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT OR IGNORE INTO metadata (id, pomo_active) VALUES (1, FALSE);
`

	if _, err := globalDB.Exec(schema); err != nil {
		return err
	}

	// Migration: Move default_device_id from spotify_tokens to config table
	if err := migrateDefaultDeviceToConfig(); err != nil {
		return err
	}

	return nil
}

func getQueries() *db.Queries {
	if globalQueries == nil {
		queries, err := GetDB()
		if err != nil {
			panic(err)
		}
		return queries
	}
	return globalQueries
}

func CloseDB() error {
	if globalDB != nil {
		err := globalDB.Close()
		globalDB = nil
		globalQueries = nil
		return err
	}
	return nil
}
