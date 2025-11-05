package core

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func SendNotification(msg string, silent bool) {
	if silent {
		fmt.Println(msg)
		fmt.Println()
	}

	// Prepend tomato emoji for pomodoro notifications
	if strings.Contains(strings.ToLower(msg), "pomo") {
		msg = "🍅 " + msg
	}

	// Only use notify-send for notifications (Linux)
	cmd := exec.Command("notify-send", msg)

	// Try to send notification but don't crash if it fails
	err := cmd.Run()
	if err != nil {
		// Log the error but don't terminate
		log.Printf("Failed to send notification: %v", err)
		// Make sure the message is still visible in the terminal
		fmt.Printf("\n%s\n", msg)
	}
}

func PauseMusic() {
	execPlayerctl("pause")
}

func PlayMusic() {
	if !IsMusicControlEnabled() {
		return
	}
	execPlayerctl("play")
}

func StopMusic() {
	if !IsMusicControlEnabled() {
		return
	}
	execPlayerctl("stop")
}

func execPlayerctl(subcmd string) {
	checkCmd := exec.Command("which", "playerctl")
	if err := checkCmd.Run(); err != nil {
		// Only warn if music control is explicitly enabled
		if IsMusicControlEnabled() {
			log.Printf("Warning: music_control_enabled is set to 'true', but playerctl is not installed or not in PATH")
			fmt.Println("⚠️  Music control is enabled in config, but playerctl is not installed.")
			fmt.Println("    Install playerctl or disable music control with: td config set music_control_enabled false")
		}
		return
	}

	err := exec.Command("playerctl", subcmd).Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			// Exit code 1 usually means no player is running, which is fine
			return
		}
		log.Printf("Failed to control media player: %v", err)
	}
}
