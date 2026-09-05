package commandline

import (
	"fmt"
	"strings"
	"unicode"

	"nextcmd/sdk"
)

// Parse converts an editable command line to a structured command. A leading !
// explicitly opts into the platform shell; every other command remains a direct
// process invocation.
func Parse(input, directory string) (sdk.Command, error) {
	trimmed := strings.TrimSpace(input)
	if strings.HasPrefix(trimmed, "!") {
		shellCommand := strings.TrimSpace(strings.TrimPrefix(trimmed, "!"))
		if shellCommand == "" {
			return sdk.Command{}, fmt.Errorf("empty shell command")
		}
		return sdk.Command{ShellCommand: shellCommand, WorkingDirectory: directory}, nil
	}
	var args []string
	var current strings.Builder
	var quote rune
	escaped := false
	flush := func() {
		if current.Len() > 0 {
			args = append(args, current.String())
			current.Reset()
		}
	}
	for _, r := range input {
		if escaped {
			current.WriteRune(r)
			escaped = false
			continue
		}
		if r == '\\' && quote != '\'' {
			escaped = true
			continue
		}
		if quote != 0 {
			if r == quote {
				quote = 0
			} else {
				current.WriteRune(r)
			}
			continue
		}
		if r == '\'' || r == '"' {
			quote = r
			continue
		}
		if unicode.IsSpace(r) {
			flush()
			continue
		}
		current.WriteRune(r)
	}
	if escaped {
		current.WriteRune('\\')
	}
	if quote != 0 {
		return sdk.Command{}, fmt.Errorf("unclosed quote")
	}
	flush()
	if len(args) == 0 {
		return sdk.Command{}, fmt.Errorf("empty command")
	}
	return sdk.Command{Executable: args[0], Args: args[1:], WorkingDirectory: directory}, nil
}
