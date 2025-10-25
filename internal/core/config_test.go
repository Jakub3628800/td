package core

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/Jakub3628800/td/internal/db"
)

func TestConfigGetSet(t *testing.T) {
	// Create a temporary database for testing
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_config.db")

	// Set the test database path
	originalPath := os.Getenv("TD_DB_PATH")
	os.Setenv("TD_DB_PATH", dbPath)
	defer os.Setenv("TD_DB_PATH", originalPath)

	// Reset global state
	globalDB = nil
	globalQueries = nil
	defer CloseDB()

	// Get database connection
	queries, err := GetDB()
	if err != nil {
		t.Fatalf("Failed to get database: %v", err)
	}

	ctx := context.Background()

	// Test 1: Get non-existent config key should return sql.ErrNoRows
	_, err = queries.GetConfig(ctx, "nonexistent_key")
	if err != sql.ErrNoRows {
		t.Errorf("Expected sql.ErrNoRows for non-existent key, got: %v", err)
	}

	// Test 2: Set a config value
	err = queries.SetConfig(ctx, db.SetConfigParams{
		Key:   "test_key",
		Value: sql.NullString{String: "test_value", Valid: true},
	})
	if err != nil {
		t.Fatalf("Failed to set config: %v", err)
	}

	// Test 3: Get the config value back
	config, err := queries.GetConfig(ctx, "test_key")
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}

	if config.Key != "test_key" {
		t.Errorf("Expected key 'test_key', got '%s'", config.Key)
	}

	if !config.Value.Valid || config.Value.String != "test_value" {
		t.Errorf("Expected value 'test_value', got '%v'", config.Value)
	}

	// Test 4: Update the config value
	err = queries.SetConfig(ctx, db.SetConfigParams{
		Key:   "test_key",
		Value: sql.NullString{String: "updated_value", Valid: true},
	})
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	// Test 5: Verify the update
	config, err = queries.GetConfig(ctx, "test_key")
	if err != nil {
		t.Fatalf("Failed to get updated config: %v", err)
	}

	if !config.Value.Valid || config.Value.String != "updated_value" {
		t.Errorf("Expected updated value 'updated_value', got '%v'", config.Value)
	}

	// Test 6: Set null value
	err = queries.SetConfig(ctx, db.SetConfigParams{
		Key:   "null_key",
		Value: sql.NullString{Valid: false},
	})
	if err != nil {
		t.Fatalf("Failed to set null config: %v", err)
	}

	config, err = queries.GetConfig(ctx, "null_key")
	if err != nil {
		t.Fatalf("Failed to get null config: %v", err)
	}

	if config.Value.Valid {
		t.Errorf("Expected null value, got valid value: '%s'", config.Value.String)
	}
}

func TestSpotifyDefaultDevice(t *testing.T) {
	// Create a temporary database for testing
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_spotify.db")

	// Set the test database path
	originalPath := os.Getenv("TD_DB_PATH")
	os.Setenv("TD_DB_PATH", dbPath)
	defer os.Setenv("TD_DB_PATH", originalPath)

	// Reset global state
	globalDB = nil
	globalQueries = nil
	defer CloseDB()

	// Create Spotify client
	client := NewSpotifyClient("test_id", "test_secret", "http://localhost:8888/callback")

	// Test 1: Get default device when none is set
	device, err := client.GetDefaultDevice()
	if err != nil {
		t.Fatalf("Failed to get default device: %v", err)
	}

	if device != "" {
		t.Errorf("Expected empty device, got '%s'", device)
	}

	// Test 2: Set default device
	err = client.SetDefaultDevice("test_device_123")
	if err != nil {
		t.Fatalf("Failed to set default device: %v", err)
	}

	// Test 3: Get default device back
	device, err = client.GetDefaultDevice()
	if err != nil {
		t.Fatalf("Failed to get default device: %v", err)
	}

	if device != "test_device_123" {
		t.Errorf("Expected device 'test_device_123', got '%s'", device)
	}

	// Test 4: Update default device
	err = client.SetDefaultDevice("new_device_456")
	if err != nil {
		t.Fatalf("Failed to update default device: %v", err)
	}

	// Test 5: Verify update
	device, err = client.GetDefaultDevice()
	if err != nil {
		t.Fatalf("Failed to get updated default device: %v", err)
	}

	if device != "new_device_456" {
		t.Errorf("Expected device 'new_device_456', got '%s'", device)
	}

	// Test 6: Clear default device (set to empty string)
	err = client.SetDefaultDevice("")
	if err != nil {
		t.Fatalf("Failed to clear default device: %v", err)
	}

	// Test 7: Verify cleared
	device, err = client.GetDefaultDevice()
	if err != nil {
		t.Fatalf("Failed to get cleared default device: %v", err)
	}

	if device != "" {
		t.Errorf("Expected empty device after clear, got '%s'", device)
	}
}
