package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/Jakub3628800/td/internal/core"
)

// Valid pomodoro durations (15-60 minutes, incrementing by 5)
var validPomodoroDurations = []int{15, 20, 25, 30, 35, 40, 45, 50, 55, 60}

type configModel struct {
	currentOption    int // 0 = music control, 1 = pomo duration
	musicEnabled     bool
	pomoDuration     int
	selectedDuration int
	showDurationMenu bool
}

func initialConfigModel() configModel {
	queries, _ := core.GetDB()
	ctx := context.Background()

	musicEnabled := false
	if value, err := core.GetConfig(ctx, queries, "music_control_enabled"); err == nil && value == "true" {
		musicEnabled = true
	}

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
	}
}

func runConfigTUI() {
	if os.Getenv("TD_TEST_MODE") == "true" {
		fmt.Println("Config TUI skipped in test mode")
		return
	}

	state, err := enableRawMode()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running config TUI: %v\n", err)
		os.Exit(1)
	}
	defer state.restore()

	model := initialConfigModel()

	for {
		renderConfig(model)
		key, err := readConfigKey()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}

		switch key {
		case "q", "ctrl+c":
			fmt.Print("\n✓ Configuration saved!\n")
			return
		case "up", "k":
			if model.showDurationMenu {
				if model.selectedDuration > 0 {
					model.selectedDuration--
				}
			} else if model.currentOption > 0 {
				model.currentOption--
			}
		case "down", "j":
			if model.showDurationMenu {
				if model.selectedDuration < len(validPomodoroDurations)-1 {
					model.selectedDuration++
				}
			} else if model.currentOption < 1 {
				model.currentOption++
			}
		case "enter", " ":
			if model.showDurationMenu {
				model.pomoDuration = validPomodoroDurations[model.selectedDuration]
				saveConfig("default_pomo_duration", fmt.Sprintf("%d", model.pomoDuration))
				model.showDurationMenu = false
			} else if model.currentOption == 0 {
				model.musicEnabled = !model.musicEnabled
				value := "false"
				if model.musicEnabled {
					value = "true"
				}
				saveConfig("music_control_enabled", value)
			} else if model.currentOption == 1 {
				model.showDurationMenu = true
			}
		case "esc":
			if model.showDurationMenu {
				model.showDurationMenu = false
				model.selectedDuration = -1
				for i, d := range validPomodoroDurations {
					if d == model.pomoDuration {
						model.selectedDuration = i
						break
					}
				}
			}
		}
	}
}

func readConfigKey() (string, error) {
	b, err := readInputByteBlocking()
	if err != nil {
		return "", err
	}

	switch b {
	case 3:
		return "ctrl+c", nil
	case '\r', '\n':
		return "enter", nil
	case 27:
		second, ok, err := readInputByteOptional()
		if err != nil {
			return "", err
		}
		if !ok || second != '[' {
			return "esc", nil
		}
		third, ok, err := readInputByteOptional()
		if err != nil {
			return "", err
		}
		if !ok {
			return "esc", nil
		}
		switch third {
		case 'A':
			return "up", nil
		case 'B':
			return "down", nil
		}
		return "esc", nil
	case 'q':
		return "q", nil
	case 'k':
		return "k", nil
	case 'j':
		return "j", nil
	case ' ':
		return " ", nil
	}
	return string(b), nil
}

func readInputByteBlocking() (byte, error) {
	buf := make([]byte, 1)
	for {
		n, err := os.Stdin.Read(buf)
		if n > 0 {
			return buf[0], nil
		}
		if err == io.EOF {
			continue
		}
		if err != nil {
			return 0, err
		}
	}
}

func readInputByteOptional() (byte, bool, error) {
	buf := make([]byte, 1)
	n, err := os.Stdin.Read(buf)
	if n > 0 {
		return buf[0], true, nil
	}
	if err == io.EOF {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return 0, false, nil
}

func renderConfig(m configModel) {
	fmt.Print("\033[H\033[2J")
	if m.showDurationMenu {
		fmt.Print(renderDurationMenu(m))
		return
	}
	fmt.Print(renderMainMenu(m))
}

func renderMainMenu(m configModel) string {
	output := "\n Configuration Menu\n\n"

	musicStatus := "OFF"
	if m.musicEnabled {
		musicStatus = "ON"
	}
	musicLine := fmt.Sprintf("Music Control: %s", musicStatus)
	if m.currentOption == 0 {
		output += selectedText("▶ "+musicLine) + "\n"
	} else {
		output += "  " + musicLine + "\n"
	}

	durLine := fmt.Sprintf("Default Pomodoro Duration: %d minutes", m.pomoDuration)
	if m.currentOption == 1 {
		output += selectedText("▶ "+durLine) + "\n"
	} else {
		output += "  " + durLine + "\n"
	}

	output += "\n" + helpText("↑/k: up  ↓/j: down  Enter: select  q: quit") + "\n"
	return output
}

func renderDurationMenu(m configModel) string {
	output := "\n Select Default Pomodoro Duration\n\n"

	for i, duration := range validPomodoroDurations {
		durationStr := fmt.Sprintf("%d minutes", duration)
		if i == m.selectedDuration {
			output += selectedText("▶ "+durationStr) + "\n"
		} else {
			output += "  " + durationStr + "\n"
		}
	}

	output += "\n" + helpText("↑/k: up  ↓/j: down  Enter: select  Esc: back  q: quit") + "\n"
	return output
}

func selectedText(s string) string {
	return "\033[32;1m" + s + "\033[0m"
}

func helpText(s string) string {
	return "\033[90m" + s + "\033[0m"
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
