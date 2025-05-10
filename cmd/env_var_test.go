package cmd

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Jakub3628800/td/core"
)

// TestTD_VAULT_LOC_Respect tests that the TD_VAULT_LOC environment variable
// is respected in all commands that interact with the filesystem
func TestTD_VAULT_LOC_Respect(t *testing.T) {
	// Skip this test as it requires user interaction
	t.Skip("Skipping integration test that requires user interaction")

	// Set up test environment with a unique temporary directory
	testVaultDir := filepath.Join(os.TempDir(), "td_vault_loc_test_"+strconv.FormatInt(time.Now().UnixNano(), 10))
	originalVaultLoc := os.Getenv("TD_VAULT_LOC")
	originalTestMode := os.Getenv("TD_TEST_MODE")

	// Clean up after test
	defer func() {
		os.Setenv("TD_VAULT_LOC", originalVaultLoc)
		os.Setenv("TD_TEST_MODE", originalTestMode)
		os.RemoveAll(testVaultDir)

		// Reset core package variables
		core.ResetForTest(originalVaultLoc, "weekly")
	}()

	// Set environment variables for testing
	os.Setenv("TD_VAULT_LOC", testVaultDir)
	os.Setenv("TD_TEST_MODE", "true") // Prevent editor from opening

	// Reset core package variables
	core.ResetForTest(testVaultDir, "weekly")

	// Create test directory if it doesn't exist
	if err := os.MkdirAll(testVaultDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Test cases for various commands
	testCases := []struct {
		name       string
		cmdSetup   func() error
		verifyFunc func() error
	}{
		{
			name: "Add command respects TD_VAULT_LOC",
			cmdSetup: func() error {
				captureCmd.SetArgs([]string{"Test task for add command"})
				return captureCmd.Execute()
			},
			verifyFunc: func() error {
				tasks, err := core.LoadLinesWithSelection(time.Now())
				if err != nil {
					return err
				}

				found := false
				for _, task := range tasks {
					if strings.Contains(task.Line, "Test task for add command") {
						found = true
						break
					}
				}

				if !found {
					t.Error("Task not found in the vault location specified by TD_VAULT_LOC")
				}
				return nil
			},
		},
		{
			name: "Pomodoro command respects TD_VAULT_LOC",
			cmdSetup: func() error {
				recordPomoSession(15, "completed")
				return nil
			},
			verifyFunc: func() error {
				// Find all files in the test vault
				var files []string
				err := filepath.Walk(testVaultDir, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if !info.IsDir() {
						files = append(files, path)
					}
					return nil
				})

				if err != nil {
					return err
				}

				if len(files) == 0 {
					t.Error("No files found in test vault directory")
					return nil
				}

				// Check if pomodoro session was recorded
				foundPomo := false
				for _, file := range files {
					content, err := os.ReadFile(file)
					if err != nil {
						t.Errorf("Failed to read file %s: %v", file, err)
						continue
					}

					if strings.Contains(string(content), "pomodoro session - 15min") {
						foundPomo = true
						break
					}
				}

				if !foundPomo {
					t.Error("Pomodoro record not found in the vault location specified by TD_VAULT_LOC")
				}
				return nil
			},
		},
	}

	// Run all test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Execute the command setup
			if err := tc.cmdSetup(); err != nil {
				t.Fatalf("Failed to execute command setup: %v", err)
			}

			// Verify the result
			if err := tc.verifyFunc(); err != nil {
				t.Fatalf("Verification failed: %v", err)
			}
		})
	}
}
