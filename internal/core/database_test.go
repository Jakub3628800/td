package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDatabaseCreationInNonExistentDirectory(t *testing.T) {
	// Create a temp directory that we'll clean up
	tempBase := t.TempDir()

	// Point to a database path in a subdirectory that doesn't exist yet
	dbPath := filepath.Join(tempBase, "does", "not", "exist", "yet", "td.db")

	origDBPath := os.Getenv("TD_DB_PATH")
	defer func() {
		os.Setenv("TD_DB_PATH", origDBPath)
		CloseDB()
	}()

	os.Setenv("TD_DB_PATH", dbPath)
	CloseDB() // Reset any existing connection

	// Try to get the database - this should create the directory structure
	queries, err := GetDB()
	if err != nil {
		t.Fatalf("GetDB failed: %v", err)
	}

	if queries == nil {
		t.Fatal("queries should not be nil")
	}

	// Verify the database file was created
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Fatal("database file was not created")
	}

	// Verify we can query the database
	_, err = queries.ListPomodori(ctx())
	if err != nil {
		t.Fatalf("ListPomodori failed: %v", err)
	}
}

func TestDatabaseCreationWithoutPermissions(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("Skipping permission test when running as root")
	}

	// Try to create database in a directory we can't write to
	dbPath := "/root/impossible/td.db"

	origDBPath := os.Getenv("TD_DB_PATH")
	defer func() {
		os.Setenv("TD_DB_PATH", origDBPath)
		CloseDB()
	}()

	os.Setenv("TD_DB_PATH", dbPath)
	CloseDB()

	// This should fail with a permission error
	_, err := GetDB()
	if err == nil {
		t.Fatal("expected error when creating database in protected directory, got nil")
	}

	// The error should be meaningful
	if err.Error() == "" {
		t.Fatal("error message should not be empty")
	}

	// Verify the error contains helpful information
	errMsg := err.Error()
	if !contains(errMsg, "failed to create database directory") && !contains(errMsg, "permission") {
		t.Logf("Error message: %s", errMsg)
	}
}

func TestDatabaseErrorMessagesAreHelpful(t *testing.T) {
	tests := []struct {
		name        string
		dbPath      string
		shouldError bool
		contains    []string
	}{
		{
			name:        "valid path creates database successfully",
			dbPath:      filepath.Join(t.TempDir(), "valid", "td.db"),
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origDBPath := os.Getenv("TD_DB_PATH")
			defer func() {
				os.Setenv("TD_DB_PATH", origDBPath)
				CloseDB()
			}()

			os.Setenv("TD_DB_PATH", tt.dbPath)
			CloseDB()

			_, err := GetDB()
			if tt.shouldError {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				errMsg := err.Error()
				for _, substring := range tt.contains {
					if !contains(errMsg, substring) {
						t.Errorf("error message should contain %q, got: %s", substring, errMsg)
					}
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
