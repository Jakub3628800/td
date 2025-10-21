package core

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Jakub3628800/td/internal/db"
)

const (
	spotifyAuthURL  = "https://accounts.spotify.com/authorize"
	spotifyTokenURL = "https://accounts.spotify.com/api/token"
	spotifyAPIBase  = "https://api.spotify.com/v1"
)

type SpotifyClient struct {
	clientID     string
	clientSecret string
	redirectURI  string
	queries      *db.Queries
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type Device struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	IsActive     bool   `json:"is_active"`
	IsRestricted bool   `json:"is_restricted"`
	VolumePercent int   `json:"volume_percent"`
}

type devicesResponse struct {
	Devices []Device `json:"devices"`
}

func NewSpotifyClient(clientID, clientSecret, redirectURI string) *SpotifyClient {
	return &SpotifyClient{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
		queries:      getQueries(),
	}
}

func (s *SpotifyClient) GetAuthURL() string {
	params := url.Values{}
	params.Add("client_id", s.clientID)
	params.Add("response_type", "code")
	params.Add("redirect_uri", s.redirectURI)
	params.Add("scope", "user-read-playback-state user-modify-playback-state")
	return fmt.Sprintf("%s?%s", spotifyAuthURL, params.Encode())
}

func (s *SpotifyClient) ExchangeCode(code string) error {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", s.redirectURI)

	return s.requestTokens(data)
}

func (s *SpotifyClient) RefreshToken() error {
	ctx := context.Background()
	tokens, err := s.queries.GetSpotifyTokens(ctx)
	if err != nil {
		return fmt.Errorf("failed to get refresh token: %w", err)
	}

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", tokens.RefreshToken)

	return s.requestTokens(data)
}

func (s *SpotifyClient) requestTokens(data url.Values) error {
	req, err := http.NewRequest("POST", spotifyTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(s.clientID + ":" + s.clientSecret))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("token request failed: %s - %s", resp.Status, string(body))
	}

	var tokenResp tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return err
	}

	expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	ctx := context.Background()
	return s.queries.UpsertSpotifyTokens(ctx, db.UpsertSpotifyTokensParams{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    expiresAt,
	})
}

func (s *SpotifyClient) getValidToken() (string, error) {
	ctx := context.Background()
	tokens, err := s.queries.GetSpotifyTokens(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not authenticated - run 'td spotify auth' first")
		}
		return "", err
	}

	if time.Now().After(tokens.ExpiresAt) {
		if err := s.RefreshToken(); err != nil {
			return "", fmt.Errorf("failed to refresh token: %w", err)
		}
		tokens, err = s.queries.GetSpotifyTokens(ctx)
		if err != nil {
			return "", err
		}
	}

	return tokens.AccessToken, nil
}

func (s *SpotifyClient) makeAPIRequest(method, endpoint string, body io.Reader) (*http.Response, error) {
	token, err := s.getValidToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, spotifyAPIBase+endpoint, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return http.DefaultClient.Do(req)
}

func (s *SpotifyClient) Play(uri string, deviceID string) error {
	endpoint := "/me/player/play"
	if deviceID != "" {
		endpoint += "?device_id=" + deviceID
	}

	var body io.Reader
	if uri != "" {
		var jsonBody string
		if strings.Contains(uri, ":playlist:") || strings.Contains(uri, ":album:") || strings.Contains(uri, ":artist:") {
			jsonBody = fmt.Sprintf(`{"context_uri": "%s"}`, uri)
		} else {
			jsonBody = fmt.Sprintf(`{"uris": ["%s"]}`, uri)
		}
		body = strings.NewReader(jsonBody)
	}

	resp, err := s.makeAPIRequest("PUT", endpoint, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("play failed: %s - %s", resp.Status, string(bodyBytes))
	}

	return nil
}

func (s *SpotifyClient) Pause(deviceID string) error {
	endpoint := "/me/player/pause"
	if deviceID != "" {
		endpoint += "?device_id=" + deviceID
	}

	resp, err := s.makeAPIRequest("PUT", endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("pause failed: %s - %s", resp.Status, string(bodyBytes))
	}

	return nil
}

func (s *SpotifyClient) GetDevices() ([]Device, error) {
	resp, err := s.makeAPIRequest("GET", "/me/player/devices", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get devices failed: %s - %s", resp.Status, string(bodyBytes))
	}

	var devicesResp devicesResponse
	if err := json.NewDecoder(resp.Body).Decode(&devicesResp); err != nil {
		return nil, err
	}

	return devicesResp.Devices, nil
}

func (s *SpotifyClient) TransferPlayback(deviceID string, play bool) error {
	jsonBody := fmt.Sprintf(`{"device_ids": ["%s"], "play": %t}`, deviceID, play)

	resp, err := s.makeAPIRequest("PUT", "/me/player", strings.NewReader(jsonBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("transfer playback failed: %s - %s", resp.Status, string(bodyBytes))
	}

	return nil
}
