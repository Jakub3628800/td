package core

import (
	"context"
	"database/sql"

	"github.com/Jakub3628800/td/internal/db"
)

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
