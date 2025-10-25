# td (To-Do ToDay)

![td logo](td-logo.svg)

To-Do ToDay is a simple, efficient tool for tracking pomodoro sessions with
configuration management. Designed for focused work sessions with session
recording and SQLite database storage.

## Features

- ⏱️ Pomodoro timer for focused work sessions
- 📊 Session recording and history tracking
- 🏷️ Tag-based session categorization
- ⚙️ Configuration management
- 💾 SQLite database storage for reliable data persistence
- 🔧 Easy-to-use CLI interface

## Getting Started

### Installation

You need [Go](https://golang.org/dl/) installed on your system. Then install td with:

```bash
go install github.com/Jakub3628800/td@latest
```

That's it! The `td` command should now be available.

**Troubleshooting**: If `td` command is not found, make sure your Go bin
directory is in your PATH:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

## Usage

To see all available commands:

```bash
td --help
```

### Common Commands

- Start a Pomodoro session (default 25 minutes):

  ```bash
  td pomo
  ```

- Start a custom duration Pomodoro:

  ```bash
  td pomo -d 10  # 10 minute session
  ```

- Add tags to a pomodoro:

  ```bash
  td pomo -t feature-dev -t refactor
  ```

- List all pomodoro sessions:

  ```bash
  td list-pomos
  ```

- View configuration options:

  ```bash
  td config
  ```

## Development

### Prerequisites

- Go 1.24 or higher
- SQLite3
- Make (for build commands)

### Local Development

Clone the repository and set up local development:

```bash
git clone https://github.com/Jakub3628800/td.git
cd td
make build      # Build the binary to bin/td
make test       # Run test suite
make lint       # Run linters
make all        # Run lint, test, and build
```

The `.envrc` file will automatically add `./bin` to your PATH when you enter
the directory (requires [direnv](https://direnv.net/)).

### Publishing a New Version

Follow these steps to publish a new release:

1. **Update the version** in `main.go`:
   ```go
   Version = "X.Y.Z"  // Update this constant
   ```

2. **Commit the version bump**:
   ```bash
   git add main.go
   git commit -m "Bump version to X.Y.Z"
   git push origin master
   ```

3. **Create and push a Git tag**:
   ```bash
   git tag -a vX.Y.Z -m "Release vX.Y.Z: <description of changes>"
   git push origin vX.Y.Z
   ```

4. **Verify the release** (wait a few seconds for Go proxy to cache):
   ```bash
   go install github.com/Jakub3628800/td@vX.Y.Z
   td --version
   ```

Users can then install the new version with:
```bash
go install github.com/Jakub3628800/td@vX.Y.Z
go install github.com/Jakub3628800/td@latest  # Always gets newest
```

**Important Notes:**
- Semantic versioning (MAJOR.MINOR.PATCH)
- Create annotated tags (`git tag -a`, not `-l`)
- Version in `main.go` must match the Git tag
- The Go module proxy caches releases; new versions may take a few seconds to appear

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file
for details.

## Acknowledgements

- [Cobra](https://github.com/spf13/cobra) for CLI interface
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) for terminal UI framework
- [SQLc](https://sqlc.dev/) for type-safe database queries
- [Lipgloss](https://github.com/charmbracelet/lipgloss) for terminal styling

---

Happy productivity tracking with td! 🍅
