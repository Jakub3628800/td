package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/Jakub3628800/td/internal/core"
)

var duration int

var pomoCmd = &cobra.Command{
	Use:   "pomo",
	Short: "Start a Pomodoro timer",
	Long:  `Start a Pomodoro timer for focused work sessions. Default duration is 25 minutes.`,
	Run: func(_ *cobra.Command, _ []string) {
		// Validate duration
		if duration <= 0 {
			fmt.Println("Error: Duration must be greater than 0 minutes")
			os.Exit(1)
		}

		m := initialPomoModel()
		p := tea.NewProgram(m)

		// Handle interrupts gracefully
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Program panicked:", r)
			}
		}()

		if _, err := p.Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(pomoCmd)
	pomoCmd.Flags().IntVarP(&duration, "duration", "d", 25, "Duration in minutes")
}

const (
	padding  = 2
	maxWidth = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

type tickMsg time.Time
type doneMsg struct{}

type pomoModel struct {
	progress  progress.Model
	start     time.Time
	duration  time.Duration
	elapsed   time.Duration
	isPaused  bool
	pauseTime time.Time
	done      bool
}

func initialPomoModel() pomoModel {
	return pomoModel{
		progress: progress.New(
			progress.WithoutPercentage(),
			progress.WithDefaultGradient(),
		),
		duration: time.Duration(duration) * time.Minute,
		start:    time.Now(),
		isPaused: false,
		done:     false,
	}
}

func (m pomoModel) Init() tea.Cmd {
	return tickCmd()
}

func (m pomoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			// Record the session even when cancelled
			recordPomoSession(duration, "cancelled")
			return m, tea.Quit
		case "p", " ":
			if m.isPaused {
				m.elapsed += time.Since(m.pauseTime)
				m.isPaused = false
				return m, tickCmd()
			}
			m.isPaused = true
			m.pauseTime = time.Now()
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:
		if m.isPaused {
			return m, nil
		}

		elapsed := time.Since(m.start) - m.elapsed
		if elapsed >= m.duration {
			if !m.done {
				m.done = true

				// Ensure completion is at 100%
				progressCmd := m.progress.SetPercent(1.0)

				// Separate goroutine for notification so it doesn't block UI
				go func() {
					core.SendNotification(fmt.Sprintf("pomo session %dm done", duration), false)
					core.PauseMusic()
					recordPomoSession(duration, "completed")
				}()

				// Add a longer delay before quitting to allow the bar to render as full
				return m, tea.Sequence(
					progressCmd,
					tea.Tick(time.Second, func(_ time.Time) tea.Msg {
						return doneMsg{}
					}),
				)
			}
			return m, tea.Quit
		}

		percentage := float64(elapsed) / float64(m.duration)
		progressCmd := m.progress.SetPercent(percentage)
		return m, tea.Batch(tickCmd(), progressCmd)

	case doneMsg:
		return m, tea.Quit

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}

	return m, nil
}

func (m pomoModel) View() string {
	var elapsed time.Duration
	if m.isPaused {
		elapsed = time.Since(m.start) - m.elapsed - time.Since(m.pauseTime)
	} else {
		elapsed = time.Since(m.start) - m.elapsed
	}

	remaining := m.duration - elapsed
	if remaining < 0 {
		remaining = 0
	}

	minutes := int(remaining.Minutes())
	seconds := int(remaining.Seconds()) % 60

	pad := strings.Repeat(" ", padding)
	status := ""
	if m.isPaused {
		status = "(Paused)"
	} else if m.done {
		status = "(Completed!)"
	}

	return "\n" +
		pad + fmt.Sprintf("%02d:%02d %s", minutes, seconds, status) + "\n" +
		pad + m.progress.View() + "\n\n" +
		pad + helpStyle("Press 'p' or space to pause/resume") + "\n" +
		pad + helpStyle("Press 'q' to quit")
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Helper function to record pomodoro sessions
func recordPomoSession(durationMinutes int, status string) {
	now := time.Now()
	_ = core.SavePomodoroLog(durationMinutes, status, now)
}
