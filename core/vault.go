package core

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var vaultLoc string
var intervalMode string // daily, weekly or monthly
var templatePath string
var skipWeekend bool

type Task struct {
	Line     string
	Selected bool
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func templateFile() string {
	return filepath.Join(vaultLoc, templatePath)
}

func init() {
	vaultLoc = getEnv("TD_VAULT_LOC", ".td")
	intervalMode = getEnv("TD_INTERVAL_MODE", "weekly")
	templatePath = getEnv("TD_TEMPLATE_PATH", ".template")
	skipWeekend = getEnv("TD_SKIP_WEEKEND", "false") == "true"
}

func getFilename(date time.Time) string {
	year, week := date.ISOWeek()
	month := date.Month().String()
	if intervalMode == "daily" {
		return filepath.Join(vaultLoc, date.Format("2006/January/02.md"))
	} else if intervalMode == "weekly" {
		return fmt.Sprintf("%s/%d/%s/week%d.md", vaultLoc, year, month, week)
	}

	return filepath.Join(vaultLoc, strconv.Itoa(year), month, month+".md")
}

func GetHeader(date time.Time) string {
	if intervalMode == "daily" {
		return date.Format("2006-01-02") + " " + date.Weekday().String() + "\n\n"
	}

	_, week := date.ISOWeek()
	return "Week " + strconv.Itoa(week) + "\n\n"
}

func NextDate(date time.Time) time.Time {
	if intervalMode == "daily" {
		next := date.AddDate(0, 0, 1)
		if skipWeekend {
			for next.Weekday() == time.Saturday || next.Weekday() == time.Sunday {
				next = next.AddDate(0, 0, 1)
			}
		}
		return next
	} else if intervalMode == "weekly" {
		return date.AddDate(0, 0, 7)
	}
	return date.AddDate(0, 1, 0) // Monthly mode
}

func PreviousDate(date time.Time) time.Time {
	if intervalMode == "daily" {
		prev := date.AddDate(0, 0, -1)
		if skipWeekend {
			for prev.Weekday() == time.Saturday || prev.Weekday() == time.Sunday {
				prev = prev.AddDate(0, 0, -1)
			}
		}
		return prev
	} else if intervalMode == "weekly" {
		return date.AddDate(0, 0, -7)
	}
	return date.AddDate(0, -1, 0) // Monthly mode
}

func openFile(date time.Time) (*os.File, error) {
	filename := getFilename(date)
	return os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func createFile(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if fileExists(templateFile()) {
		// Use filepath.Clean to sanitize the path
		cleanTemplatePath := filepath.Clean(templateFile())
		cleanDestPath := filepath.Clean(path)

		// Use absolute paths to avoid path traversal
		absTemplatePath, err := filepath.Abs(cleanTemplatePath)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %w", err)
		}

		absDestPath, err := filepath.Abs(cleanDestPath)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %w", err)
		}

		cmd := exec.Command("cp", absTemplatePath, absDestPath)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to copy file: %w", err)
		}
	} else {
		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer file.Close()
	}
	return nil
}

func AddTask(date time.Time, line string) error {
	line = "- [ ] " + line

	filename := getFilename(date)
	if !fileExists(filename) {
		if err := createFile(filename); err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
	}
	file, err := openFile(date)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(line + "\n")
	if err != nil {
		return err
	}

	return nil
}

func isLineCheckbox(line string) (bool, bool) {
	pattern := regexp.MustCompile(`^([\t ]*)- \[( |x)\] ?(.*)$`)
	if matches := pattern.FindStringSubmatch(line); matches != nil {
		return true, matches[2] == "x"
	}
	return false, false
}

func linesWithSelection(filename string) ([]Task, error) {
	var tasks []Task

	if !fileExists(filename) {
		if fileExists(templateFile()) {
			return linesWithSelection(templateFile())
		}
		return tasks, nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		if trimmedLine == "" || trimmedLine == "- [ ]" || trimmedLine == "- [x]" {
			continue
		}
		isCheck, selected := isLineCheckbox(line)
		if isCheck {
			tasks = append(tasks, Task{line, selected})
		}
	}

	if err := scanner.Err(); err != nil {
		return tasks, fmt.Errorf("error reading file: %v", err)
	}
	return tasks, nil
}

