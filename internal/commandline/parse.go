package commandline

import (
	"fmt"
	"strings"
	"unicode"

	"nextcmd/sdk"
)

// Parse converts an editable command line to a structured command without invoking a shell.
func Parse(input, directory string) (sdk.Command, error) {
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
