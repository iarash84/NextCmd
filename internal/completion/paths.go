package completion

import (
	"os"
	"path/filepath"
	"strings"

	"nextcmd/sdk"
)

func completeBuiltinPath(input, directory string) ([]sdk.Suggestion, bool) {
	trimmed := strings.TrimSpace(input)
	lower := strings.ToLower(trimmed)
	command := ""
	for _, candidate := range []string{":cd", "cd", ":ls", ":mkdir", ":del", ":trash"} {
		if lower == candidate {
			return nil, true
		}
		if strings.HasPrefix(lower, candidate+" ") {
			command = candidate
			break
		}
	}
	if command == "" {
		return nil, false
	}

	remainder := strings.TrimSpace(trimmed[len(command):])
	flags := []string{}
	if command == ":del" {
		for {
			switch {
			case strings.HasPrefix(remainder, "--dry-run "):
				flags = append(flags, "--dry-run")
				remainder = strings.TrimSpace(strings.TrimPrefix(remainder, "--dry-run"))
			case strings.HasPrefix(remainder, "--permanent "):
				flags = append(flags, "--permanent")
				remainder = strings.TrimSpace(strings.TrimPrefix(remainder, "--permanent"))
			default:
				goto flagsDone
			}
		}
	}

flagsDone:
	remainder = strings.Trim(remainder, `"'`)
	parent, prefix := splitPathPrefix(directory, remainder)
	entries, err := os.ReadDir(parent)
	if err != nil {
		return nil, true
	}
	dirsOnly := command == "cd" || command == ":cd" || command == ":mkdir"
	items := []sdk.Suggestion{}
	for _, entry := range entries {
		if dirsOnly && !entry.IsDir() {
			continue
		}
		if !strings.HasPrefix(strings.ToLower(entry.Name()), strings.ToLower(prefix)) {
			continue
		}
		suggestedPath := filepath.Join(filepath.Dir(remainder), entry.Name())
		if remainder == "" || filepath.Dir(remainder) == "." {
			suggestedPath = entry.Name()
		}
		if entry.IsDir() && (command == ":ls" || command == ":del" || command == ":trash") {
			suggestedPath += string(os.PathSeparator)
		}
		args := append([]string{}, flags...)
		args = append(args, suggestedPath)
		items = append(items, sdk.Suggestion{
			Command:     sdk.Command{Executable: command, Args: args},
			Title:       "Complete path",
			Description: "Complete a path from the active working directory",
			Reason:      "Built-in NextCmd path completion",
			Kind:        sdk.Completion,
			Priority:    100,
			Relevance:   100,
			Risk:        sdk.Safe,
			Source:      "nextcmd",
		})
	}
	return items, true
}

func splitPathPrefix(directory, requested string) (string, string) {
	if requested == "" {
		return directory, ""
	}
	if requested == "~" || strings.HasPrefix(requested, "~/") || strings.HasPrefix(requested, `~\`) {
		home, err := os.UserHomeDir()
		if err == nil {
			if requested == "~" {
				return home, ""
			}
			requested = filepath.Join(home, requested[2:])
		}
	}
	parent := filepath.Dir(requested)
	prefix := filepath.Base(requested)
	if parent == "." {
		return directory, prefix
	}
	if !filepath.IsAbs(parent) {
		parent = filepath.Join(directory, parent)
	}
	if abs, err := filepath.Abs(parent); err == nil {
		parent = abs
	}
	return filepath.Clean(parent), prefix
}
