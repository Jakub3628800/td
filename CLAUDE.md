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
- `cmd/` - Cobra CLI command definitions (root, day, add, pomo commands)
- `core/` - Core business logic split into three main components:
  - `vault.go` - Task storage and file management (markdown-based)
  - `day.go` - Day logging system (JSON-based Collins Score tracking)
  - `session.go` - Notification and media control utilities

### Key Architectural Patterns

#### Task Storage (vault.go)

- Tasks stored as markdown checkbox format (`- [ ]` / `- [x]`)
- File organization by date/week/month based on `TD_INTERVAL_MODE`
- Template system for new file creation
- Environment-driven configuration (TD\_\* variables)

#### Day Logging (day.go)

- Separate JSON-based tracking for Collins Score methodology
- Stores daily start/end times, goals, ratings, and pomodoro sessions
- File structure: `{vault}/{year}/{month}/{day}.json`

#### CLI Design

- Default command is `day` (shows daily view)
- Argument rewriting: `--start`/`--end` → `day --start`/`day --end`
- Version handling in main.go before Cobra execution

### Environment Configuration

Key environment variables that affect behavior:

- `TD_VAULT_LOC` - Storage location (default: `.td`)
- `TD_INTERVAL_MODE` - daily/weekly/monthly file organization
- `TD_TEMPLATE_PATH` - Template file for new entries
- `TD_SKIP_WEEKEND` - Skip weekends in daily mode
- `TD_COPY_PREVIOUS` - Copy previous period's content
- `TD_TEST_MODE` - Disable editor launching during tests

### Dependencies

- Cobra for CLI framework
- Bubble Tea for TUI components
- Lipgloss for terminal styling
- Standard library for file operations and JSON handling
