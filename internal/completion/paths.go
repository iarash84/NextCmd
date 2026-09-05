package completion

import (
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"nextcmd/sdk"
)

func completeCommandPath(input, directory string) []sdk.Suggestion {
	trimmed := strings.TrimLeftFunc(input, unicode.IsSpace)
	if strings.HasPrefix(trimmed, "!") || strings.HasPrefix(trimmed, ":") {
		return nil
	}
	completed, current, ok := splitCompletionInput(input)
	if !ok || len(completed) == 0 || strings.HasPrefix(current, "-") {
		return nil
	}
	parent, prefix := splitPathPrefix(directory, current)
	entries, err := os.ReadDir(parent)
	if err != nil {
		return nil
	}
	items := make([]sdk.Suggestion, 0)
	for _, entry := range entries {
		if !strings.HasPrefix(strings.ToLower(entry.Name()), strings.ToLower(prefix)) {
			continue
		}
		suggestedPath := completePathBase(current, entry.Name())
		kind := "file"
		if entry.IsDir() {
			suggestedPath += string(os.PathSeparator)
			kind = "directory"
		}
		args := append([]string(nil), completed[1:]...)
		args = append(args, suggestedPath)
		items = append(items, sdk.Suggestion{
			Command:     sdk.Command{Executable: completed[0], Args: args},
			Title:       "Complete path",
			Description: "Complete a " + kind + " path from the active working directory",
			Reason:      "Matches a local " + kind,
			Kind:        sdk.Completion,
			Priority:    30,
			Relevance:   70,
			Risk:        sdk.Safe,
			Source:      "nextcmd",
		})
	}
	return items
}

func completePathBase(current, name string) string {
	if current == "~" {
		return "~" + string(os.PathSeparator) + name
	}
	if separator := strings.LastIndexAny(current, `/\`); separator >= 0 {
		return current[:separator+1] + name
	}
	return name
}

// splitCompletionInput returns complete tokens separately from the token being
// edited. Unlike command execution parsing, an open quote is valid here.
func splitCompletionInput(input string) ([]string, string, bool) {
	var tokens []string
	var current strings.Builder
	var quote rune
	escaped := false
	tokenStarted := false
	endedWithSpace := false
	flush := func() {
		if tokenStarted {
			tokens = append(tokens, current.String())
			current.Reset()
			tokenStarted = false
		}
	}
	for _, char := range input {
		if escaped {
			current.WriteRune(char)
			tokenStarted = true
			escaped = false
			endedWithSpace = false
			continue
		}
		if char == '\\' && quote != '\'' && os.PathSeparator != '\\' {
			escaped = true
			tokenStarted = true
			endedWithSpace = false
			continue
		}
		if quote != 0 {
			if char == quote {
				quote = 0
			} else {
				current.WriteRune(char)
			}
			tokenStarted = true
			endedWithSpace = false
			continue
		}
		if char == '\'' || char == '"' {
			quote = char
			tokenStarted = true
			endedWithSpace = false
			continue
		}
		if unicode.IsSpace(char) {
			flush()
			endedWithSpace = true
			continue
		}
		current.WriteRune(char)
		tokenStarted = true
		endedWithSpace = false
	}
	if escaped {
		current.WriteRune('\\')
		tokenStarted = true
	}
	flush()
	if endedWithSpace {
		return tokens, "", len(tokens) > 0
	}
	if len(tokens) < 2 {
		return nil, "", false
	}
	return tokens[:len(tokens)-1], tokens[len(tokens)-1], true
}

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
			case strings.HasPrefix(remainder, "--yes "):
				flags = append(flags, "--yes")
				remainder = strings.TrimSpace(strings.TrimPrefix(remainder, "--yes"))
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
			if requested == "~" || requested == "~/" || requested == `~\` {
				return home, ""
			}
			requested = filepath.Join(home, requested[2:])
		}
	}
	if strings.HasSuffix(requested, string(os.PathSeparator)) {
		parent := requested
		if !filepath.IsAbs(parent) {
			parent = filepath.Join(directory, parent)
		}
		return filepath.Clean(parent), ""
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
