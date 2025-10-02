# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with
code in this repository.

## Common Development Commands

### Building and Running

- `make build` - Run tests then build the binary to `bin/td`
- `make run` - Run the application directly with `go run main.go`
- `make install` - Build and install to `~/.local/bin/td`
- `make clean` - Remove built binary

### Testing and Quality

- `make test` or `go test -v ./...` - Run full test suite
- `make lint` - Run golangci-lint (auto-installs if needed)
- `make all` - Run lint, test, and build in sequence

### Single Test Execution

- `go test -v ./cmd -run TestSpecificFunction` - Run specific test
- `go test -v ./core` - Run tests for core package only

## Architecture Overview

### Project Structure

- `main.go` - Entry point with version handling and CLI argument rewriting
- `cmd/` - Cobra CLI command definitions (root, day, pomo commands)
- `internal/core/` - Core business logic:
  - `database.go` - SQLite database connection management
  - `day.go` - Day logging system (SQLite-based Collins Score tracking)
  - `session.go` - Notification and media control utilities
- `internal/db/` - Generated sqlc database queries and models
- `sqlc/` - SQL schema and query definitions

### Key Architectural Patterns

#### Data Storage (SQLite)

- SQLite database with sqlc for type-safe queries
- Database location: `{vault}/td.db` (where vault defaults to `.td`)
- Tables: `days`, `pomodori`, `metadata`
- Environment variable: `TD_VAULT_LOC` controls database location

#### Day Logging (day.go)

- SQLite-based tracking for Collins Score methodology
- Stores daily start/end times, goals, ratings, and pomodoro sessions
- UPSERT operations for day start/end data

#### CLI Design

- Default command is `day` (shows daily view)
- Argument rewriting: `--start`/`--end` → `day --start`/`day --end`
- Version handling in main.go before Cobra execution

### Environment Configuration

Key environment variables that affect behavior:

- `TD_VAULT_LOC` - Database storage location (default: `.td`)
- `TD_TEST_MODE` - Disable TUI interactions during tests

### Dependencies

- Cobra for CLI framework
- Bubble Tea for TUI components
- Lipgloss for terminal styling
- SQLite for data storage
- sqlc for type-safe database queries
