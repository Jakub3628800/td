/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"os"

	"github.com/Jakub3628800/td/cmd"
)

// These variables are set during build time
var (
	// Version is the current version of the application
	Version = "0.1.0"
	// Commit is the git commit hash of the build
	Commit = "unknown"
	// BuildDate is the date when the binary was built
	BuildDate = "unknown"
)

func main() {
	// Check for version flag
	for _, arg := range os.Args {
		if arg == "-v" || arg == "--version" {
			fmt.Printf("td version %s (commit: %s, built: %s)\n", Version, Commit, BuildDate)
			os.Exit(0)
		}
	}

	// Intercept --start or --end as first argument and rewrite to 'day --start' or 'day --end'
	if len(os.Args) > 1 && (os.Args[1] == "--start" || os.Args[1] == "--end") {
		// Insert 'day' as the first argument after the program name
		os.Args = append([]string{os.Args[0], "day"}, os.Args[1:]...)
	}

	cmd.Execute()
}
