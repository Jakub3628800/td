package core

import (
	"database/sql"
	"fmt"
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
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		dbPath = filepath.Join(homeDir, ".local", "share", "td", "td.db")
	}

	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory '%s': %w\n\nPlease ensure you have write permissions to this location or set TD_DB_PATH to a different location", dbDir, err)
	}

	database, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database at '%s': %w", dbPath, err)
	}

	globalDB = database
	globalQueries = db.New(database)

	if err := initializeDB(); err != nil {
		return nil, fmt.Errorf("failed to initialize database schema: %w\n\nThe database file may be corrupted. Try removing '%s' and running the command again", err, dbPath)
	}

	return globalQueries, nil
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

INSERT OR IGNORE INTO metadata (id, pomo_active) VALUES (1, FALSE);
`

	_, err := globalDB.Exec(schema)
	return err
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
