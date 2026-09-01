package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func parseMakeDirectory(input string) (requested string, handled bool, err error) {
	trimmed := strings.TrimSpace(input)
	lower := strings.ToLower(trimmed)
	if lower == ":mkdir" {
		return "", true, fmt.Errorf("usage: :mkdir <path>")
	}
	if !strings.HasPrefix(lower, ":mkdir ") {
		return "", false, nil
	}
	requested = strings.TrimSpace(trimmed[len(":mkdir"):])
	if requested == "" {
		return "", true, fmt.Errorf("usage: :mkdir <path>")
	}
	if requested[0] == '\'' || requested[0] == '"' {
		quote := requested[0]
		if len(requested) < 2 || requested[len(requested)-1] != quote {
			return "", true, fmt.Errorf("directory path has an unclosed quote")
		}
		requested = requested[1 : len(requested)-1]
	}
	return requested, true, nil
}

func makeDirectory(current, requested string) (string, error) {
	requested = strings.TrimSpace(requested)
	if requested == "" {
		return "", fmt.Errorf("usage: :mkdir <path>")
	}
	if requested == "~" || strings.HasPrefix(requested, "~/") || strings.HasPrefix(requested, `~\`) {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("resolve home directory: %w", err)
		}
		if requested == "~" {
			requested = home
		} else {
			requested = filepath.Join(home, requested[2:])
		}
	}
	if !filepath.IsAbs(requested) {
		requested = filepath.Join(current, requested)
	}
	resolved, err := filepath.Abs(requested)
	if err != nil {
		return "", fmt.Errorf("resolve directory: %w", err)
	}
	if err := os.MkdirAll(filepath.Clean(resolved), 0o755); err != nil {
		return "", fmt.Errorf("create directory %q: %w", resolved, err)
	}
	return filepath.Clean(resolved), nil
}
