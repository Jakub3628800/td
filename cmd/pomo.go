package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/Jakub3628800/td/internal/core"
)

var duration int
var tags []string

const (
	padding  = 2
	maxWidth = 80
)

var pomoHelp = `Start a Pomodoro timer for focused work sessions. Default duration is 25 minutes.

Usage:
  td pomo [-d minutes] [-t tag]

Flags:
  -d, --duration int   Duration in minutes (default 25)
  -t, --tag string     Add tags to the pomodoro (can be repeated)
`

func runPomoCommand(args []string) error {
	duration = 25
	tags = []string{}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "-h" || arg == "--help":
			fmt.Print(pomoHelp)
			return nil
		case arg == "-d" || arg == "--duration":
			if i+1 >= len(args) {
				return fmt.Errorf("%s requires a value", arg)
			}
			i++
			value, err := strconv.Atoi(args[i])
			if err != nil {
				return fmt.Errorf("invalid duration %q: %w", args[i], err)
			}
			duration = value
		case strings.HasPrefix(arg, "--duration="):
			value, err := strconv.Atoi(strings.TrimPrefix(arg, "--duration="))
			if err != nil {
				return fmt.Errorf("invalid duration %q: %w", arg, err)
			}
			duration = value
		case strings.HasPrefix(arg, "-d="):
			value, err := strconv.Atoi(strings.TrimPrefix(arg, "-d="))
			if err != nil {
				return fmt.Errorf("invalid duration %q: %w", arg, err)
			}
			duration = value
		case arg == "-t" || arg == "--tag":
			if i+1 >= len(args) {
				return fmt.Errorf("%s requires a value", arg)
			}
			i++
			tags = append(tags, args[i])
		case strings.HasPrefix(arg, "--tag="):
			tags = append(tags, strings.TrimPrefix(arg, "--tag="))
		case strings.HasPrefix(arg, "-t="):
			tags = append(tags, strings.TrimPrefix(arg, "-t="))
		default:
			return fmt.Errorf("unknown pomo argument %q\n\n%s", arg, pomoHelp)
		}
	}

	if duration <= 0 {
		return fmt.Errorf("duration must be greater than 0 minutes")
	}

	hasRunning, err := core.HasRunningPomodoro()
	if err != nil {
		return fmt.Errorf("error checking for running pomodoro: %w", err)
	}
	if hasRunning {
		return fmt.Errorf("a pomodoro is already running. Please wait for it to finish or use 'td list-pomos' to see active sessions")
	}

	return runPomoTimer()
}

type terminalState struct {
	state string
}

func enableRawMode() (*terminalState, error) {
	stateCmd := exec.Command("stty", "-g")
	stateCmd.Stdin = os.Stdin
	stateBytes, err := stateCmd.Output()
	if err != nil {
		return nil, err
	}
	state := strings.TrimSpace(string(stateBytes))
	cmd := exec.Command("stty", "raw", "-echo", "min", "0", "time", "1")
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return &terminalState{state: state}, nil
}

func (s *terminalState) restore() {
	if s == nil || s.state == "" {
		return
	}
	// #nosec G204 -- state is captured from `stty -g` immediately before raw mode.
	cmd := exec.Command("stty", s.state)
	cmd.Stdin = os.Stdin
	_ = cmd.Run()
}

