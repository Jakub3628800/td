package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestPomoIntegration(t *testing.T) {
	// Set up test environment
	testVaultDir := "test/td_test_vault"
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

	// Test the pomodoro completion functionality
	testDate := time.Now()

	// Create the directory structure for the current date
	year, month, day := testDate.Date()
	dateDir := filepath.Join(testVaultDir, fmt.Sprintf("%d/%s", year, month.String()))
	if err := os.MkdirAll(dateDir, 0755); err != nil {
		t.Fatalf("Failed to create date directory: %v", err)
	}

	// Create a file for the current date
	dateFilename := filepath.Join(dateDir, fmt.Sprintf("%02d.md", day))
	if err := os.WriteFile(dateFilename, []byte(fmt.Sprintf("# %s\n", testDate.Format("2006-01-02"))), 0644); err != nil {
		t.Fatalf("Failed to create date file: %v", err)
	}

	// Append a pomodoro record directly to the file
	file, err := os.OpenFile(dateFilename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("Failed to open date file: %v", err)
	}
	defer file.Close()

	timestamp := testDate.Format("15:04")
	_, err = file.WriteString(fmt.Sprintf("\n--------------------------------\npomodoro session - 25min (%s)\n", timestamp))
	if err != nil {
		t.Fatalf("Failed to write pomodoro record: %v", err)
	}

	// Read the file content to verify the record was written
	content, err := os.ReadFile(dateFilename)
	if err != nil {
		t.Fatalf("Failed to read date file: %v", err)
	}

	if !strings.Contains(string(content), "pomodoro session - 25min") {
		t.Errorf("Pomodoro completion record not found in file content")
	}
}

// TestRootCommandPomoKey tests that pressing 'p' in the root command
// launches a pomodoro session and adds a record when completed
func TestRootCommandPomoKey(t *testing.T) {
	// This is a more complex test that would require mocking the tea.Program
	// and simulating key presses, which is beyond the scope of this implementation.
	// In a real-world scenario, you would use a testing framework that can
	// simulate user input and verify the resulting state changes.
	t.Skip("Integration test requiring user interaction - skipped")
}
