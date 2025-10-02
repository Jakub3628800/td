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
	execPlayerctl("play")
}

func execPlayerctl(subcmd string) {
	checkCmd := exec.Command("which", "playerctl")
	if err := checkCmd.Run(); err != nil {
		log.Printf("playerctl not found, media control unavailable")
		return
	}

	err := exec.Command("playerctl", subcmd).Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return
		}
		log.Printf("Failed to control media player: %v", err)
	}
}
