package cmd

import (
	"io"
	"os"
	"testing"
)

func TestHelpSubcommandCompatibility(t *testing.T) {
	restore := discardStdout(t)
	defer restore()

	if err := execute([]string{"help", "pomo"}); err != nil {
		t.Fatalf("execute help pomo failed: %v", err)
	}
}

func TestPomoShortFlagEqualsCompatibility(t *testing.T) {
	restore := discardStdout(t)
	defer restore()

	if err := runPomoCommand([]string{"-d=1", "-t=work", "--help"}); err != nil {
		t.Fatalf("runPomoCommand failed: %v", err)
	}
	if duration != 1 {
		t.Fatalf("duration = %d, want 1", duration)
	}
	if len(tags) != 1 || tags[0] != "work" {
		t.Fatalf("tags = %v, want [work]", tags)
	}
}

func discardStdout(t *testing.T) func() {
	t.Helper()

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe failed: %v", err)
	}
	os.Stdout = w

	return func() {
		t.Helper()

		_ = w.Close()
		os.Stdout = old
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}
}
