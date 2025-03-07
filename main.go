/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"td/cmd"
)

// Version information
var (
	Version   = "0.1.0"
	Commit    = "unknown"
	BuildDate = "unknown"
)

func main() {
	// Check if version flag is provided
	versionFlag := flag.Bool("version", false, "Print version information")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("td version %s (commit: %s, built: %s)\n", Version, Commit, BuildDate)
		os.Exit(0)
	}

	cmd.Execute()
}
