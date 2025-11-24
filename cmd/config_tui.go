package cmd

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Jakub3628800/td/internal/core"
)

// Valid pomodoro durations (15-60 minutes, incrementing by 5)
var validPomodoroDurations = []int{15, 20, 25, 30, 35, 40, 45, 50, 55, 60}

type configModel struct {
	// State
	currentOption   int // 0 = music control, 1 = pomo duration
	musicEnabled    bool
	pomoDuration    int
	selectedDuration int
	showDurationMenu bool

	// Style
	selectedStyle lipgloss.Style
	normalStyle   lipgloss.Style
	helpStyle     lipgloss.Style
}

func initialConfigModel() configModel {
	queries, _ := core.GetDB()
	ctx := context.Background()

	// Get current music control setting
	musicEnabled := false
	if value, err := core.GetConfig(ctx, queries, "music_control_enabled"); err == nil && value == "true" {
		musicEnabled = true
	}

	// Get current pomo duration
	pomoDuration := core.GetDefaultPromoDuration()
	selectedDuration := pomoDuration
	for i, d := range validPomodoroDurations {
		if d == pomoDuration {
			selectedDuration = i
			break
		}
	}

	return configModel{
		currentOption:    0,
		musicEnabled:     musicEnabled,
		pomoDuration:     pomoDuration,
		selectedDuration: selectedDuration,
		showDurationMenu: false,
		selectedStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00")).Bold(true),
		normalStyle:      lipgloss.NewStyle(),
		helpStyle:        lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")),
	}
}

func (m configModel) Init() tea.Cmd {
	return nil
}

func (m configModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "up", "k":
			if m.showDurationMenu {
				if m.selectedDuration > 0 {
					m.selectedDuration--
				}
			} else {
				if m.currentOption > 0 {
					m.currentOption--
				}
			}
			return m, nil

		case "down", "j":
			if m.showDurationMenu {
				if m.selectedDuration < len(validPomodoroDurations)-1 {
					m.selectedDuration++
				}
			} else {
				if m.currentOption < 1 {
					m.currentOption++
				}
			}
			return m, nil

		case "enter", " ":
			if m.showDurationMenu {
				// Save the selected duration
				m.pomoDuration = validPomodoroDurations[m.selectedDuration]
				saveConfig("default_pomo_duration", fmt.Sprintf("%d", m.pomoDuration))
				m.showDurationMenu = false
			} else if m.currentOption == 0 {
				// Toggle music control
				m.musicEnabled = !m.musicEnabled
				value := "false"
				if m.musicEnabled {
					value = "true"
				}
				saveConfig("music_control_enabled", value)
			} else if m.currentOption == 1 {
				// Enter duration selection menu
				m.showDurationMenu = true
			}
			return m, nil

		case "esc":
			if m.showDurationMenu {
				m.showDurationMenu = false
				m.selectedDuration = -1
				for i, d := range validPomodoroDurations {
					if d == m.pomoDuration {
						m.selectedDuration = i
						break
					}
				}
			}
			return m, nil
		}
	}
	return m, nil
}

func (m configModel) View() string {
	if m.showDurationMenu {
		return m.renderDurationMenu()
	}
	return m.renderMainMenu()
}

func (m configModel) renderMainMenu() string {
	output := "\n Configuration Menu\n\n"

	// Music Control option
	musicStatus := "OFF"
	if m.musicEnabled {
		musicStatus = "ON"
	}
	musicLine := fmt.Sprintf("Music Control: %s", musicStatus)
	if m.currentOption == 0 {
		output += m.selectedStyle.Render("▶ " + musicLine) + "\n"
	} else {
		output += m.normalStyle.Render("  " + musicLine) + "\n"
	}

	// Pomo Duration option
	durLine := fmt.Sprintf("Default Pomodoro Duration: %d minutes", m.pomoDuration)
	if m.currentOption == 1 {
		output += m.selectedStyle.Render("▶ " + durLine) + "\n"
	} else {
		output += m.normalStyle.Render("  " + durLine) + "\n"
	}

	output += "\n" + m.helpStyle.Render("↑/k: up  ↓/j: down  Enter: select  q: quit") + "\n"

	return output
}

func (m configModel) renderDurationMenu() string {
	output := "\n Select Default Pomodoro Duration\n\n"

	for i, duration := range validPomodoroDurations {
		durationStr := fmt.Sprintf("%d minutes", duration)
		if i == m.selectedDuration {
			output += m.selectedStyle.Render("▶ " + durationStr) + "\n"
		} else {
			output += m.normalStyle.Render("  " + durationStr) + "\n"
		}
	}

	output += "\n" + m.helpStyle.Render("↑/k: up  ↓/j: down  Enter: select  Esc: back  q: quit") + "\n"

	return output
}

func saveConfig(key, value string) {
	queries, err := core.GetDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
		return
	}

	ctx := context.Background()
	if err := core.SetConfig(ctx, queries, key, value); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
		return
	}
}

func runConfigTUI() {
	if os.Getenv("TD_TEST_MODE") == "true" {
		fmt.Println("Config TUI skipped in test mode")
		return
	}

	m := initialConfigModel()
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running config TUI: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✓ Configuration saved!")
}

