package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"td/core"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "td",
	Short: "A simple, efficient Text User Interface (TUI) app for tracking tasks",
	Long: `To-Do ToDay (td) is a simple, efficient Text User Interface (TUI) app for tracking tasks 
with a focus on daily workflow. Seamlessly add and check off tasks while the backend 
stores your progress in easy-to-read markdown files.

Features:
- 📝 Quick task addition and management
- ✅ Simple checkbox-style task completion
- 📁 Markdown file storage for easy version control and portability
- 📆 Daily, weekly, and monthly view options
- 🖥️ Clean and intuitive TUI for distraction-free productivity`,
	Run: func(_ *cobra.Command, _ []string) {
		p := tea.NewProgram(initialModel())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type model struct {
	cursor int
	tasks  []core.Task
	date   time.Time
}

func initialModel() model {
	tasks, _ := core.LoadLinesWithSelection(time.Now())
	return model{
		tasks: tasks,
		date:  time.Now(),
	}
}

func (m model) Save() {
	// Implement save functionality if needed
}

func (m *model) Refresh() {
	tasks, _ := core.LoadLinesWithSelection(time.Now())
	m.tasks = tasks
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "left", "h":
			m.date = core.PreviousDate(m.date)
			a := &m
			a.Refresh()
		case "right", "l":
			m.date = core.NextDate(m.date)
			a := &m
			a.Refresh()
		case "down", "j":
			if m.cursor < len(m.tasks)-1 {
				m.cursor++
			}
		case "p":
			// Run pomodoro session
			cmd := exec.Command("td", "pomo")
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			return m, tea.Sequence(
				tea.Quit,
				func() tea.Msg {
					err := cmd.Run()
					if err != nil {
						return err
					}

					// The pomodoro session will record itself when completed
					// We don't need to record anything here
					return nil
				},
			)
		case "e":
			lineNumber, _ := core.ContainsLine(m.date, m.tasks[m.cursor].Line)
			if err := core.OpenEditor(m.date, lineNumber, false); err != nil {
				return m, tea.Quit
			}
			a := &m
			a.Refresh()
		case "enter", " ":
			selected := m.tasks[m.cursor].Selected
			if selected {
				m.tasks[m.cursor].Selected = false
				if err := core.UpdateTaskStatus(false, m.tasks[m.cursor].Line, m.date); err != nil {
					return m, tea.Quit
				}
			} else {
				m.tasks[m.cursor].Selected = true
				if err := core.UpdateTaskStatus(true, m.tasks[m.cursor].Line, m.date); err != nil {
					return m, tea.Quit
				}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	s := core.GetHeader(m.date)

	for i, task := range m.tasks {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if m.tasks[i].Selected {
			checked = "x"
		}

		replaced := strings.ReplaceAll(task.Line, "- [ ]", "")
		replaced = strings.ReplaceAll(replaced, "- [x]", "")
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, replaced)
	}

	// Use the existing helpStyle from pomo.go
	s += "\n" + helpStyle("Press q to quit, p to start a pomodoro session.")

	return s
}