func runPomoTimer() error {
	state, err := enableRawMode()
	if err != nil {
		return fmt.Errorf("error enabling raw terminal mode: %w", err)
	}
	defer state.restore()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Program panicked:", r)
			hasRunning, err := core.HasRunningPomodoro()
			if err == nil && hasRunning {
				fmt.Println("Cancelling running pomodoro session...")
				core.StopMusic()
				recordPomoSession(duration, "cancelled", tags)
			}
		}
	}()

	keyCh := make(chan byte, 8)
	doneKeys := make(chan struct{})
	go readKeys(keyCh, doneKeys)
	defer close(doneKeys)

	core.PlayMusic()

	start := time.Now()
	pomodoroDuration := time.Duration(duration) * time.Minute
	paused := false
	pauseTime := time.Time{}
	pausedElapsed := time.Duration(0)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	renderPomo(start, pomodoroDuration, pausedElapsed, paused, pauseTime, false)

	for {
		select {
		case key := <-keyCh:
			switch key {
			case 'q', 3: // ctrl+c
				core.StopMusic()
				recordPomoSession(duration, "cancelled", tags)
				fmt.Print("\n")
				return nil
			case 'p', ' ':
				if paused {
					pausedElapsed += time.Since(pauseTime)
					paused = false
				} else {
					paused = true
					pauseTime = time.Now()
				}
				renderPomo(start, pomodoroDuration, pausedElapsed, paused, pauseTime, false)
			}
		case <-ticker.C:
			if paused {
				continue
			}
			elapsed := time.Since(start) - pausedElapsed
			if elapsed >= pomodoroDuration {
				renderPomo(start, pomodoroDuration, pausedElapsed, paused, pauseTime, true)
				core.StopMusic()
				core.SendNotification(fmt.Sprintf("pomo session %dm done", duration), false)
				recordPomoSession(duration, "completed", tags)
				time.Sleep(time.Second)
				fmt.Print("\n")
				return nil
			}
			renderPomo(start, pomodoroDuration, pausedElapsed, paused, pauseTime, false)
		}
	}
}

func readKeys(keyCh chan<- byte, done <-chan struct{}) {
	buf := make([]byte, 1)
	for {
		select {
		case <-done:
			return
		default:
		}
		n, err := os.Stdin.Read(buf)
		if n == 0 {
			if err == io.EOF {
				continue
			}
			if err != nil {
				return
			}
			continue
		}
		if err != nil && err != io.EOF {
			return
		}
		select {
		case keyCh <- buf[0]:
		case <-done:
			return
		}
	}
}

func renderPomo(start time.Time, sessionDuration, pausedElapsed time.Duration, paused bool, pauseTime time.Time, done bool) {
	elapsed := time.Since(start) - pausedElapsed
	if paused {
		elapsed -= time.Since(pauseTime)
	}
	if elapsed < 0 {
		elapsed = 0
	}

	remaining := sessionDuration - elapsed
	if remaining < 0 {
		remaining = 0
	}

	minutes := int(remaining.Minutes())
	seconds := int(remaining.Seconds()) % 60

	status := ""
	if paused {
		status = "(Paused)"
	} else if done {
		status = "(Completed!)"
	}

	percent := float64(elapsed) / float64(sessionDuration)
	if done || percent > 1 {
		percent = 1
	}

	pad := strings.Repeat(" ", padding)
	barWidth := terminalWidth() - padding*2 - 4
	if barWidth > maxWidth {
		barWidth = maxWidth
	}
	if barWidth < 10 {
		barWidth = 10
	}

	fmt.Print("\033[H\033[2J")
	fmt.Print("\n")
	fmt.Printf("%s%02d:%02d %s\n", pad, minutes, seconds, status)
	fmt.Printf("%s%s\n\n", pad, renderProgressBar(percent, barWidth))
	fmt.Printf("%s\033[90mPress 'p' or space to pause/resume\033[0m\n", pad)
	fmt.Printf("%s\033[90mPress 'q' to quit\033[0m", pad)
}

func renderProgressBar(percent float64, width int) string {
	filled := int(percent * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}

func terminalWidth() int {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return maxWidth
	}
	fields := strings.Fields(string(out))
	if len(fields) != 2 {
		return maxWidth
	}
	width, err := strconv.Atoi(fields[1])
	if err != nil {
		return maxWidth
	}
	return width
}

func recordPomoSession(durationMinutes int, status string, tags []string) {
	now := time.Now()
	if err := core.SavePomodoroLog(durationMinutes, status, tags, now); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving pomodoro: %v\n", err)
	}
}
