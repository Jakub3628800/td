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
	// Check if playerctl exists first
	checkCmd := exec.Command("which", "playerctl")
	if err := checkCmd.Run(); err != nil {
		// playerctl is not installed or not in PATH, just log and return
		log.Printf("playerctl not found, media control unavailable")
		return
	}

	err := exec.Command("playerctl", subcmd).Run()
	if err != nil {
		// Suppress error output if exit status is 1 (no player running)
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return
		}
		// Log other errors
		log.Printf("Failed to control media player: %v", err)
	}
}
