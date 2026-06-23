package cmd

import (
	"fmt"
	"os"
)

var rootHelp = `To-Do ToDay (td) is a simple tool for tracking pomodoro sessions with configuration management.

Usage:
  td <command> [flags]

Commands:
  pomo         Start a Pomodoro timer
  config       Manage configuration options
  list-pomos   List all pomodoro sessions
  self-update  Check for and install updates

Flags:
  -h, --help   Show help
`

func Execute() {
	if err := execute(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func execute(args []string) error {
	if len(args) == 0 {
		fmt.Print(rootHelp)
		return nil
	}

	switch args[0] {
	case "-h", "--help", "help":
		return runHelpCommand(args[1:])
	case "pomo":
		return runPomoCommand(args[1:])
	case "config":
		return runConfigCommand(args[1:])
	case "list-pomos":
		return runListPomosCommand(args[1:])
	case "self-update":
		return runSelfUpdateCommand(args[1:])
	default:
		return fmt.Errorf("unknown command %q\n\n%s", args[0], rootHelp)
	}
}

func runHelpCommand(args []string) error {
	if len(args) == 0 {
		fmt.Print(rootHelp)
		return nil
	}

	switch args[0] {
	case "pomo":
		fmt.Print(pomoHelp)
	case "config":
		fmt.Print(configHelp)
	case "list-pomos":
		fmt.Print(listPomosHelp)
	case "self-update":
		fmt.Print(selfUpdateHelp)
	default:
		return fmt.Errorf("unknown help topic %q\n\n%s", args[0], rootHelp)
	}
	return nil
}

func SetVersion(version string) {
	Version = version
}
