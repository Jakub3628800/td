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

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/td.git
   cd td
   ```

2. Build the application:
   ```bash
   go build -o td main.go
   ```

3. (Optional) Move the binary to a location in your PATH for easy access:
   ```bash
   sudo mv td /usr/local/bin/
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

### Testing GitHub Actions Locally

You can test GitHub Actions workflows locally using [act](https://github.com/nektos/act), a tool that runs GitHub Actions locally.

#### Prerequisites

1. Install act:
   ```bash
   # macOS
   brew install act
   
   # Linux
   curl -s https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash
   ```

2. Create a `.secrets` file with your GitHub token:
   ```
   GITHUB_TOKEN=your_github_token_here
   ```
   
   To create a token with the necessary permissions:
   - Go to https://github.com/settings/tokens
   - Click "Generate new token" (classic)
   - Give it a name like "TD Release Token"
   - Select the "repo" scope (Full control of private repositories)
   - Click "Generate token"
   - Copy the token and paste it in your `.secrets` file

#### Testing Options

1. **Simple local release test** (recommended for testing):
   ```bash
   ./scripts/test-release-locally.sh
   ```
   This script simulates a release without making any GitHub API calls.

2. **Test with act in dry-run mode**:
   ```bash
   ./scripts/test-github-actions.sh --dry-run
   ```
   This will show what would happen without making any actual API calls.

3. **Test with act using a simplified workflow**:
   ```bash
   ./scripts/test-github-actions.sh -w local-release-test.yml
   ```
   This uses a simplified workflow that's easier to test locally.

4. **Full workflow test** (requires valid GitHub token with proper permissions):
   ```bash
   ./scripts/test-github-actions.sh
   ```

#### Troubleshooting

If you encounter a 403 Forbidden error when creating releases:
1. Make sure your GitHub token has the "repo" scope
2. Try using the simplified test workflow instead
3. Consider creating the release directly on GitHub instead of testing locally

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgements

- [Cobra](https://github.com/spf13/cobra) for CLI interface
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) for terminal UI framework

---

Happy task managing with td! 🎉