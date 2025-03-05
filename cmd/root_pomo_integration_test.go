package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"td/core"

	tea "github.com/charmbracelet/bubbletea"
)

// TestRootPomoIntegration tests the integration between the root command and pomo command
// This test simulates what happens when a user presses 'p' in the root command
// and completes a pomodoro session
func TestRootPomoIntegration(t *testing.T) {
	t.Skip("Skipping integration test that requires more complex setup")

	// Set up test environment
	testVaultDir := "test/td_test_vault_integration"
	originalVaultLoc := os.Getenv("TD_VAULT_LOC")
	originalTestMode := os.Getenv("TD_TEST_MODE")
	originalIntervalMode := os.Getenv("TD_INTERVAL_MODE")

	// Set environment variables for testing
	os.Setenv("TD_VAULT_LOC", testVaultDir)
	os.Setenv("TD_TEST_MODE", "true")      // Prevent editor from opening
	os.Setenv("TD_INTERVAL_MODE", "daily") // Use daily mode for easier testing

	// Clean up after test
	defer func() {
		os.Setenv("TD_VAULT_LOC", originalVaultLoc)
		os.Setenv("TD_TEST_MODE", originalTestMode)
		os.Setenv("TD_INTERVAL_MODE", originalIntervalMode)
		os.RemoveAll(testVaultDir)
	}()

	// Create test directory if it doesn't exist
	if err := os.MkdirAll(testVaultDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create a test template file
	templatePath := filepath.Join(testVaultDir, ".template")
	if err := os.WriteFile(templatePath, []byte("# Template File\n- [ ] Task 1\n- [ ] Task 2\n"), 0644); err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	// Simulate what happens when a pomodoro session completes
	testDate := time.Now()

	// Create the directory structure for the current date
	year, month, day := testDate.Date()
	dateDir := filepath.Join(testVaultDir, strconv.Itoa(year), month.String())
	if err := os.MkdirAll(dateDir, 0755); err != nil {
		t.Fatalf("Failed to create date directory: %v", err)
	}

	// Create a file for the current date
	dateFilename := filepath.Join(dateDir, fmt.Sprintf("%02d.md", day))
	if err := os.WriteFile(dateFilename, []byte(fmt.Sprintf("# %s\n", testDate.Format("2006-01-02"))), 0644); err != nil {
		t.Fatalf("Failed to create date file: %v", err)
	}

	// Simulate adding a pomodoro completion record
	err := core.AddTask(testDate, "Completed pomodoro session")
	if err != nil {
		t.Fatalf("Failed to add pomodoro task: %v", err)
	}

	// Check if the record was added correctly
	tasks, err := core.LoadLinesWithSelection(testDate)
	if err != nil {
		t.Fatalf("Failed to load tasks: %v", err)
	}

	// Verify that the pomodoro record exists
	found := false
	for _, task := range tasks {
		if strings.Contains(task.Line, "Completed pomodoro session") {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Pomodoro completion record not found in tasks")
	}

	// Read the file content to verify the record was written
	content, err := os.ReadFile(dateFilename)
	if err != nil {
		t.Fatalf("Failed to read date file: %v", err)
	}

	if !strings.Contains(string(content), "Completed pomodoro session") {
		t.Errorf("Pomodoro completion record not found in file content")
	}
}

// TestRootCommandPomoKeyHandler tests that the 'p' key handler is properly set up
// in the root command's Update method
func TestRootCommandPomoKeyHandler(t *testing.T) {
	// Create a model instance
	m := initialModel()

	// Create a key message for 'p'
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}}

	// Call the Update method with the 'p' key message
	_, cmd := m.Update(msg)

	// Verify that the command is not nil (indicating that a command was returned)
	if cmd == nil {
		t.Errorf("Expected a command to be returned when pressing 'p', got nil")
	}
}
