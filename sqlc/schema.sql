CREATE TABLE metadata (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    pomo_active BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE pomodori (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP,
    duration_minutes INTEGER NOT NULL,
    completed BOOLEAN DEFAULT FALSE,
    tags TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE days (
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
