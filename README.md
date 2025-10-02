# td (To-Do ToDay)

![td logo](td-logo.svg)

To-Do ToDay is a simple, efficient Text User Interface (TUI) app for tracking daily productivity
with a focus on day logging and pomodoro sessions. Uses the Collins Score methodology to track
daily well-being and productivity with SQLite database storage.

## Features

- Daily productivity tracking with Collins Score methodology
- Pomodoro timer with session recording
- SQLite database storage for reliable data persistence
- Clean and intuitive TUI for distraction-free productivity
- Day ratings, goals, and focus time tracking

## Getting Started

### Installation

You need [Go](https://golang.org/dl/) installed on your system. Then install td with:

```bash
go install github.com/Jakub3628800/td@latest
```

That's it! The `td` command should now be available.

**Troubleshooting**: If `td` command is not found, make sure your Go bin directory is in your PATH:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

## Usage

To see all available commands:

```bash
td --help
```

### Common Commands

- Start your day (set shutdown time and daily goal):

  ```bash
  td day --start
  ```

- End your day (rate your day and record focus hours):

  ```bash
  td day --end
  ```

- Start a Pomodoro session:

  ```bash
  td pomo
  ```

- Default command shows day interface:

  ```bash
  td
  ```

## About the Collins Score

This technique was developed by Jim Collins, author of "Good to Great," as a method to track and improve daily well-being and productivity.

**How it works:**

- **Morning Setup**: Set your work shutdown time and define what would make today a "+1 day"
- **Evening Review**: Rate your day from -2 (terrible) to +2 (excellent), explain why, and log your focus hours
- **Continuous Improvement**: Track patterns over time to optimize your daily routine

td makes this process simple with guided prompts and automatic data storage.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file
for details.

## Acknowledgements

- [Cobra](https://github.com/spf13/cobra) for CLI interface
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) for terminal UI framework
- Jim Collins for the Collins Score methodology

---

Happy productivity tracking with td!
