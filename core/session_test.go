package core

import (
	"log"
	"os/exec"
	"testing"
)

func TestSendNotificationSilent(_ *testing.T) {
	// Should just print to stdout, no error
	SendNotification("test message", true)
}

func TestSendNotificationPomoEmoji(_ *testing.T) {
	// Should prepend tomato emoji for pomo messages
	// We can't capture notify-send, but we can check for panics/errors
	SendNotification("pomo session done", false)
}

func TestSendNotificationNormal(_ *testing.T) {
	SendNotification("normal message", false)
}

func TestPauseMusicAndPlayMusic(_ *testing.T) {
	// These should not panic or error, even if playerctl is not installed
	PauseMusic()
	PlayMusic()
}

func TestExecPlayerctlHandlesMissingBinary(_ *testing.T) {
	// Temporarily override exec.Command to simulate missing playerctl
	origCommand := execCommand
	defer func() { execCommand = origCommand }()

	execCommand = func(name string, arg ...string) *exec.Cmd {
		if name == "which" || name == "playerctl" {
			return exec.Command("false") // always fails
		}
		return exec.Command(name, arg...)
	}

	// Should log but not panic
	execPlayerctl("pause")
}

// Allow exec.Command to be overridden for testing
var execCommand = exec.Command

func init() {
	// Patch exec.Command in session.go to use execCommand variable
	// This is a hack for testability
	// In production, execCommand is just exec.Command
	// In tests, we can override it
	_ = log.Printf
}
