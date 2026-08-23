package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ResolveDirectory(current, requested string) (string, error) {
	requested = strings.TrimSpace(requested)
	if requested == "" || requested == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("resolve home directory: %w", err)
		}
		requested = home
	} else if strings.HasPrefix(requested, "~/") || strings.HasPrefix(requested, `~\`) {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("resolve home directory: %w", err)
		}
		requested = filepath.Join(home, requested[2:])
	}
	if !filepath.IsAbs(requested) {
		requested = filepath.Join(current, requested)
	}
	resolved, err := filepath.Abs(requested)
	if err != nil {
		return "", fmt.Errorf("resolve working directory: %w", err)
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return "", fmt.Errorf("open working directory %q: %w", resolved, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("working directory %q is not a directory", resolved)
	}
	return filepath.Clean(resolved), nil
}

func parseChangeDirectory(input string) (requested string, handled bool, err error) {
	trimmed := strings.TrimSpace(input)
	lower := strings.ToLower(trimmed)
	commandLength := 0
	switch {
	case lower == "cd":
		return "", true, nil
	case lower == ":cd":
		return "", true, nil
	case strings.HasPrefix(lower, "cd "):
		commandLength = len("cd")
	case strings.HasPrefix(lower, ":cd "):
		commandLength = len(":cd")
	default:
		return "", false, nil
	}
	requested = strings.TrimSpace(trimmed[commandLength:])
	if requested == "" {
		return "", true, nil
	}
	if requested[0] == '\'' || requested[0] == '"' {
		quote := requested[0]
		if len(requested) < 2 || requested[len(requested)-1] != quote {
			return "", true, fmt.Errorf("working directory has an unclosed quote")
		}
		requested = requested[1 : len(requested)-1]
	}
	return requested, true, nil
}

func isPrintDirectoryCommand(input string) bool {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "pwd", ":pwd":
		return true
	default:
		return false
	}
}
