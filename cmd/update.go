package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var Version string

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check for and install updates",
	Long:  `Check for the latest version of td and optionally update to it.`,
	Run: func(_ *cobra.Command, _ []string) {
		checkAndUpdate()
	},
}

func checkAndUpdate() {
	latestVersion, err := getLatestVersion()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking for updates: %v\n", err)
		os.Exit(1)
	}

	cmp := compareVersions(Version, latestVersion)
	if cmp == 0 {
		fmt.Printf("You are already on the latest version: %s\n", Version)
		return
	}

	if cmp > 0 {
		fmt.Printf("You are on a newer version (%s) than the latest tag (%s)\n", Version, latestVersion)
		return
	}

	fmt.Printf("Current version: %s\n", Version)
	fmt.Printf("Latest version:  %s\n", latestVersion)
	fmt.Print("\nUpdate to the latest version? (y/n): ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	if response != "y" && response != "yes" {
		fmt.Println("Update cancelled.")
		return
	}

	performUpdate(latestVersion)
}

// compareVersions compares two semantic versions (major.minor.patch)
// Returns: -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func compareVersions(v1, v2 string) int {
	parts1 := parseVersion(v1)
	parts2 := parseVersion(v2)

	for i := 0; i < 3; i++ {
		if parts1[i] < parts2[i] {
			return -1
		}
		if parts1[i] > parts2[i] {
			return 1
		}
	}
	return 0
}

// parseVersion extracts major, minor, patch from a version string (with optional 'v' prefix)
// Returns [major, minor, patch] as integers, defaults to 0 for missing parts
func parseVersion(v string) [3]int {
	v = strings.TrimPrefix(v, "v")
	parts := strings.Split(v, ".")
	result := [3]int{0, 0, 0}

	for i := 0; i < len(parts) && i < 3; i++ {
		if num, err := strconv.Atoi(parts[i]); err == nil {
			result[i] = num
		}
	}
	return result
}

func getLatestVersion() (string, error) {
	// Query GitHub API for all tags and find the latest version
	resp, err := http.Get("https://api.github.com/repos/Jakub3628800/td/tags")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch tags: status %d", resp.StatusCode)
	}

	var tags []struct {
		Name string `json:"name"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(body, &tags); err != nil {
		return "", err
	}

	if len(tags) == 0 {
		return "", fmt.Errorf("no tags found")
	}

	// Return the first tag (latest)
	return tags[0].Name, nil
}

func performUpdate(version string) {
	fmt.Printf("Downloading and installing version %s...\n", version)

	cmd := exec.Command("go", "install", fmt.Sprintf("github.com/Jakub3628800/td@%s", version))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error during update: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nSuccessfully updated to version %s!\n", version)
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
