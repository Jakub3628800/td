package core

import (
	"os/exec"
	"testing"
)

// mockExecCommand returns a command that always succeeds without doing anything
func mockExecCommand(name string, arg ...string) *exec.Cmd {
	return exec.Command("true")
}

func init() {
	// Mock execCommand globally for all tests to prevent playerctl from running
	execCommand = mockExecCommand
}

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
	// execCommand is mocked, so no actual playerctl runs
	PauseMusic()
	PlayMusic()
}

func TestExecPlayerctlHandlesMissingBinary(_ *testing.T) {
	// Override exec.Command to simulate missing playerctl
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
