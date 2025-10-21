package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Jakub3628800/td/internal/core"
)

var spotifyClientID string
var spotifyClientSecret string
var spotifyRedirectURI string

var spotifyCmd = &cobra.Command{
	Use:   "spotify",
	Short: "Control Spotify playback",
	Long:  `Control Spotify playback including playing, pausing, and managing devices.`,
}

var spotifyAuthCmd = &cobra.Command{
	Use:   "auth [authorization-code]",
	Short: "Authenticate with Spotify",
	Long: `Authenticate with Spotify using OAuth 2.0.
Run without arguments to get the authorization URL, then run again with the code from the callback.`,
	Run: func(_ *cobra.Command, args []string) {
		client := getSpotifyClient()

		if len(args) == 0 {
			authURL := client.GetAuthURL()
			fmt.Println("Please visit this URL to authorize:")
			fmt.Println(authURL)
			fmt.Println("\nAfter authorizing, run:")
			fmt.Println("  td spotify auth <code>")
			return
		}

		code := args[0]
		if err := client.ExchangeCode(code); err != nil {
			fmt.Fprintf(os.Stderr, "Error exchanging code: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Successfully authenticated with Spotify!")
	},
}

var spotifyPlayCmd = &cobra.Command{
	Use:   "play [uri]",
	Short: "Play a track, album, or playlist",
	Long: `Play a track, album, or playlist on Spotify.
If no URI is provided, resumes current playback.

Examples:
  td spotify play
  td spotify play spotify:track:6rqhFgbbKwnb9MLmUQDhG6
  td spotify play spotify:playlist:37i9dQZF1DXcBWIGoYBM5M
  td spotify play --device abc123 spotify:album:1DFixLWuPkv3KT3TnV35m3`,
	Run: func(cmd *cobra.Command, args []string) {
		client := getSpotifyClient()

		var uri string
		if len(args) > 0 {
			uri = args[0]
		}

		deviceID, _ := cmd.Flags().GetString("device")

		if err := client.Play(uri, deviceID); err != nil {
			fmt.Fprintf(os.Stderr, "Error playing: %v\n", err)
			os.Exit(1)
		}

		if uri == "" {
			fmt.Println("Resumed playback")
		} else {
			fmt.Printf("Playing: %s\n", uri)
		}
	},
}

var spotifyPauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause Spotify playback",
	Long:  `Pause the current Spotify playback.`,
	Run: func(cmd *cobra.Command, _ []string) {
		client := getSpotifyClient()

		deviceID, _ := cmd.Flags().GetString("device")

		if err := client.Pause(deviceID); err != nil {
			fmt.Fprintf(os.Stderr, "Error pausing: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Playback paused")
	},
}

var spotifyDevicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "List available Spotify devices",
	Long:  `List all available Spotify Connect devices.`,
	Run: func(_ *cobra.Command, _ []string) {
		client := getSpotifyClient()

		devices, err := client.GetDevices()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting devices: %v\n", err)
			os.Exit(1)
		}

		if len(devices) == 0 {
			fmt.Println("No devices found")
			return
		}

		fmt.Println("Available devices:")
		for _, device := range devices {
			status := ""
			if device.IsActive {
				status = " (active)"
			}
			if device.IsRestricted {
				status += " (restricted)"
			}

			fmt.Printf("  [%s] %s - %s%s\n", device.ID, device.Name, device.Type, status)
		}
	},
}

func getSpotifyClient() *core.SpotifyClient {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	if clientID == "" {
		clientID = spotifyClientID
	}
	if clientID == "" {
		fmt.Fprintln(os.Stderr, "Error: SPOTIFY_CLIENT_ID not set")
		fmt.Fprintln(os.Stderr, "Set it via environment variable or use --client-id flag")
		os.Exit(1)
	}

	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	if clientSecret == "" {
		clientSecret = spotifyClientSecret
	}
	if clientSecret == "" {
		fmt.Fprintln(os.Stderr, "Error: SPOTIFY_CLIENT_SECRET not set")
		fmt.Fprintln(os.Stderr, "Set it via environment variable or use --client-secret flag")
		os.Exit(1)
	}

	redirectURI := os.Getenv("SPOTIFY_REDIRECT_URI")
	if redirectURI == "" {
		redirectURI = spotifyRedirectURI
	}
	if redirectURI == "" {
		redirectURI = "http://localhost:8888/callback"
	}

	return core.NewSpotifyClient(clientID, clientSecret, redirectURI)
}

func init() {
	rootCmd.AddCommand(spotifyCmd)

	spotifyCmd.PersistentFlags().StringVar(&spotifyClientID, "client-id", "", "Spotify Client ID")
	spotifyCmd.PersistentFlags().StringVar(&spotifyClientSecret, "client-secret", "", "Spotify Client Secret")
	spotifyCmd.PersistentFlags().StringVar(&spotifyRedirectURI, "redirect-uri", "http://localhost:8888/callback", "OAuth redirect URI")

	spotifyCmd.AddCommand(spotifyAuthCmd)
	spotifyCmd.AddCommand(spotifyPlayCmd)
	spotifyCmd.AddCommand(spotifyPauseCmd)
	spotifyCmd.AddCommand(spotifyDevicesCmd)

	spotifyPlayCmd.Flags().StringP("device", "d", "", "Device ID to play on")
	spotifyPauseCmd.Flags().StringP("device", "d", "", "Device ID to pause")
}
