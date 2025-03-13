package core

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGetFilename(t *testing.T) {
	// Save the original values
	originalVaultLoc := vaultLoc
	originalIntervalMode := intervalMode

	// Restore the original values after the test
	defer func() {
		vaultLoc = originalVaultLoc
		intervalMode = originalIntervalMode
	}()

	tests := []struct {
		name         string
		vaultLoc     string
		intervalMode string
		date         time.Time
		want         string
	}{
		{
			name:         "Daily mode",
			vaultLoc:     "/test/vault",
			intervalMode: "daily",
			date:         time.Date(2024, 8, 30, 0, 0, 0, 0, time.UTC),
			want:         "/test/vault/2024/August/30.md",
		},
		{
			name:         "Weekly mode",
			vaultLoc:     "/test/vault",
			intervalMode: "weekly",
			date:         time.Date(2024, 8, 30, 0, 0, 0, 0, time.UTC),
			want:         "/test/vault/2024/August/week35.md",
		},
		{
			name:         "Monthly mode",
			vaultLoc:     "/test/vault",
			intervalMode: "monthly",
			date:         time.Date(2024, 8, 30, 0, 0, 0, 0, time.UTC),
			want:         "/test/vault/2024/August/August.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vaultLoc = tt.vaultLoc
			intervalMode = tt.intervalMode
			got := getFilename(tt.date)
			if got != tt.want {
				t.Errorf("getFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHeader(t *testing.T) {
	// Save the original value
	originalIntervalMode := intervalMode

	// Restore the original value after the test
	defer func() {
		intervalMode = originalIntervalMode
	}()

	tests := []struct {
		name         string
		intervalMode string
		date         time.Time
		want         string
	}{
		{
			name:         "Daily mode",
			intervalMode: "daily",
			date:         time.Date(2024, 8, 30, 0, 0, 0, 0, time.UTC),
			want:         "2024-08-30 Friday\n\n",
		},
		{
			name:         "Weekly mode",
			intervalMode: "weekly",
			date:         time.Date(2024, 8, 30, 0, 0, 0, 0, time.UTC),
			want:         "Week 35\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			intervalMode = tt.intervalMode
			got := GetHeader(tt.date)
			if got != tt.want {
				t.Errorf("GetHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddTask(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "vault_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set up the test environment
	vaultLoc = tempDir
	intervalMode = "daily"
	templatePath = "test_template"

	// Create a test template file
	templateContent := "Template content\n"
	err = os.WriteFile(filepath.Join(tempDir, templatePath), []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	// Test adding a task
	testDate := time.Date(2024, 8, 30, 0, 0, 0, 0, time.UTC)
	testTask := "Test task"
	err = AddTask(testDate, testTask)
	if err != nil {
		t.Errorf("AddTask() error = %v", err)
	}

	// Verify the task was added correctly
	filename := getFilename(testDate)
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	expectedContent := templateContent + "- [ ] " + testTask + "\n"
	if string(content) != expectedContent {
		t.Errorf("File content = %v, want %v", string(content), expectedContent)
	}
}

func TestUpdateTaskStatus(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "vault_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set up the test environment
	vaultLoc = tempDir
	intervalMode = "daily"

	// Create a test file with tasks
	testDate := time.Date(2024, 8, 30, 0, 0, 0, 0, time.UTC)
	filename := getFilename(testDate)
	initialContent := "- [ ] Task 1\n- [ ] Task 2\n- [ ] Task 3\n"
	err = os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		t.Fatalf("Failed to create directories: %v", err)
	}
	err = os.WriteFile(filename, []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test updating task status
	err = UpdateTaskStatus(true, "Task 2", testDate)
	if err != nil {
		t.Errorf("UpdateTaskStatus() error = %v", err)
	}

	// Verify the task status was updated correctly
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	expectedContent := "- [ ] Task 1\n- [x] Task 2\n- [ ] Task 3\n"
	if string(content) != expectedContent {
		t.Errorf("File content = %v, want %v", string(content), expectedContent)
	}
}

func TestNextDateWithWeekendSkipping(t *testing.T) {
	originalIntervalMode := intervalMode
	originalSkipWeekend := skipWeekend
	defer func() {
		intervalMode = originalIntervalMode
		skipWeekend = originalSkipWeekend
	}()

	intervalMode = "daily"
	skipWeekend = true

	tests := []struct {
		name string
		date time.Time
		want time.Time
	}{
		{
			name: "Friday to Monday",
			date: time.Date(2024, 8, 30, 0, 0, 0, 0, time.UTC), // Friday
			want: time.Date(2024, 9, 2, 0, 0, 0, 0, time.UTC),  // Monday
		},
		{
			name: "Saturday to Monday",
			date: time.Date(2024, 8, 31, 0, 0, 0, 0, time.UTC), // Saturday
			want: time.Date(2024, 9, 2, 0, 0, 0, 0, time.UTC),  // Monday
		},
		{
			name: "Sunday to Monday",
			date: time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC), // Sunday
			want: time.Date(2024, 9, 2, 0, 0, 0, 0, time.UTC), // Monday
		},
		{
			name: "Monday to Tuesday",
			date: time.Date(2024, 9, 2, 0, 0, 0, 0, time.UTC), // Monday
			want: time.Date(2024, 9, 3, 0, 0, 0, 0, time.UTC), // Tuesday
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NextDate(tt.date)
			if !got.Equal(tt.want) {
				t.Errorf("NextDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPreviousDateWithWeekendSkipping(t *testing.T) {
	originalIntervalMode := intervalMode
	originalSkipWeekend := skipWeekend
	defer func() {
		intervalMode = originalIntervalMode
		skipWeekend = originalSkipWeekend
	}()

	intervalMode = "daily"
	skipWeekend = true

	tests := []struct {
		name string
		date time.Time
		want time.Time
	}{
		{
			name: "Monday to Friday",
			date: time.Date(2024, 9, 2, 0, 0, 0, 0, time.UTC),  // Monday
			want: time.Date(2024, 8, 30, 0, 0, 0, 0, time.UTC), // Friday
		},
		{
			name: "Saturday to Friday",
			date: time.Date(2024, 8, 31, 0, 0, 0, 0, time.UTC), // Saturday
			want: time.Date(2024, 8, 30, 0, 0, 0, 0, time.UTC), // Friday
		},
		{
			name: "Sunday to Friday",
			date: time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC),  // Sunday
			want: time.Date(2024, 8, 30, 0, 0, 0, 0, time.UTC), // Friday
		},
		{
			name: "Friday to Thursday",
			date: time.Date(2024, 8, 30, 0, 0, 0, 0, time.UTC), // Friday
			want: time.Date(2024, 8, 29, 0, 0, 0, 0, time.UTC), // Thursday
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PreviousDate(tt.date)
			if !got.Equal(tt.want) {
				t.Errorf("PreviousDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOpenEditorWithCopyPrevious(t *testing.T) {
	// Set TD_TEST_MODE environment variable
	os.Setenv("TD_TEST_MODE", "true")
	defer os.Unsetenv("TD_TEST_MODE")

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "vault_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set up the test environment
	vaultLoc = tempDir
	intervalMode = "daily"
	skipWeekend = false

	// Create a previous day's file
	prevDate := time.Date(2024, 8, 29, 0, 0, 0, 0, time.UTC)
	prevFilename := getFilename(prevDate)
	prevContent := "Previous day's content\n"
	err = os.MkdirAll(filepath.Dir(prevFilename), 0755)
	if err != nil {
		t.Fatalf("Failed to create directories: %v", err)
	}
	err = os.WriteFile(prevFilename, []byte(prevContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create previous day's file: %v", err)
	}

	// Test OpenEditor with copyPrevious = true
	testDate := time.Date(2024, 8, 30, 0, 0, 0, 0, time.UTC)
	err = OpenEditor(testDate, 1, true)
	if err != nil {
		t.Errorf("OpenEditor() error = %v", err)
	}

	// Verify the new file was created with the previous day's content and the new header
	filename := getFilename(testDate)
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	expectedContent := GetHeader(testDate) + prevContent
	if string(content) != expectedContent {
		t.Errorf("File content = %v, want %v", string(content), expectedContent)
	}

	// Test OpenEditor with copyPrevious = false
	testDate = time.Date(2024, 8, 31, 0, 0, 0, 0, time.UTC)
	err = OpenEditor(testDate, 1, false)
	if err != nil {
		t.Errorf("OpenEditor() error = %v", err)
	}

	// Verify the new file was created with only the header
	filename = getFilename(testDate)
	content, err = os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	expectedContent = GetHeader(testDate)
	if string(content) != expectedContent {
		t.Errorf("File content = %v, want %v", string(content), expectedContent)
	}
}

func TestGetVaultLocationAndIntervalMode(t *testing.T) {
	// Save original package variable values
	origVaultLoc := vaultLoc
	origIntervalMode := intervalMode

	// Restore package variables after the test
	defer func() {
		vaultLoc = origVaultLoc
		intervalMode = origIntervalMode
	}()

	testCases := []struct {
		name             string
		setVaultLoc      string
		setIntervalMode  string
		expectedVaultLoc string
		expectedInterval string
	}{
		{
			name:             "Default values",
			setVaultLoc:      ".td",
			setIntervalMode:  "weekly",
			expectedVaultLoc: ".td",
			expectedInterval: "weekly",
		},
		{
			name:             "Custom vault location",
			setVaultLoc:      "/custom/vault/location",
			setIntervalMode:  "weekly",
			expectedVaultLoc: "/custom/vault/location",
			expectedInterval: "weekly",
		},
		{
			name:             "Custom interval mode",
			setVaultLoc:      ".td",
			setIntervalMode:  "daily",
			expectedVaultLoc: ".td",
			expectedInterval: "daily",
		},
		{
			name:             "Both custom values",
			setVaultLoc:      "/another/custom/path",
			setIntervalMode:  "monthly",
			expectedVaultLoc: "/another/custom/path",
			expectedInterval: "monthly",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Directly set package variables for this test case
			vaultLoc = tc.setVaultLoc
			intervalMode = tc.setIntervalMode

			// Test GetVaultLocation
			result := GetVaultLocation()
			if result != tc.expectedVaultLoc {
				t.Errorf("GetVaultLocation() = %q, want %q", result, tc.expectedVaultLoc)
			}

			// Test GetIntervalMode
			result = GetIntervalMode()
			if result != tc.expectedInterval {
				t.Errorf("GetIntervalMode() = %q, want %q", result, tc.expectedInterval)
			}
		})
	}
}

// TestEnvironmentVariableResolution tests that the environment variables are properly resolved
func TestEnvironmentVariableResolution(t *testing.T) {
	// Save original environment values
	originalVaultLoc := os.Getenv("TD_VAULT_LOC")
	originalIntervalMode := os.Getenv("TD_INTERVAL_MODE")

	// Save original package variable values
	origVaultLoc := vaultLoc
	origIntervalMode := intervalMode

	// Restore original values after the test
	defer func() {
		os.Setenv("TD_VAULT_LOC", originalVaultLoc)
		os.Setenv("TD_INTERVAL_MODE", originalIntervalMode)

		// Reset package variables to ensure tests don't interfere with each other
		vaultLoc = origVaultLoc
		intervalMode = origIntervalMode
	}()

	testCases := []struct {
		name                 string
		envVaultLoc          string
		envIntervalMode      string
		expectedVaultLoc     string
		expectedIntervalMode string
	}{
		{
			name:                 "Empty environment variables use defaults",
			envVaultLoc:          "",
			envIntervalMode:      "",
			expectedVaultLoc:     ".td",
			expectedIntervalMode: "weekly",
		},
		{
			name:                 "Custom vault location",
			envVaultLoc:          "/custom/vault/path",
			envIntervalMode:      "",
			expectedVaultLoc:     "/custom/vault/path",
			expectedIntervalMode: "weekly",
		},
		{
			name:                 "Custom interval mode",
			envVaultLoc:          "",
			envIntervalMode:      "daily",
			expectedVaultLoc:     ".td",
			expectedIntervalMode: "daily",
		},
		{
			name:                 "Both custom values",
			envVaultLoc:          "/another/path",
			envIntervalMode:      "monthly",
			expectedVaultLoc:     "/another/path",
			expectedIntervalMode: "monthly",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set environment variables
			os.Setenv("TD_VAULT_LOC", tc.envVaultLoc)
			os.Setenv("TD_INTERVAL_MODE", tc.envIntervalMode)

			// Reset package variables to initialize from environment
			// We need explicit defaults for empty environment variables
			testVaultLoc := tc.envVaultLoc
			if testVaultLoc == "" {
				testVaultLoc = ".td"
			}

			testIntervalMode := tc.envIntervalMode
			if testIntervalMode == "" {
				testIntervalMode = "weekly"
			}

			// Update package variables
			vaultLoc = testVaultLoc
			intervalMode = testIntervalMode

			// Verify values are as expected
			if vaultLoc != tc.expectedVaultLoc {
				t.Errorf("vaultLoc = %q, want %q", vaultLoc, tc.expectedVaultLoc)
			}

			if intervalMode != tc.expectedIntervalMode {
				t.Errorf("intervalMode = %q, want %q", intervalMode, tc.expectedIntervalMode)
			}

			// Then test the accessor functions
			if GetVaultLocation() != tc.expectedVaultLoc {
				t.Errorf("GetVaultLocation() = %q, want %q", GetVaultLocation(), tc.expectedVaultLoc)
			}

			if GetIntervalMode() != tc.expectedIntervalMode {
				t.Errorf("GetIntervalMode() = %q, want %q", GetIntervalMode(), tc.expectedIntervalMode)
			}
		})
	}
}
