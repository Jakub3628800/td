# 📅 td (To-Do ToDay)

<img src="td-logo.svg" alt="td logo" height="100">

To-Do ToDay is a simple, efficient Text User Interface (TUI) app for tracking tasks with a focus on daily workflow. Seamlessly add and check off tasks while the backend stores your progress in easy-to-read markdown files.

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

You can download the latest release from the [GitHub Releases page](https://github.com/yourusername/td/releases).

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

1. When changes are pushed to the `master` branch, a new release is automatically created.
2. The release version is determined by the version in `main.go`.
3. To create a test release without pushing to master, you can manually trigger the "Test Release" workflow from the GitHub Actions tab.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgements

- [Cobra](https://github.com/spf13/cobra) for CLI interface
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) for terminal UI framework

---

Happy task managing with td! 🎉