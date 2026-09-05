package sdk

import (
	"strconv"
	"strings"
)

// Command is a shell-independent process invocation.
type Command struct {
	Executable       string            `json:"executable"`
	Args             []string          `json:"args,omitempty"`
	ShellCommand     string            `json:"shellCommand,omitempty"`
	WorkingDirectory string            `json:"workingDirectory,omitempty"`
	Environment      map[string]string `json:"environment,omitempty"`
}

func (c Command) Display() string {
	if c.ShellCommand != "" {
		return "! " + c.ShellCommand
	}
	parts := append([]string{quote(c.Executable)}, mapStrings(c.Args, quote)...)
	return strings.Join(parts, " ")
}

func quote(s string) string {
	if s != "" && !strings.ContainsAny(s, " \t\n\r\"'") {
		return s
	}
	return strconv.Quote(s)
}

func mapStrings(in []string, f func(string) string) []string {
	out := make([]string, len(in))
	for i, value := range in {
		out[i] = f(value)
	}
	return out
}
