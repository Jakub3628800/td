package core

import (
	"fmt"
	"log"
	"os/exec"
)

func SendNotification(msg string, silent bool) {
	if silent {
		fmt.Println(msg)
		fmt.Println()
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
		// Don't crash if playerctl fails, just log it
		log.Printf("Failed to control media player: %v", err)
	}
}
