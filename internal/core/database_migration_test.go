package core

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/Jakub3628800/td/internal/db"
)

func TestDefaultDeviceMigration(t *testing.T) {
	// Create a temporary database for testing
	dbPath := "/tmp/test_migration.db"
	defer os.Remove(dbPath)

	os.Setenv("TD_DB_PATH", dbPath)
	defer os.Unsetenv("TD_DB_PATH")

	// Initialize global DB
	_, err := GetDB()
	if err != nil {
		t.Fatalf("Failed to init DB: %v", err)
	}
	defer CloseDB()

	// Set up the old schema with default_device_id column
	oldSchema := `
	DROP TABLE IF EXISTS spotify_tokens;
	CREATE TABLE spotify_tokens (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		access_token TEXT NOT NULL,
		refresh_token TEXT NOT NULL,
		expires_at TIMESTAMP NOT NULL,
		default_device_id TEXT,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := globalDB.Exec(oldSchema); err != nil {
		t.Fatalf("Failed to create old schema: %v", err)
	}

	// Insert test data with a default device ID
	testDeviceID := "test-device-123"
	_, err = globalDB.Exec(`
		INSERT INTO spotify_tokens (id, access_token, refresh_token, expires_at, default_device_id)
		VALUES (1, 'test-access', 'test-refresh', ?, ?)
	`, time.Now().Add(1*time.Hour), testDeviceID)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Close and reinitialize to trigger migration
	CloseDB()
	globalQueries = nil // Reset global state
	queries, err := GetDB()
	if err != nil {
		t.Fatalf("Failed to reinit DB for migration: %v", err)
	}

	// Verify the config table was created and data migrated
	ctx := context.Background()

	config, err := queries.GetConfig(ctx, "spotify_default_device")
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}

	if !config.Value.Valid || config.Value.String != testDeviceID {
		t.Errorf("Expected default device %s, got %v", testDeviceID, config.Value)
	}

	// Verify the old column no longer exists
	var columnExists bool
	row := globalDB.QueryRow(`
		SELECT COUNT(*) > 0
		FROM pragma_table_info('spotify_tokens')
		WHERE name = 'default_device_id'
	`)
	if err := row.Scan(&columnExists); err != nil {
		t.Fatalf("Failed to check column existence: %v", err)
	}

	if columnExists {
		t.Error("default_device_id column should have been removed")
	}
}

func TestDefaultDeviceMigrationNoExistingData(t *testing.T) {
	// Test that migration works when there's no existing default_device_id
	dbPath := "/tmp/test_migration_empty.db"
	defer os.Remove(dbPath)

	os.Setenv("TD_DB_PATH", dbPath)
	defer os.Unsetenv("TD_DB_PATH")

	queries, err := GetDB()
	if err != nil {
		t.Fatalf("Failed to init DB: %v", err)
	}
	defer CloseDB()

	// Verify config table exists but has no spotify_default_device entry
	ctx := context.Background()

	_, err = queries.GetConfig(ctx, "spotify_default_device")
	if err != sql.ErrNoRows {
		t.Errorf("Expected ErrNoRows, got %v", err)
	}
}

func TestSetAndGetDefaultDevice(t *testing.T) {
	dbPath := "/tmp/test_default_device.db"
	defer os.Remove(dbPath)

	os.Setenv("TD_DB_PATH", dbPath)
	defer os.Unsetenv("TD_DB_PATH")

	queries, err := GetDB()
	if err != nil {
		t.Fatalf("Failed to init DB: %v", err)
	}
	defer CloseDB()

	client := NewSpotifyClient("test-id", "test-secret", "http://localhost:8888/callback")

	// Test setting default device
	testDeviceID := "test-device-456"
	if err := client.SetDefaultDevice(testDeviceID); err != nil {
		t.Fatalf("Failed to set default device: %v", err)
	}

	// Test getting default device
	deviceID, err := client.GetDefaultDevice()
	if err != nil {
		t.Fatalf("Failed to get default device: %v", err)
	}

	if deviceID != testDeviceID {
		t.Errorf("Expected device ID %s, got %s", testDeviceID, deviceID)
	}

	// Test updating default device
	newDeviceID := "test-device-789"
	if err := client.SetDefaultDevice(newDeviceID); err != nil {
		t.Fatalf("Failed to update default device: %v", err)
	}

	deviceID, err = client.GetDefaultDevice()
	if err != nil {
		t.Fatalf("Failed to get updated default device: %v", err)
	}

	if deviceID != newDeviceID {
		t.Errorf("Expected device ID %s, got %s", newDeviceID, deviceID)
	}

	// Test getting when no default is set returns empty string
	ctx := context.Background()
	if err := queries.SetConfig(ctx, db.SetConfigParams{
		Key:   "spotify_default_device",
		Value: sql.NullString{Valid: false},
	}); err != nil {
		t.Fatalf("Failed to clear default device: %v", err)
	}

	deviceID, err = client.GetDefaultDevice()
	if err != nil {
		t.Fatalf("Failed to get cleared default device: %v", err)
	}

	if deviceID != "" {
		t.Errorf("Expected empty device ID, got %s", deviceID)
	}
}
