package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/Jakub3628800/td/internal/core"
)

var startDay bool
var endDay bool

var grayStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render

var dayCmd = &cobra.Command{
	Use:   "day",
	Short: "Start or end your day, or log day data interactively.",
	Long:  `Start or end your day, or log day data interactively.`,
	Run: func(_ *cobra.Command, _ []string) {
		now := time.Now()
		log, err := core.LoadDayLog(now)
		if startDay {
			p := tea.NewProgram(initialStartDayModel())
			if m, err := p.Run(); err == nil {
				if final, ok := m.(startDayModel); ok && final.confirmed {
					saveDayStart(final.shutdownTimeMinutes, final.dayGoal)
					notifyDayStart(final.shutdownTimeMinutes)
				}
			} else {
				fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
			}
			return
		}
		if endDay {
			p := tea.NewProgram(initialDayModel())
			if m, err := p.Run(); err == nil {
				if final, ok := m.(dayModel); ok && final.confirmed {
					saveDayEnd(final.rating, final.reason, final.focusHours)
				}
			} else {
				fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
			}
			return
		}
		// No --start or --end: decide what to do
		if err != nil || log.Start.StartedAt.IsZero() {
			// No log for today or not started: start the day
			p := tea.NewProgram(initialStartDayModel())
			if m, err := p.Run(); err == nil {
				if final, ok := m.(startDayModel); ok && final.confirmed {
					saveDayStart(final.shutdownTimeMinutes, final.dayGoal)
					notifyDayStart(final.shutdownTimeMinutes)
				}
			} else {
				fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
			}
			return
		}
		if log.End.FinishedAt.IsZero() {
			// Not finished: end the day
			p := tea.NewProgram(initialDayModel())
			if m, err := p.Run(); err == nil {
				if final, ok := m.(dayModel); ok && final.confirmed {
					saveDayEnd(final.rating, final.reason, final.focusHours)
				}
			} else {
				fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
			}
			return
		}
		fmt.Println("🌱 You have already finished work, go touch grass!")
		fmt.Println("To edit your response, use:")
		fmt.Println(grayStyle("td --start"))
		fmt.Println(grayStyle("td --end"))
	},
}

type dayModel struct {
	rating       int
	confirmed    bool
	askingReason bool
	askingFocus  bool
	reason       string
	input        textinput.Model
	focusHours   float64
}

func initialDayModel() dayModel {
	ti := textinput.New()
	ti.Placeholder = "Type your reason and press Enter..."
	ti.Focus()
	return dayModel{rating: 0, confirmed: false, askingReason: false, askingFocus: false, input: ti, focusHours: 0}
}

func (m dayModel) Init() tea.Cmd {
	return nil
}

