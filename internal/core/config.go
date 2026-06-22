package core

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/Jakub3628800/td/internal/db"
)

// AllowedConfigKeys defines the config keys that can be set via CLI
var AllowedConfigKeys = map[string]string{
	"music_control_enabled":  "Enable/disable music control during pomodoro sessions (true/false)",
	"spotify_default_device": "Default Spotify device ID for playback control",
	"default_pomo_duration":  "Default duration for pomodoro sessions in minutes (15-60, increment by 5)",
}

// ValidateConfigKey checks if a config key is allowed
func ValidateConfigKey(key string) error {
	if _, ok := AllowedConfigKeys[key]; !ok {
		return fmt.Errorf("unknown config key '%s'. Allowed keys: %v", key, getAllowedKeys())
	}
	return nil
}

// getAllowedKeys returns a slice of all allowed config keys
func getAllowedKeys() []string {
	keys := make([]string, 0, len(AllowedConfigKeys))
	for k := range AllowedConfigKeys {
		keys = append(keys, k)
	}
	return keys
}

// GetConfig retrieves a config value by key
func GetConfig(ctx context.Context, queries *db.Queries, key string) (string, error) {
	config, err := queries.GetConfig(ctx, key)
	if err != nil {
		return "", err
	}
	if config.Value.Valid {
		return config.Value.String, nil
	}
	return "", nil
}

// SetConfig stores a config value by key
func SetConfig(ctx context.Context, queries *db.Queries, key, value string) error {
	return queries.SetConfig(ctx, db.SetConfigParams{
		Key:   key,
		Value: sql.NullString{String: value, Valid: true},
	})
}

// IsMusicControlEnabled checks if music control is enabled in config
func IsMusicControlEnabled() bool {
	ctx := context.Background()
	queries, err := GetDB()
	if err != nil {
		return false
	}

	value, err := GetConfig(ctx, queries, "music_control_enabled")
	if err != nil {
		return false
	}

	return value == "true"
}

// SpotifyClient wraps Spotify-related configuration
type SpotifyClient struct {
	clientID     string
	clientSecret string
	redirectURI  string
}

// NewSpotifyClient creates a new Spotify client
func NewSpotifyClient(clientID, clientSecret, redirectURI string) *SpotifyClient {
	return &SpotifyClient{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
	}
}

// GetDefaultDevice retrieves the default Spotify device ID from config
func (sc *SpotifyClient) GetDefaultDevice() (string, error) {
	ctx := context.Background()
	queries, err := GetDB()
	if err != nil {
		return "", err
	}

	device, err := GetConfig(ctx, queries, "spotify_default_device")
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return device, nil
}

// SetDefaultDevice stores the default Spotify device ID in config
func (sc *SpotifyClient) SetDefaultDevice(deviceID string) error {
	ctx := context.Background()
	queries, err := GetDB()
	if err != nil {
		return err
	}

	if deviceID == "" {
		return queries.SetConfig(ctx, db.SetConfigParams{
			Key:   "spotify_default_device",
			Value: sql.NullString{Valid: false},
		})
	}

	return SetConfig(ctx, queries, "spotify_default_device", deviceID)
}

// GetDefaultPromoDuration retrieves the default pomodoro duration from config
// Returns 25 as default if not set or invalid
func GetDefaultPromoDuration() int {
	ctx := context.Background()
	queries, err := GetDB()
	if err != nil {
		return 25
	}

	value, err := GetConfig(ctx, queries, "default_pomo_duration")
	if err != nil || value == "" {
		return 25
	}

	duration, err := strconv.Atoi(value)
	if err != nil || duration < 15 || duration > 60 {
		return 25
	}

	return duration
}
