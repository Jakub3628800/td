# �� td (To-Do ToDay)

![td logo](td-logo.svg)

To-Do ToDay is a simple, efficient Text User Interface (TUI) app for tracking tasks
with a focus on daily workflow. Seamlessly add and check off tasks while the backend
stores your progress in easy-to-read markdown files.

## 🌟 Features

- 📝 Quick task addition and management
- ✅ Simple checkbox-style task completion
- 📁 Markdown file storage for easy version control and portability
- 📆 Daily, weekly, and monthly view options
- 🖥️ Clean and intuitive TUI for distraction-free productivity

## 🚀 Getting Started

### Prerequisites

- Go 1.16 or higher

### Installation

#### From Releases

You can download the latest release from the
[GitHub Releases page](https://github.com/yourusername/td/releases).

#### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/td.git
cd td

# Build and install
make install
```

## 🎯 Usage

To see all available commands:

```bash
td --help
```

### Common Commands

- Add a task:

  ```bash
  td add "Complete project proposal"
  ```

- List tasks:

  ```bash
  td list
  ```

- Start a Pomodoro session:

  ```bash
  td pomo
  ```

## 🛠️ Development

### Run Locally

To run the application without building:

```bash
go run main.go
```

### Testing

Run the test suite:

```bash
go test -v ./...
```

### Version Management

The project uses semantic versioning. To bump the version, use the provided script:

```bash
# Bump the patch version (0.1.0 -> 0.1.1)
./scripts/bump-version.sh patch

# Bump the minor version (0.1.0 -> 0.2.0)
./scripts/bump-version.sh minor

# Bump the major version (0.1.0 -> 1.0.0)
./scripts/bump-version.sh major
```

### Release Process

The project has an automated release process:

1. When changes are pushed to the `master` branch, a new release is automatically
   created.
2. The release version is determined by the version in `main.go`.
3. To create a test release without pushing to master, you can manually trigger
   the "Test Release" workflow
   from the GitHub Actions tab
   on GitHub.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file
for details.

## 🙏 Acknowledgements

- [Cobra](https://github.com/spf13/cobra) for CLI interface
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) for terminal UI framework

---

Happy task managing with td! 🎉

The Collins Score
This technique was developed by Jim Collins, the author of "Good to Great," as a
method to track and improve daily well-being and productivity. [37:06, 38:03]

How it Works (Original Version):

At the end of each day, rate your day on a scale of -2 to +2: [38:40]
+2: Excellent day
+1: Good day
0: Languishing (just getting by)
-1: Bad day
-2: Very bad day
Write a sentence or two explaining why you gave that score. [46:50]
Track the amount of creative or deep work time you had. [38:20]

Modified Version (as discussed in the video for daily use):

In the morning, ask yourself:

What time will I shut down work today? [47:28]
What do I need to achieve today for it to be a "+1 day"? [47:33]
At the end of the day (EOD), ask yourself:

How would you rate your day (using the -2 to +2 scale)? [47:42]
Why would you rate it that way? [47:49]
How many hours of creative/focused time did you have? [47:56]