func (m dayModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	msgTyped, ok := msg.(tea.KeyMsg)
	if ok {
		if m.askingReason {
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			if msgTyped.Type == tea.KeyEnter {
				m.reason = m.input.Value()
				m.askingReason = false
				m.askingFocus = true
				return m, nil
			}
			return m, cmd
		}
		if m.askingFocus {
			s := msgTyped.String()
			switch s {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				m.focusHours += 0.5
			case "down", "j":
				if m.focusHours > 0 {
					m.focusHours -= 0.5
					if m.focusHours < 0 {
						m.focusHours = 0
					}
				}
			case "enter", " ":
				m.confirmed = true
				return m, tea.Quit
			}
			return m, nil
		}
		s := msgTyped.String()
		switch s {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.rating < 2 {
				m.rating++
			}
		case "down", "j":
			if m.rating > -2 {
				m.rating--
			}
		case "enter", " ":
			m.askingReason = true
			m.input.SetValue("")
			m.input.Focus()
			return m, nil
		}
	}
	if m.askingReason {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
	return m, nil
}

var ratingColors = map[int]string{
	-2: "#e57373", // lighter red
	-1: "#e57373", // lighter red
	0:  "#1976d2", // blue
	1:  "#81c784", // green
	2:  "#1b5e20", // intense green
}

var ratingLabels = map[int]string{
	-2: "Terrible",
	-1: "Bad",
	0:  "Neutral",
	1:  "Good",
	2:  "Great",
}

func (m dayModel) View() string {
	if m.askingReason {
		return "📝 Why would you rate it that way?\n\n" + m.input.View() + "\n\n(Enter to confirm, q to quit)"
	}
	if m.askingFocus {
		color := "#1976d2"
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Bold(true)
		prompt := fmt.Sprintf("🕰️ How many hours of focus time did you have after this? [%s]", style.Render(fmt.Sprintf("%.1f", m.focusHours)))
		return prompt + "\n\nUse ↑/↓ or j/k to change, Enter to confirm, q to quit."
	}
	color := ratingColors[m.rating]
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Bold(true)
	prompt := fmt.Sprintf("⭐ How would you rate your day [%s]", style.Render(fmt.Sprintf("%d", m.rating)))

	s := prompt + "\n\nUse ↑/↓ or j/k to change, Enter to confirm, q to quit."
	if m.confirmed {
		s += fmt.Sprintf("\n\n🙌 You rated your day: %d (%s)\n📝 Reason: %s\n🕰️ Focus hours: %.1f", m.rating, ratingLabels[m.rating], m.reason, m.focusHours)
	}
	return s
}

func saveDayStart(shutdownTimeMinutes int, dayGoal string) {
	now := time.Now()
	year, month, day := now.Date()
	hour := shutdownTimeMinutes / 60
	minute := shutdownTimeMinutes % 60
	shutdownTime := time.Date(year, month, day, hour, minute, 0, 0, now.Location())
	start := core.DayStart{
		ShutdownTime: shutdownTime,
		DayGoal:      dayGoal,
		StartedAt:    now,
	}
	_ = core.SaveDayStart(now, start)
}

func saveDayEnd(rating int, reason string, focusHours float64) {
	now := time.Now()
	end := core.DayEnd{
		Rating:     rating,
		Reason:     reason,
		FocusHours: focusHours,
		FinishedAt: now,
	}
	_ = core.SaveDayEnd(now, end)
}

func init() {
	dayCmd.Flags().BoolVar(&startDay, "start", false, "Start the day and answer planning questions (override)")
	dayCmd.Flags().BoolVar(&endDay, "end", false, "End the day and answer end questions (override)")
	rootCmd.AddCommand(dayCmd)
}

type startDayModel struct {
	step                int
	confirmed           bool
	shutdownTimeMinutes int // minutes since midnight
	dayGoal             string
	input               textinput.Model
}

func initialStartDayModel() startDayModel {
	ti := textinput.New()
	ti.Placeholder = "Describe your day goal..."
	return startDayModel{step: 0, confirmed: false, shutdownTimeMinutes: 17 * 60, input: ti}
}

func (m startDayModel) Init() tea.Cmd {
	return nil
}

func (m startDayModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.step == 0 {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				if m.shutdownTimeMinutes < 23*60+45 {
					m.shutdownTimeMinutes += 15
				}
			case "down", "j":
				if m.shutdownTimeMinutes > 0 {
					m.shutdownTimeMinutes -= 15
				}
			case "enter", " ":
				m.step = 1
				m.input.SetValue("")
				m.input.Focus()
				return m, nil
			}
			return m, nil
		}
		if m.step == 1 {
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			if msg.Type == tea.KeyEnter {
				m.dayGoal = m.input.Value()
				m.confirmed = true
				return m, tea.Quit
			}
			if msg.String() == "ctrl+c" || msg.String() == "q" {
				return m, tea.Quit
			}
			return m, cmd
		}
	}
	return m, nil
}

func (m startDayModel) View() string {
	if m.step == 0 {
		hour := m.shutdownTimeMinutes / 60
		minute := m.shutdownTimeMinutes % 60
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("#1976d2")).Bold(true)
		return fmt.Sprintf("🗓️ What time will I shut down work today?\n\n[%s]\n\nUse ↑/↓ or j/k to change, Enter to confirm, q to quit.", style.Render(fmt.Sprintf("%02d:%02d", hour, minute)))
	}
	if m.step == 1 {
		return "📝 What do I need to achieve today for it to be a \"+1 day\"?\n\n" + m.input.View() + "\n\n(Enter to confirm, q to quit)"
	}
	if m.confirmed {
		hour := m.shutdownTimeMinutes / 60
		minute := m.shutdownTimeMinutes % 60
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("#1976d2")).Bold(true)
		return fmt.Sprintf("Shutdown time: %s\nDay goal: %s", style.Render(fmt.Sprintf("%02d:%02d", hour, minute)), m.dayGoal)
	}
	return ""
}

func notifyDayStart(shutdownTimeMinutes int) {
	now := time.Now()
	year, month, day := now.Date()
	hour := shutdownTimeMinutes / 60
	minute := shutdownTimeMinutes % 60
	shutdownTime := time.Date(year, month, day, hour, minute, 0, 0, now.Location())

	// Notification 1: 1 hour before shutdown
	oneHourBefore := shutdownTime.Add(-1 * time.Hour)
	scheduleNotification(oneHourBefore, "Your desired time for finishing work is 1hr from now.")

	// Notification 2: at shutdown time
	scheduleNotification(shutdownTime, "You said you'd ideally finish work now.")
}

func scheduleNotification(t time.Time, content string) {
	datetime := t.Format("2006-01-02 15:04")
	// #nosec G204 - datetime is controlled format and content is from internal string constants
	cmd := exec.Command("systemd-run", "--user", "--on-calendar="+datetime, "notify-send", content)
	_ = cmd.Run() // ignore errors for now
}