func LoadLinesWithSelection(date time.Time) ([]Task, error) {
	filename := getFilename(date)
	return linesWithSelection(filename)
}

var copyPreviousEnv bool

func init() {
	copyPreviousEnv = getEnv("TD_COPY_PREVIOUS", "false") == "true"
}

func OpenEditor(date time.Time, lineNumber int, copyPrevious bool) error {
	filename := getFilename(date)

	// Ensure the directory exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create the file if it doesn't exist
	if !fileExists(filename) {
		var content string

		// Check if template file exists and use it as a base
		if fileExists(templateFile()) && !copyPrevious && !copyPreviousEnv {
			// Read template content
			templateContent, err := os.ReadFile(templateFile())
			if err != nil {
				return fmt.Errorf("failed to read template file: %w", err)
			}
			content = string(templateContent)
		} else if copyPrevious || copyPreviousEnv {
			// If copy from previous is requested, use that instead of template
			prevDate := PreviousDate(date)
			prevFilename := getFilename(prevDate)
			if fileExists(prevFilename) {
				prevContent, err := os.ReadFile(prevFilename)
				if err != nil {
					return fmt.Errorf("failed to read previous file: %w", err)
				}
				content = string(prevContent)
			}
		}

		// Add the header
		header := GetHeader(date)
		content = header + content

		if err := os.WriteFile(filename, []byte(content), 0600); err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
	}

	// Skip launching the editor during tests
	if os.Getenv("TD_TEST_MODE") == "true" {
		return nil
	}

	// Get the editor from the EDITOR environment variable, defaulting to "vim"
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	// Prepare the command
	args := []string{fmt.Sprintf("+%d", lineNumber), filename}
	cmd := exec.Command(editor, args...)

	// Set the command's standard input, output, and error to the current program's ones
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command and wait for it to finish
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error running editor: %w", err)
	}

	return nil
}

func UpdateTaskStatus(selected bool, taskDescription string, date time.Time) error {
	filename := getFilename(date)
	if !fileExists(filename) {
		if err := createFile(filename); err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
	}

	file, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	lines := strings.Split(string(file), "\n")

	lineUpdated := false
	for i, line := range lines {
		if strings.Contains(line, taskDescription) {
			if selected {
				lines[i] = strings.Replace(line, "- [ ]", "- [x]", 1)
			} else {
				lines[i] = strings.Replace(line, "- [x]", "- [ ]", 1)
			}
			lineUpdated = true
			break
		}
	}

	if !lineUpdated {
		return fmt.Errorf("task not found in the file")
	}

	updatedContent := strings.Join(lines, "\n")

	err = os.WriteFile(filename, []byte(updatedContent), 0600)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}

func ContainsLine(date time.Time, searchLine string) (int, error) {
	filename := getFilename(date)

	if !fileExists(filename) {
		return containsLine(templateFile(), searchLine)
	}
	return containsLine(filename, searchLine)
}

func containsLine(filename string, searchLine string) (int, error) {
	if !fileExists(filename) {
		return 0, nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimPrefix(strings.TrimPrefix(line, "- [ ]"), "- [x]")
		if strings.TrimSpace(trimmedLine) == strings.TrimSpace(searchLine) {
			return lineNumber, nil
		}
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return 0, nil
}

// GetVaultLocation returns the configured vault location
func GetVaultLocation() string {
	return vaultLoc
}

// GetIntervalMode returns the configured interval mode
func GetIntervalMode() string {
	return intervalMode
}

// ResetForTest allows test code to reset the package variables
// This should only be used in tests, not in production code
func ResetForTest(newVaultLoc, newIntervalMode string) {
	vaultLoc = newVaultLoc
	intervalMode = newIntervalMode
}
