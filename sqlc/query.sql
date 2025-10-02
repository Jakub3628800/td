-- name: GetMetadata :one
SELECT * FROM metadata WHERE id = 1;

-- name: UpdatePomoActive :exec
UPDATE metadata SET pomo_active = ? WHERE id = 1;

-- name: InsertPomodoro :one
INSERT INTO pomodori (start_time, end_time, duration_minutes, completed)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: UpdatePomodoroEndTime :exec
UPDATE pomodori SET end_time = ?, completed = ? WHERE id = ?;

-- name: GetActivePomodoro :one
SELECT * FROM pomodori WHERE end_time IS NULL ORDER BY start_time DESC LIMIT 1;

-- name: ListPomodori :many
SELECT * FROM pomodori ORDER BY start_time DESC;

-- name: CreateDay :one
INSERT INTO days (date, shutdown_time, day_goal, started_at)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: UpdateDayEnd :exec
UPDATE days SET rating = ?, reason = ?, focus_hours = ?, finished_at = ? WHERE date = ?;

-- name: UpdateDayStart :exec
UPDATE days SET shutdown_time = ?, day_goal = ?, started_at = ? WHERE date = ?;

-- name: GetDayByDate :one
SELECT * FROM days WHERE date = ?;

-- name: ListDays :many
SELECT * FROM days ORDER BY date DESC;
